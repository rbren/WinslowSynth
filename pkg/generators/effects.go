package generators

func TheWorks(name string, shape OscillatorShape) Generator {
	osc := Oscillator{
		Info: Info{
			Name:  name,
			Group: name,
		},
		Shape: shape,
		SubGenerators: map[string]Generator{
			"Frequency": GetHarmonicConstant(name),
			"Amplitude": GetLFO(name, GetADSR(name)),
		},
	}
	var inst Generator = osc
	inst = NoiseFilter{
		SubGenerators: map[string]Generator{"Input": inst},
	}
	inst = Reverb{
		SubGenerators: map[string]Generator{"Input": inst},
	}
	inst = Delay{
		SubGenerators: map[string]Generator{"Input": inst},
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
	return ADSR{}.Initialize(name).(ADSR)
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
				SubGenerators: map[string]Generator{
					"Bias": Constant{Value: 1.0},
					"Amplitude": Constant{
						Info:  Info{Group: name, Subgroup: "LFO", Name: "Strength"},
						Value: 0.0,
						Min:   0.0,
						Max:   2.0,
					},
					"Frequency": Constant{
						Info:  Info{Group: name, Subgroup: "LFO", Name: "Freq"},
						Value: 2.0,
						Min:   0.0,
						Max:   20.0,
					},
				},
			},
		},
	}
}
