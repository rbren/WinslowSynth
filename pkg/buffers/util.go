package buffers

func MixBuffers(bufs [][]float32) []float32 {
	output := make([]float32, len(bufs[0]))
	for _, buf := range bufs {
		for i := range buf {
			output[i] += buf[i]
		}
	}
	return output
}
