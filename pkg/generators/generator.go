package generators

type Generator interface {
	GetValue(t uint64) float32
}

func GetSamples(g Generator, samples []float32, t uint64) {
	for idx := range samples {
		samples[idx] = g.GetValue(t + uint64(idx))
	}
}
