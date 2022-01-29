package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/rbren/midi/pkg/input"
	"github.com/rbren/midi/pkg/logger"
	"github.com/rbren/midi/pkg/music"
)

var sendInterval = 10 * time.Millisecond

var upgrader = websocket.Upgrader{} // use default options

type MessageIn struct {
	Key    string
	Action string
}

type MessageOut struct {
	Time uint64
}

type Server struct {
	notes      chan input.InputKey
	connection *websocket.Conn
	Player     *music.MusicPlayer
}

func (s Server) StartListening() (chan input.InputKey, error) {
	s.notes = make(chan input.InputKey, 20)
	http.HandleFunc("/connect", s.connect)
	http.Handle("/", http.FileServer(http.Dir("web")))
	go http.ListenAndServe(":8080", nil)
	go s.startReadLoop()
	go s.startWriteLoop()
	fmt.Println("started listening")
	return s.notes, nil
}

func (s Server) Close() error {
	return nil
}

func (s Server) connect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.ForceLog("upgrade:", err)
		return
	}
	defer c.Close()
	s.connection = c
}

func (s Server) startReadLoop() {
	for {
		if s.connection == nil {
			continue
		}
		msg := MessageIn{}
		err := s.connection.ReadJSON(&msg)
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

func (s Server) startWriteLoop() {
	ticker := time.NewTicker(sendInterval)
	for {
		select {
		case <-ticker.C:
			if s.connection == nil {
				continue
			}
			msg := MessageOut{
				Time: s.Player.CurrentSample,
			}
			err := s.connection.WriteJSON(msg)
			if err != nil {
				logger.ForceLog("server write error:", err)
				break
			}
		}
	}
}
