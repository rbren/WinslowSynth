package generators

import (
	"math"

	"github.com/rbren/midi/pkg/config"
)

type Reverb struct {
	Info     Info
	Strength Generator
	Delay    Generator
	Repeats  Generator
	Input    Generator
}

func (r Reverb) SubGenerators() []Generator {
	return []Generator{r.Input, r.Strength, r.Delay, r.Repeats}
}

func (r Reverb) Initialize(group string) Generator {
	var maxDelay float32 = 1000.0
	maxRepeats := float32(math.Floor(float64(historyMs-1) / float64(maxDelay)))
	if r.Input == nil {
		panic("Reverb has no input")
	}
	if r.Strength == nil {
		r.Strength = Constant{
			Info: Info{
				Name:     "Strength",
				Group:    group,
				Subgroup: "Reverb",
			},
			Value: 1.0,
			Min:   0.0,
			Max:   1.0,
		}
	}
	if r.Delay == nil {
		r.Delay = Constant{
			Info: Info{
				Name:     "Delay",
				Group:    group,
				Subgroup: "Reverb",
			},
			Value: 250,
			Min:   0,
			Max:   maxDelay,
		}
	}
	if r.Repeats == nil {
		r.Repeats = Constant{
			Info: Info{
				Name:     "Repeats",
				Group:    group,
				Subgroup: "Reverb",
			},
			Value: 3,
			Min:   0,
			Max:   maxRepeats,
			Step:  1,
		}
	}
	r.Input = r.Input.Initialize(group)
	r.Input = r.Input.Copy(UseDefaultHistoryLength)
	r.Delay = r.Delay.Initialize(group)
	r.Strength = r.Strength.Initialize(group)
	r.Repeats = r.Repeats.Initialize(group)
	return r
}

func (d Reverb) GetValue(t, r uint64) float32 {
	val := GetValue(d.Input, t, r)

	numRepeats := int(GetValue(d.Repeats, t, r))
	delayMs := GetValue(d.Delay, t, r)
	startAmplitude := GetValue(d.Strength, t, r)
	if startAmplitude == 0 || delayMs == 0 || numRepeats == 0 {
		return val
	}

	samplesPerMs := config.MainConfig.SampleRate / 1000
	for repetition := 0; repetition < numRepeats; repetition++ {
		delaySamples := uint64(int(delayMs) * (repetition + 1) * samplesPerMs)
		amplitude := startAmplitude * (1.0 - float32(repetition)/float32(numRepeats))
		oldTime := t - delaySamples
		if oldTime > 0 {
			oldVal := GetValue(d.Input, oldTime, r)
			val += amplitude * oldVal
		}
	}
	return val
}

func (d Reverb) GetInfo() Info { return d.Info }
func (d Reverb) Copy(historyLen int) Generator {
	d.Info = d.Info.Copy(historyLen)
	d.Strength = d.Strength.Copy(CopyExistingHistoryLength)
	d.Delay = d.Delay.Copy(CopyExistingHistoryLength)
	d.Repeats = d.Repeats.Copy(CopyExistingHistoryLength)
	d.Input = d.Input.Copy(CopyExistingHistoryLength)
	return d
}
