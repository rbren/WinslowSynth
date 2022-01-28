package output

import (
	"errors"
	"fmt"
	"sync"
	"time"

	oto "github.com/hajimehoshi/oto/v2"

	"github.com/rbren/midi/pkg/logger"
)

type OutputLine struct {
	sampleRate int
	Line       AudioReaderWriter
	Player     oto.Player
}

func NewOutputLine(sampleRate int) (*OutputLine, error) {
	line := NewAudioReaderWriter(sampleRate * 4 * 10)
	logger.Log("create output", sampleRate, len(line.buffer))
	ctx, _, err := oto.NewContext(sampleRate, 2, 2)
	if err != nil {
		return nil, err
	}
	player := ctx.NewPlayer(line)

	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				logger.Log("Unplayed:", player.UnplayedBufferSize())
				if err := player.Err(); err != nil {
					logger.Log("player had an error:", err)
				}
			}
		}
	}()

	return &OutputLine{
		sampleRate: sampleRate,
		Line:       line,
		Player:     player,
	}, nil
}

type AudioReaderWriter struct {
	buffer   []byte
	ReadPos  *int
	WritePos *int
	player   oto.Player
	lock     sync.Mutex
}

func NewAudioReaderWriter(capacity int) AudioReaderWriter {
	ReadPos := 0
	WritePos := 0
	return AudioReaderWriter{
		buffer:   make([]byte, capacity),
		ReadPos:  &ReadPos,
		WritePos: &WritePos,
	}
}

func (m AudioReaderWriter) GetBufferDelay() int {
	if *m.ReadPos <= *m.WritePos {
		return *m.WritePos - *m.ReadPos
	}
	return len(m.buffer) - (*m.ReadPos - *m.WritePos)
}

func (m *AudioReaderWriter) incrementReadPos() {
	*m.ReadPos++
	if *m.ReadPos >= len(m.buffer) {
		fmt.Println("read full buffer")
		*m.ReadPos = 0
	}
}

func (m *AudioReaderWriter) incrementWritePos() {
	*m.WritePos++
	if *m.WritePos >= len(m.buffer) {
		fmt.Println("wrote full buffer")
		*m.WritePos = 0
	}
}

func (m AudioReaderWriter) Read(p []byte) (n int, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	// oto tries to read up to .5 seconds at a time
	// i.e. at 48kHz, 1 chan, 1 deep, 24000 samples
	// i.e. at 48kHz, 2 chan, 2 deep, 96000 samples
	numRead := 0
	numNonZero := 0
	for idx := range p {
		if *m.ReadPos == *m.WritePos {
			break
		}
		p[idx] = m.buffer[*m.ReadPos]
		if p[idx] != 0 {
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

func (m AudioReaderWriter) Write(p []byte) (n int, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	numWritten := 0
	for _, b := range p {
		curReadPos := *m.ReadPos
		m.buffer[*m.WritePos] = b
		m.incrementWritePos()
		numWritten++
		if curReadPos == *m.WritePos {
			return numWritten, errors.New("Caught up to the reader!")
		}
	}
	//logger.Log("write", len(p), numWritten)
	return numWritten, nil
}

func (m AudioReaderWriter) WriteAudio(left []float64, right []float64) (n int, err error) {
	if len(left) != len(right) {
		panic("Two different sized channels!")
	}

	buf := make([]byte, 2*2*len(left))
	channels := [][]float64{left, right}
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
}
