package input

import (
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

const timeToKeyUp = 250 * time.Millisecond

type pressEvent struct {
	lastTime time.Time
	inputKey InputKey
}

type QwertyKeyboard struct{}

func (k QwertyKeyboard) StartListening() (chan InputKey, error) {
	err := termbox.Init()
	if err != nil {
		return nil, err
	}
	onKeys := map[termbox.Key]*pressEvent{}
	notes := getOutputChannel()

	ticker := time.NewTicker(10 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				tickTime := time.Now()
				toDelete := []termbox.Key{}
				for key, evt := range onKeys {
					diff := tickTime.Sub(evt.lastTime)
					if diff > timeToKeyUp {
						evt.inputKey.Action = "channel.NoteOff"
						notes <- evt.inputKey
						toDelete = append(toDelete, key)
					}
				}
				for _, del := range toDelete {
					delete(onKeys, del)
				}
			}
		}
	}()

	go func() {
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				if ev.Key == termbox.KeyCtrlC {
					k.Close()
					os.Exit(0)
				}
				if press, ok := onKeys[ev.Key]; ok {
					press.lastTime = time.Now()
				} else {
					press := pressEvent{
						lastTime: time.Now(),
						inputKey: InputKey{Action: "channel.NoteOn"},
					}
					onKeys[ev.Key] = &press
					notes <- press.inputKey
				}
				termbox.Flush()
			case termbox.EventInterrupt:
				k.Close()
				os.Exit(0)
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
