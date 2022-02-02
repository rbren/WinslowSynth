package generators

import (
	"github.com/sirupsen/logrus"
)

type Info struct {
	Name            string
	Group           string
	History         []float32
	HistoryPosition int
}

type Generator interface {
	GetInfo() *Info
	SetInfo(Info)
	GetValue(elapsed uint64, releasedAt uint64) float32
}

type Instrument interface {
	Generator
}

func GetDefaultInstrument() Instrument {
	return Mega()
}

func SetUpInstrument(i Instrument) {
	name := ""
	if info := i.GetInfo(); info != nil {
		name = info.Name
	}
	i.SetInfo(Info{
		Name:    name,
		History: getEmptyHistory(),
	})
}

func copyInfo(dest *Info, src Info) {
	if dest == nil {
		logrus.Error("Tried to set info on a nil instrument")
		return
	}
	dest.Name = src.Name
	dest.Group = src.Group
	dest.History = src.History
	dest.HistoryPosition = src.HistoryPosition
}
