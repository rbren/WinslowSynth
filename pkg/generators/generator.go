package generators

import (
	"github.com/rbren/midi/pkg/config"
)

type Info struct {
	Name     string
	Group    string
	Subgroup string
	History  *History
}

type History struct {
	samples  []float32
	Position int
	Time     uint64
}

type Generator interface {
	Initialize(group string) Generator
	GetInfo() Info
	Copy(historyLen int) Generator
	SubGenerators() []Generator
	GetValue(elapsed uint64, releasedAt uint64) float32
}

func getTimeInSamples(ms float32) uint64 {
	samplesPerMs := config.MainConfig.SampleRate / 1000
	return uint64(int(ms) * samplesPerMs)
}
