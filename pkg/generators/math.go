package generators

import ()

type Average struct {
	Info       Info
	Generators []Generator
}

type Multiply struct {
	Info       Info
	Generators []Generator
}

func (s Average) GetValue(t, releasedAt uint64) float32 {
	if len(s.Generators) == 0 {
		return 0.0
	}
	var val float32 = 0.0
	for _, gen := range s.Generators {
		val += GetValue(gen, t, releasedAt)
	}
	return val / float32(len(s.Generators))
}

func (m Multiply) GetValue(t, releasedAt uint64) float32 {
	var val float32 = 1.0
	for _, gen := range m.Generators {
		val *= GetValue(gen, t, releasedAt)
	}
	return val
}

func (s Average) GetInfo() Info  { return s.Info }
func (m Multiply) GetInfo() Info { return m.Info }
func (s Average) Copy(historyLen int) Generator {
	s.Info = s.Info.Copy(historyLen)
	return s
}
func (m Multiply) Copy(historyLen int) Generator {
	m.Info = m.Info.Copy(historyLen)
	return m
}
