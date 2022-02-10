package generators

var Library = map[string]Generator{}

func init() {
	Library = map[string]Generator{
		"winslow": Mega().Initialize("Winslow"),
		"sine":    BasicSine().Initialize("Sine"),
		"saw":     BasicSaw().Initialize("Saw"),
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
		}.Initialize("Average"),
		"square": BasicSquare().Initialize("Square"),
		"reverb": Reverb{
			Input: BasicSine(),
		}.Initialize("Reverb"),
		"warbler":     Warbler().Initialize("Warbler"),
		"dirty":       DirtySawWave().Initialize("Dirty Saw"),
		"harmonic":    HarmonicOscillator().Initialize("Harmonic"),
		"noiseFilter": NoisySineWave().Initialize("Noisy Sine"),
	}
}

func GetDefaultInstrument() Generator {
	return Mega()
}
