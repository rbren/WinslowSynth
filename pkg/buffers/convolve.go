package buffers

func Convolve(input, kernel []float32) []float32 {
	if len(input) <= len(kernel) {
		panic("provided buffer is not greater than the filter weights")
	}

	output := make([]float32, len(input))
	for i := 0; i < len(output); i++ {
		output[i] = 0.0
		// note: this assumes kernel == 0 outside the provided range
		for m := 0; m < len(kernel); m++ {
			output[i] = output[i] + input[i-m]*kernel[m]
		}
	}
	return output
}
