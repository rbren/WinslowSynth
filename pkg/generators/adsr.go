package generators

type ADSR struct {
	PeakLevel    float32
	SustainLevel float32
	AttackTime   uint64
	DecayTime    uint64
	ReleaseTime  uint64
}

func (a ADSR) GetValue(t, r uint64) float32 {
	if r == 0 {
		if t < a.AttackTime {
			return a.Attack(t)
		}
		if t < a.DecayTime {
			return a.Decay(t)
		}
		return a.SustainLevel
	}
	return a.Release(t, r)
}

func (a ADSR) Attack(t uint64) float32 {
	percentDone := float32(t) / float32(a.AttackTime)
	return a.PeakLevel * percentDone
}

func (a ADSR) Decay(t uint64) float32 {
	timeInDecay := t - a.AttackTime
	percentDone := float32(timeInDecay) / float32(a.DecayTime)
	levelDiff := a.PeakLevel - a.SustainLevel
	return a.SustainLevel + float32(levelDiff)*(1.0-percentDone)
}

func (a ADSR) Release(t, r uint64) float32 {
	timeSinceRelease := t - r
	if timeSinceRelease > a.ReleaseTime {
		return 0.0
	}
	baseVal := a.SustainLevel
	if t < a.AttackTime {
		baseVal = a.Attack(r)
	} else if t < a.DecayTime {
		baseVal = a.Decay(r)
	}
	percentDone := float32(timeSinceRelease) / float32(a.ReleaseTime)
	return baseVal * (1.0 - percentDone)
}
