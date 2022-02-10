package generators

import (
	"github.com/rbren/midi/pkg/config"
)

type Delay struct {
	Info   Info
	Amount Generator
	Input  Generator
}

func NewDelay(input Generator, amt Generator) Delay {
	// TODO: ensure input has history being tracked
	return Delay{
		Input:  input.Copy(UseDefaultHistoryLength),
		Amount: amt,
	}
}

func (d Delay) GetValue(t, r uint64) float32 {
	samplesPerMs := config.MainConfig.SampleRate / 1000
	amtMs := GetValue(d.Amount, t, r)
	amtSamples := int(amtMs) * samplesPerMs
	//GetValue(d.Input, t, r) // Ignore current value, but store it in history
	return GetValue(d.Input, t-uint64(amtSamples), r)
}

func (d Delay) GetInfo() Info { return d.Info }
func (d Delay) Copy(historyLen int) Generator {
	d.Info = d.Info.Copy(historyLen)
	return d
}
