package generators

import (
	"math/rand"
	"time"
)

type Noise struct {
	Info   Info
	Amount Generator
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (n Noise) SubGenerators() []Generator {
	return []Generator{n.Amount}
}

func (n Noise) Initialize(name string) Generator {
	if n.Amount == nil {
		n.Amount = Constant{
			Info:  Info{Group: name, Name: "Noise"},
			Value: .1,
			Min:   0.0,
			Max:   1.0,
		}
	}
	n.Amount = n.Amount.Initialize(name)
	return n
}

func (n Noise) GetValue(t, r uint64) float32 {
	random := rand.Float32()
	amt := GetValue(n.Amount, t, r)
	max := 1.0 + amt
	min := 1.0 - amt
	return min + random*(max-min)
}

func (n Noise) GetInfo() Info { return n.Info }
func (n Noise) Copy(historyLen int) Generator {
	n.Info = n.Info.Copy(historyLen)
	n.Amount = n.Amount.Copy(CopyExistingHistoryLength)
	return n
}
