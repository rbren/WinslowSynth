package server

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"

	"github.com/rbren/midi/pkg/generators"
	"github.com/rbren/midi/pkg/input"
	"github.com/rbren/midi/pkg/logger"
	"github.com/rbren/midi/pkg/music"
)

var sendInterval = 50 * time.Millisecond

var upgrader = websocket.Upgrader{} // use default options

type MessageIn struct {
	Key    string
	Action string
}

type MessageOut struct {
	Time       uint64
	Instrument generators.Instrument
}

type Server struct {
	Name       string
	notes      chan input.InputKey
	connection *websocket.Conn
	generator  generators.Generator
	Player     *music.MusicPlayer
}

func (s *Server) Initialize() {
	s.notes = make(chan input.InputKey, 20)
	go s.startReadLoop()
	go s.startWriteLoop()
	http.HandleFunc("/connect", s.connect)
	http.Handle("/", http.FileServer(http.Dir("web")))
	go http.ListenAndServe(":8080", nil)
	logger.Log("started listening")
}

func (s Server) StartListening() (chan input.InputKey, error) {
	return s.notes, nil
}

func (s Server) Close() error {
	s.connection.Close()
	return nil
}

func (s *Server) connect(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.ForceLog("upgrade:", err)
		return
	}
	s.connection = c
	logger.ForceLog("connected")
}

func (s *Server) startReadLoop() {
	for {
		if s.connection == nil {
			continue
		}
		msg := MessageIn{}
		err := s.connection.ReadJSON(&msg)
		if err != nil {
			logger.ForceLog("server read error:", err)
			s.connection = nil
			s.Player.Clear()
			continue
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
		s.notes <- inputKey
	}
}

func (s *Server) startWriteLoop() {
	ticker := time.NewTicker(sendInterval)
	for {
		select {
		case <-ticker.C:
			if s.connection == nil {
				continue
			}
			msg := MessageOut{
				Time:       s.Player.CurrentSample,
				Instrument: s.Player.Instrument,
			}
			err := s.connection.WriteJSON(msg)
			if err != nil {
				logger.ForceLog("server write error:", err)
				break
			}
		}
	}
}
