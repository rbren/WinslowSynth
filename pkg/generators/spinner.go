package generators

import (
	"math"
)

type Spinner struct {
	Amplitude Generator
	Frequency Generator
	Phase     Generator
	Bias      Generator
}

func NewSpinner(a, f, p float32) Spinner {
	return Spinner{
		Amplitude: Constant{Value: a},
		Frequency: Constant{Value: f},
		Phase:     Constant{Value: p},
		Bias:      Constant{Value: 0.0},
	}
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
	s.initialize()
	pos := 2.0 * math.Pi * GetPhasePosition(s.Frequency, s.Phase, time, releasedAt)
	amp := s.Amplitude.GetValue(time, releasedAt)
	return s.Bias.GetValue(time, releasedAt) + amp*float32(math.Sin(float64(pos)))
}
