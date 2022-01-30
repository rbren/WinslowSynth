package generators

type Generator interface {
	GetValue(elapsed uint64, releasedAt uint64) float32
}

type Instrument interface {
	Generator
	SetFrequency(freq float32) Instrument
}

func GetDefaultInstrument() Instrument {
	return Warbler()
}
