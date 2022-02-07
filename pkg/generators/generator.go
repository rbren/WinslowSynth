package generators

import (
	_ "github.com/sirupsen/logrus"
)

type Info struct {
	Name    string
	Group   string
	History *History
}

type History struct {
	samples  []float32
	Position int
	Time     uint64
}

type Generator interface {
	GetInfo() Info
	Copy(historyLen int) Generator
	GetValue(elapsed uint64, releasedAt uint64) float32
}
