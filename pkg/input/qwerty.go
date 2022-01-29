package input

import (
	"fmt"
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

var QwertyToMidi = map[string]int64{
	"a": 58, // A#
	"s": 59, // B
	"d": 60, // middle C
	"r": 61, // C#
	"f": 62, // D
	"t": 63, // D#
	"g": 64, // E
	"h": 65, // F
	"u": 66, // F#
	"j": 67, // G
	"i": 68, // G#
	"k": 69, // A
	"o": 70, // A#
	"l": 71, // B
	";": 72, // C
}

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
				fmt.Println("Key", string(ev.Ch))
				if ev.Key == termbox.KeyCtrlC {
					k.Close()
					os.Exit(0)
				}
				if press, ok := onKeys[ev.Key]; ok {
					press.lastTime = time.Now()
				} else {
					midi := QwertyToMidi[string(ev.Ch)]
					note := MidiNotes[midi]
					press := pressEvent{
						lastTime: time.Now(),
						inputKey: InputKey{
							Action:    "channel.NoteOn",
							Key:       midi,
							Frequency: note.Frequency,
						},
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
