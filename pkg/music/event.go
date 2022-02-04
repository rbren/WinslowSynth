package music

type Event struct {
	Key         int64
	Frequency   float32
	Velocity    int64
	AttackTime  uint64
	ReleaseTime uint64
}

type EventType int

const (
	AttackEvent EventType = iota
	ReleaseEvent
)

func (e Event) getRelativeTime(absoluteTime uint64) (uint64, uint64) {
	elapsed := absoluteTime - e.AttackTime
	var releasedAt uint64 = 0
	if e.ReleaseTime != 0 {
		elapsedSinceRelease := absoluteTime - e.ReleaseTime
		releasedAt = elapsed - elapsedSinceRelease
	}
	return elapsed, releasedAt
}
