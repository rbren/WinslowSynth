package generators

import (
	"fmt"
)

type Sum struct {
	Generators []Generator
}

func setListFrequency(gs []Generator, freq float32) []Generator {
	set := make([]Generator, len(gs))
	for idx, g := range gs {
		if i, ok := g.(Instrument); ok {
			set[idx] = i.SetFrequency(freq)
		} else {
			set[idx] = g
		}
	}
	return set
}

func (s Sum) GetValue(t, releasedAt uint64) float32 {
	var val float32 = 0.0
	for _, gen := range s.Generators {
		val += gen.GetValue(t, releasedAt)
	}
	return val
}

func (s Sum) SetFrequency(freq float32) Instrument {
	s.Generators = setListFrequency(s.Generators, freq)
	return s
}

type Multiply struct {
	Generators []Generator
}

func (m Multiply) GetValue(t, releasedAt uint64) float32 {
	var val float32 = 1.0
	for _, gen := range m.Generators {
		val *= gen.GetValue(t, releasedAt)
	}
	return val
}

func (m Multiply) SetFrequency(freq float32) Instrument {
	m.Generators = setListFrequency(m.Generators, freq)
	fmt.Println("set generators", m.Generators)
	return m
}
