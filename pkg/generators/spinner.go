package generators

import (
	"math"

	"github.com/rbren/midi/pkg/config"
)

// A * exp(i * (wt + phi)) = Amplitude * exp(i * (Frequency + Phase))
type Spinner struct {
	Amplitude float32
	Phase     float32
	Frequency float32
}

func (s Spinner) GetValue(time uint64) float32 {
	samplesPerPeriod := float32(config.MainConfig.SampleRate) / s.Frequency
	sampleLoc := int(time % uint64(samplesPerPeriod))
	pos := 2.0 * math.Pi * (float32(sampleLoc) / samplesPerPeriod) // pos is in [0, 2pi]
	return float32(math.Sin(float64(pos)))
}

func (s Spinner) Multiply(mul Spinner) Spinner {
	// A exp(i(wt+p)) * B exp(i(vt+q)) = ABexp(i((w+v)t+p+q))
	out := Spinner{}
	out.Amplitude = s.Amplitude * mul.Amplitude
	out.Phase = s.Phase + mul.Phase
	out.Frequency = s.Frequency + mul.Frequency
	return out
}
