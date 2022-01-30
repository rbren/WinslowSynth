package generators

import (
	"github.com/rbren/midi/pkg/config"
)

var Library = map[string]Instrument{
	"warbler":  Warbler(),
	"sine":     BasicSine(),
	"saw":      BasicSaw(),
	"dirty":    DirtySawWave(),
	"harmonic": HarmonicSpinner(),
}

// Instruments are Generators/Spinners with Frequency=nil
func Warbler() Spinner {
	adsr := BasicADSR()
	adsrInner := adsr
	adsrInner.SustainLevel = Constant{"Warble Amt", 20.0, 0.0, 100.0}
	return Spinner{
		Amplitude:     adsr,
		DropOnRelease: true,
		Frequency: Spinner{
			// setting Bias on this sets the overall freq
			Amplitude: adsrInner,
			Frequency: Constant{"Warble Speed", 4, 0.0, 20.0},
		},
	}
}

func BasicADSR() ADSR {
	samplesPerMs := config.MainConfig.SampleRate / 1000
	attackMs := 50
	decayMs := 1000
	releaseMs := 1000
	return ADSR{
		PeakLevel:    Constant{Value: 1.0},
		SustainLevel: Constant{"Sustain", 0.5, 0.0, 1.0},
		AttackTime:   uint64(attackMs * samplesPerMs),
		DecayTime:    uint64(decayMs * samplesPerMs),
		ReleaseTime:  uint64(releaseMs * samplesPerMs),
	}
}

func BasicSine() Spinner {
	return Spinner{
		Amplitude: BasicADSR(),
	}
}

func BasicSaw() SawWave {
	return SawWave{
		Amplitude: BasicADSR(),
	}
}

func DirtySawWave() Instrument {
	base := BasicSaw()
	base.Amplitude = Multiply{
		Generators: []Generator{base.Amplitude, Noise{Min: .5, Max: 1.5}},
	}
	return base
}

func HarmonicSpinner() Instrument {
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
