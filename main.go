package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"

	"github.com/rbren/midi/pkg/input"
	"github.com/rbren/midi/pkg/logger"
	"github.com/rbren/midi/pkg/music"
	"github.com/rbren/midi/pkg/output"
	"github.com/rbren/midi/pkg/server"
)

const useServer = true

func startServer() {
	s := server.Server{}
	s.Initialize()
	browserNotes, err := s.StartListening()
	closeOnExit(s)
	must(err)

	notes := browserNotes
	var keyListener input.InputDevice = &input.MidiKeyboard{}
	defer keyListener.Close()
	midiNotes, err := keyListener.StartListening()
	if err != nil {
		fmt.Println("Couldn't find MIDI keyboard, using browser")
	} else {
		notes = midiNotes
	}

	player, out := startOutput(notes)
	defer out.Close()
	s.Player = player

	fmt.Println("Ready!")
	done := make(chan bool)
	<-done
}

func startLocal() {
	inputDevice, notes, err := input.StartBestInputDevice()
	defer inputDevice.Close()
	closeOnExit(inputDevice)
	must(err)
	_, out := startOutput(notes)
	defer out.Close()

	fmt.Println("Ready!")
	done := make(chan bool)
	<-done
}

func main() {
	if useServer {
		startServer()
	} else {
		startLocal()
	}
}

func startOutput(notes chan input.InputKey) (*music.MusicPlayer, *output.PortAudioOutput) {
	out, err := output.NewPortAudioOutput()
	must(err)
	logrus.Info("created output line")

	musicPlayer := music.NewMusicPlayer(out.Buffer)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				logger.Recover("main", e)
			}
		}()
		musicPlayer.Start(notes)
		logrus.Info("started music player")
	}()
	out.Start()
	return &musicPlayer, out
}

func closeOnExit(inputDevice input.InputDevice) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)

	go func() {
		select {
		case sig := <-c:
			fmt.Printf("Got %s signal. Aborting...\n", sig)
			inputDevice.Close()
			os.Exit(1)
		}
	}()
}

func must(err error) {
	if err != nil {
		panic(err.Error())
	}
}
