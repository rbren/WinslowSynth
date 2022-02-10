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

func NewNoiseFilter(input Generator, amt Generator) NoiseFilter {
	return NoiseFilter{
		Info: Info{
			History: getEmptyHistory(),
		},
		Input:  input,
		Amount: amt,
	}
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
