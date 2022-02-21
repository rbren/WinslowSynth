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

func (a Average) SubGenerators() []Generator  { return a.Generators }
func (m Multiply) SubGenerators() []Generator { return m.Generators }

func (a Average) Initialize(group string) Generator {
	for idx := range a.Generators {
		a.Generators[idx] = a.Generators[idx].Initialize(group)
	}
	return a
}

func (m Multiply) Initialize(group string) Generator {
	for idx := range m.Generators {
		m.Generators[idx] = m.Generators[idx].Initialize(group)
	}
	return m
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
func (s Average) Copy(historyLen int, storeFrequencies bool) Generator {
	s.Info = s.Info.Copy(historyLen, storeFrequencies)
	for idx := range s.Generators {
		s.Generators[idx] = s.Generators[idx].Copy(CopyExistingHistoryLength, false)
	}
	return s
}
func (m Multiply) Copy(historyLen int, storeFrequencies bool) Generator {
	m.Info = m.Info.Copy(historyLen, storeFrequencies)
	for idx := range m.Generators {
		m.Generators[idx] = m.Generators[idx].Copy(CopyExistingHistoryLength, false)
	}
	return m
}
