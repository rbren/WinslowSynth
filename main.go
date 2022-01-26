package main

import (
	"fmt"

	"github.com/rbren/midi/pkg/midi"
	"github.com/rbren/midi/pkg/music"
	"github.com/rbren/midi/pkg/output"
)

const SampleRate = 48000
const SampleBuffer = 10000
const MidiBuffer = 1000

func main() {
	notes := make(chan midi.MidiNote, MidiBuffer)
	done := make(chan bool)
	go func() {
		err := midi.StartDriver(notes, done)
		must(err)
	}()

	out, err := output.NewOutputLine(SampleRate)
	must(err)
	out.Player.Play()

	startPlaying := func() {
		for {
			buf := music.GenerateFrequency(440.0, SampleRate, SampleBuffer)
			out.Line.Write(buf)
		}
	}

	for {
		select {
		case note := <-notes:
			if note.Action == "channel.NoteOn" {
				startPlaying()
			}
			fmt.Println("rec note", note)
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
