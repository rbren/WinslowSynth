package generators

type Harmonic struct {
	Info          Info
	Modes         []Mode
	Oscillator    Oscillator
	SubGenerators SubGenerators
}

type Mode struct {
	Amplitude float32
	Frequency float32
}

func (h Harmonic) GetInfo() Info                   { return h.Info }
func (h Harmonic) GetSubGenerators() SubGenerators { return h.SubGenerators }
func (h Harmonic) Copy(historyLen int) Generator {
	h.Info = h.Info.Copy(historyLen)
	h.SubGenerators = h.SubGenerators.Copy()
	return h
}

func (h Harmonic) Initialize(name string) Generator {
	if h.SubGenerators == nil {
		h.SubGenerators = make(map[string]Generator)
	}
	toAverage := []Generator{h.Oscillator}
	for _, mode := range h.Modes {
		amp := Multiply{
			Generators: []Generator{h.Oscillator.SubGenerators["Amplitude"], Constant{Value: mode.Amplitude}},
		}
		freq := Multiply{
			Generators: []Generator{h.Oscillator.SubGenerators["Frequency"], Constant{Value: mode.Frequency}},
		}
		modeGenerator := Oscillator{
			SubGenerators: map[string]Generator{
				"Amplitude": amp,
				"Frequency": freq,
				"Phase":     h.Oscillator.SubGenerators["Phase"],
			},
		}
		toAverage = append(toAverage, modeGenerator)
	}
	h.SubGenerators["Average"] = Average{Generators: toAverage}
	h.SubGenerators["Average"] = h.SubGenerators["Average"].Initialize(name).(Average)
	h.Oscillator = h.Oscillator.Initialize(name).(Oscillator)
	return h
}

func (h Harmonic) GetValue(t, r uint64) float32 {
	return GetValue(h.SubGenerators["Average"], t, r)
}
