package main

import (
	"flag"
	"fmt"

	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/generators"
	"github.com/rbren/midi/pkg/performance"
)

const numTrials = 3
const valuesPerTrial = 100000
const releasesPerTrial = 10

var instrumentToTest string

func init() {
	flag.StringVar(&instrumentToTest, "instrument", "winslow", "instrument to test")
}

func main() {
	samplesPerSec := config.MainConfig.SampleRate
	samplesPerMicroSec := float64(samplesPerSec) / 1000 / 1000
	fmt.Printf("Min is %.02f us per sample\n", 1/samplesPerMicroSec)
	fmt.Printf("Target is %.02f us per sample\n", 1/samplesPerMicroSec/100)
	for key, inst := range generators.Library {
		if instrumentToTest != "" && key != instrumentToTest {
			continue
		}
		fmt.Println("Testing " + key)
		var avg float64
		for trial := 0; trial < numTrials; trial++ {
			duration := performance.TestPerformance(inst, valuesPerTrial)
			fmt.Printf("  trial %d: %.2fµs\n", trial, duration)
			avg += duration
		}
		fmt.Printf("  average is %.2fµs\n", avg/numTrials)
	}
}
