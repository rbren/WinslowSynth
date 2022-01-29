package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/rbren/midi/pkg/input"
	"github.com/rbren/midi/pkg/logger"
	"github.com/rbren/midi/pkg/music"
	"github.com/rbren/midi/pkg/output"
)

const SampleRate = 48000

func main() {
	inputDevice, notes, err := input.StartBestInputDevice()
	logger.Log("started input")
	defer inputDevice.Close()
	must(err)

	out, err := output.NewPortAudioOutput(SampleRate)
	must(err)
	defer out.Close()
	logger.Log("created output line")

	musicPlayer := music.NewMusicPlayer(SampleRate, out.Buffer)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("PANIC!!! %v \n", r)
			}
		}()
		musicPlayer.Start(notes)
		logger.Log("started music player")
	}()
	out.Start()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	go func() {
		select {
		case sig := <-c:
			fmt.Printf("Got %s signal. Aborting...\n", sig)
			inputDevice.Close()
			out.Close()
			os.Exit(1)
		}
	}()

	fmt.Println("Ready!")
	done := make(chan bool)
	<-done
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
