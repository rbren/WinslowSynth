package generators

import (
	"math/rand"
	"time"
)

type Noise struct {
	Info          Info
	SubGenerators SubGenerators
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (n Noise) GetInfo() Info                   { return n.Info }
func (n Noise) GetSubGenerators() SubGenerators { return n.SubGenerators }
func (n Noise) Copy(historyLen int) Generator {
	n.Info = n.Info.Copy(historyLen)
	n.SubGenerators = n.SubGenerators.Copy()
	return n
}

func (n Noise) Initialize(name string) Generator {
	if n.SubGenerators == nil {
		n.SubGenerators = make(map[string]Generator)
	}
	if n.SubGenerators["Amount"] == nil {
		n.SubGenerators["Amount"] = Constant{
			Info:  Info{Group: name, Name: "Noise"},
			Value: .1,
			Min:   0.0,
			Max:   1.0,
		}
	}
	n.SubGenerators["Amount"] = n.SubGenerators["Amount"].Initialize(name)
	return n
}

func (n Noise) GetValue(t, r uint64) float32 {
	random := rand.Float32()
	amt := GetValue(n.SubGenerators["Amount"], t, r)
	max := 1.0 + amt
	min := 1.0 - amt
	return min + random*(max-min)
}
