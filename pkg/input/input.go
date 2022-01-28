package input

import (
	"fmt"
)

const inputBufferSize = 20

func getOutputChannel() chan InputKey {
	return make(chan InputKey, inputBufferSize)
}

func StartBestInputDevice() (InputDevice, chan InputKey, error) {
	var keyListener InputDevice = &MidiKeyboard{}
	notes, err := keyListener.StartListening()
	if err != nil {
		fmt.Println("Couldn't find MIDI keyboard, connecting to QWERTY")
		keyListener = &QwertyKeyboard{}
		notes, err = keyListener.StartListening()
		if err != nil {
			return keyListener, nil, err
		}
	}
	return keyListener, notes, nil
}

type InputDevice interface {
	StartListening() (chan InputKey, error)
	Close() error
}
