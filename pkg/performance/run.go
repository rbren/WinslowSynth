package performance

import (
	"math/rand"
	"time"

	"github.com/rbren/midi/pkg/generators"
)

const numTrials = 3
const valuesPerTrial = 10000
const releasesPerTrial = 50

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestPerformance(g generators.Generator, useHist bool) float64 {
	avgDuration := 0.0
	generators.SetFrequency(g, 440.0)
	for trial := 0; trial < numTrials; trial++ {
		histSize := 0
		if useHist {
			histSize = generators.UseDefaultHistoryLength
		}
		start := time.Now()
		total := 0
		for rel := 0; rel < releasesPerTrial; rel++ {
			doRelease := rand.Float64() < .5
			gTrial := g.Copy(histSize)
			for time := 0; time < valuesPerTrial; time++ {
				total++
				releaseTime := 0
				if time > 0 && doRelease {
					releaseTime = rand.Intn(time)
				}
				generators.GetValue(gTrial, uint64(time), uint64(releaseTime))
			}
		}
		duration := float64(time.Since(start).Microseconds()) / float64(total)
		avgDuration += duration
		//fmt.Printf("  trial %d took %fus per sample\n", trial, float64(duration.Microseconds())/float64(total))
	}
	return avgDuration / numTrials
}
