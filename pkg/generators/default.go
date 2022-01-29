package generators

import (
	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/input"
)

func GetDefaultGenerator(key input.InputKey) Generator {
	samplesPerMs := config.MainConfig.SampleRate / 1000
	attackMs := 50
	decayMs := 1000
	releaseMs := 1000
	return SawWave{
		Frequency: Constant{key.Frequency},
		Amplitude: ADSR{
			PeakLevel:    1.0,
			SustainLevel: 0.5,
			AttackTime:   uint64(attackMs * samplesPerMs),
			DecayTime:    uint64(decayMs * samplesPerMs),
			ReleaseTime:  uint64(releaseMs * samplesPerMs),
		},
		/*
			Amplitude: Multiply{
				Generators: []Generator{amplitudeRamp(), Noise{Min: .5, Max: 1.5}},
			},
		*/
	}
}

func harmonicSpinner(key input.InputKey) Generator {
	base := simpleRamper(key)
	return Harmonic{
		Spinner: base,
		Modes: []Mode{
			Mode{Frequency: 1.5, Amplitude: .25},
			Mode{Frequency: 2.0, Amplitude: .15},
			Mode{Frequency: 4.0, Amplitude: .1},
		},
	}
}

func amplitudeRamp() Ramp {
	samplesPerMs := config.MainConfig.SampleRate / 1000
	rampUpMs := 100
	rampDownMs := 500
	return Ramp{
		RampUp:   uint64(rampUpMs * samplesPerMs),
		RampDown: uint64(rampDownMs * samplesPerMs),
		Target:   1.0,
	}
}

func simpleRamper(key input.InputKey) Spinner {
	g := Spinner{
		Amplitude: amplitudeRamp(),
		Frequency: Constant{Value: key.Frequency},
		Phase:     Constant{Value: 0.0},
	}
	return g
}
