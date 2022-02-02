package generators

func GetADSR(name string) ADSR {
	return ADSR{
		PeakLevel: Constant{
			Info:  &Info{Group: name, Name: "Level"},
			Value: 1.0,
			Min:   0.0,
			Max:   1.0,
		},
		SustainLevel: Constant{
			Info:  &Info{Group: name, Name: "Sustain"},
			Value: 0.8,
			Min:   0.0,
			Max:   1.0,
		},
		AttackTime: Constant{
			Info:  &Info{Group: name, Name: "Attack"},
			Value: 300,
			Min:   0.0,
			Max:   1000,
			Step:  1.0,
		},
		DecayTime: Constant{
			Info:  &Info{Group: name, Name: "Decay"},
			Value: 500,
			Min:   0.0,
			Max:   1000,
			Step:  1.0,
		},
		ReleaseTime: Constant{
			Info:  &Info{Group: name, Name: "Release"},
			Value: 1000,
			Min:   0.0,
			Max:   3000,
			Step:  1.0,
		},
	}
}

func GetHarmonicConstant(name string) Instrument {
	return Multiply{
		Generators: []Generator{
			frequencyConst(),
			Constant{
				Info:  &Info{Group: name, Name: "Harmonic"},
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
					Info:  &Info{Group: name, Name: "LFO strength"},
					Value: 0.0,
					Min:   0.0,
					Max:   2.0,
				},
				Frequency: Constant{
					Info:  &Info{Group: name, Name: "LFO frequency"},
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
			Info:  &Info{Group: name, Name: "Noise"},
			Value: 0.0,
			Min:   0.0,
			Max:   1.0,
		},
	}
}
