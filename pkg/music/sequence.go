package music

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"

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
	lock          *sync.Mutex
	Events        []*Event
	LastFrequency float32
}

func NewSequence() Sequence {
	return Sequence{
		Events: []*Event{},
		lock:   &sync.Mutex{},
	}
}

func (s *Sequence) Add(note input.InputKey, time uint64) {
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

func (s *Sequence) attack(key input.InputKey, time uint64) {
	logrus.Infof("attack %d %d", key.Key, time)
	s.release(key, time)
	s.LastFrequency = key.Frequency
	s.Events = append(s.Events, &Event{
		AttackTime:  time,
		ReleaseTime: 0,
		Frequency:   key.Frequency,
		Key:         key.Key,
		Velocity:    key.Velocity,
	})
}

func (s *Sequence) release(key input.InputKey, time uint64) {
	logrus.Infof("release %d %d", key.Key, time)
	for _, evt := range s.Events {
		if evt.Key == key.Key && evt.ReleaseTime == 0 {
			evt.ReleaseTime = time
		}
	}
}

func (s *Sequence) ClearOldEvents(absoluteTime uint64) {
	s.Events = funk.Filter(s.Events, func(event *Event) bool {
		if event.ReleaseTime == 0 {
			return true
		}
		elapsedSinceRelease := absoluteTime - event.ReleaseTime
		return elapsedSinceRelease <= uint64(maxReleaseTimeSamples)
	}).([]*Event)
}

func (s *Sequence) GetSamples(inst generators.Instrument, absoluteTime uint64, numSamples int) []float32 {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.ClearOldEvents(absoluteTime) // TODO: put this on its own loop
	allSamples := [][]float32{}
	for _, event := range s.Events {
		eventSamples := make([]float32, numSamples)
		withFreq := generators.SetFrequency(inst, event.Frequency)
		if withFreq.GetInfo() != nil {
			//withFreq.SetInfo(generators.Info{}) // Make sure not to overwrite the history of inst
		}
		t, r := event.getRelativeTime(absoluteTime)
		for idx := range eventSamples {
			eventSamples[idx] = generators.GetValue(withFreq, t+uint64(idx), r)
		}
		allSamples = append(allSamples, eventSamples)
	}
	var output []float32
	if len(allSamples) == 0 {
		output = make([]float32, numSamples)
	} else {
		output = buffers.MixBuffers(allSamples)
	}
	generators.AddHistory(inst, absoluteTime, output)
	return output
}
