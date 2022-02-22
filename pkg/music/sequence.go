package music

import (
	"fmt"
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
	eventsCopy := []*Event{}
	for _, event := range s.Events {
		// events may change during loop, but that's OK as long as we're
		// working with the same set for the duration of GetSamples
		eventsCopy = append(eventsCopy, event)
	}

	startLen := len(eventsCopy)
	//logrus.Infof("%d generators", len(eventsCopy))
	allSamples := make([][]float32, len(eventsCopy))
	var wg sync.WaitGroup
	for eventIdx, event := range eventsCopy {
		wg.Add(1)
		go func(eventIdx int, event *Event) {
			defer wg.Done()
			eventSamples := event.GetSamples(absoluteTime, numSamples, s.SampleRateHandicap)
			allSamples[eventIdx] = eventSamples
		}(eventIdx, event)
	}
	wg.Wait()
	if len(eventsCopy) != startLen {
		panic(fmt.Errorf("number of events changed from %d to %d during GetSamples", startLen, len(eventsCopy)))
	}
	var output []float32
	if len(allSamples) == 0 {
		output = make([]float32, numSamples)
	} else {
		output = buffers.MixBuffers(allSamples, buffers.NaiveSumMix, .5)
	}
	generators.AddHistory(s.Instrument, absoluteTime, output)
	duration := time.Since(start)
	s.AdjustSampleRateHandicap(numSamples, duration)
	return output
}

func (s *Sequence) AdjustSampleRateHandicap(numSamples int, duration time.Duration) {
	samplesPerMs := config.MainConfig.SampleRate / 1000
	samplesPerSprint := numSamples
	msPerSprint := samplesPerSprint / samplesPerMs
	durationMs := float32(duration.Microseconds()) / 1000
	ratio := durationMs / float32(msPerSprint)

	if ratio > sampleRateRatioMax {
		if s.SampleRateHandicap >= maxSampleRateHandicap {
			logrus.Errorf("FULLY DOWNSAMPLED: with %d generators, ratio was %f", len(s.Events), ratio)
		}
		s.SampleRateHandicap = float32(math.Min(maxSampleRateHandicap, float64(s.SampleRateHandicap+sampleRateHandicapJump)))
	} else if s.SampleRateHandicap != 0.0 && ratio < sampleRateRatioMin {
		s.SampleRateHandicap = float32(math.Max(0, float64(s.SampleRateHandicap-sampleRateHandicapJump)))
	}
	if rand.Float32() < sampleRateHandicapJump*5 {
		logrus.Infof("with %d generators, took %.02f ms to compute %d ms. Ratio was %.02f, sample handicap is %.02f",
			len(s.Events), durationMs, msPerSprint, ratio, s.SampleRateHandicap)
	}
}
