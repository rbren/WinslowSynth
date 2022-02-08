package generators

import (
	"math"

	"github.com/rbren/midi/pkg/config"
)

type Reverb struct {
	Info     Info
	Strength Generator
	Delay    Generator
	Decay    Generator
	Input    Generator
}

func NewReverb(group string, input Generator) Reverb {
	return Reverb{
		Input: input.Copy(UseDefaultHistoryLength),
		Strength: Constant{
			Info: Info{
				Name:  "Reverb Strength",
				Group: group,
			},
			Value: 1.0,
			Min:   0.0,
			Max:   1.0,
		},
		Delay: Constant{
			Info: Info{
				Name:  "Reverb Delay",
				Group: group,
			},
			Value: 250,
			Min:   0,
			Max:   1000,
		},
		Decay: Constant{
			Info: Info{
				Name:  "Reverb Decay",
				Group: group,
			},
			Value: 1000,
			Min:   0,
			Max:   float32(historyMs - 1),
		},
	}
}

func (d Reverb) GetValue(t, r uint64) float32 {
	val := GetValue(d.Input, t, r)

	samplesPerMs := config.MainConfig.SampleRate / 1000
	decayTimeMs := GetValue(d.Decay, t, r)
	delayMs := GetValue(d.Delay, t, r)
	numRepeats := int(math.Floor(float64(decayTimeMs / delayMs)))
	startAmplitude := GetValue(d.Strength, t, r)
	for repetition := 0; repetition < numRepeats; repetition++ {
		delaySamples := uint64(int(delayMs) * (repetition + 1) * samplesPerMs)
		amplitude := startAmplitude * (1.0 - float32(repetition)/float32(numRepeats))
		oldTime := t - delaySamples
		if oldTime > 0 {
			oldVal := GetValueCached(d.Input, oldTime)
			if oldVal == nil {
				panic("Couldn't get cached value for reverb")
			}
			val += amplitude * (*oldVal)
		}
	}
	return val
}

func (d Reverb) GetInfo() Info { return d.Info }
func (d Reverb) Copy(historyLen int) Generator {
	d.Info = d.Info.Copy(historyLen)
	return d
}
