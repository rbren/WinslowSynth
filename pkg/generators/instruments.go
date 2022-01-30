package generators

import (
	"github.com/rbren/midi/pkg/config"
)

// Instruments are Generators/Spinners with Frequency=nil
func Warbler() Spinner {
	return Spinner{
		Frequency: Spinner{
			// setting Bias on this sets the overall freq
			Amplitude: Constant{20.0},
			Frequency: Constant{4},
		},
	}
}

func BasicSineWave() Spinner {
	return Spinner{}
}

func RampSawWave() SawWave {
	samplesPerMs := config.MainConfig.SampleRate / 1000
	attackMs := 50
	decayMs := 1000
	releaseMs := 1000
	return SawWave{
		Amplitude: ADSR{
			PeakLevel:    1.0,
			SustainLevel: 0.5,
			AttackTime:   uint64(attackMs * samplesPerMs),
			DecayTime:    uint64(decayMs * samplesPerMs),
			ReleaseTime:  uint64(releaseMs * samplesPerMs),
		},
	}
}

func DirtySawWave() Generator {
	base := RampSawWave()
	base.Amplitude = Multiply{
		Generators: []Generator{base.Amplitude, Noise{Min: .5, Max: 1.5}},
	}
	return base
}

func HarmonicSpinner() Generator {
	base := SimpleRamper()
	return Harmonic{
		Spinner: base,
		Modes: []Mode{
			Mode{Frequency: 1.5, Amplitude: .25},
			Mode{Frequency: 2.0, Amplitude: .15},
			Mode{Frequency: 4.0, Amplitude: .1},
		},
	}
}

func AmplitudeRamp() Ramp {
	samplesPerMs := config.MainConfig.SampleRate / 1000
	rampUpMs := 100
	rampDownMs := 500
	return Ramp{
		RampUp:   uint64(rampUpMs * samplesPerMs),
		RampDown: uint64(rampDownMs * samplesPerMs),
		Target:   1.0,
	}
}

func SimpleRamper() Spinner {
	g := Spinner{
		Amplitude: AmplitudeRamp(),
		Phase:     Constant{Value: 0.0},
	}
	return g
}
