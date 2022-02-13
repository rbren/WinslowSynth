package generators

import (
	"math/rand"
	"time"
)

type NoiseFilter struct {
	Info          Info
	SubGenerators SubGenerators
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (n NoiseFilter) GetInfo() Info                   { return n.Info }
func (n NoiseFilter) GetSubGenerators() SubGenerators { return n.SubGenerators }
func (n NoiseFilter) Copy(historyLen int) Generator {
	n.Info = n.Info.Copy(historyLen)
	n.SubGenerators = n.SubGenerators.Copy()
	return n
}

func (n NoiseFilter) Initialize(name string) Generator {
	n.Info.History = getEmptyHistory()
	if n.SubGenerators["Input"] == nil {
		panic("NoiseFilter has no input")
	}
	if n.SubGenerators["Amount"] == nil {
		n.SubGenerators["Amount"] = Constant{
			Info:  Info{Group: name, Name: "Noise"},
			Value: 0.0,
			Min:   0.0,
			Max:   0.5,
		}
	}
	n.SubGenerators["Input"] = n.SubGenerators["Input"].Initialize(name)
	n.SubGenerators["Amount"] = n.SubGenerators["Amount"].Initialize(name)
	return n
}

func (n NoiseFilter) GetValue(t, r uint64) float32 {
	random := rand.Float32()
	amt := GetValue(n.SubGenerators["Amount"], t, r)
	hold := random < amt && t != 0
	if !hold {
		return GetValue(n.SubGenerators["Input"], t, r)
	}
	return GetValue(n.SubGenerators["Input"], t-1, r)
}
