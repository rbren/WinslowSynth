package generators

import (
	"fmt"
	"math"
)

type OscillatorShape int

const (
	WaveShape OscillatorShape = iota
	SawShape
	SquareShape
)

type Oscillator struct {
	Info          *Info
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
	pos := 2.0 * math.Pi * GetPhasePosition(s.Frequency, s.Phase, t, r)
	amp := s.Amplitude.GetValue(t, r)
	return s.Bias.GetValue(t, r) + amp*float32(math.Sin(float64(pos)))
}

func (s Oscillator) GetSaw(t, r uint64) float32 {
	return s.Amplitude.GetValue(t, r) * GetPhasePosition(s.Frequency, s.Phase, t, r)
}

func (s Oscillator) GetSquare(t, r uint64) float32 {
	phasePos := GetPhasePosition(s.Frequency, s.Phase, t, r)
	var val float32 = 1.0
	if phasePos > .5 {
		val = -1.0
	}
	return val * s.Amplitude.GetValue(t, r)
}

func (s Oscillator) GetInfo() *Info    { return s.Info }
func (s Oscillator) SetInfo(info Info) { copyInfo(s.Info, info) }
