package generators

type Constant struct {
	Value float32
}

func (c Constant) GetValue(t uint64) float32 {
	return c.Value
}
