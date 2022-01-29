package music

import (
	"math"
)

func GenerateFrequency(hz float32, sampleRate int, numSamples int, startSample uint64) []float32 {
	output := make([]float32, numSamples)
	samplesPerPeriod := int(float32(sampleRate) / hz)
	for idx := range output {
		sampleLoc := int((uint64(idx) + startSample) % uint64(samplesPerPeriod))
		pos := 2.0 * math.Pi * (float32(sampleLoc) / float32(samplesPerPeriod)) // pos is in [0, 2pi]
		output[idx] = float32(math.Sin(float64(pos)))
	}
	return output
}
