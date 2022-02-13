package generators

import (
	"github.com/rbren/midi/pkg/config"
)

type Delay struct {
	Info          Info
	SubGenerators SubGenerators
}

func (d Delay) GetInfo() Info                   { return d.Info }
func (d Delay) GetSubGenerators() SubGenerators { return d.SubGenerators }
func (d Delay) Copy(historyLen int) Generator {
	d.Info = d.Info.Copy(historyLen)
	d.SubGenerators = d.SubGenerators.Copy()
	return d
}

func (d Delay) Initialize(name string) Generator {
	if d.SubGenerators == nil || d.SubGenerators["Input"] == nil {
		panic("Delay has no input")
	}
	if d.SubGenerators["Amount"] == nil {
		d.SubGenerators["Amount"] = Constant{
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
	d.SubGenerators["Input"] = d.SubGenerators["Input"].Initialize(name)
	d.SubGenerators["Input"] = d.SubGenerators["Input"].Copy(UseDefaultHistoryLength)
	d.SubGenerators["Amount"] = d.SubGenerators["Amount"].Initialize(name)
	return d
}

func (d Delay) GetValue(t, r uint64) float32 {
	samplesPerMs := config.MainConfig.SampleRate / 1000
	amtMs := GetValue(d.SubGenerators["Amount"], t, r)
	amtSamples := int(amtMs) * samplesPerMs
	return GetValue(d.SubGenerators["Input"], t-uint64(amtSamples), r)
}
