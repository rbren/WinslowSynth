package logger

import (
	"fmt"
	"os"
)

func Log(arg ...interface{}) {
	if os.Getenv("LOG") == "true" {
		ForceLog(arg...)
	}
}

func ForceLog(arg ...interface{}) {
	fmt.Println(arg...)
}
