package generators

type SawWave struct {
	Amplitude Generator
	Frequency Generator
	Phase     Generator
}

func (s *SawWave) initialize() {
	if s.Amplitude == nil {
		s.Amplitude = Constant{Value: 1.0}
	}
	if s.Phase == nil {
		s.Phase = Constant{Value: 0.0}
	}
}

func (s SawWave) GetValue(time, releasedAt uint64) float32 {
	s.initialize()
	return s.Amplitude.GetValue(time, releasedAt) * GetPhasePosition(s.Frequency, s.Phase, time, releasedAt)
}
