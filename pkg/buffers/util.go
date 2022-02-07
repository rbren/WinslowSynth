package buffers

import (
	"math"
)

func MixBuffers(bufs [][]float32) []float32 {
	if len(bufs) == 0 {
		panic("tried to mix empty list of buffers")
	}
	output := make([]float32, len(bufs[0]))
	for i := range output {
		for _, buf := range bufs {
			output[i] += buf[i]
		}
	}
	return output
}

/*
x1 .5
x2 .5
x3 .25

1: x1
2: x1 + (1 - x1)x2
=  x1 + x2 - x1x2

3: x1 + x2 - x1x2 + x3(1 - (x1 + x2 - x1x2))
=  x1 + x2 - x1x2 + x3 - x3(x1 + x2 - x1x2)
=  x1 + x2 + x3 - x1x2 - x3x1 - x3x2 + x1x2x3
*/

func MixBuffersNoOverflow(bufs [][]float32) []float32 {
	if len(bufs) == 0 {
		panic("tried to mix empty list of buffers")
	}
	output := make([]float32, len(bufs[0]))
	for i := range output {
		var total float32
		for _, buf := range bufs {
			total += float32(1.0-math.Abs(float64(total))) * buf[i]
		}
		output[i] = total
	}
	return output
}
