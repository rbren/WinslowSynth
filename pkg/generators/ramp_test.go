package generators

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func assertAboutEqual(t *testing.T, expected, actual float32) {
	assert.InDelta(t, expected, actual, .00001)
}

func TestRamp(t *testing.T) {
	r := Ramp{RampUp: 10, RampDown: 5, Target: 100.0}

	// steady state
	assertAboutEqual(t, float32(100.0), r.GetValue(100, 0))
	assertAboutEqual(t, float32(100.0), r.GetValue(1000, 0))

	// ramp up
	assertAboutEqual(t, float32(0.0), r.GetValue(0, 0))
	assertAboutEqual(t, float32(50.0), r.GetValue(5, 0))
	assertAboutEqual(t, float32(90.0), r.GetValue(9, 0))
	assertAboutEqual(t, float32(100.0), r.GetValue(10, 0))

	// ramp down
	assertAboutEqual(t, float32(100.0), r.GetValue(100, 100))
	assertAboutEqual(t, float32(80.0), r.GetValue(101, 100))
	assertAboutEqual(t, float32(20.0), r.GetValue(104, 100))
	assertAboutEqual(t, float32(0.0), r.GetValue(105, 100))
	assertAboutEqual(t, float32(0.0), r.GetValue(110, 100))

	// ramp down without reaching peak
	assertAboutEqual(t, float32(50.0), r.GetValue(5, 0))
	assertAboutEqual(t, float32(40.0), r.GetValue(6, 5))
	assertAboutEqual(t, float32(30.0), r.GetValue(7, 5))
	assertAboutEqual(t, float32(20.0), r.GetValue(8, 5))
	assertAboutEqual(t, float32(10.0), r.GetValue(9, 5))
	assertAboutEqual(t, float32(0.0), r.GetValue(10, 5))
	assertAboutEqual(t, float32(0.0), r.GetValue(11, 5))

}
