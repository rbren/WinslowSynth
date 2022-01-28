package input

import (
	"errors"
	"fmt"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/rtmididrv"
)

type MidiKeyboard struct {
	closeFunc func() error
}

func (m MidiKeyboard) StartListening() (chan InputKey, error) {
	drv, err := rtmididrv.New()
	if err != nil {
		return nil, err
	}

	ins, err := drv.Ins()
	if err != nil {
		return nil, err
	}

	outs, err := drv.Outs()
	if err != nil {
		return nil, err
	}

	if len(ins) < 1 || len(outs) < 1 {
		return nil, errors.New("No midi device")
	}
	in, out := ins[0], outs[0]

	err = in.Open()
	if err != nil {
		return nil, err
	}

	err = out.Open()
	if err != nil {
		return nil, err
	}

	m.closeFunc = func() error {
		drv.Close()
		in.Close()
		out.Close()
		return nil
	}

	notes := getOutputChannel()

	// to disable logging, pass mid.NoLogger() as option
	rd := reader.New(
		// reader.NoLogger(),
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			note, err := ParseMidiNote(msg.String())
			if err != nil {
				fmt.Println("error parsing note", err)
				panic(err)
			}
			notes <- note
		}),
	)

	err = rd.ListenTo(in)
	if err != nil {
		return nil, err
	}
	return notes, nil
}

func (m MidiKeyboard) Close() error {
	return m.closeFunc()
}
