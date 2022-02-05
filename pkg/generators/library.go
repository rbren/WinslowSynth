package generators

import (
	"github.com/rbren/midi/pkg/config"
)

var historyMs = 50
var historyLength int
var Library = map[string]Instrument{}

func init() {
	historyLength = historyMs * (config.MainConfig.SampleRate / 1000) // store 1s
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
