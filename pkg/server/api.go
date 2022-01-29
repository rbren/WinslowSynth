package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/rbren/midi/pkg/input"
	"github.com/rbren/midi/pkg/logger"
)

var upgrader = websocket.Upgrader{} // use default options

type Message struct {
	Key    string
	Action string
}

type Server struct {
	notes chan input.InputKey
}

func (s Server) StartListening() (chan input.InputKey, error) {
	s.notes = make(chan input.InputKey, 20)
	http.HandleFunc("/echo", s.echo)
	http.Handle("/", http.FileServer(http.Dir("web")))
	go http.ListenAndServe(":8080", nil)
	fmt.Println("started listening")
	return s.notes, nil
}

func (s Server) Close() error {
	return nil
}

func (s Server) echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.ForceLog("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		msg := Message{}
		err := c.ReadJSON(&msg)
		fmt.Println("msg", msg)
		if err != nil {
			logger.ForceLog("server read error:", err)
			break
		}

		midi, ok := input.QwertyToMidi[msg.Key]
		if !ok {
			continue
		}
		note := input.MidiNotes[midi]
		action := "channel.NoteOn"
		if msg.Action == "up" {
			action = "channel.NoteOff"
		}
		inputKey := input.InputKey{
			Action:    action,
			Key:       midi,
			Frequency: note.Frequency,
		}
		fmt.Println("send", inputKey)
		s.notes <- inputKey
	}
}
