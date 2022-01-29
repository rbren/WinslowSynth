package generators

import (
	"math"

	"github.com/rbren/midi/pkg/config"
)

func GenerateSine(hz float32, numSamples int, startSample uint64) []float32 {
	output := make([]float32, numSamples)
	samplesPerPeriod := int(float32(config.MainConfig.SampleRate) / hz)
	for idx := range output {
		sampleLoc := int((uint64(idx) + startSample) % uint64(samplesPerPeriod))
		pos := 2.0 * math.Pi * (float32(sampleLoc) / float32(samplesPerPeriod)) // pos is in [0, 2pi]
		output[idx] = float32(math.Sin(float64(pos)))
	}
	return output
}
