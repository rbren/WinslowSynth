package main

import (
	"fmt"

	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/generators"
	"github.com/rbren/midi/pkg/performance"
)

func main() {
	samplesPerSec := config.MainConfig.SampleRate
	samplesPerMicroSec := float64(samplesPerSec) / 1000 / 1000
	fmt.Printf("Min is %.02f us per sample\n", 1/samplesPerMicroSec)
	fmt.Printf("Target is %.02f us per sample\n", 1/samplesPerMicroSec/100)
	for key, inst := range generators.Library {
		fmt.Println("Testing " + key)
		var avg float64
		avg = performance.TestPerformance(inst, true)
		fmt.Printf("  with history:    average is %.2fµs\n", avg)
		if key == "noiseFilter" {
			continue // needs history
		}
		avg = performance.TestPerformance(inst, false)
		fmt.Printf("  without history: average is %.2fµs\n", avg)
	}
}
