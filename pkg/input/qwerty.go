package input

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/nsf/termbox-go"
)

const timeToKeyUp = 250 * time.Millisecond

type pressEvent struct {
	lastTime time.Time
	inputKey InputKey
	sustain  bool
}

type QwertyKeyboard struct {
	lock   sync.Mutex
	onKeys map[rune]*pressEvent
	notes  chan InputKey
}

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
	k.onKeys = map[rune]*pressEvent{}
	k.notes = getOutputChannel()

	ticker := time.NewTicker(10 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				k.clearNotes()
			}
		}
	}()

	go func() {
		for {
			switch ev := termbox.PollEvent(); ev.Type {
			case termbox.EventKey:
				k.handleKeyPress(ev)
				termbox.Flush()
			case termbox.EventInterrupt:
				k.Close()
				os.Exit(0)
			case termbox.EventError:
				panic(ev.Err)
			}
		}
	}()
	return k.notes, nil
}

func (k QwertyKeyboard) Close() error {
	termbox.Close()
	return nil
}

func (k QwertyKeyboard) handleKeyPress(ev termbox.Event) {
	fmt.Println("Key", string(ev.Ch))
	k.lock.Lock()
	defer k.lock.Unlock()
	if ev.Key == termbox.KeyCtrlC {
		k.Close()
		os.Exit(0)
	}
	char := strings.ToLower(string(ev.Ch))
	sustain := string(ev.Ch) != char
	if !sustain {
		for _, pressEvent := range k.onKeys {
			pressEvent.sustain = false
		}
	}
	if press, ok := k.onKeys[ev.Ch]; ok {
		press.lastTime = time.Now()
	} else {
		midi := QwertyToMidi[char]
		note := MidiNotes[midi]
		press := pressEvent{
			lastTime: time.Now(),
			sustain:  sustain,
			inputKey: InputKey{
				Action:    "channel.NoteOn",
				Key:       midi,
				Frequency: note.Frequency,
			},
		}
		k.onKeys[ev.Ch] = &press
		k.notes <- press.inputKey
	}
}

func (k QwertyKeyboard) clearNotes() {
	tickTime := time.Now()
	toDelete := []rune{}
	k.lock.Lock()
	defer k.lock.Unlock()
	for key, evt := range k.onKeys {
		if evt.sustain {
			continue
		}
		diff := tickTime.Sub(evt.lastTime)
		if diff > timeToKeyUp {
			evt.inputKey.Action = "channel.NoteOff"
			k.notes <- evt.inputKey
			toDelete = append(toDelete, key)
		}
	}
	for _, del := range toDelete {
		delete(k.onKeys, del)
	}
}
