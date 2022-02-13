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
	Shape         OscillatorShape
	DropOnRelease bool
	SubGenerators SubGenerators
}

func (s Oscillator) GetInfo() Info                   { return s.Info }
func (s Oscillator) GetSubGenerators() SubGenerators { return s.SubGenerators }
func (s Oscillator) Copy(historyLen int) Generator {
	s.Info = s.Info.Copy(historyLen)
	s.SubGenerators = s.SubGenerators.Copy()
	return s
}

func (s Oscillator) Initialize(group string) Generator {
	if s.SubGenerators == nil {
		s.SubGenerators = make(map[string]Generator)
	}
	if s.SubGenerators["Frequency"] == nil {
		s.SubGenerators["Frequency"] = Constant{
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
	if s.SubGenerators["Amplitude"] == nil {
		s.SubGenerators["Amplitude"] = Constant{
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
	if s.SubGenerators["Phase"] == nil {
		s.SubGenerators["Phase"] = Constant{
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
	if s.SubGenerators["Bias"] == nil {
		s.SubGenerators["Bias"] = Constant{
			Info: Info{
				Group:    group,
				Subgroup: "Oscillator",
			},
			Value: 0.0,
		}
	}
	for key, g := range s.SubGenerators {
		s.SubGenerators[key] = g.Initialize(group)
	}
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
	amp := GetValue(s.SubGenerators["Amplitude"], t, r)
	return GetValue(s.SubGenerators["Bias"], t, r) + amp*float32(math.Sin(float64(pos)))
}

func (s Oscillator) GetSaw(t, r uint64) float32 {
	fraction := s.GetPhasePosition(t, r)
	return GetValue(s.SubGenerators["Amplitude"], t, r) * (fraction*2.0 - 1.0)
}

func (s Oscillator) GetSquare(t, r uint64) float32 {
	phasePos := s.GetPhasePosition(t, r)
	var val float32 = 1.0
	if phasePos > .5 {
		val = -1.0
	}
	return val * GetValue(s.SubGenerators["Amplitude"], t, r)
}

// GetPhasePosition returns the current position as a fraction of a full period
func (s Oscillator) GetPhasePosition(time, releasedAt uint64) float32 {
	samplesPerPeriod := float32(config.MainConfig.SampleRate) / s.SubGenerators["Frequency"].GetValue(time, releasedAt)
	phase := s.SubGenerators["Phase"].GetValue(time, releasedAt)
	sampleLoc := int((time + uint64(phase)) % uint64(samplesPerPeriod))
	return float32(sampleLoc) / samplesPerPeriod
}
