package generators

type ADSR struct {
	Info          Info
	SubGenerators SubGenerators
}

func (a ADSR) GetSubGenerators() SubGenerators { return a.SubGenerators }

func (a ADSR) Initialize(name string) Generator {
	if a.SubGenerators == nil {
		a.SubGenerators = map[string]Generator{}
	}
	if a.SubGenerators["PeakLevel"] == nil {
		a.SubGenerators["PeakLevel"] = Constant{
			Info:  Info{Group: name, Subgroup: "ADSR", Name: "Peak"},
			Value: 1.0,
			Min:   0.0,
			Max:   1.0,
		}
	}
	if a.SubGenerators["SustainLevel"] == nil {
		a.SubGenerators["SustainLevel"] = Constant{
			Info:  Info{Group: name, Subgroup: "ADSR", Name: "Sustain"},
			Value: 0.8,
			Min:   0.0,
			Max:   1.0,
		}
	}
	if a.SubGenerators["AttackTime"] == nil {
		a.SubGenerators["AttackTime"] = Constant{
			Info:  Info{Group: name, Subgroup: "ADSR", Name: "Attack"},
			Value: 300,
			Min:   0.0,
			Max:   1000,
			Step:  1.0,
		}
	}
	if a.SubGenerators["DecayTime"] == nil {
		a.SubGenerators["DecayTime"] = Constant{
			Info:  Info{Group: name, Subgroup: "ADSR", Name: "Decay"},
			Value: 500,
			Min:   0.0,
			Max:   1000,
			Step:  1.0,
		}
	}
	if a.SubGenerators["ReleaseTime"] == nil {
		a.SubGenerators["ReleaseTime"] = Constant{
			Info:  Info{Group: name, Subgroup: "ADSR", Name: "Release"},
			Value: 1000,
			Min:   0.0,
			Max:   3000,
			Step:  1.0,
		}
	}
	for key, g := range a.SubGenerators {
		a.SubGenerators[key] = g.Initialize(name)
	}
	return a
}

func (a ADSR) GetValue(t, r uint64) float32 {
	attackTime := getTimeInSamples(GetValue(a.SubGenerators["AttackTime"], t, r))
	if t < attackTime {
		return a.Attack(t, r)
	}
	decayTime := getTimeInSamples(GetValue(a.SubGenerators["DecayTime"], t, r))
	if t < attackTime+decayTime {
		return a.Decay(t, r)
	}
	if r == 0 {
		return GetValue(a.SubGenerators["SustainLevel"], t, r)
	}
	return a.Release(t, r)
}

func (a ADSR) Attack(t, r uint64) float32 {
	percentDone := float32(t) / float32(getTimeInSamples(GetValue(a.SubGenerators["AttackTime"], t, r)))
	return GetValue(a.SubGenerators["PeakLevel"], t, r) * percentDone
}

func (a ADSR) Decay(t, r uint64) float32 {
	timeInDecay := t - getTimeInSamples(GetValue(a.SubGenerators["AttackTime"], t, r))
	percentDone := float32(timeInDecay) / float32(getTimeInSamples(GetValue(a.SubGenerators["DecayTime"], t, r)))
	levelDiff := GetValue(a.SubGenerators["PeakLevel"], t, r) - GetValue(a.SubGenerators["SustainLevel"], t, r)
	return GetValue(a.SubGenerators["SustainLevel"], t, r) + float32(levelDiff)*(1.0-percentDone)
}

func (a ADSR) Release(t, r uint64) float32 {
	minTimeOfRelease := getTimeInSamples(GetValue(a.SubGenerators["AttackTime"], t, r) + GetValue(a.SubGenerators["DecayTime"], t, r))
	timeOfRelease := r
	if timeOfRelease < minTimeOfRelease {
		timeOfRelease = minTimeOfRelease
	}
	timeSinceRelease := t - timeOfRelease
	desiredReleaseTime := getTimeInSamples(GetValue(a.SubGenerators["ReleaseTime"], t, r))
	if timeSinceRelease > desiredReleaseTime {
		return 0.0
	}
	baseVal := GetValue(a.SubGenerators["SustainLevel"], t, r)
	percentDone := float32(timeSinceRelease) / float32(desiredReleaseTime)
	return baseVal * (1.0 - percentDone)
}
func (a ADSR) GetInfo() Info { return a.Info }
func (a ADSR) Copy(historyLen int) Generator {
	a.Info = a.Info.Copy(historyLen)
	for key, g := range a.SubGenerators {
		a.SubGenerators[key] = g.Copy(CopyExistingHistoryLength)
	}
	return a
}
