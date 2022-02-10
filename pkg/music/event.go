package music

import (
	"math"
	"math/rand"

	"github.com/sirupsen/logrus"

	"github.com/rbren/midi/pkg/buffers"
	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/generators"
)

const maxReleaseTimeMs = 10000
const minZeroedTimeMs = 1000
const zeroThreshold = 0.01

var maxReleaseTimeSamples int
var minZeroedTimeSamples int

func init() {
	maxReleaseTimeSamples = maxReleaseTimeMs * config.MainConfig.SampleRate / 1000
	minZeroedTimeSamples = minZeroedTimeMs * config.MainConfig.SampleRate / 1000
}

type Event struct {
	Key         int64
	Frequency   float32
	Velocity    int64
	AttackTime  uint64
	ReleaseTime uint64
	Generator   generators.Generator
	Zeroed      bool
}

type EventType int

const (
	AttackEvent EventType = iota
	ReleaseEvent
)

func (e Event) getRelativeTime(absoluteTime uint64) (uint64, uint64) {
	elapsed := absoluteTime - e.AttackTime
	var releasedAt uint64 = 0
	if e.ReleaseTime != 0 {
		elapsedSinceRelease := absoluteTime - e.ReleaseTime
		releasedAt = elapsed - elapsedSinceRelease
	}
	return elapsed, releasedAt
}

func (e Event) StillActive(absoluteTime uint64) bool {
	if e.ReleaseTime == 0 {
		return true
	}
	elapsedSinceRelease := absoluteTime - e.ReleaseTime
	if e.Zeroed && elapsedSinceRelease > uint64(minZeroedTimeSamples) {
		return false
	}
	return elapsedSinceRelease <= uint64(maxReleaseTimeSamples)
}

func (e *Event) GetSamples(absoluteTime uint64, numSamples, handicapModulus int) []float32 {
	eventSamples := make([]float32, numSamples)
	t, r := e.getRelativeTime(absoluteTime)
	zeroed := true
	lastIdxCalculated := -1
	numCalculated := 0
	for idx := range eventSamples {
		if rand.Intn(handicapModulus) == 0 || idx == 0 || idx == numSamples-1 {
			val := generators.GetValue(e.Generator, t+uint64(idx), r)
			eventSamples[idx] = val
			if math.Abs(float64(val)) > zeroThreshold {
				zeroed = false
			}
			if lastIdxCalculated != -1 {
				buffers.InterpolateEvents(eventSamples, lastIdxCalculated, idx)
			}
			numCalculated++
			lastIdxCalculated = idx
		}
	}
	logrus.Infof("Calculated %.02f%% of  values", 100*float32(numCalculated)/float32(len(eventSamples)))
	e.Zeroed = zeroed
	return eventSamples
}
