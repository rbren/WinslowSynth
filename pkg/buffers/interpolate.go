package buffers

func InterpolateEvents(eventSamples []float32, startIdx, endIdx int) {
	startVal := eventSamples[startIdx]
	endVal := eventSamples[endIdx]
	valDiff := endVal - startVal
	numSamplesToInterpolate := endIdx - startIdx
	for i := startIdx + 1; i < endIdx; i++ {
		// at i=startIdx, pct = 0
		// at i=endIdx, pct = 1
		pct := float32(i-startIdx) / float32(numSamplesToInterpolate)
		eventSamples[i] = startVal + pct*valDiff
	}
}
