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

func (s *Oscillator) initialize() {
	if s.Amplitude == nil {
		s.Amplitude = Constant{Value: 1.0}
	}
	if s.Phase == nil {
		s.Phase = Constant{Value: 0.0}
	}
	if s.Bias == nil {
		s.Bias = Constant{Value: 0.0}
	}
}

func (s Oscillator) GetValue(t, r uint64) float32 {
	if s.DropOnRelease && r != 0 {
		return 0.0
	}
	s.initialize()
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
	phaseVal := s.Phase.GetValue(time, releasedAt)
	phaseScaled := (samplesPerPeriod * phaseVal) / (2.0 * math.Pi)
	sampleLoc := int((time + uint64(phaseScaled)) % uint64(samplesPerPeriod))
	return float32(sampleLoc) / samplesPerPeriod
}

func (s Oscillator) GetInfo() Info { return s.Info }
func (s Oscillator) Copy(historyLen int) Generator {
	s.Info = s.Info.Copy(historyLen)
	return s
}
