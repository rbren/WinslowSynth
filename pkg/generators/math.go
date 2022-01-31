package generators

import ()

type Sum struct {
	Info       *Info
	Generators []Generator
}

type Multiply struct {
	Info       *Info
	Generators []Generator
}

func (s Sum) GetValue(t, releasedAt uint64) float32 {
	var val float32 = 0.0
	for _, gen := range s.Generators {
		val += gen.GetValue(t, releasedAt)
	}
	return val
}

func (m Multiply) GetValue(t, releasedAt uint64) float32 {
	var val float32 = 1.0
	for _, gen := range m.Generators {
		val *= gen.GetValue(t, releasedAt)
	}
	return val
}

func (s Sum) GetInfo() *Info { return s.Info }

func (m Multiply) GetInfo() *Info { return m.Info }
