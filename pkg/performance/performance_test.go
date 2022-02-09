package main

import (
	"fmt"
	"testing"
	"time"

	"github.com/rbren/midi/pkg/generators"
)

const numTrials = 10
const valuesPerTrial = 10000

func TestPerformance(t *testing.T) {
	for key, inst := range generators.Library {
		fmt.Println("Testing " + key)
		testPerformance(inst)
	}
}

func testPerformance(g generators.Generator) {
	for trial := 0; trial < numTrials; trial++ {
		start := time.Now()
		for time := 0; time < valuesPerTrial; time++ {
			g.GetValue(uint64(time), 0) // TODO: test release times
		}
		duration := time.Since(start)
		fmt.Printf("  trial %d took %dus\n", trial, duration.Microseconds())
	}
}
