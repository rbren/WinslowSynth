package generators

type ADSR struct {
	Info         Info
	PeakLevel    Generator
	AttackTime   Generator
	DecayTime    Generator
	SustainLevel Generator
	ReleaseTime  Generator
}

func (a ADSR) SubGenerators() []Generator {
	return []Generator{a.PeakLevel, a.SustainLevel, a.AttackTime, a.DecayTime, a.ReleaseTime}
}

func (a ADSR) Initialize(name string) Generator {
	if a.PeakLevel == nil {
		a.PeakLevel = Constant{
			Info:  Info{Group: name, Subgroup: "ADSR", Name: "Peak"},
			Value: 1.0,
			Min:   0.0,
			Max:   1.0,
		}
	}
	if a.SustainLevel == nil {
		a.SustainLevel = Constant{
			Info:  Info{Group: name, Subgroup: "ADSR", Name: "Sustain"},
			Value: 0.8,
			Min:   0.0,
			Max:   1.0,
		}
	}
	if a.AttackTime == nil {
		a.AttackTime = Constant{
			Info:  Info{Group: name, Subgroup: "ADSR", Name: "Attack"},
			Value: 300,
			Min:   0.0,
			Max:   1000,
			Step:  1.0,
		}
	}
	if a.DecayTime == nil {
		a.DecayTime = Constant{
			Info:  Info{Group: name, Subgroup: "ADSR", Name: "Decay"},
			Value: 500,
			Min:   0.0,
			Max:   1000,
			Step:  1.0,
		}
	}
	if a.ReleaseTime == nil {
		a.ReleaseTime = Constant{
			Info:  Info{Group: name, Subgroup: "ADSR", Name: "Release"},
			Value: 1000,
			Min:   0.0,
			Max:   3000,
			Step:  1.0,
		}
	}
	a.PeakLevel = a.PeakLevel.Initialize(name)
	a.SustainLevel = a.SustainLevel.Initialize(name)
	a.AttackTime = a.AttackTime.Initialize(name)
	a.DecayTime = a.DecayTime.Initialize(name)
	a.ReleaseTime = a.ReleaseTime.Initialize(name)
	return a
}

func (a ADSR) GetValue(t, r uint64) float32 {
	attackTime := getTimeInSamples(GetValue(a.AttackTime, t, r))
	if t < attackTime {
		return a.Attack(t, r)
	}
	decayTime := getTimeInSamples(GetValue(a.DecayTime, t, r))
	if t < attackTime+decayTime {
		return a.Decay(t, r)
	}
	if r == 0 {
		return GetValue(a.SustainLevel, t, r)
	}
	return a.Release(t, r)
}

func (a ADSR) Attack(t, r uint64) float32 {
	percentDone := float32(t) / float32(getTimeInSamples(GetValue(a.AttackTime, t, r)))
	return GetValue(a.PeakLevel, t, r) * percentDone
}

func (a ADSR) Decay(t, r uint64) float32 {
	timeInDecay := t - getTimeInSamples(GetValue(a.AttackTime, t, r))
	percentDone := float32(timeInDecay) / float32(getTimeInSamples(GetValue(a.DecayTime, t, r)))
	levelDiff := GetValue(a.PeakLevel, t, r) - GetValue(a.SustainLevel, t, r)
	return GetValue(a.SustainLevel, t, r) + float32(levelDiff)*(1.0-percentDone)
}

func (a ADSR) Release(t, r uint64) float32 {
	minTimeOfRelease := getTimeInSamples(GetValue(a.AttackTime, t, r) + GetValue(a.DecayTime, t, r))
	timeOfRelease := r
	if timeOfRelease < minTimeOfRelease {
		timeOfRelease = minTimeOfRelease
	}
	timeSinceRelease := t - timeOfRelease
	desiredReleaseTime := getTimeInSamples(GetValue(a.ReleaseTime, t, r))
	if timeSinceRelease > desiredReleaseTime {
		return 0.0
	}
	baseVal := GetValue(a.SustainLevel, t, r)
	percentDone := float32(timeSinceRelease) / float32(desiredReleaseTime)
	return baseVal * (1.0 - percentDone)
}
func (a ADSR) GetInfo() Info { return a.Info }
func (a ADSR) Copy(historyLen int) Generator {
	a.Info = a.Info.Copy(historyLen)
	a.PeakLevel = a.PeakLevel.Copy(CopyExistingHistoryLength)
	a.AttackTime = a.AttackTime.Copy(CopyExistingHistoryLength)
	a.DecayTime = a.DecayTime.Copy(CopyExistingHistoryLength)
	a.SustainLevel = a.SustainLevel.Copy(CopyExistingHistoryLength)
	a.ReleaseTime = a.ReleaseTime.Copy(CopyExistingHistoryLength)
	return a
}
