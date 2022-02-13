package generators

import (
	"math"

	"github.com/rbren/midi/pkg/config"
)

type Reverb struct {
	Info          Info
	SubGenerators SubGenerators
}

func (r Reverb) GetInfo() Info                   { return r.Info }
func (r Reverb) GetSubGenerators() SubGenerators { return r.SubGenerators }
func (r Reverb) Copy(historyLen int) Generator {
	r.Info = r.Info.Copy(historyLen)
	r.SubGenerators = r.SubGenerators.Copy()
	return r
}

func (r Reverb) Initialize(group string) Generator {
	var maxDelay float32 = 1000.0
	maxRepeats := float32(math.Floor(float64(historyMs-1) / float64(maxDelay)))
	if r.SubGenerators["Input"] == nil {
		panic("Reverb has no input")
	}
	if r.SubGenerators["Strength"] == nil {
		r.SubGenerators["Strength"] = Constant{
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
	if r.SubGenerators["Delay"] == nil {
		r.SubGenerators["Delay"] = Constant{
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
	if r.SubGenerators["Repeats"] == nil {
		r.SubGenerators["Repeats"] = Constant{
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
	r.SubGenerators["Input"] = r.SubGenerators["Input"].Initialize(group)
	r.SubGenerators["Input"] = r.SubGenerators["Input"].Copy(UseDefaultHistoryLength)
	r.SubGenerators["Delay"] = r.SubGenerators["Delay"].Initialize(group)
	r.SubGenerators["Strength"] = r.SubGenerators["Strength"].Initialize(group)
	r.SubGenerators["Repeats"] = r.SubGenerators["Repeats"].Initialize(group)
	return r
}

func (d Reverb) GetValue(t, r uint64) float32 {
	val := GetValue(d.SubGenerators["Input"], t, r)

	numRepeats := int(GetValue(d.SubGenerators["Repeats"], t, r))
	delayMs := GetValue(d.SubGenerators["Delay"], t, r)
	startAmplitude := GetValue(d.SubGenerators["Strength"], t, r)
	if startAmplitude == 0 || delayMs == 0 || numRepeats == 0 {
		return val
	}

	samplesPerMs := config.MainConfig.SampleRate / 1000
	for repetition := 0; repetition < numRepeats; repetition++ {
		delaySamples := uint64(int(delayMs) * (repetition + 1) * samplesPerMs)
		amplitude := startAmplitude * (1.0 - float32(repetition)/float32(numRepeats))
		oldTime := t - delaySamples
		if oldTime > 0 {
			oldVal := GetValue(d.SubGenerators["Input"], oldTime, r)
			val += amplitude * oldVal
		}
	}
	return val
}
