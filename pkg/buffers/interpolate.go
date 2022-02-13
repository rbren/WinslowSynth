package buffers

func InterpolateEvents(eventSamples []float32, startIdx, endIdx int) {
	startVal := eventSamples[startIdx]
	endVal := eventSamples[endIdx]
	valDiff := endVal - startVal
	numSamples := len(eventSamples)
	numSamplesToInterpolate := endIdx - startIdx
	if numSamplesToInterpolate < 0 {
		numSamplesToInterpolate += len(eventSamples)
	}
	if numSamplesToInterpolate == 0 {
		return
	}
	interpolationPos := 1
	for i := (startIdx + 1) % numSamples; i != endIdx; i = (i + 1) % numSamples {
		// at i=startIdx, pct = 0
		// at i=endIdx, pct = 1
		pct := float32(interpolationPos) / float32(numSamplesToInterpolate)
		eventSamples[i] = startVal + pct*valDiff
		interpolationPos++
	}
}
