package output

import (
	"errors"

	oto "github.com/hajimehoshi/oto/v2"
)

type OutputLine struct {
	sampleRate int
	Line       MusicReaderWriter
	Player     oto.Player
}

func NewOutputLine(sampleRate int) (*OutputLine, error) {
	line := NewMusicReaderWriter(sampleRate * 10)
	ctx, _, err := oto.NewContext(sampleRate, 2, 2)
	if err != nil {
		return nil, err
	}
	return &OutputLine{
		sampleRate: sampleRate,
		Line:       line,
		Player:     ctx.NewPlayer(line),
	}, nil
}

type MusicReaderWriter struct {
	buffer   []byte
	readPos  *int
	writePos *int
	player   oto.Player
}

func NewMusicReaderWriter(capacity int) MusicReaderWriter {
	readPos := 0
	writePos := 0
	return MusicReaderWriter{
		buffer:   make([]byte, capacity),
		readPos:  &readPos,
		writePos: &writePos,
	}
}

func (m *MusicReaderWriter) incrementReadPos() {
	*m.readPos++
	if *m.readPos >= len(m.buffer) {
		*m.readPos = 0
	}
}

func (m *MusicReaderWriter) incrementWritePos() {
	*m.writePos++
	if *m.writePos >= len(m.buffer) {
		*m.writePos = 0
	}
}

func (m MusicReaderWriter) Read(p []byte) (n int, err error) {
	numRead := 0
	for idx := range p {
		if *m.readPos == *m.writePos {
			break
		}
		p[idx] = m.buffer[*m.readPos]
		m.incrementReadPos()
		numRead++
	}
	return numRead, nil
}

func (m MusicReaderWriter) Write(p []byte) (n int, err error) {
	numWritten := 0
	for _, b := range p {
		m.buffer[*m.writePos] = b
		m.incrementWritePos()
		numWritten++
		if *m.readPos == *m.writePos {
			return numWritten, errors.New("Caught up to the reader!")
		}
	}
	return numWritten, nil
}
