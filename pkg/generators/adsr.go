package generators

import (
	"github.com/rbren/midi/pkg/config"
)

type ADSR struct {
	Info         *Info
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
	attackTime := getTimeInSamples(getValue(a.AttackTime, t, r))
	if t < attackTime {
		return a.Attack(t, r)
	}
	decayTime := getTimeInSamples(getValue(a.DecayTime, t, r))
	if t < attackTime+decayTime {
		return a.Decay(t, r)
	}
	if r == 0 {
		return getValue(a.SustainLevel, t, r)
	}
	return a.Release(t, r)
}

func (a ADSR) Attack(t, r uint64) float32 {
	percentDone := float32(t) / float32(getTimeInSamples(getValue(a.AttackTime, t, r)))
	return getValue(a.PeakLevel, t, r) * percentDone
}

func (a ADSR) Decay(t, r uint64) float32 {
	timeInDecay := t - getTimeInSamples(getValue(a.AttackTime, t, r))
	percentDone := float32(timeInDecay) / float32(getTimeInSamples(getValue(a.DecayTime, t, r)))
	levelDiff := getValue(a.PeakLevel, t, r) - getValue(a.SustainLevel, t, r)
	return getValue(a.SustainLevel, t, r) + float32(levelDiff)*(1.0-percentDone)
}

func (a ADSR) Release(t, r uint64) float32 {
	minTimeOfRelease := getTimeInSamples(getValue(a.AttackTime, t, r) + getValue(a.DecayTime, t, r))
	timeOfRelease := r
	if timeOfRelease < minTimeOfRelease {
		timeOfRelease = minTimeOfRelease
	}
	timeSinceRelease := t - timeOfRelease
	desiredReleaseTime := getTimeInSamples(getValue(a.ReleaseTime, t, r))
	if timeSinceRelease > desiredReleaseTime {
		return 0.0
	}
	baseVal := getValue(a.SustainLevel, t, r)
	percentDone := float32(timeSinceRelease) / float32(desiredReleaseTime)
	return baseVal * (1.0 - percentDone)
}
func (a ADSR) GetInfo() *Info    { return a.Info }
func (a ADSR) SetInfo(info Info) { copyInfo(a.Info, info) }
