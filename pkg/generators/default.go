package generators

import (
	"github.com/rbren/midi/pkg/input"
)

func GetDefaultGenerator(key input.InputKey) Generator {
	g := NewSpinner(1.0, key.Frequency, 0.0)
	return g
}
