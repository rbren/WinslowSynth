package performance

import (
	"math/rand"
	"time"

	"github.com/rbren/midi/pkg/generators"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func CheckPerformance(g generators.Generator, duration int) float64 {
	g = g.Copy(generators.UseDefaultHistoryLength, true)
	generators.SetFrequency(g, 440.0)

	start := time.Now()
	releaseTime := rand.Intn(duration)
	for time := 0; time < duration; time++ {
		r := 0
		if releaseTime >= time {
			r = releaseTime
		}
		generators.GetValue(g, uint64(time), uint64(r))
	}
	return float64(time.Since(start).Microseconds()) / float64(duration)
}
