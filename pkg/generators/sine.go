package generators

import (
	"math"

	"github.com/rbren/midi/pkg/config"
)

type SineWave struct {
	Frequency float32
}

func (s SineWave) GetValue(startSample uint64) float32 {
	samplesPerPeriod := int(float32(config.MainConfig.SampleRate) / s.Frequency)
	sampleLoc := int(startSample % uint64(samplesPerPeriod))
	pos := 2.0 * math.Pi * (float32(sampleLoc) / float32(samplesPerPeriod)) // pos is in [0, 2pi]
	return float32(math.Sin(float64(pos)))
}
