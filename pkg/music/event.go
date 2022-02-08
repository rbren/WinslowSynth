package music

import (
	"math"

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
	for idx := range eventSamples {
		if idx%handicapModulus == 0 || idx == numSamples-1 {
			val := generators.GetValue(e.Generator, t+uint64(idx), r)
			eventSamples[idx] = val
			if math.Abs(float64(val)) > zeroThreshold {
				zeroed = false
			}
		}
	}
	e.Zeroed = zeroed
	var prev, next float32
	for idx := range eventSamples {
		remainder := idx % handicapModulus
		if remainder == 0 {
			prev = eventSamples[idx]
			nextIdx := idx + handicapModulus
			if nextIdx >= len(eventSamples) {
				nextIdx = len(eventSamples) - 1
			}
			next = eventSamples[nextIdx]
		} else {
			weightNext := float32(remainder) / float32(handicapModulus)
			weightPrev := 1.0 - weightNext
			eventSamples[idx] = weightPrev*prev + weightNext*next
		}
	}
	return eventSamples
}
