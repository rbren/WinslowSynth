package generators

import (
	"fmt"

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

const numMegaOscillators = 1

func Mega() Instrument {
	waves := []Generator{}

	if numMegaOscillators > 1 {
		for i := 0; i < 30; i++ {
			wave1 := TheWorks(fmt.Sprintf("Sine%d", i), WaveShape)
			wave2 := TheWorks(fmt.Sprintf("Saw%d", i), SawShape)
			wave3 := TheWorks(fmt.Sprintf("Square%d", i), SquareShape)
			waves = append(waves, wave1, wave2, wave3)
		}
	} else {
		wave1 := TheWorks("Sine", WaveShape)
		wave2 := TheWorks("Saw", SawShape)
		wave3 := TheWorks("Square", SquareShape)
		waves = append(waves, wave1, wave2, wave3)
	}

	return Average{
		Info: &Info{
			Name:    "Mega",
			History: getEmptyHistory(),
		},
		Generators: waves,
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
