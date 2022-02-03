package generators

import (
	"github.com/rbren/midi/pkg/config"
)

func frequencyConst() Constant {
	return Constant{
		Info:  &Info{Group: "", Name: "Frequency"},
		Value: 440.0,
		Min:   20.0,
		Max:   20000.0,
	}
}

func BasicSine() Oscillator {
	return Oscillator{
		Frequency: frequencyConst(),
		Info: &Info{
			Name: "Basic Sine",
		},
		Amplitude: GetADSR("adsr"),
	}
}

func BasicSaw() Oscillator {
	return Oscillator{
		Info:      &Info{Name: "Basic Saw"},
		Shape:     SawShape,
		Frequency: frequencyConst(),
		Amplitude: GetADSR("adsr"),
	}
}

func BasicSquare() Oscillator {
	return Oscillator{
		Info:      &Info{Name: "Basic Square"},
		Shape:     SquareShape,
		Frequency: frequencyConst(),
		Amplitude: GetADSR("adsr"),
	}
}

func Mega() Instrument {
	oscSin := BasicSine()
	oscSaw := BasicSaw()
	oscSqr := BasicSquare()

	oscSin.Frequency = GetHarmonicConstant("Sine")
	oscSaw.Frequency = GetHarmonicConstant("Saw")
	oscSqr.Frequency = GetHarmonicConstant("Square")

	oscSin.Amplitude = GetADSR("Sine")
	oscSaw.Amplitude = GetADSR("Saw")
	oscSqr.Amplitude = GetADSR("Square")

	oscSin.Amplitude = GetLFO("Sine", oscSin.Amplitude)
	oscSaw.Amplitude = GetLFO("Saw", oscSaw.Amplitude)
	oscSqr.Amplitude = GetLFO("Square", oscSqr.Amplitude)

	wave1 := AddNoise("Sine", oscSin)
	wave2 := AddNoise("Saw", oscSaw)
	wave3 := AddNoise("Square", oscSqr)

	wave1 = AddLevel("Sine", wave1)
	wave2 = AddLevel("Saw", wave2)
	wave3 = AddLevel("Square", wave3)

	return Average{
		Info: &Info{
			Name:    "Mega",
			History: getEmptyHistory(),
		},
		Generators: []Generator{wave1, wave2, wave3},
	}
}

func NoisySineWave() Instrument {
	base := BasicSine()
	return NoiseFilter{
		Info:  &Info{Name: "Noisy Sine"},
		Input: base,
		Amount: Constant{
			Info:  &Info{Group: "", Name: "Noise"},
			Value: .2,
			Min:   0.0,
			Max:   1.0,
		},
	}
}

func Warbler() Oscillator {
	adsr := GetADSR("adsr")
	adsrInner := adsr
	adsrInner.SustainLevel = Constant{
		Info:  &Info{Group: "", Name: "Warble Amt"},
		Value: 20.0,
		Min:   0.0,
		Max:   100.0,
	}
	return Oscillator{
		Info:          &Info{Name: "Warbler"},
		Amplitude:     adsr,
		DropOnRelease: true,
		Frequency: Oscillator{
			Amplitude: adsrInner,
			Frequency: Constant{
				Info:  &Info{Group: "", Name: "Warble Speed"},
				Value: 4,
				Min:   0.0,
				Max:   20.0,
			},
			// setting Bias on this sets the overall freq
			Bias: frequencyConst(),
		},
	}
}

func DirtySawWave() Instrument {
	base := BasicSaw()
	base.Amplitude = Multiply{
		Info: &Info{Name: "Dirty Saw"},
		Generators: []Generator{
			base.Amplitude,
			Noise{
				Amount: Constant{
					Info:  &Info{Group: "", Name: "Noise"},
					Value: .1,
					Min:   0.0,
					Max:   1.0,
				},
			},
		},
	}
	return base
}

func HarmonicOscillator() Instrument {
	base := SimpleRamper()
	return Harmonic{
		Info:       &Info{Name: "Harmonic"},
		Oscillator: base,
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
		Info:     &Info{Name: "Ramper"},
		RampUp:   uint64(rampUpMs * samplesPerMs),
		RampDown: uint64(rampDownMs * samplesPerMs),
		Target:   1.0,
	}
}

func SimpleRamper() Oscillator {
	g := Oscillator{
		Info:      &Info{Name: "Simple Ramper"},
		Frequency: frequencyConst(),
		Amplitude: AmplitudeRamp(),
		Phase:     Constant{Value: 0.0},
	}
	return g
}
