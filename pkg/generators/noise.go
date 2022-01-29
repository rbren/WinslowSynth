package generators

import (
	"math/rand"
	"time"
)

type Noise struct {
	Max float32
	Min float32
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func (n Noise) GetValue(t, r uint64) float32 {
	random := rand.Float32()
	return n.Min + random*(n.Max-n.Min)
}
