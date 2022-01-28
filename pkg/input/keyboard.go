package input

import (
	"fmt"
	"os"

	"github.com/eiannone/keyboard"
)

type QwertyKeyboard struct{}

func (k QwertyKeyboard) StartListening() error {
	keysEvents, err := keyboard.GetKeys(10)
	if err != nil {
		return err
	}
	defer func() {
		err := keyboard.Close()
		if err != nil {
			panic(err)
		}
	}()

	for {
		event := <-keysEvents
		if event.Key == keyboard.KeyEsc || event.Key == keyboard.KeyCtrlC {
			break
		}
		fmt.Printf("You pressed: rune %q, key %X\r\n", event.Rune, event.Key)
		if event.Err != nil {
			panic(event.Err)
		}
	}
	err = keyboard.Close()
	if err != nil {
		panic(err)
	}
	os.Exit(0)
	return nil
}
