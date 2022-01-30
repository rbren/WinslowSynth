package generators

import (
	"github.com/rbren/midi/pkg/input"
)

func GetDefaultGenerator(key input.InputKey) Generator {
	return BasicSineWave()
}
