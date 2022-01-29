package generators

import (
	"math"

	"github.com/rbren/midi/pkg/config"
)

// GetPhasePosition returns the current position as a fraction of a full period
func GetPhasePosition(freq Generator, phase Generator, time, releasedAt uint64) float32 {
	samplesPerPeriod := float32(config.MainConfig.SampleRate) / freq.GetValue(time, releasedAt)
	phaseVal := phase.GetValue(time, releasedAt)
	phaseScaled := (samplesPerPeriod * phaseVal) / (2.0 * math.Pi)
	sampleLoc := int((time + uint64(phaseScaled)) % uint64(samplesPerPeriod))
	return float32(sampleLoc) / samplesPerPeriod
}
