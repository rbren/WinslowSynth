package generators

import (
	"github.com/sirupsen/logrus"
)

type Info struct {
	Name    string
	Group   string
	History History
}

type History struct {
	Samples  []float32
	Position int
	Time     uint64
}

func copyInfo(dest *Info, src Info) {
	if dest == nil {
		panic("Tried to set info on an instrument with no existing info")
	}
	dest.Name = src.Name
	dest.Group = src.Group
	dest.History = src.History
}

func getEmptyHistory() History {
	return History{
		Samples: make([]float32, historyLength),
	}
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
	logrus.Info("set up instrument", i)
	info := Info{}
	if existingInfo := i.GetInfo(); existingInfo != nil {
		info = *existingInfo
	}
	info.History = getEmptyHistory()
	logrus.Info("setting info", i)
	i.SetInfo(info)
}

func GetValue(g Generator, t, r uint64) float32 {
	// TODO: use history as a cache
	val := g.GetValue(t, r)
	//AddHistory(g, t, []float32{val})

	return val
}

func AddHistory(g Generator, startTime uint64, history []float32) {
	i := g.GetInfo()
	if i == nil || i.History.Samples == nil {
		return
	}
	for idx, val := range history {
		i.History.Samples[i.History.Position] = val
		i.History.Position = (i.History.Position + 1) % len(i.History.Samples)
		i.History.Time = startTime + uint64(idx)
	}
}
