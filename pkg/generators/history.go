package generators

import (
	"math"
	"math/cmplx"
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

type History struct {
	samples       []float32
	frequencyBins []complex64
	Position      int
	Time          uint64
}

func getEmptyHistory() *History {
	return &History{
		samples: make([]float32, historyLength),
	}
}

func getEmptyHistoryWithFrequencies(numBins int) *History {
	h := getEmptyHistory()
	h.frequencyBins = make([]complex64, numBins)
	return h
}

func AddHistory(g Generator, startTime uint64, history []float32) {
	i := g.GetInfo()
	if i.History == nil || i.History.samples == nil {
		return
	}
	i.History.Add(startTime, history)
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
			idx := (info.History.Position - int(timeDiff)) % len(info.History.samples)
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

func (h *History) Add(startTime uint64, samples []float32) {
	origPos := h.Position
	timeSinceLastSample := startTime - h.Time
	earliestNewPos := (h.Position + int(timeSinceLastSample)) % len(h.samples)
	removedSamples := []float32{}
	newSamples := []float32{}
	for idx := range samples {
		idxTime := startTime + uint64(idx)
		if h.Time != 0 && idxTime <= h.Time {
			// we've already filled this spot
			continue
		}
		idxPos := (earliestNewPos + idx) % len(h.samples)
		removedSamples = append(removedSamples, h.samples[idxPos])
		newSamples = append(newSamples, samples[idx])
		h.samples[idxPos] = samples[idx]
		h.Position = idxPos
		h.Time = idxTime
	}
	if timeSinceLastSample > 1 {
		buffers.InterpolateEvents(h.samples, origPos, earliestNewPos)
	}
	h.UpdateFrequencies(removedSamples, newSamples)
}

func (h *History) UpdateFrequencies(removedSamples, newSamples []float32) {
	if len(h.frequencyBins) == 0 {
		return
	}
	for sampleIdx := range removedSamples {
		oldSample := removedSamples[sampleIdx]
		newSample := newSamples[sampleIdx]
		for binIdx, binValue := range h.frequencyBins {
			coeff := cmplx.Exp(2 * math.Pi * complex(0, 1) * complex(float64(binIdx), 0)) // TODO: precompute
			h.frequencyBins[binIdx] = complex64(coeff) * (binValue + complex64(complex(newSample-oldSample, 0)))
		}
	}
}

func (h History) GetOrdered(numSamples int) []float32 {
	startIdx := (h.Position + 1) % len(h.samples)
	ordered := append(h.samples[startIdx:], h.samples[0:startIdx]...)
	if numSamples != -1 {
		return ordered[len(ordered)-numSamples:]
	}
	return ordered
}

func (h History) GetOrderedComplex() []complex128 {
	out := make([]complex128, len(h.samples))
	startIdx := (h.Position + 1) % len(h.samples)
	for i := 0; i < len(out); i++ {
		idx := (startIdx + i) % len(h.samples)
		out[i] = complex(float64(h.samples[idx]), 0)
	}
	return out
}
