package generators

type Harmonic struct {
	Modes   []Mode
	Spinner Spinner
	Sum     *Sum
}

type Mode struct {
	Amplitude float32
	Frequency float32
}

func (h *Harmonic) initialize() {
	if h.Sum != nil {
		return
	}
	toSum := []Generator{h.Spinner}
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
		toSum = append(toSum, modeGenerator)
	}
	h.Sum = &Sum{Generators: toSum}
}

func (h Harmonic) GetValue(t, r uint64) float32 {
	h.initialize()
	return h.Sum.GetValue(t, r)
}