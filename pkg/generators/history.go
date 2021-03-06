package generators

import (
	"math"
	"math/cmplx"
	"os"

	dsp "github.com/mjibson/go-dsp/fft"
	"github.com/sirupsen/logrus"

	"github.com/rbren/midi/pkg/buffers"
	"github.com/rbren/midi/pkg/config"
)

var useHistory bool
var useFrequencies bool

const CopyExistingHistoryLength = -1
const UseDefaultHistoryLength = -2

var frequencyCoefficients []complex64
var historyLength int

func init() {
	historyLength = config.MainConfig.HistoryMs * (config.MainConfig.SampleRate / 1000)
	logrus.Infof("Using %d bins by default", config.MainConfig.FrequencyBins)
	useHistory = os.Getenv("NO_HISTORY") == ""
	useFrequencies = os.Getenv("NO_FREQUENCIES") == ""

	precomputeFrequencyCoefficients()
}

func precomputeFrequencyCoefficients() {
	twoPiI := 2 * math.Pi * complex(0, 1)
	frequencyCoefficients = make([]complex64, config.MainConfig.FrequencyBins)
	for i := 0; i < config.MainConfig.FrequencyBins; i++ {
		frac := float64(i) / float64(config.MainConfig.FrequencyBins)
		frequencyCoefficients[i] = complex64(cmplx.Exp(twoPiI * complex(frac, 0)))
	}
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

func getEmptyHistoryWithFrequencies() *History {
	h := getEmptyHistory()
	h.frequencyBins = make([]complex64, config.MainConfig.FrequencyBins)
	return h
}

func AddHistory(g Generator, startTime uint64, history []float32) {
	if !useHistory {
		return
	}
	i := g.GetInfo()
	if i.History == nil || i.History.samples == nil {
		return
	}
	i.History.Add(startTime, history)
}

func GetValue(g Generator, t, r uint64) float32 {
	cached := GetValueCached(g, t)
	if cached != nil {
		return *cached
	}
	val := g.GetValue(t, r)
	AddHistory(g, t, []float32{val})
	return val
}

func GetValueCached(g Generator, t uint64) *float32 {
	if !useHistory {
		return nil
	}
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

func (i Info) Copy(historyLen int, useFrequencies bool) Info {
	initialLength := 0
	if i.History != nil {
		initialLength = len(i.History.samples)
	}

	if useFrequencies {
		i.History = getEmptyHistoryWithFrequencies()
	} else {
		i.History = getEmptyHistory()
	}
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
	doInterpolation := h.Time != 0 && timeSinceLastSample > 1
	historyLength := len(h.samples)
	earliestNewPos := buffers.Modulus(h.Position+int(timeSinceLastSample), historyLength)
	for idx := range samples {
		idxTime := startTime + uint64(idx)
		if h.Time != 0 && idxTime <= h.Time {
			// we've already filled this spot
			continue
		}
		idxPos := buffers.Modulus((earliestNewPos + idx), historyLength)
		h.samples[idxPos] = samples[idx]
		h.Position = idxPos
		h.Time = idxTime
	}
	if doInterpolation {
		buffers.InterpolateEvents(h.samples, origPos, earliestNewPos)
	}
	h.UpdateFrequencies(origPos, h.Position)
}

func (h *History) UpdateFrequencies(startPos, endPos int) {
	if !useFrequencies {
		return
	}
	numBins := len(h.frequencyBins)
	if numBins == 0 {
		return
	}
	historyLength := len(h.samples)
	startIdx := buffers.Modulus(startPos+1, historyLength)
	endIdx := buffers.Modulus(endPos+1, historyLength)
	for posIdx := startIdx; posIdx != endIdx; posIdx = buffers.Modulus(posIdx+1, historyLength) {
		newSample := h.samples[posIdx]
		oldIdx := buffers.Modulus(posIdx-numBins, historyLength)
		oldSample := h.samples[oldIdx]
		diff := complex64(complex(newSample-oldSample, 0))
		for binIdx, binValue := range h.frequencyBins {
			h.frequencyBins[binIdx] = frequencyCoefficients[binIdx] * (binValue + diff)
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

func (h History) GetFrequencies() []float32 {
	out := make([]float32, len(h.frequencyBins))
	for idx := range out {
		out[idx] = float32(cmplx.Abs(complex128(h.frequencyBins[idx])))
	}
	return out
}

func (h History) CalculateFrequencies() []float32 {
	samples := h.GetOrdered(len(h.samples))
	samples64 := make([]float64, len(samples))
	for idx := range samples {
		samples64[idx] = float64(samples[idx])
	}
	transformed := dsp.FFTReal(samples64)
	freqs := make([]float32, len(transformed))
	for idx := range freqs {
		freqs[idx] = float32(cmplx.Abs(transformed[idx]))
	}
	return freqs
}
