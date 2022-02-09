package generators

var Library = map[string]Generator{}

func init() {
	Library = map[string]Generator{
		"winslow": Mega(),
		"sine":    BasicSine(),
		"saw":     BasicSaw(),
		"average": Average{
			Generators: []Generator{
				BasicSaw(),
				BasicSaw(),
				BasicSaw(),
				BasicSaw(),
				BasicSaw(),
				BasicSaw(),
				BasicSaw(),
				BasicSaw(),
			},
		},
		"square":      BasicSquare(),
		"reverb":      AddReverb("sine", BasicSine()),
		"warbler":     Warbler(),
		"dirty":       DirtySawWave(),
		"harmonic":    HarmonicOscillator(),
		"noiseFilter": NoisySineWave(),
	}
}

func GetDefaultInstrument() Generator {
	return Mega()
}
