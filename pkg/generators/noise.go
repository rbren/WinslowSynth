package generators

import (
	"math/rand"
	"time"
)

type Noise struct {
	Info   *Info
	Amount Generator
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (n Noise) GetValue(t, r uint64) float32 {
	random := rand.Float32()
	amt := getValue(n.Amount, t, r)
	max := 1.0 + amt
	min := 1.0 - amt
	return min + random*(max-min)
}

func (n Noise) GetInfo() *Info    { return n.Info }
func (n Noise) SetInfo(info Info) { copyInfo(n.Info, info) }
