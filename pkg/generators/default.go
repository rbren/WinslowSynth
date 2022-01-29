package generators

import (
	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/input"
)

func GetDefaultGenerator(key input.InputKey) Generator {
	return SawWave{
		Frequency: Constant{key.Frequency},
		Amplitude: amplitudeRamp(),
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
