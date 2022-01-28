package main

import (
	"fmt"

	"github.com/rbren/midi/pkg/input"
	"github.com/rbren/midi/pkg/music"
	"github.com/rbren/midi/pkg/output"
)

const SampleRate = 48000
const SampleBuffer = 10000

func main() {
	inputDevice, notes, err := input.StartBestInputDevice()
	fmt.Println("started input")
	defer inputDevice.Close()
	must(err)

	out, err := output.NewOutputLine(SampleRate)
	must(err)
	fmt.Println("created output line")

	musicPlayer := music.NewMusicPlayer(SampleRate, out.Line)
	go musicPlayer.Start()
	fmt.Println("started music player")
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
			if err := out.Player.Err(); err != nil {
				fmt.Println("there was an error!", err)
				//out.Player.Play()
			}
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
