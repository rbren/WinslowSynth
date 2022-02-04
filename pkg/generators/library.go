package generators

import (
	"github.com/rbren/midi/pkg/config"
)

var historyLength int
var Library = map[string]Instrument{}

func init() {
	historyLength = 50 * (config.MainConfig.SampleRate / 1000) // store 1s
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
