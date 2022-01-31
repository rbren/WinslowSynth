package generators

import (
	"math/rand"
	"time"
)

type NoiseFilter struct {
	Info   *Info
	Amount Generator
	Input  Instrument
	buffer []float32
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (n *NoiseFilter) initialize() {
	if n.buffer == nil {
		n.buffer = make([]float32, 1)
	}
}

func (n NoiseFilter) GetValue(t, r uint64) float32 {
	n.initialize()
	random := rand.Float32()
	amt := n.Amount.GetValue(t, r)
	hold := random < amt
	val := n.buffer[0]
	if !hold {
		val = n.Input.GetValue(t, r)
	}
	n.buffer[0] = val
	return val
}
func (n NoiseFilter) GetInfo() *Info { return n.Info }
