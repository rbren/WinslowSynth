package generators

import ()

type Info struct {
	Name            string
	Group           string
	History         []float32
	HistoryPosition int
}

type Generator interface {
	GetInfo() *Info
	GetValue(elapsed uint64, releasedAt uint64) float32
}

type Instrument interface {
	Generator
}

func GetDefaultInstrument() Instrument {
	return BasicSine()
}
