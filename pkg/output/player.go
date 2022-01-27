package output

import (
	"errors"

	oto "github.com/hajimehoshi/oto/v2"

	"github.com/rbren/midi/pkg/logger"
)

type OutputLine struct {
	sampleRate int
	Line       AudioReaderWriter
	Player     oto.Player
}

func NewOutputLine(sampleRate int) (*OutputLine, error) {
	line := NewAudioReaderWriter(sampleRate * 10)
	logger.Log("create output", sampleRate, len(line.buffer))
	ctx, _, err := oto.NewContext(sampleRate, 1, 1)
	if err != nil {
		return nil, err
	}
	return &OutputLine{
		sampleRate: sampleRate,
		Line:       line,
		Player:     ctx.NewPlayer(line),
	}, nil
}

type AudioReaderWriter struct {
	buffer   []byte
	ReadPos  *int
	WritePos *int
	player   oto.Player
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
		*m.ReadPos = 0
	}
}

func (m *AudioReaderWriter) incrementWritePos() {
	*m.WritePos++
	if *m.WritePos >= len(m.buffer) {
		*m.WritePos = 0
	}
}

func (m AudioReaderWriter) Read(p []byte) (n int, err error) {
	numRead := 0
	for idx := range p {
		if *m.ReadPos == *m.WritePos {
			break
		}
		p[idx] = m.buffer[*m.ReadPos]
		m.incrementReadPos()
		numRead++
	}
	return numRead, nil
}

func (m AudioReaderWriter) Write(p []byte) (n int, err error) {
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
	return numWritten, nil
}
