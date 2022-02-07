package music

import (
	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/generators"
)

const maxReleaseTimeMs = 10000
const minZeroedTimeMs = 1000

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
