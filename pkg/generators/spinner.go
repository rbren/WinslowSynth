package generators

import (
	"math"
)

type Spinner struct {
	Amplitude Generator
	Frequency Generator
	Phase     Generator
	Value     float32
}

func NewSpinner(a, f, p float32) Spinner {
	return Spinner{
		Amplitude: Constant{Value: a},
		Frequency: Constant{Value: f},
		Phase:     Constant{Value: p},
	}
}

func (s Spinner) GetValue(time, releasedAt uint64) float32 {
	pos := 2.0 * math.Pi * GetPhasePosition(s.Frequency, s.Phase, time, releasedAt)
	amp := s.Amplitude.GetValue(time, releasedAt)
	return s.Value + amp*float32(math.Sin(float64(pos)))
}
