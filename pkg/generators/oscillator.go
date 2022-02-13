package generators

import (
	"fmt"
	"math"

	"github.com/rbren/midi/pkg/config"
)

type OscillatorShape int

const (
	WaveShape OscillatorShape = iota
	SawShape
	SquareShape
)

type Oscillator struct {
	Info          Info
	Amplitude     Generator
	Frequency     Generator
	Phase         Generator
	Bias          Generator
	Shape         OscillatorShape
	DropOnRelease bool
}

func (s Oscillator) SubGenerators() []Generator {
	return []Generator{s.Amplitude, s.Frequency, s.Phase, s.Bias}
}

func (s Oscillator) Initialize(group string) Generator {
	if s.Frequency == nil {
		s.Frequency = Constant{
			Info: Info{
				Name:     "Frequency",
				Group:    group,
				Subgroup: "Oscillator",
			},
			Value: 440,
			Min:   20,
			Max:   20000,
		}
	}
	if s.Amplitude == nil {
		s.Amplitude = Constant{
			Info: Info{
				Name:     "Amplitude",
				Group:    group,
				Subgroup: "Oscillator",
			},
			Value: 1.0,
			Min:   0.0,
			Max:   1.0,
		}
	}
	if s.Phase == nil {
		s.Phase = Constant{
			Info: Info{
				Name:     "Phase",
				Group:    group,
				Subgroup: "Oscillator",
			},
			Value: 0.0,
			Min:   0.0,
			Max:   2.0 * math.Pi,
		}
	}
	if s.Bias == nil {
		s.Bias = Constant{
			Info: Info{
				Group:    group,
				Subgroup: "Oscillator",
			},
			Value: 0.0,
		}
	}
	s.Frequency = s.Frequency.Initialize(group)
	s.Amplitude = s.Amplitude.Initialize(group)
	s.Phase = s.Phase.Initialize(group)
	s.Bias = s.Bias.Initialize(group)
	return s
}

func (s Oscillator) GetValue(t, r uint64) float32 {
	if s.DropOnRelease && r != 0 {
		return 0.0
	}
	if s.Shape == WaveShape {
		return s.GetWave(t, r)
	}
	if s.Shape == SawShape {
		return s.GetSaw(t, r)
	}
	if s.Shape == SquareShape {
		return s.GetSquare(t, r)
	}
	panic(fmt.Errorf("Unknown shape %d", s.Shape))
	return 0.0
}

func (s Oscillator) GetWave(t, r uint64) float32 {
	pos := 2.0 * math.Pi * s.GetPhasePosition(t, r)
	amp := GetValue(s.Amplitude, t, r)
	return GetValue(s.Bias, t, r) + amp*float32(math.Sin(float64(pos)))
}

func (s Oscillator) GetSaw(t, r uint64) float32 {
	fraction := s.GetPhasePosition(t, r)
	return GetValue(s.Amplitude, t, r) * (fraction*2.0 - 1.0)
}

func (s Oscillator) GetSquare(t, r uint64) float32 {
	phasePos := s.GetPhasePosition(t, r)
	var val float32 = 1.0
	if phasePos > .5 {
		val = -1.0
	}
	return val * GetValue(s.Amplitude, t, r)
}

// GetPhasePosition returns the current position as a fraction of a full period
func (s Oscillator) GetPhasePosition(time, releasedAt uint64) float32 {
	samplesPerPeriod := float32(config.MainConfig.SampleRate) / s.Frequency.GetValue(time, releasedAt)
	phase := s.Phase.GetValue(time, releasedAt)
	sampleLoc := int((time + uint64(phase)) % uint64(samplesPerPeriod))
	return float32(sampleLoc) / samplesPerPeriod
}

func (s Oscillator) GetInfo() Info { return s.Info }
func (s Oscillator) Copy(historyLen int) Generator {
	s.Info = s.Info.Copy(historyLen)
	s.Amplitude = s.Amplitude.Copy(CopyExistingHistoryLength)
	s.Frequency = s.Frequency.Copy(CopyExistingHistoryLength)
	s.Phase = s.Phase.Copy(CopyExistingHistoryLength)
	s.Bias = s.Bias.Copy(CopyExistingHistoryLength)
	return s
}
