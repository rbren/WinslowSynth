package generators

import (
	"math/rand"
	"time"
)

type NoiseFilter struct {
	Info   *Info
	Amount Generator
	Input  Instrument
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func NewNoiseFilter(input Generator, amt Generator) NoiseFilter {
	return NoiseFilter{
		Info: &Info{
			History: getEmptyHistory(),
		},
		Input:  input,
		Amount: amt,
	}
}

func (n NoiseFilter) GetValue(t, r uint64) float32 {
	random := rand.Float32()
	amt := n.Amount.GetValue(t, r)
	hold := random < amt
	if !hold {
		return n.Input.GetValue(t, r)
	}
	lastValueIndex := n.Info.HistoryPosition - 1
	if lastValueIndex < 0 {
		lastValueIndex = len(n.Info.History) - 1
	}
	val := n.Info.History[lastValueIndex]
	return val
}

func (n NoiseFilter) GetInfo() *Info    { return n.Info }
func (n NoiseFilter) SetInfo(info Info) { copyInfo(n.Info, info) }
