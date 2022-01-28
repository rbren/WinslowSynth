package input

import (
	"os"

	"github.com/nsf/termbox-go"
)

type QwertyKeyboard struct{}

func (k QwertyKeyboard) StartListening() (chan InputKey, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}
	notes := getOutputChannel()
	go func() {
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				if ev.Key == termbox.KeyCtrlC {
					k.Close()
					os.Exit(0)
				}
				notes <- InputKey{Action: "channel.NoteOn"}
				termbox.Flush()
			case termbox.EventError:
				panic(ev.Err)
			}
		}
	}()
	return notes, nil
}

func (k QwertyKeyboard) Close() error {
	termbox.Close()
	return nil
}
