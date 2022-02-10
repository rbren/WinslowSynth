package generators

func TheWorks(name string, shape OscillatorShape) Generator {
	osc := Oscillator{
		Info: Info{
			Name:  name,
			Group: name,
		},
		Frequency: GetHarmonicConstant(name),
		Amplitude: GetLFO(name, GetADSR(name)),
		Shape:     shape,
	}
	var inst Generator = osc
	inst = NoiseFilter{
		Input: inst,
	}
	inst = Reverb{
		Input: inst,
	}
	inst = Delay{
		Input: inst,
	}
	inst = AddLevel(name, inst)
	return inst.Initialize(name)
}

func AddLevel(name string, inst Generator) Generator {
	return Multiply{
		Info: Info{Group: name},
		Generators: []Generator{
			Constant{
				Info: Info{
					Name:  "Level",
					Group: name,
				},
				Value: 1.0,
				Min:   0.0,
				Max:   1.0,
			},
			inst,
		},
	}
}

func GetADSR(name string) ADSR {
	return ADSR{
		PeakLevel: Constant{
			Info:  Info{Group: name, Subgroup: "ADSR", Name: "Peak"},
			Value: 1.0,
			Min:   0.0,
			Max:   1.0,
		},
		SustainLevel: Constant{
			Info:  Info{Group: name, Subgroup: "ADSR", Name: "Sustain"},
			Value: 0.8,
			Min:   0.0,
			Max:   1.0,
		},
		AttackTime: Constant{
			Info:  Info{Group: name, Subgroup: "ADSR", Name: "Attack"},
			Value: 300,
			Min:   0.0,
			Max:   1000,
			Step:  1.0,
		},
		DecayTime: Constant{
			Info:  Info{Group: name, Subgroup: "ADSR", Name: "Decay"},
			Value: 500,
			Min:   0.0,
			Max:   1000,
			Step:  1.0,
		},
		ReleaseTime: Constant{
			Info:  Info{Group: name, Subgroup: "ADSR", Name: "Release"},
			Value: 1000,
			Min:   0.0,
			Max:   3000,
			Step:  1.0,
		},
	}
}

func GetHarmonicConstant(name string) Generator {
	return Multiply{
		Generators: []Generator{
			frequencyConst(),
			Constant{
				Info:  Info{Group: name, Name: "Harmonic"},
				Value: 1.0,
				Min:   1.0,
				Max:   8.0,
				Step:  1.0,
			},
		},
	}
}

func GetLFO(name string, amplitude Generator) Generator {
	return Multiply{
		Generators: []Generator{
			amplitude,
			Oscillator{
				Bias: Constant{Value: 1.0},
				Amplitude: Constant{
					Info:  Info{Group: name, Subgroup: "LFO", Name: "Strength"},
					Value: 0.0,
					Min:   0.0,
					Max:   2.0,
				},
				Frequency: Constant{
					Info:  Info{Group: name, Subgroup: "LFO", Name: "Freq"},
					Value: 2.0,
					Min:   0.0,
					Max:   20.0,
				},
			},
		},
	}
}
