package generators

import (
	"container/list"
)

type Info struct {
	Name    string
	Group   string
	History *list.List
}

type Generator interface {
	GetInfo() *Info
	GetValue(elapsed uint64, releasedAt uint64) float32
}

type Instrument interface {
	Generator
}

func GetDefaultInstrument() Instrument {
	return Mega()
}
