package logger

import (
	"fmt"
	"os"
	"runtime/debug"
)

func Log(arg ...interface{}) {
	if os.Getenv("LOG") == "true" {
		ForceLog(arg...)
	}
}

func ForceLog(arg ...interface{}) {
	fmt.Println(arg...)
}

func Recover(msg string, r interface{}) {
	fmt.Printf("PANIC at %s %v \n", msg, r)
	debug.PrintStack()
}
