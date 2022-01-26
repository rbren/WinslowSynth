package music

import (
	"math"
)

func GenerateFrequency(hz float64, sampleRate int, samples int) []byte {
	out := make([]byte, samples)
	samplesPerPeriod := int(float64(sampleRate) / hz)
	period := getSinePeriod(samplesPerPeriod)
	for idx := range out {
		out[idx] = period[idx%len(period)]
	}
	return out
}

func getSinePeriod(numSamples int) []byte {
	samples := make([]byte, numSamples)

	for idx := range samples {
		pos := 2.0 * math.Pi * (float64(idx) / float64(numSamples)) // pos is in [0, 2pi]
		val := math.Sin(pos)
		samples[idx] = byte(int(255 * val))
	}
	return samples
}
