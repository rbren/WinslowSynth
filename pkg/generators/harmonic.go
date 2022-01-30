package generators

type Harmonic struct {
	Modes   []Mode
	Spinner Spinner
	Sum     Sum
}

type Mode struct {
	Amplitude float32
	Frequency float32
}

func (h *Harmonic) initialize(force bool) {
	if !force && len(h.Sum.Generators) == len(h.Modes) {
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
	h.Sum = Sum{Generators: toSum}
}

func (h Harmonic) GetValue(t, r uint64) float32 {
	h.initialize(false)
	return h.Sum.GetValue(t, r)
}

func (h Harmonic) SetFrequency(f float32) Generator {
	h.Spinner.SetFrequency(f)
	h.initialize(true)
	return h
}
