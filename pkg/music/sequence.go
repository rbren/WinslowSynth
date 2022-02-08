package music

import (
	"math"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"

	"github.com/rbren/midi/pkg/buffers"
	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/generators"
	"github.com/rbren/midi/pkg/input"
)

type Sequence struct {
	lock               *sync.Mutex
	Instrument         generators.Generator
	Events             []*Event
	SampleRateHandicap float32
	LastFrequency      float32
}

func NewSequence() Sequence {
	return Sequence{
		Instrument: generators.GetDefaultInstrument(),
		Events:     []*Event{},
		lock:       &sync.Mutex{},
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
		Generator:   generators.SetFrequency(s.Instrument, key.Frequency),
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
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Events = funk.Filter(s.Events, func(event *Event) bool {
		return event.StillActive(absoluteTime)
	}).([]*Event)
}

func (s *Sequence) GetSamples(absoluteTime uint64, numSamples int) []float32 {
	s.lock.Lock()
	defer s.lock.Unlock()
	start := time.Now()
	samplesPerMs := config.MainConfig.SampleRate / 1000
	samplesPerSprint := numSamples
	msPerSprint := samplesPerSprint / samplesPerMs
	handicapModulus := int(math.Max(1.0, math.Ceil(float64(s.SampleRateHandicap))))

	//logrus.Infof("%d generators", len(s.Events))
	allSamples := make([][]float32, len(s.Events))
	var wg sync.WaitGroup
	for eventIdx, event := range s.Events {
		wg.Add(1)
		go func(eventIdx int, event *Event) {
			defer wg.Done()
			eventSamples := event.GetSamples(absoluteTime, numSamples, handicapModulus)
			allSamples[eventIdx] = eventSamples
		}(eventIdx, event)
	}
	wg.Wait()
	var output []float32
	if len(allSamples) == 0 {
		output = make([]float32, numSamples)
	} else {
		output = buffers.MixBuffers(allSamples)
	}
	generators.AddHistory(s.Instrument, absoluteTime, output)
	duration := time.Since(start)
	ratio := float32(duration.Milliseconds()) / float32(msPerSprint)
	if ratio > 1 {
		s.SampleRateHandicap = float32(math.Max(float64(s.SampleRateHandicap), 1.0)) + 1
		logrus.Warningf("DOWNSAMPLE: with %d generators, ratio was %f, increased sample handicap to %f", len(s.Events), ratio, s.SampleRateHandicap)
	} else if ratio < .25 && s.SampleRateHandicap >= 1 {
		s.SampleRateHandicap = s.SampleRateHandicap - 1
		logrus.Warningf("UPSAMPLE: with %d generators, ratio was %f, decreased sample handicap to %f", len(s.Events), ratio, s.SampleRateHandicap)
	}
	return output
}
