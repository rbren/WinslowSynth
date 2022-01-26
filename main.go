package main

import (
	"fmt"

	"github.com/rbren/midi/pkg/midi"
	"github.com/rbren/midi/pkg/output"
)

func main() {
	notes := make(chan midi.MidiNote, 1000)
	done := make(chan bool)
	go func() {
		err := midi.StartDriver(notes, done)
		must(err)
	}()

	out, err := output.NewOutputLine(48000)
	must(err)
	out.Player.Play()

	for {
		select {
		case note := <-notes:
			fmt.Println("rec note", note)
		}
	}
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
