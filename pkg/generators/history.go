package generators

import (
	"github.com/rbren/midi/pkg/config"
)

const CopyExistingHistoryLength = -1
const UseDefaultHistoryLength = -2

var historyMs = 1000
var historyLength int

func init() {
	historyLength = historyMs * (config.MainConfig.SampleRate / 1000)
}

func getEmptyHistory() *History {
	return &History{
		Samples: make([]float32, historyLength),
	}
}

func AddHistory(g Generator, startTime uint64, history []float32) {
	i := g.GetInfo()
	if i.History == nil || i.History.Samples == nil {
		return
	}
	for idx, val := range history {
		idxTime := startTime + uint64(idx)
		if i.History.Time >= idxTime {
			// we've already filled this spot
			continue
		}
		i.History.Samples[i.History.Position] = val
		i.History.Position = (i.History.Position + 1) % len(i.History.Samples)
		i.History.Time = idxTime
	}
}

func GetValue(g Generator, t, r uint64) float32 {
	// TODO: use history as a cache
	val := g.GetValue(t, r)
	AddHistory(g, t, []float32{val})
	return val
}

func (i Info) Copy(historyLen int) Info {
	initialLength := 0
	if i.History != nil {
		initialLength = len(i.History.Samples)
	}

	i.History = getEmptyHistory()
	if historyLen == UseDefaultHistoryLength {
		return i
	}

	if historyLen == CopyExistingHistoryLength {
		historyLen = initialLength
	}
	if historyLen == 0 {
		i.History.Samples = nil
	} else {
		i.History.Samples = make([]float32, historyLen)
	}
	return i
}

func (h History) GetOrdered() []float32 {
	return append(h.Samples[h.Position:], h.Samples[0:h.Position]...)
}
