package generators

import (
	"github.com/rbren/midi/pkg/config"
)

type Delay struct {
	Info   *Info
	Amount Generator
	Input  Instrument
}

func NewDelay(input Generator, amt Generator) Delay {
	input.SetInfo(Info{
		History: getEmptyHistory(),
	})
	return Delay{
		Info:   &Info{},
		Input:  input,
		Amount: amt,
	}
}

func (d Delay) GetValue(t, r uint64) float32 {
	samplesPerMs := config.MainConfig.SampleRate / 1000
	amtMs := getValue(d.Amount, t, r)
	amtSamples := int(amtMs) * samplesPerMs
	inputInfo := d.Input.GetInfo()
	valueIndex := inputInfo.HistoryPosition - 1 - amtSamples
	if valueIndex < 0 {
		valueIndex = len(inputInfo.History) + valueIndex
	}
	val := inputInfo.History[valueIndex]
	return val
}

func (d Delay) GetInfo() *Info    { return d.Info }
func (d Delay) SetInfo(info Info) { copyInfo(d.Info, info) }
