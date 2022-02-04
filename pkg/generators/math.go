package generators

import ()

type Average struct {
	Info       *Info
	Generators []Generator
}

type Multiply struct {
	Info       *Info
	Generators []Generator
}

func (s Average) GetValue(t, releasedAt uint64) float32 {
	if len(s.Generators) == 0 {
		return 0.0
	}
	var val float32 = 0.0
	for _, gen := range s.Generators {
		val += getValue(gen, t, releasedAt)
	}
	return val / float32(len(s.Generators))
}

func (m Multiply) GetValue(t, releasedAt uint64) float32 {
	var val float32 = 1.0
	for _, gen := range m.Generators {
		val *= getValue(gen, t, releasedAt)
	}
	return val
}

func (s Average) GetInfo() *Info     { return s.Info }
func (s Average) SetInfo(info Info)  { copyInfo(s.Info, info) }
func (m Multiply) GetInfo() *Info    { return m.Info }
func (m Multiply) SetInfo(info Info) { copyInfo(m.Info, info) }
