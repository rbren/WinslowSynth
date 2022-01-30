package generators

import (
	"fmt"
	"math"
)

type Spinner struct {
	Amplitude     Generator
	Frequency     Generator
	Phase         Generator
	Bias          Generator
	DropOnRelease bool
}

func (s *Spinner) initialize() {
	if s.Amplitude == nil {
		s.Amplitude = Constant{Name: "Amplitude", Value: 1.0}
	}
	if s.Phase == nil {
		s.Phase = Constant{Name: "Phase", Value: 0.0}
	}
	if s.Bias == nil {
		s.Bias = Constant{Name: "Bias", Value: 0.0}
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

func (s Spinner) SetFrequency(freq float32) Instrument {
	ret := s
	if ret.Frequency == nil {
		ret.Frequency = Constant{Value: freq}
		return ret
	} else if c, ok := ret.Frequency.(Constant); ok {
		c.Value = freq
		ret.Frequency = c
	} else if f, ok := ret.Frequency.(Spinner); ok {
		ret.Frequency = f.SetBias(freq)
	} else if i, ok := ret.Frequency.(Instrument); ok {
		ret.Frequency = i.SetFrequency(freq)
	} else {
		fmt.Printf("tried to set frequency %#v\n", ret.Frequency)
		panic("can't set frequency")
	}
	return ret
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
