package output

import (
	"errors"

	oto "github.com/hajimehoshi/oto/v2"
)

type OutputLine struct {
	sampleRate int
	Line       AudioReaderWriter
	Player     oto.Player
}

func NewOutputLine(sampleRate int) (*OutputLine, error) {
	line := NewAudioReaderWriter(sampleRate * 10)
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
	readPos  *int
	writePos *int
	player   oto.Player
}

func NewAudioReaderWriter(capacity int) AudioReaderWriter {
	readPos := 0
	writePos := 0
	return AudioReaderWriter{
		buffer:   make([]byte, capacity),
		readPos:  &readPos,
		writePos: &writePos,
	}
}

func (m *AudioReaderWriter) incrementReadPos() {
	*m.readPos++
	if *m.readPos >= len(m.buffer) {
		*m.readPos = 0
	}
}

func (m *AudioReaderWriter) incrementWritePos() {
	*m.writePos++
	if *m.writePos >= len(m.buffer) {
		*m.writePos = 0
	}
}

func (m AudioReaderWriter) Read(p []byte) (n int, err error) {
	//fmt.Println("try read", len(p))
	numRead := 0
	for idx := range p {
		if *m.readPos == *m.writePos {
			break
		}
		p[idx] = m.buffer[*m.readPos]
		m.incrementReadPos()
		numRead++
	}
	if numRead > 0 {
		//fmt.Println("read", numRead)
	}
	return numRead, nil
}

func (m AudioReaderWriter) Write(p []byte) (n int, err error) {
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
