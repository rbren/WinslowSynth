package buffers

func MixBuffers(bufs [][]float32) []float32 {
	if len(bufs) == 0 {
		panic("tried to mix empty list of buffers")
	}
	output := make([]float32, len(bufs[0]))
	for i := range output {
		for _, buf := range bufs {
			output[i] += buf[i]
		}
		//output[i] = output[i] / float32(len(bufs))
	}
	return output
}
