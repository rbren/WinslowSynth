package generators

type ADSR struct {
	Info         *Info
	PeakLevel    Constant
	SustainLevel Constant
	AttackTime   uint64
	DecayTime    uint64
	ReleaseTime  uint64
}

func (a ADSR) GetValue(t, r uint64) float32 {
	if r == 0 {
		if t < a.AttackTime {
			return a.Attack(t, r)
		}
		if t < a.DecayTime {
			return a.Decay(t, r)
		}
		return a.SustainLevel.GetValue(t, r)
	}
	return a.Release(t, r)
}

func (a ADSR) Attack(t, r uint64) float32 {
	percentDone := float32(t) / float32(a.AttackTime)
	return a.PeakLevel.GetValue(t, r) * percentDone
}

func (a ADSR) Decay(t, r uint64) float32 {
	timeInDecay := t - a.AttackTime
	percentDone := float32(timeInDecay) / float32(a.DecayTime)
	levelDiff := a.PeakLevel.GetValue(t, r) - a.SustainLevel.GetValue(t, r)
	return a.SustainLevel.GetValue(t, r) + float32(levelDiff)*(1.0-percentDone)
}

func (a ADSR) Release(t, r uint64) float32 {
	timeSinceRelease := t - r
	if timeSinceRelease > a.ReleaseTime {
		return 0.0
	}
	baseVal := a.SustainLevel.GetValue(t, r)
	if t < a.AttackTime {
		baseVal = a.Attack(r, r)
	} else if t < a.DecayTime {
		baseVal = a.Decay(r, r)
	}
	percentDone := float32(timeSinceRelease) / float32(a.ReleaseTime)
	return baseVal * (1.0 - percentDone)
}

func (a ADSR) GetInfo() *Info { return a.Info }
