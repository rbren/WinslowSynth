package generators

type Constant struct {
	Name  string
	Value float32
	Min   float32
	Max   float32
}

func (c Constant) GetValue(t, r uint64) float32 {
	return c.Value
}
