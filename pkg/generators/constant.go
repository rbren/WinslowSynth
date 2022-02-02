package generators

type Constant struct {
	Info  *Info
	Value float32
	Min   float32
	Max   float32
}

func (c Constant) GetValue(t, r uint64) float32 {
	return c.Value
}

func (c Constant) GetInfo() *Info    { return c.Info }
func (c Constant) SetInfo(info Info) { copyInfo(c.Info, info) }
