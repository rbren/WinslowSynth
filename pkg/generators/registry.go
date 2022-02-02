package generators

import (
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/rbren/midi/pkg/buffers"
	"github.com/rbren/midi/pkg/config"
)

var maxReleaseTimeSamples int
var historyLength int

func init() {
	historyLength = config.MainConfig.SampleRate // store 1 second
	maxReleaseTimeSamples = config.MainConfig.SampleRate * 3
}

type Registry struct {
	Events map[int64]*Event
	lock   sync.Mutex
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
	r.lock.Lock()
	defer r.lock.Unlock()
	logrus.Infof("attack %d %d", key, time)
	r.Events[key] = &Event{
		Generator:   g,
		AttackTime:  time,
		ReleaseTime: 0,
	}
}

func (r Registry) Release(key int64, time uint64, g Generator) {
	r.lock.Lock()
	defer r.lock.Unlock()
	logrus.Infof("release %d %d", key, time)
	existing, ok := r.Events[key]
	if !ok {
		logrus.Error("Released key without attack!", key)
		return
	}
	existing.ReleaseTime = time
}

func (r Registry) ClearOldEvents(absoluteTime uint64) {
	r.lock.Lock()
	defer r.lock.Unlock()
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
	r.ClearOldEvents(absoluteTime) // TODO: put this on its own loop
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
			addHistory(event.Generator, eventSamples[idx])
		}
		samples = buffers.MixBuffers([][]float32{samples, eventSamples})
	}
	return samples
}

func getEmptyHistory() []float32 {
	return make([]float32, historyLength)
}

func addHistory(g Generator, val float32) {
	i := g.GetInfo()
	if i == nil || i.History == nil {
		return
	}
	i.History[i.HistoryPosition] = val
	i.HistoryPosition = (i.HistoryPosition + 1) % historyLength
}
