package music

import (
	"math"
)

func GenerateFrequency(hz float32, sampleRate int, numSamples int) []float32 {
	samplesPerPeriod := int(float32(sampleRate) / hz)
	period := getSinePeriod(samplesPerPeriod)
	output := make([]float32, numSamples)
	for idx := range output {
		output[idx] = period[idx%len(period)]
	}
	return output
}

func getSinePeriod(numSamples int) []float32 {
	samples := make([]float32, numSamples)

	for idx := range samples {
		pos := 2.0 * math.Pi * (float32(idx) / float32(numSamples)) // pos is in [0, 2pi]
		samples[idx] = float32(math.Sin(float64(pos)))
	}
	return samples
}
