package generators

var Library = map[string]Generator{}

func init() {
	Library = map[string]Generator{
		"winslow":     Mega(),
		"warbler":     Warbler(),
		"sine":        AddDelay("sine", BasicSine()),
		"saw":         BasicSaw(),
		"square":      BasicSquare(),
		"dirty":       DirtySawWave(),
		"harmonic":    HarmonicOscillator(),
		"noiseFilter": NoisySineWave(),
	}
}

func GetDefaultInstrument() Generator {
	return Mega()
}
