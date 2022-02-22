package generators

var Library = map[string]Generator{}

func init() {
	Library = map[string]Generator{
		"oscillator": Oscillator{}.Initialize("Oscillator"),
		"winslow":    Mega().Initialize("Winslow"),
		"sine":       Sine().Initialize("Sine"),
		"saw":        EnvSaw().Initialize("Saw"),
		"average": Average{
			Generators: []Generator{
				EnvSaw(),
				EnvSaw(),
				EnvSaw(),
				EnvSaw(),
				EnvSaw(),
				EnvSaw(),
				EnvSaw(),
				EnvSaw(),
			},
		}.Initialize("Average"),
		"square": EnvSquare().Initialize("Square"),
		"reverb": Reverb{
			Input: EnvSine(),
		}.Initialize("Reverb"),
		"warbler":     Warbler().Initialize("Warbler"),
		"dirty":       DirtySawWave().Initialize("Dirty Saw"),
		"harmonic":    HarmonicOscillator().Initialize("Harmonic"),
		"noiseFilter": NoisySineWave().Initialize("Noisy Sine"),
	}
}

func GetDefaultInstrument() Generator {
	return Mega().Copy(UseDefaultHistoryLength, true)
}
