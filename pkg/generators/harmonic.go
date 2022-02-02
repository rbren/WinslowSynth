package generators

type Harmonic struct {
	Info    *Info
	Modes   []Mode
	Spinner Spinner
	Average     Average
}

type Mode struct {
	Amplitude float32
	Frequency float32
}

func (h *Harmonic) initialize(force bool) {
	if !force && len(h.Average.Generators) == len(h.Modes) {
		return
	}
	toAverage := []Generator{h.Spinner}
	for _, mode := range h.Modes {
		amp := Multiply{
			Generators: []Generator{h.Spinner.Amplitude, Constant{Value: mode.Amplitude}},
		}
		freq := Multiply{
			Generators: []Generator{h.Spinner.Frequency, Constant{Value: mode.Frequency}},
		}
		modeGenerator := Spinner{
			Amplitude: amp,
			Frequency: freq,
			Phase:     h.Spinner.Phase,
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
