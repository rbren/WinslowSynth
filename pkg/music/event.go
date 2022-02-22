package music

import (
	"math"
	"math/rand"

	"github.com/sirupsen/logrus"

	"github.com/rbren/midi/pkg/buffers"
	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/generators"
)

var maxReleaseTimeSamples int
var minZeroedTimeSamples int

func init() {
	maxReleaseTimeSamples = config.MainConfig.MaxReleaseTimeMs * config.MainConfig.SampleRate / 1000
	minZeroedTimeSamples = config.MainConfig.MinZeroedTimeMs * config.MainConfig.SampleRate / 1000
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

func (e *Event) GetSamples(absoluteTime uint64, numSamples int, handicap float32) []float32 {
	eventSamples := make([]float32, numSamples)
	t, r := e.getRelativeTime(absoluteTime)
	zeroed := true
	lastIdxCalculated := -1
	numCalculated := 0
	for idx := range eventSamples {
		if rand.Float32() > handicap || idx == 0 || idx == numSamples-1 {
			val := generators.GetValue(e.Generator, t+uint64(idx), r)
			eventSamples[idx] = val
			if math.Abs(float64(val)) > float64(config.MainConfig.ZeroSampleThreshold) {
				zeroed = false
			}
			if lastIdxCalculated != idx-1 {
				buffers.InterpolateEvents(eventSamples, lastIdxCalculated, idx)
			}
			numCalculated++
			lastIdxCalculated = idx
		}
	}
	logrus.Debugf("Calculated %.02f%% of  values", 100*float32(numCalculated)/float32(len(eventSamples)))
	e.Zeroed = zeroed
	return eventSamples
}
