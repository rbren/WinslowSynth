package buffers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpolation(t *testing.T) {
	events := []float32{1.0, 0.0, 2.0}
	InterpolateEvents(events, 0, 2)
	assert.Equal(t, []float32{1.0, 1.5, 2.0}, events)

	events = []float32{-1.0, 0.0, 1.0}
	InterpolateEvents(events, 0, 2)
	assert.Equal(t, float32(0.0), events[1])

	events = []float32{-1.0, 0.0, 1.0}
	InterpolateEvents(events, 2, 2)
	assert.Equal(t, []float32{-1.0, 0.0, 1.0}, events)

	events = []float32{-1.0, 0.0, 0.0, 0.0, 1.0}
	InterpolateEvents(events, 0, 4)
	assert.Equal(t, float32(-.5), events[1])
	assert.Equal(t, float32(0.0), events[2])
	assert.Equal(t, float32(.5), events[3])

	events = []float32{-1.0, 1.0}
	InterpolateEvents(events, 0, 1)
	assert.Equal(t, float32(-1.0), events[0])
	assert.Equal(t, float32(1.0), events[1])
}

func TestInterpolationCircular(t *testing.T) {
	events := []float32{1.0, 0.0, 2.0}
	InterpolateEvents(events, 2, 0)
	assert.Equal(t, []float32{1.0, 0.0, 2.0}, events)

	events = []float32{0.2, 0.0, 0.0, 1.0, 0.0}
	InterpolateEvents(events, 3, 0)
	assert.Equal(t, []float32{0.2, 0.0, 0.0, 1.0, .6}, events)
}
