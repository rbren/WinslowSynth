package generators

type SquareWave struct {
	Amplitude Generator
	Frequency Generator
	Phase     Generator
}

func (s *SquareWave) initialize() {
	if s.Amplitude == nil {
		s.Amplitude = Constant{Value: 1.0}
	}
	if s.Phase == nil {
		s.Phase = Constant{Value: 0.0}
	}
}

func (s SquareWave) GetValue(time, releasedAt uint64) float32 {
	s.initialize()
	phasePos := GetPhasePosition(s.Frequency, s.Phase, time, releasedAt)
	var val float32 = 1.0
	if phasePos > .5 {
		val = -1.0
	}
	return val * s.Amplitude.GetValue(time, releasedAt)
}
