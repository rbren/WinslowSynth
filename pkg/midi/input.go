package midi

import (
	"fmt"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/rtmididrv"
)

func StartDriver(notes chan MidiNote, done chan bool) error {
	drv, err := rtmididrv.New()
	if err != nil {
		return err
	}
	defer drv.Close()

	ins, err := drv.Ins()
	if err != nil {
		return err
	}

	outs, err := drv.Outs()
	if err != nil {
		return err
	}

	if len(ins) < 1 || len(outs) < 1 {
		panic("No midi device!")
	}
	in, out := ins[0], outs[0]

	err = in.Open()
	if err != nil {
		return err
	}

	err = out.Open()
	if err != nil {
		return err
	}

	defer in.Close()
	defer out.Close()

	// to disable logging, pass mid.NoLogger() as option
	rd := reader.New(
		reader.NoLogger(),
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
		return err
	}
	<-done
	return nil
}
