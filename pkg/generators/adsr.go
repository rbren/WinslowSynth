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
	attackTime := getTimeInSamples(a.AttackTime.GetValue(t, r))
	if t < attackTime {
		return a.Attack(t, r)
	}
	if t < attackTime+getTimeInSamples(a.DecayTime.GetValue(t, r)) {
		return a.Decay(t, r)
	}
	if r == 0 {
		return a.SustainLevel.GetValue(t, r)
	}
	return a.Release(t, r)
}

func (a ADSR) Attack(t, r uint64) float32 {
	percentDone := float32(t) / float32(getTimeInSamples(a.AttackTime.GetValue(t, r)))
	return a.PeakLevel.GetValue(t, r) * percentDone
}

func (a ADSR) Decay(t, r uint64) float32 {
	timeInDecay := t - getTimeInSamples(a.AttackTime.GetValue(t, r))
	percentDone := float32(timeInDecay) / float32(getTimeInSamples(a.DecayTime.GetValue(t, r)))
	levelDiff := a.PeakLevel.GetValue(t, r) - a.SustainLevel.GetValue(t, r)
	return a.SustainLevel.GetValue(t, r) + float32(levelDiff)*(1.0-percentDone)
}

func (a ADSR) Release(t, r uint64) float32 {
	timeSinceRelease := t - r
	desiredReleaseTime := getTimeInSamples(a.ReleaseTime.GetValue(t, r))
	if timeSinceRelease > desiredReleaseTime {
		return 0.0
	}
	baseVal := a.SustainLevel.GetValue(t, r)
	percentDone := float32(timeSinceRelease) / float32(desiredReleaseTime)
	return baseVal * (1.0 - percentDone)
}
func (a ADSR) GetInfo() *Info    { return a.Info }
func (a ADSR) SetInfo(info Info) { copyInfo(a.Info, info) }
