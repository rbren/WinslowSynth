package music

import (
	"math"
)

func GenerateFrequency(hz float64, sampleRate int, numSamples int) []float64 {
	samplesPerPeriod := int(float64(sampleRate) / hz)
	period := getSinePeriod(samplesPerPeriod)
	output := make([]float64, numSamples)
	for idx := range output {
		output[idx] = period[idx%len(period)]
	}
	return output
}

func getSinePeriod(numSamples int) []float64 {
	samples := make([]float64, numSamples)

	for idx := range samples {
		pos := 2.0 * math.Pi * (float64(idx) / float64(numSamples)) // pos is in [0, 2pi]
		samples[idx] = math.Sin(pos)
	}
	return samples
}
