package main

import (
	"fmt"
	"io"
	"time"

	oto "github.com/hajimehoshi/oto/v2"

	"github.com/rbren/midi/pkg/midi"
)

type MusicPlayer struct {
	buffer   []byte
	readPos  *int
	writePos *int
}

func NewMusicPlayer(capacity int) MusicPlayer {
	zero := 0
	return MusicPlayer{
		buffer:  make([]byte, capacity),
		readPos: &zero,
	}
}

func (m *MusicPlayer) incrementReadPos() {
	*m.readPos++
	if *m.readPos >= len(m.buffer) {
		*m.readPos = 0
	}
}

func (m *MusicPlayer) incrementWritePos() {
	*m.writePos++
	if *m.writePos >= len(m.buffer) {
		*m.writePos = 0
	}
}

func (m MusicPlayer) Read(p []byte) (n int, err error) {
	numRead := 0
	for idx := range p {
		p[idx] = m.buffer[*m.readPos]
		m.incrementReadPos()
		numRead++
	}
	return numRead, nil
}

func (m MusicPlayer) Write(p []byte) (n int, err error) {
	numWritten := 0
	for _, b := range p {
		m.buffer[*m.writePos] = b
		m.incrementWritePos()
		numWritten++
	}
	return numWritten, nil
}

func startOutputBuffer() io.Writer {
	sampleRate := 48000
	bufferSize := sampleRate * 10
	ctx, _, err := oto.NewContext(sampleRate, 2, 2)
	must(err)
	musicPlayer := NewMusicPlayer(bufferSize)
	player := ctx.NewPlayer(musicPlayer)
	player.Play()
	return musicPlayer
}

// This example reads from the first input and and writes to the first output port
func main() {
	notes := make(chan midi.MidiNote, 1000)
	done := make(chan bool)
	go func() {
		err := midi.StartDriver(notes, done)
		must(err)
	}()
	startOutputBuffer()

	for {
		select {
		case note := <-notes:
			fmt.Println("rec note", note)
		}
	}
	time.Sleep(10000 * time.Second)
	// Output: got channel.NoteOn channel 0 key 60 velocity 100
	// got channel.NoteOff channel 0 key 60
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
