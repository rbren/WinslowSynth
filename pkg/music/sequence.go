package music

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/thoas/go-funk"

	"github.com/rbren/midi/pkg/buffers"
	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/generators"
	"github.com/rbren/midi/pkg/input"
)

const maxSampleRateHandicap = 0.9
const sampleRateHandicapJump = .01
const sampleRateRatioMin = .3
const sampleRateRatioMax = .7

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
	start := time.Now()
	samplesPerMs := config.MainConfig.SampleRate / 1000
	samplesPerSprint := numSamples
	msPerSprint := samplesPerSprint / samplesPerMs

	//logrus.Infof("%d generators", len(s.Events))
	allSamples := make([][]float32, len(s.Events))
	var wg sync.WaitGroup
	for eventIdx, event := range s.Events {
		wg.Add(1)
		go func(eventIdx int, event *Event) {
			defer wg.Done()
			eventSamples := event.GetSamples(absoluteTime, numSamples, s.SampleRateHandicap)
			allSamples[eventIdx] = eventSamples
		}(eventIdx, event)
	}
	wg.Wait()
	var output []float32
	if len(allSamples) == 0 {
		output = make([]float32, numSamples)
	} else {
		output = buffers.MixBuffers(allSamples, buffers.NaiveSumMix, .5)
	}
	generators.AddHistory(s.Instrument, absoluteTime, output)
	duration := time.Since(start)
	ratio := float32(duration.Milliseconds()) / float32(msPerSprint)
	s.AdjustSampleRateHandicap(ratio)
	return output
}

func (s *Sequence) AdjustSampleRateHandicap(ratio float32) {
	if ratio > sampleRateRatioMax {
		if s.SampleRateHandicap >= maxSampleRateHandicap {
			logrus.Errorf("FULLY DOWNSAMPLED: with %d generators, ratio was %f", len(s.Events), ratio)
		}
		s.SampleRateHandicap = float32(math.Min(maxSampleRateHandicap, float64(s.SampleRateHandicap+sampleRateHandicapJump)))
	} else if s.SampleRateHandicap != 0.0 && ratio < sampleRateRatioMin {
		s.SampleRateHandicap = float32(math.Max(0, float64(s.SampleRateHandicap-sampleRateHandicapJump)))
	}
	if rand.Float32() < sampleRateHandicapJump*5 {
		logrus.Infof("with %d generators, ratio was %f, sample handicap is %f", len(s.Events), ratio, s.SampleRateHandicap)
	}
}
