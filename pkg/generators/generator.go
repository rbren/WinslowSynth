package generators

import (
	"github.com/sirupsen/logrus"
)

type Info struct {
	Name            string
	Group           string
	History         []float32
	HistoryPosition int
	HistoryTime     uint64
}

func copyInfo(dest *Info, src Info) {
	if dest == nil {
		panic("Tried to set info on an instrument with no existing info")
	}
	dest.Name = src.Name
	dest.Group = src.Group
	dest.History = src.History
	dest.HistoryTime = src.HistoryTime
}

func getEmptyHistory() []float32 {
	return make([]float32, historyLength)
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
	AddHistory(g, t, []float32{val})

	return val
}

func AddHistory(g Generator, startTime uint64, history []float32) {
	i := g.GetInfo()
	if i == nil || i.History == nil {
		return
	}
	for idx, val := range history {
		i.History[i.HistoryPosition] = val
		i.HistoryPosition = (i.HistoryPosition + 1) % len(i.History)
		i.HistoryTime = startTime + uint64(idx)
	}
}
