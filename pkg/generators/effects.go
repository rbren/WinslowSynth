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
		PeakLevel:    Constant{name, name + " Level", 1.0, 0.0, 1.0},
		SustainLevel: Constant{name, name + " Sustain", 0.8, 0.0, 1.0},
		AttackTime:   uint64(attackMs * samplesPerMs),
		DecayTime:    uint64(decayMs * samplesPerMs),
		ReleaseTime:  uint64(releaseMs * samplesPerMs),
	}
}

func GetHarmonicConstant(name string) Instrument {
	return Multiply{
		Generators: []Generator{
			frequencyConst(),
			Constant{name, name + " Harmonic", 1.0, .5, 4.0},
		},
	}
}

func GetLFO(name string, amplitude Generator) Instrument {
	return Multiply{
		Generators: []Generator{
			amplitude,
			Spinner{
				Bias:      Constant{Value: 1.0},
				Amplitude: Constant{name, name + " LFO strength", 0.0, 0.0, 2.0},
				Frequency: Constant{name, name + " LFO frequency", 2.0, 0.0, 20.0},
			},
		},
	}
}

func AddNoise(name string, base Generator) Instrument {
	return NoiseFilter{
		Input:  base,
		Amount: Constant{name, name + " Noise", 0.0, 0.0, 1.0},
	}
}
