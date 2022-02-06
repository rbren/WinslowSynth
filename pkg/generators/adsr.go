package generators

import (
	"github.com/rbren/midi/pkg/config"
)

type ADSR struct {
	Info         Info
	PeakLevel    Constant
	AttackTime   Constant
	DecayTime    Constant
	SustainLevel Constant
	ReleaseTime  Constant
}

func getTimeInSamples(ms float32) uint64 {
	samplesPerMs := config.MainConfig.SampleRate / 1000
	return uint64(int(ms) * samplesPerMs)
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
	return a
}
