package generators

var Library = map[string]Generator{}

func init() {
	Library = map[string]Generator{
		"warbler":     Warbler(),
		"sine":        AddDelay("sine", BasicSine()),
		"saw":         BasicSaw(),
		"square":      BasicSquare(),
		"dirty":       DirtySawWave(),
		"harmonic":    HarmonicOscillator(),
		"noiseFilter": NoisySineWave(),
		"mega":        Mega(),
	}
}

func GetDefaultInstrument() Generator {
	return Mega()
}
