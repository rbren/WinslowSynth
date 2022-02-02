package input

import (
	"strconv"
	"strings"
)

type InputKey struct {
	Action        string
	Velocity      int64
	Key           int64
	Channel       int64
	AbsoluteValue int64
	Frequency     float32
}

func ParseMidiNote(s string) (InputKey, error) {
	note := InputKey{}
	parts := strings.Split(s, " ")
	for idx, val := range parts {
		if idx == 0 {
			note.Action = val
			continue
		}
		if idx == len(parts)-1 {
			break
		}
		intVal, intErr := strconv.ParseInt(parts[idx+1], 10, 0)
		if val == "key" {
			if intErr != nil {
				return note, intErr
			}
			note.Key = intVal
		}
		if val == "velocity" {
			if intErr != nil {
				return note, intErr
			}
			note.Velocity = intVal
		}
		if val == "absValue" {
			if intErr != nil {
				return note, intErr
			}
			note.AbsoluteValue = intVal
		}
		if val == "channel" {
			if intErr != nil {
				return note, intErr
			}
			note.Channel = intVal
		}
	}
	note.Frequency = MidiNotes[note.Key].Frequency
	return note, nil
}
