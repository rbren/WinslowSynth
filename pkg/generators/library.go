package generators

var Library = map[string]Instrument{}

func init() {
	Library = map[string]Instrument{
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
