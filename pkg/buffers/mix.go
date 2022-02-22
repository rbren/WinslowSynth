package buffers

import (
	"fmt"
	"math"
)

type MixStrategy int

const (
	NaiveSumMix MixStrategy = iota
	TanhMix
	BackoffMix
)

func MixBuffers(bufs [][]float32, strategy MixStrategy, scale float32) []float32 {
	if strategy == NaiveSumMix {
		return MixBuffersNaiveSum(bufs, scale)
	}
	if strategy == TanhMix {
		return MixBuffersTanh(bufs, scale)
	}
	if strategy == BackoffMix {
		return MixBuffersBackoff(bufs, scale)
	}
	panic("Unknown strategy")
	return nil
}

func MixBuffersNaiveSum(bufs [][]float32, scale float32) []float32 {
	if len(bufs) == 0 {
		panic("tried to mix empty list of buffers")
	}
	for idx, buf := range bufs {
		if len(buf) != len(bufs[0]) {
			panic(fmt.Errorf("tried to mix buffers of different size. Got %d, expected %d at idx %d of %d", len(buf), len(bufs[0]), idx, len(bufs)))
		}
	}
	output := make([]float32, len(bufs[0]))
	for i := range output {
		for _, buf := range bufs {
			output[i] += scale * buf[i]
		}
	}
	return output
}

func MixBuffersTanh(bufs [][]float32, scale float32) []float32 {
	if len(bufs) == 0 {
		panic("tried to mix empty list of buffers")
	}
	output := make([]float32, len(bufs[0]))
	for i := range output {
		for _, buf := range bufs {
			output[i] += scale * buf[i]
		}
		output[i] = float32(math.Tanh(float64(output[i])))
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

func MixBuffersBackoff(bufs [][]float32, scale float32) []float32 {
	if len(bufs) == 0 {
		panic("tried to mix empty list of buffers")
	}
	output := make([]float32, len(bufs[0]))
	for i := range output {
		var total float32
		for _, buf := range bufs {
			total += float32(1.0-math.Abs(float64(total))) * scale * buf[i]
		}
		output[i] = total
	}
	return output
}
