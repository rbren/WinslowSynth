package generators

import (
	"fmt"
)

type Average struct {
	Info          Info
	Generators    []Generator
	SubGenerators SubGenerators
}

type Multiply struct {
	Info          Info
	Generators    []Generator
	SubGenerators SubGenerators
}

func (a Average) GetInfo() Info                    { return a.Info }
func (m Multiply) GetInfo() Info                   { return m.Info }
func (a Average) GetSubGenerators() SubGenerators  { return a.SubGenerators }
func (m Multiply) GetSubGenerators() SubGenerators { return m.SubGenerators }
func (a Average) Copy(historyLen int) Generator {
	a.Info = a.Info.Copy(historyLen)
	a.SubGenerators = a.SubGenerators.Copy()
	return a
}
func (m Multiply) Copy(historyLen int) Generator {
	m.Info = m.Info.Copy(historyLen)
	m.SubGenerators = m.SubGenerators.Copy()
	return m
}

func (a Average) Initialize(group string) Generator {
	a.SubGenerators = map[string]Generator{}
	for idx := range a.Generators {
		a.SubGenerators[fmt.Sprintf("%d", idx)] = a.Generators[idx].Initialize(group)
	}
	return a
}

func (m Multiply) Initialize(group string) Generator {
	m.SubGenerators = map[string]Generator{}
	for idx := range m.Generators {
		m.SubGenerators[fmt.Sprintf("%d", idx)] = m.Generators[idx].Initialize(group)
	}
	return m
}

func (s Average) GetValue(t, releasedAt uint64) float32 {
	if len(s.Generators) == 0 {
		return 0.0
	}
	var val float32 = 0.0
	for _, gen := range s.SubGenerators {
		val += GetValue(gen, t, releasedAt)
	}
	return val / float32(len(s.Generators))
}

func (m Multiply) GetValue(t, releasedAt uint64) float32 {
	var val float32 = 1.0
	for _, gen := range m.SubGenerators {
		val *= GetValue(gen, t, releasedAt)
	}
	return val
}
