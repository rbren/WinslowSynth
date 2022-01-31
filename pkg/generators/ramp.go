package generators

import (
	_ "github.com/rbren/midi/pkg/logger"
)

type Ramp struct {
	Info     *Info
	RampUp   uint64
	RampDown uint64
	Target   float32
}

func (r Ramp) GetValue(t, releasedAt uint64) float32 {
	if releasedAt == 0 && t < r.RampUp {
		return r.RampUpVal(t)
	}
	if releasedAt != 0 {
		return r.RampDownVal(t, releasedAt)
	}
	return r.Target
}

func (r Ramp) RampUpVal(t uint64) float32 {
	percentDone := float32(t) / float32(r.RampUp)
	return r.Target * percentDone
}

func (r Ramp) RampDownVal(t, releasedAt uint64) float32 {
	timeSinceRelease := t - releasedAt
	if timeSinceRelease > r.RampDown {
		return 0.0
	}
	startVal := r.Target
	if t < r.RampUp {
		startVal = r.RampUpVal(releasedAt)
	}

	percentDone := float32(timeSinceRelease) / float32(r.RampDown)
	return startVal * (1.0 - percentDone)
}

func (r Ramp) GetInfo() *Info { return r.Info }
