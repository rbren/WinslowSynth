package generators

import (
	"github.com/rbren/midi/pkg/config"
)

func GetADSR(name string) ADSR {
	samplesPerMs := config.MainConfig.SampleRate / 1000
	attackMs := 50
	decayMs := 1000
	releaseMs := 1000
	return ADSR{
		PeakLevel: Constant{
			Info:  &Info{Group: name, Name: name + " Level"},
			Value: 1.0,
			Min:   0.0,
			Max:   1.0,
		},
		SustainLevel: Constant{
			Info:  &Info{Group: name, Name: name + " Sustain"},
			Value: 0.8,
			Min:   0.0,
			Max:   1.0,
		},
		AttackTime:  uint64(attackMs * samplesPerMs),
		DecayTime:   uint64(decayMs * samplesPerMs),
		ReleaseTime: uint64(releaseMs * samplesPerMs),
	}
}

func GetHarmonicConstant(name string) Instrument {
	return Multiply{
		Generators: []Generator{
			frequencyConst(),
			Constant{
				Info:  &Info{Group: name, Name: name + " Harmonic"},
				Value: 1.0,
				Min:   .5,
				Max:   4.0,
			},
		},
	}
}

func GetLFO(name string, amplitude Generator) Instrument {
	return Multiply{
		Generators: []Generator{
			amplitude,
			Spinner{
				Bias: Constant{Value: 1.0},
				Amplitude: Constant{
					Info:  &Info{Group: name, Name: name + " LFO strength"},
					Value: 0.0,
					Min:   0.0,
					Max:   2.0,
				},
				Frequency: Constant{
					Info:  &Info{Group: name, Name: name + " LFO frequency"},
					Value: 2.0,
					Min:   0.0,
					Max:   20.0,
				},
			},
		},
	}
}

func AddNoise(name string, base Generator) Instrument {
	return NoiseFilter{
		Input: base,
		Amount: Constant{
			Info:  &Info{Group: name, Name: name + " Noise"},
			Value: 0.0,
			Min:   0.0,
			Max:   1.0,
		},
	}
}
