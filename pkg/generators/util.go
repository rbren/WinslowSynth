package generators

import (
	"math"

	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/logger"
)

// GetPhasePosition returns the current position as a fraction of a full period
func GetPhasePosition(freq Generator, phase Generator, time, releasedAt uint64) float32 {
	samplesPerPeriod := float32(config.MainConfig.SampleRate) / freq.GetValue(time, releasedAt)
	phaseVal := phase.GetValue(time, releasedAt)
	phaseScaled := (samplesPerPeriod * phaseVal) / (2.0 * math.Pi)
	sampleLoc := int((time + uint64(phaseScaled)) % uint64(samplesPerPeriod))
	return float32(sampleLoc) / samplesPerPeriod
}

func SetFrequency(inst Generator, freq float32) Generator {
	if inst == nil {
		return Constant{freq}
	}
	if c, ok := inst.(Constant); ok {
		c.Value = freq
		return c
	}
	if s, ok := inst.(Spinner); ok {
		s.Frequency = SetFrequency(s.Frequency, freq)
		return s
	}
	if s, ok := inst.(SawWave); ok {
		s.Frequency = SetFrequency(s.Frequency, freq)
		return s
	}
	if s, ok := inst.(SquareWave); ok {
		s.Frequency = SetFrequency(s.Frequency, freq)
		return s
	}
	logger.ForceLog("No ability to set frequency for", inst)
	return nil
}
