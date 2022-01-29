package generators

import (
	"github.com/rbren/midi/pkg/config"
)

type SawWave struct {
	Frequency float32
}

func (s SawWave) GetValue(time, releasedAt uint64) float32 {
	if releasedAt != 0 {
		return 0.0
	}
	samplesPerPeriod := float32(config.MainConfig.SampleRate) / s.Frequency
	sampleLoc := int(time % uint64(samplesPerPeriod))
	return float32(sampleLoc) / samplesPerPeriod
}
