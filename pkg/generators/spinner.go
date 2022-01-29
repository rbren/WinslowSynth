package generators

import (
	"math"

	"github.com/rbren/midi/pkg/config"
)

// A * exp(i * (wt + phi)) = Amplitude * exp(i * (Frequency + Phase))
type Spinner struct {
	Amplitude Generator
	Frequency Generator
	Phase     Generator
}

func NewSpinner(a, f, p float32) Spinner {
	return Spinner{
		Amplitude: Constant{Value: a},
		Frequency: Constant{Value: f},
		Phase:     Constant{Value: p},
	}
}

func (s Spinner) GetValue(time, releasedAt uint64) float32 {
	samplesPerPeriod := float32(config.MainConfig.SampleRate) / s.Frequency.GetValue(time, releasedAt)
	sampleLoc := int(time % uint64(samplesPerPeriod))
	pos := 2.0 * math.Pi * (float32(sampleLoc) / samplesPerPeriod) // pos is in [0, 2pi]
	amp := s.Amplitude.GetValue(time, releasedAt)
	return amp * float32(math.Sin(float64(pos)+float64(s.Phase.GetValue(time, releasedAt))))
}
