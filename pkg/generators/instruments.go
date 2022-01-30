package generators

import (
	"github.com/rbren/midi/pkg/config"
)

var Library = map[string]Instrument{
	"warbler":     Warbler(),
	"sine":        BasicSine(),
	"saw":         BasicSaw(),
	"square":      BasicSquare(),
	"dirty":       DirtySawWave(),
	"harmonic":    HarmonicSpinner(),
	"noiseFilter": NoisySineWave(),
	"mega":        Mega(),
}

func frequencyConst() Constant {
	return Constant{"Frequency", 440.0, 20.0, 20000.0}
}

func BasicSine() Spinner {
	return Spinner{
		Frequency: frequencyConst(),
		Amplitude: BasicADSR(),
	}
}

func BasicSaw() SawWave {
	return SawWave{
		Frequency: frequencyConst(),
		Amplitude: BasicADSR(),
	}
}

func BasicSquare() SquareWave {
	return SquareWave{
		Frequency: frequencyConst(),
		Amplitude: BasicADSR(),
	}
}

func GetHarmonicConstant(name string, freq Generator) Instrument {
	return Multiply{
		Generators: []Generator{
			frequencyConst(),
			Constant{name + " Harmonic", 1.0, .5, 4.0},
		},
	}
}

func Mega() Instrument {
	wave1 := BasicSine()
	wave1.Frequency = GetHarmonicConstant("Sine", wave1.Frequency)
	wave2 := BasicSaw()
	wave3 := BasicSquare()
	return Sum{
		Generators: []Generator{wave1, wave2, wave3},
	}
}

func NoisySineWave() Instrument {
	base := BasicSine()
	return NoiseFilter{
		Input:  base,
		Amount: Constant{"Noise", .2, 0.0, 1.0},
	}
}

func Warbler() Spinner {
	adsr := BasicADSR()
	adsrInner := adsr
	adsrInner.SustainLevel = Constant{"Warble Amt", 20.0, 0.0, 100.0}
	return Spinner{
		Amplitude:     adsr,
		DropOnRelease: true,
		Frequency: Spinner{
			Amplitude: adsrInner,
			Frequency: Constant{"Warble Speed", 4, 0.0, 20.0},
			// setting Bias on this sets the overall freq
			Bias: frequencyConst(),
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
		SustainLevel: Constant{Value: 0.5},
		AttackTime:   uint64(attackMs * samplesPerMs),
		DecayTime:    uint64(decayMs * samplesPerMs),
		ReleaseTime:  uint64(releaseMs * samplesPerMs),
	}
}

func DirtySawWave() Instrument {
	base := BasicSaw()
	base.Amplitude = Multiply{
		Generators: []Generator{
			base.Amplitude,
			Noise{
				Amount: Constant{"Noise", .1, 0.0, 1.0},
			},
		},
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
		Frequency: frequencyConst(),
		Amplitude: AmplitudeRamp(),
		Phase:     Constant{Value: 0.0},
	}
	return g
}
