package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	oto "github.com/hajimehoshi/oto/v2"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/midi/writer"
	"gitlab.com/gomidi/rtmididrv"
)

type MidiNote struct {
	Action        string
	Velocity      int64
	Key           int64
	Channel       int64
	AbsoluteValue int64
}

func ParseMidiNote(s string) MidiNote {
	note := MidiNote{}
	parts := strings.Split(s, " ")
	for idx, val := range parts {
		if idx == 0 {
			note.Action = val
			continue
		}
		if idx == len(parts)-1 {
			break
		}
		intVal, intErr := strconv.ParseInt(parts[idx+1], 10, 0)
		if val == "key" {
			must(intErr)
			note.Key = intVal
		}
		if val == "velocity" {
			must(intErr)
			note.Velocity = intVal
		}
		if val == "absValue" {
			must(intErr)
			note.AbsoluteValue = intVal
		}
		if val == "channel" {
			must(intErr)
			note.Channel = intVal
		}
	}
	if strings.Contains(s, "channel.NoteOn") {
		note.Action = "NoteOn"
	} else if strings.Contains(s, "channel.NoteOff") {
		note.Action = "NoteOff"
	}
	return note
}

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
	drv, err := rtmididrv.New()
	if err != nil {
		panic(err)
	}

	startOutputBuffer()

	// make sure to close all open ports at the end
	defer drv.Close()

	ins, err := drv.Ins()
	must(err)

	outs, err := drv.Outs()
	must(err)

	in, out := ins[0], outs[0]

	must(in.Open())
	must(out.Open())

	defer in.Close()
	defer out.Close()

	// the writer we are writing to
	wr := writer.New(out)

	// to disable logging, pass mid.NoLogger() as option
	rd := reader.New(
		reader.NoLogger(),
		// write every message to the out port
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			note := ParseMidiNote(msg.String())
			fmt.Printf("got %#v\n", note)
		}),
	)

	// listen for MIDI
	err = rd.ListenTo(in)
	must(err)
	fmt.Println("listening...")
	time.Sleep(100000 * time.Second)

	err = writer.NoteOn(wr, 60, 100)
	must(err)

	err = writer.NoteOff(wr, 60)

	must(err)
	// Output: got channel.NoteOn channel 0 key 60 velocity 100
	// got channel.NoteOff channel 0 key 60
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
