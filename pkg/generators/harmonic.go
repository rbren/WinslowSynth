package generators

type Harmonic struct {
	Info       Info
	Modes      []Mode
	Oscillator Oscillator
	Average    Average
}

type Mode struct {
	Amplitude float32
	Frequency float32
}

func (h Harmonic) SubGenerators() []Generator {
	return []Generator{h.Average}
}

func (h Harmonic) Initialize(name string) Generator {
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
	h.Average = h.Average.Initialize(name).(Average)
	h.Oscillator = h.Oscillator.Initialize(name).(Oscillator)
	return h
}

func (h Harmonic) GetValue(t, r uint64) float32 {
	return GetValue(h.Average, t, r)
}

func (h Harmonic) GetInfo() Info { return h.Info }
func (h Harmonic) Copy(historyLen int) Generator {
	h.Info = h.Info.Copy(historyLen)
	h.Average = h.Average.Copy(CopyExistingHistoryLength).(Average)
	h.Oscillator = h.Oscillator.Copy(CopyExistingHistoryLength).(Oscillator)
	return h
}
