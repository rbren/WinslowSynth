package generators

import (
	"github.com/rbren/midi/pkg/config"
)

type SquareWave struct {
	Frequency float32
}

func (s SquareWave) GetValue(time, releasedAt uint64) float32 {
	if releasedAt != 0 {
		return 0.0
	}
	samplesPerPeriod := float32(config.MainConfig.SampleRate) / s.Frequency
	sampleLoc := int(time % uint64(samplesPerPeriod))
	midpoint := int(samplesPerPeriod) / 2
	if sampleLoc < midpoint {
		return 1.0
	}
	return -1.0
}
