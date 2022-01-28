package main

import (
	"fmt"

	"github.com/rbren/midi/pkg/input"
	"github.com/rbren/midi/pkg/music"
	"github.com/rbren/midi/pkg/output"
)

const SampleRate = 48000

func main() {
	inputDevice, notes, err := input.StartBestInputDevice()
	fmt.Println("started input")
	defer inputDevice.Close()
	must(err)

	out := output.PortAudioOutput{}
	err = out.Start(SampleRate)
	must(err)
	defer out.Close()
	fmt.Println("created output line")

	musicPlayer := music.NewMusicPlayer(SampleRate, out.Buffer)
	go musicPlayer.Start(notes)
	fmt.Println("started music player")

	fmt.Println("Ready!")
	done := make(chan bool)
	<-done
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
