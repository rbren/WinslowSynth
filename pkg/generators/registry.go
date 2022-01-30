package generators

import (
	"github.com/rbren/midi/pkg/buffers"
	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/logger"
)

var maxReleaseTimeSamples int

func init() {
	maxReleaseTimeSamples = config.MainConfig.SampleRate * 5
}

type Registry struct {
	Events map[int64]*Event
}

func NewRegistry() Registry {
	return Registry{
		Events: map[int64]*Event{},
	}
}

type Event struct {
	Time        uint64
	Generator   Generator
	AttackTime  uint64
	ReleaseTime uint64
}

type EventType int

const (
	AttackEvent EventType = iota
	ReleaseEvent
)

func (r Registry) Attack(key int64, time uint64, g Generator) {
	logger.Log("attack", key)
	r.Events[key] = &Event{
		Generator:   g,
		AttackTime:  time,
		ReleaseTime: 0,
	}
}

func (r Registry) Release(key int64, time uint64, g Generator) {
	logger.Log("release", key)
	existing, ok := r.Events[key]
	if !ok {
		logger.ForceLog("Released key without attack!", key)
		return
	}
	existing.ReleaseTime = time
}

func (r Registry) ClearOldEvents(absoluteTime uint64) {
	remove := []int64{}
	for key, event := range r.Events {
		if event.ReleaseTime == 0 {
			continue
		}
		elapsedSinceRelease := absoluteTime - event.ReleaseTime
		if elapsedSinceRelease > uint64(maxReleaseTimeSamples) {
			remove = append(remove, key)
		}
	}
	for _, key := range remove {
		delete(r.Events, key)
	}
}

func (r Registry) GetSamples(absoluteTime uint64, numSamples int) []float32 {
	r.ClearOldEvents(absoluteTime)
	samples := make([]float32, numSamples)
	for _, event := range r.Events {
		eventSamples := make([]float32, numSamples)
		elapsed := absoluteTime - event.AttackTime
		var releasedAt uint64 = 0
		if event.ReleaseTime != 0 {
			elapsedSinceRelease := absoluteTime - event.ReleaseTime
			releasedAt = elapsed - elapsedSinceRelease
		}
		for idx := range eventSamples {
			eventSamples[idx] = event.Generator.GetValue(elapsed+uint64(idx), releasedAt)
		}
		samples = buffers.MixBuffers([][]float32{samples, eventSamples})
	}
	return samples
}
