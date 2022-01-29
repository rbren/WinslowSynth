package output

import (
	"errors"
	"fmt"
	"sync"

	"github.com/rbren/midi/pkg/logger"
)

type CircularAudioBuffer struct {
	left     []float32
	right    []float32
	ReadPos  *int
	WritePos *int
	lock     sync.Mutex
}

func NewCircularAudioBuffer(capacity int) *CircularAudioBuffer {
	ReadPos := 0
	WritePos := 0
	return &CircularAudioBuffer{
		right:    make([]float32, capacity),
		left:     make([]float32, capacity),
		ReadPos:  &ReadPos,
		WritePos: &WritePos,
	}
}

func (m CircularAudioBuffer) GetCapacity() int {
	return len(m.left)
}

func (m CircularAudioBuffer) GetBufferDelay() int {
	if *m.ReadPos <= *m.WritePos {
		return *m.WritePos - *m.ReadPos
	}
	return len(m.left) - (*m.ReadPos - *m.WritePos)
}

func (m *CircularAudioBuffer) incrementReadPos() {
	*m.ReadPos++
	if *m.ReadPos >= len(m.left) {
		fmt.Println("read full buffer")
		*m.ReadPos = 0
	}
}

func (m *CircularAudioBuffer) incrementWritePos() {
	*m.WritePos++
	if *m.WritePos >= len(m.left) {
		fmt.Println("wrote full buffer")
		*m.WritePos = 0
	}
}

func (m CircularAudioBuffer) Read(p []byte) (n int, err error) {
	panic("unimplemented")
	return 0, nil
}

func (m CircularAudioBuffer) ReadChannels(p [][]float32) (n int, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	// oto tries to read up to .5 seconds at a time
	// i.e. at 48kHz, 1 chan, 1 deep, 24000 samples
	// i.e. at 48kHz, 2 chan, 2 deep, 96000 samples
	numRead := 0
	numNonZero := 0
	for idx := range p[0] {
		if *m.ReadPos == *m.WritePos {
			break
		}
		p[0][idx] = m.left[*m.ReadPos]
		p[1][idx] = m.right[*m.ReadPos]
		if p[0][idx] != 0.0 || p[1][idx] != 0.0 {
			numNonZero++
		}
		m.incrementReadPos()
		numRead++
	}
	if numNonZero > 0 {
		logger.Log("read", numRead, numNonZero)
	}
	return numRead, nil
}

func (m CircularAudioBuffer) Write(left []float32, right []float32) (n int, err error) {
	if len(left) != len(right) {
		panic("Two different sized channels!")
	}
	m.lock.Lock()
	defer m.lock.Unlock()
	numWritten := 0
	for idx := range left {
		curReadPos := *m.ReadPos
		m.left[*m.WritePos] = left[idx]
		m.right[*m.WritePos] = right[idx]
		m.incrementWritePos()
		numWritten++
		if curReadPos == *m.WritePos {
			return numWritten, errors.New("Caught up to the reader!")
		}
	}
	//logger.Log("write", len(p), numWritten)
	return numWritten, nil
}

func (m CircularAudioBuffer) WriteAudio(left []float32, right []float32) (n int, err error) {
	return m.Write(left, right)

	/*
		buf := make([]byte, 2*2*len(left))
		channels := [][]float32{left, right}
		for c := range channels {
			for i := range channels[c] {
				val := channels[c][i]
				if val < -1 {
					val = -1
				}
				if val > +1 {
					val = +1
				}
				valInt16 := int16(val * (1<<15 - 1))
				low := byte(valInt16)
				high := byte(valInt16 >> 8)
				buf[i*4+c*2+0] = low
				buf[i*4+c*2+1] = high
			}
		}
		return m.Write(buf)
	*/
}

func (m CircularAudioBuffer) ReadAudio(out [][]float32) {
	m.ReadChannels(out)
}