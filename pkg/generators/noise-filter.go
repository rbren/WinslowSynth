package generators

import (
	"math/rand"
	"time"
)

type NoiseFilter struct {
	Info   Info
	Amount Generator
	Input  Generator
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (n NoiseFilter) SubGenerators() []Generator {
	return []Generator{n.Input, n.Amount}
}

func (n NoiseFilter) Initialize(name string) Generator {
	n.Info.History = getEmptyHistory()
	if n.Input == nil {
		panic("NoiseFilter has no input")
	}
	if n.Amount == nil {
		n.Amount = Constant{
			Info:  Info{Group: name, Name: "Noise"},
			Value: 0.0,
			Min:   0.0,
			Max:   0.5,
		}
	}
	n.Input = n.Input.Initialize(name)
	n.Amount = n.Amount.Initialize(name)
	return n
}

func (n NoiseFilter) GetValue(t, r uint64) float32 {
	random := rand.Float32()
	amt := GetValue(n.Amount, t, r)
	hold := random < amt && t != 0
	if !hold {
		return GetValue(n.Input, t, r)
	}
	return GetValue(n.Input, t-1, r)
}

func (n NoiseFilter) GetInfo() Info { return n.Info }
func (n NoiseFilter) Copy(historyLen int) Generator {
	n.Info = n.Info.Copy(historyLen)
	return n
}
