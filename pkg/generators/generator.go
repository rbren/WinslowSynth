package generators

type Generator interface {
	GetValue(elapsed uint64, releasedAt uint64) float32
}

type Instrument interface {
	Generator
}

func GetDefaultInstrument() Instrument {
	return Mega()
}
