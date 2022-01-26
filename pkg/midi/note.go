package midi

import (
	"strconv"
	"strings"
)

type MidiNote struct {
	Action        string
	Velocity      int64
	Key           int64
	Channel       int64
	AbsoluteValue int64
}

func ParseMidiNote(s string) (MidiNote, error) {
	note := MidiNote{}
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
	return note, nil
}
