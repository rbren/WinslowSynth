package music

import (
	"sync"

	"github.com/sirupsen/logrus"

	"github.com/rbren/midi/pkg/buffers"
	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/generators"
	"github.com/rbren/midi/pkg/input"
)

var maxReleaseTimeSamples int

func init() {
	maxReleaseTimeSamples = config.MainConfig.SampleRate * 10
}

type Sequence struct {
	Events map[int64]*Event
	lock   *sync.Mutex
}

func NewSequence() Sequence {
	return Sequence{
		Events: map[int64]*Event{},
		lock:   &sync.Mutex{},
	}
}

func (s Sequence) Add(note input.InputKey, time uint64) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if note.Action == "channel.NoteOn" {
		s.attack(note, time)
	} else if note.Action == "channel.NoteOff" {
		s.release(note, time)
	} else {
		logrus.Info("No action for " + note.Action)
	}
}

func (s Sequence) attack(key input.InputKey, time uint64) {
	logrus.Infof("attack %d %d", key.Key, time)
	// TODO: allow more than one event per key simultaneously
	s.Events[key.Key] = &Event{
		AttackTime:  time,
		ReleaseTime: 0,
		Frequency:   key.Frequency,
		Key:         key.Key,
		Velocity:    key.Velocity,
	}
}

func (s Sequence) release(key input.InputKey, time uint64) {
	logrus.Infof("release %d %d", key.Key, time)
	existing, ok := s.Events[key.Key]
	if !ok {
		logrus.Error("Released key without attack!", key)
		return
	}
	existing.ReleaseTime = time
}

func (s Sequence) ClearOldEvents(absoluteTime uint64) {
	remove := []int64{}
	for key, event := range s.Events {
		if event.ReleaseTime == 0 {
			continue
		}
		elapsedSinceRelease := absoluteTime - event.ReleaseTime
		if elapsedSinceRelease > uint64(maxReleaseTimeSamples) {
			remove = append(remove, key)
		}
	}
	for _, key := range remove {
		delete(s.Events, key)
	}
}

func (s Sequence) GetSamples(inst generators.Instrument, absoluteTime uint64, numSamples int) []float32 {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.ClearOldEvents(absoluteTime) // TODO: put this on its own loop
	if len(s.Events) == 0 {
		return make([]float32, numSamples)
	}
	allSamples := [][]float32{}
	for _, event := range s.Events {
		logrus.Infof("process event %#v", event)
		eventSamples := make([]float32, numSamples)
		withFreq := generators.SetFrequency(inst, event.Frequency)
		t, r := event.getRelativeTime(absoluteTime)
		for idx := range eventSamples {
			eventSamples[idx] = generators.GetValue(withFreq, t+uint64(idx), r)
		}
		allSamples = append(allSamples, eventSamples)
	}
	return buffers.MixBuffers(allSamples)
}
