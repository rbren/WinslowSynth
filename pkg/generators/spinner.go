package generators

import (
	"math"
)

type Spinner struct {
	Info          *Info
	Amplitude     Generator
	Frequency     Generator
	Phase         Generator
	Bias          Generator
	DropOnRelease bool
}

func (s *Spinner) initialize() {
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

func (s Spinner) GetValue(time, releasedAt uint64) float32 {
	if s.DropOnRelease && releasedAt != 0 {
		return 0.0
	}
	s.initialize()
	pos := 2.0 * math.Pi * GetPhasePosition(s.Frequency, s.Phase, time, releasedAt)
	amp := s.Amplitude.GetValue(time, releasedAt)
	return s.Bias.GetValue(time, releasedAt) + amp*float32(math.Sin(float64(pos)))
}

func (s Spinner) SetBias(bias float32) Generator {
	ret := s
	if ret.Bias == nil {
		ret.Bias = Constant{Value: bias}
	} else if b, ok := ret.Bias.(Spinner); ok {
		ret.Bias = b.SetBias(bias)
	} else {
		panic("can't set bias")
	}
	return ret
}

func (s Spinner) GetInfo() *Info { return s.Info }
