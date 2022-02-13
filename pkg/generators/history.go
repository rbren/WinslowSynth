package generators

import (
	"os"

	_ "github.com/sirupsen/logrus"

	"github.com/rbren/midi/pkg/buffers"
	"github.com/rbren/midi/pkg/config"
)

var useHistory bool

const CopyExistingHistoryLength = -1
const UseDefaultHistoryLength = -2

var historyMs = 5000
var historyLength int

func init() {
	historyLength = historyMs * (config.MainConfig.SampleRate / 1000)
	useHistory = os.Getenv("NO_HISTORY") == ""
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
	origPos := i.History.Position
	timeSinceLastSample := startTime - i.History.Time
	earliestNewPos := (i.History.Position + int(timeSinceLastSample)) % len(i.History.samples)
	for idx := range history {
		idxTime := startTime + uint64(idx)
		if i.History.Time != 0 && idxTime <= i.History.Time {
			// we've already filled this spot
			continue
		}
		idxPos := (earliestNewPos + idx) % len(i.History.samples)
		i.History.samples[idxPos] = history[idx]
		i.History.Position = idxPos
		i.History.Time = idxTime
	}
	if timeSinceLastSample > 1 {
		buffers.InterpolateEvents(i.History.samples, origPos, earliestNewPos)
	}
}

func GetValue(g Generator, t, r uint64) float32 {
	if useHistory {
		cached := GetValueCached(g, t)
		if cached != nil {
			return *cached
		}
	}
	val := g.GetValue(t, r)
	if useHistory {
		AddHistory(g, t, []float32{val})
	}
	return val
}

func GetValueCached(g Generator, t uint64) *float32 {
	info := g.GetInfo()
	if info.History != nil && info.History.samples != nil {
		if timeDiff := info.History.Time - t; timeDiff >= 0 && timeDiff < uint64(len(info.History.samples)) {
			idx := (info.History.Position - 1 - int(timeDiff)) % len(info.History.samples)
			if idx < 0 {
				idx = idx + len(info.History.samples)
			}
			val := info.History.samples[idx]
			return &val
		}
	}
	return nil
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
