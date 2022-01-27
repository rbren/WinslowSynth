package main

import (
	"fmt"

	"github.com/rbren/midi/pkg/midi"
	"github.com/rbren/midi/pkg/music"
	"github.com/rbren/midi/pkg/output"
)

const SampleRate = 4800
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
	musicPlayer := music.NewMusicPlayer(SampleRate, out.Line)

	go musicPlayer.Start()
	out.Player.Play()
	fmt.Println("Ready!")

	for {
		select {
		case note := <-notes:
			fmt.Println("note", note)
			if note.Action == "channel.NoteOn" {
				musicPlayer.ActiveKeys[note.Key] = music.Note{
					Frequency: 440.0,
					Velocity:  note.Velocity,
				}
			} else if note.Action == "channel.NoteOff" {
				delete(musicPlayer.ActiveKeys, note.Key)
			} else {
				fmt.Println("No action for " + note.Action)
			}
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
