package server

import (
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"

	"github.com/rbren/midi/pkg/config"
	"github.com/rbren/midi/pkg/generators"
	"github.com/rbren/midi/pkg/input"
	"github.com/rbren/midi/pkg/music"
)

var sendInterval time.Duration
var samplesPerSend int

func init() {
	msPerSend := config.MainConfig.ServerSendIntervalMs
	samplesPerMs := config.MainConfig.SampleRate / 1000
	samplesPerSend = samplesPerMs * msPerSend
	sendInterval = time.Duration(msPerSend) * time.Millisecond
}

var upgrader = websocket.Upgrader{} // use default options

type MessageIn struct {
	Key    string
	Action string
	Value  float32
}

type MessageOut struct {
	Time               uint64
	Frequency          float32
	Instrument         generators.Generator
	HistoryUpdate      []float32
	Frequencies        []float32
	Instruments        []string
	Constants          []generators.Constant
	Config             config.Config
	SampleRateHandicap float32
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
	go func() {
		logrus.Fatal(http.ListenAndServe(":8080", nil))
	}()
	logrus.Info("started listening")
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
		logrus.Errorf("error connecting: %v", err)
		return
	}
	s.connection = c
	logrus.Info("connected")
}

func (s *Server) startReadLoop() {
	for {
		if s.connection == nil {
			continue
		}
		msg := MessageIn{}
		err := s.connection.ReadJSON(&msg)
		if err != nil {
			logrus.Errorf("server read error: %v", err)
			s.connection = nil
			s.Player.Clear()
			continue
		}
		if msg.Action == "up" || msg.Action == "down" {
			s.NoteAction(msg)
		} else if msg.Action == "set" {
			s.SetAction(msg)
		} else if msg.Action == "choose" {
			s.ChooseAction(msg)
		} else {
			logrus.Error("unknown action", msg.Action)
		}
	}
}

func (s Server) ChooseAction(msg MessageIn) {
	if inst, ok := generators.Library[msg.Key]; ok {
		if inst == nil {
			panic("instrument not found:" + msg.Key)
		}
		s.Player.Sequence.Instrument = inst.Copy(generators.UseDefaultHistoryLength, true)
		logrus.Info("set instrument", s.Player.Sequence.Instrument.GetInfo().Name)
	} else {
		logrus.Error("instrument not found:", msg.Key)
	}
}

func (s Server) SetAction(msg MessageIn) {
	logrus.Infof("Set %s to %f", msg.Key, msg.Value)
	parts := strings.Split(msg.Key, "/")
	group := parts[0]
	name := parts[1]
	s.Player.Sequence.Instrument = generators.SetConstant(s.Player.Sequence.Instrument, group, name, msg.Value)
}

func (s Server) NoteAction(msg MessageIn) {
	logrus.Infof("Note %s %s", msg.Key, msg.Action)
	midi, ok := input.QwertyToMidi[msg.Key]
	if !ok {
		return
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

func (s *Server) startWriteLoop() {
	ticker := time.NewTicker(sendInterval)
	for {
		select {
		case <-ticker.C:
			if s.connection == nil {
				continue
			}
			instruments := []string{}
			for k := range generators.Library {
				instruments = append(instruments, k)
			}
			sort.Strings(instruments)
			msg := MessageOut{
				Time:               s.Player.CurrentSample,
				Instrument:         s.Player.Sequence.Instrument,
				HistoryUpdate:      s.Player.Sequence.Instrument.GetInfo().History.GetOrdered(samplesPerSend),
				Frequencies:        s.Player.Sequence.Instrument.GetInfo().History.GetFrequencies(),
				SampleRateHandicap: s.Player.Sequence.SampleRateHandicap,
				Instruments:        instruments,
				Frequency:          s.Player.Sequence.LastFrequency,
				Config:             config.MainConfig,
			}
			err := s.connection.WriteJSON(msg)
			if err != nil {
				logrus.Error("server write error:", err)
				break
			}
		}
	}
}
