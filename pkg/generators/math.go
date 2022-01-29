package generators

type Sum struct {
	Generators []Generator
}

func (s Sum) GetValue(t uint64) float32 {
	var val float32 = 0.0
	for _, gen := range s.Generators {
		val += gen.GetValue(t)
	}
	return val
}
