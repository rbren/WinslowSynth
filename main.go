package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/rbren/midi/pkg/input"
	"github.com/rbren/midi/pkg/logger"
	"github.com/rbren/midi/pkg/music"
	"github.com/rbren/midi/pkg/output"
	"github.com/rbren/midi/pkg/server"
)

const useServer = true

func main() {
	var inputDevice input.InputDevice
	var notes chan input.InputKey
	var err error
	if useServer {
		s := server.Server{}
		notes, err = s.StartListening()
		inputDevice = s
	} else {
		inputDevice, notes, err = input.StartBestInputDevice()
	}
	logger.Log("started input", notes)
	defer inputDevice.Close()
	must(err)

	out, err := output.NewPortAudioOutput()
	must(err)
	defer out.Close()
	logger.Log("created output line")

	musicPlayer := music.NewMusicPlayer(out.Buffer)
	go func() {
		defer func() {
			if e := recover(); e != nil {
				logger.Recover("main", e)
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
