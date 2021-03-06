package output

import (
	"errors"

	"github.com/sirupsen/logrus"
)

type CircularAudioBuffer struct {
	left     []float32
	right    []float32
	ReadPos  *int
	WritePos *int
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
		logrus.Debug("read full buffer")
		*m.ReadPos = 0
	}
}

func (m *CircularAudioBuffer) incrementWritePos() {
	*m.WritePos++
	if *m.WritePos >= len(m.left) {
		logrus.Debug("wrote full buffer")
		*m.WritePos = 0
	}
}

func (m CircularAudioBuffer) Read(p []byte) (n int, err error) {
	panic("unimplemented")
	return 0, nil
}

func (m CircularAudioBuffer) ReadChannels(p [][]float32) (n int, err error) {
	numRead := 0
	numNonZero := 0
	for idx := range p[0] {
		if *m.ReadPos == *m.WritePos {
			logrus.Info("CAUGHT UP TO WRITER")
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
	for i := numRead; i < len(p[0]); i++ {
		p[0][i] = 0.0
		p[1][i] = 0.0
	}
	return numRead, nil
}

func (m CircularAudioBuffer) Write(left []float32, right []float32) (n int, err error) {
	if len(left) != len(right) {
		panic("Two different sized channels!")
	}
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
	return numWritten, nil
}

func (m CircularAudioBuffer) WriteAudio(left []float32, right []float32) (n int, err error) {
	return m.Write(left, right)
}

func (m CircularAudioBuffer) ReadAudio(out [][]float32) {
	m.ReadChannels(out)
}
