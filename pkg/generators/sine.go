package generators

import (
	"math"

	"github.com/rbren/midi/pkg/config"
)

type SineWave struct {
	Frequency float32
}

func (s SineWave) GetValue(time, releasedAt uint64) float32 {
	if releasedAt != 0 {
		return 0.0
	}
	samplesPerPeriod := float32(config.MainConfig.SampleRate) / s.Frequency
	sampleLoc := int(time % uint64(samplesPerPeriod))
	pos := 2.0 * math.Pi * (float32(sampleLoc) / samplesPerPeriod) // pos is in [0, 2pi]
	return float32(math.Sin(float64(pos)))
}
