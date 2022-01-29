package generators

import (
	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/input"
)

func GetDefaultGenerator(key input.InputKey) Generator {
	rampUpMs := 100
	rampDownMs := 500
	samplesPerMs := config.MainConfig.SampleRate / 1000
	g := Spinner{
		Amplitude: Ramp{
			RampUp:   uint64(rampUpMs * samplesPerMs),
			RampDown: uint64(rampDownMs * samplesPerMs),
			Target:   1.0,
		},
		Frequency: Constant{Value: key.Frequency},
		Phase:     Constant{Value: 0.0},
	}
	return g
}
