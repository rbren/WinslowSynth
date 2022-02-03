package generators

type Harmonic struct {
	Info       *Info
	Modes      []Mode
	Oscillator Oscillator
	Average    Average
}

type Mode struct {
	Amplitude float32
	Frequency float32
}

func (h *Harmonic) initialize(force bool) {
	if !force && len(h.Average.Generators) == len(h.Modes) {
		return
	}
	toAverage := []Generator{h.Oscillator}
	for _, mode := range h.Modes {
		amp := Multiply{
			Generators: []Generator{h.Oscillator.Amplitude, Constant{Value: mode.Amplitude}},
		}
		freq := Multiply{
			Generators: []Generator{h.Oscillator.Frequency, Constant{Value: mode.Frequency}},
		}
		modeGenerator := Oscillator{
			Amplitude: amp,
			Frequency: freq,
			Phase:     h.Oscillator.Phase,
		}
		toAverage = append(toAverage, modeGenerator)
	}
	h.Average = Average{Generators: toAverage}
}

func (h Harmonic) GetValue(t, r uint64) float32 {
	h.initialize(false)
	return h.Average.GetValue(t, r)
}

func (h Harmonic) GetInfo() *Info    { return h.Info }
func (h Harmonic) SetInfo(info Info) { copyInfo(h.Info, info) }
