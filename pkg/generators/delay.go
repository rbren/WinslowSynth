package generators

import (
	"github.com/rbren/midi/pkg/config"
)

type Delay struct {
	Info   Info
	Amount Generator
	Input  Generator
}

func (d Delay) SubGenerators() []Generator {
	return []Generator{d.Input, d.Amount}
}

func (d Delay) Initialize(name string) Generator {
	if d.Input == nil {
		panic("Delay has no input")
	}
	if d.Amount == nil {
		d.Amount = Constant{
			Info: Info{
				Name:  "Delay",
				Group: name,
			},
			Value: 10,
			Min:   0,
			Max:   500,
			Step:  10,
		}
	}
	d.Input = d.Input.Initialize(name)
	d.Input = d.Input.Copy(UseDefaultHistoryLength, false)
	d.Amount = d.Amount.Initialize(name)
	return d
}

func (d Delay) GetValue(t, r uint64) float32 {
	samplesPerMs := config.MainConfig.SampleRate / 1000
	amtMs := GetValue(d.Amount, t, r)
	amtSamples := int(amtMs) * samplesPerMs
	return GetValue(d.Input, t-uint64(amtSamples), r)
}

func (d Delay) GetInfo() Info { return d.Info }
func (d Delay) Copy(historyLen int, storeFrequencies bool) Generator {
	d.Info = d.Info.Copy(historyLen, storeFrequencies)
	d.Amount = d.Amount.Copy(CopyExistingHistoryLength, false)
	d.Input = d.Input.Copy(CopyExistingHistoryLength, false)
	return d
}
