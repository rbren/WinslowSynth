package main

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/generators"
)

const numTrials = 3
const valuesPerTrial = 10000
const releasesPerTrial = 50

func init() {
	rand.Seed(time.Now().UnixNano())
}

func TestPerformance(t *testing.T) {
	for key, inst := range generators.Library {
		fmt.Println("Testing " + key)
		testPerformance(inst)
	}

	samplesPerSec := config.MainConfig.SampleRate
	samplesPerMicroSec := float64(samplesPerSec) / 1000 / 1000
	fmt.Printf("Min is %.02f us per sample\n", 1/samplesPerMicroSec)
	fmt.Printf("Target is %.02f us per sample\n", 1/samplesPerMicroSec/100)
}

func testPerformance(g generators.Generator) {
	generators.SetFrequency(g, 440.0)
	for trial := 0; trial < numTrials; trial++ {
		gTrial := g.Copy(generators.UseDefaultHistoryLength)
		start := time.Now()
		total := 0
		for time := 0; time < valuesPerTrial; time++ {
			for rel := 0; rel < releasesPerTrial; rel++ {
				total++
				releaseTime := 0
				if time > 0 && rand.Float64() < .5 {
					releaseTime = rand.Intn(time)
				}
				generators.GetValue(gTrial, uint64(time), uint64(releaseTime))
			}
		}
		duration := time.Since(start)
		fmt.Printf("  trial %d took %fus per sample\n", trial, float64(duration.Microseconds())/float64(total))
	}
}
