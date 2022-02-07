package generators

import (
	_ "github.com/sirupsen/logrus"

	"github.com/rbren/midi/pkg/config"
)

const CopyExistingHistoryLength = -1
const UseDefaultHistoryLength = -2

var historyMs = 5000
var historyLength int

func init() {
	historyLength = historyMs * (config.MainConfig.SampleRate / 1000)
}

func getEmptyHistory() *History {
	return &History{
		samples: make([]float32, historyLength),
	}
}

func AddHistory(g Generator, startTime uint64, history []float32) {
	i := g.GetInfo()
	if i.History == nil || i.History.samples == nil {
		return
	}
	for idx, val := range history {
		idxTime := startTime + uint64(idx)
		if i.History.Time >= idxTime {
			// we've already filled this spot
			continue
		}
		i.History.samples[i.History.Position] = val
		i.History.Position = (i.History.Position + 1) % len(i.History.samples)
		i.History.Time = idxTime
	}
}

func GetValue(g Generator, t, r uint64) float32 {
	// TODO: use history as a cache
	info := g.GetInfo()
	if info.History != nil && info.History.samples != nil {
		if timeDiff := info.History.Time - t; timeDiff >= 0 && timeDiff < uint64(len(info.History.samples)) {
			idx := (info.History.Position - 1 - int(timeDiff)) % len(info.History.samples)
			if idx < 0 {
				idx = idx + len(info.History.samples)
			}
			return info.History.samples[idx]
		}
	}
	val := g.GetValue(t, r)
	AddHistory(g, t, []float32{val})
	return val
}

func (i Info) Copy(historyLen int) Info {
	initialLength := 0
	if i.History != nil {
		initialLength = len(i.History.samples)
	}

	i.History = getEmptyHistory()
	if historyLen == UseDefaultHistoryLength {
		return i
	}

	if historyLen == CopyExistingHistoryLength {
		historyLen = initialLength
	}
	if historyLen == 0 {
		i.History.samples = nil
	} else {
		i.History.samples = make([]float32, historyLen)
	}
	return i
}

func (h History) GetOrdered(numSamples int) []float32 {
	ordered := append(h.samples[h.Position:], h.samples[0:h.Position]...)
	if numSamples != -1 {
		return ordered[len(ordered)-numSamples:]
	}
	return ordered
}
