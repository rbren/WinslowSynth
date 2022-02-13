package generators

type Constant struct {
	Info  Info
	Value float32
	Min   float32
	Max   float32
	Step  float32
}

func (c Constant) GetInfo() Info                   { return c.Info }
func (c Constant) GetSubGenerators() SubGenerators { return map[string]Generator{} }
func (c Constant) Copy(historyLen int) Generator {
	c.Info = c.Info.Copy(historyLen)
	return c
}

func (c Constant) Initialize(group string) Generator { return c }

func (c Constant) GetValue(t, r uint64) float32 {
	return c.Value
}
