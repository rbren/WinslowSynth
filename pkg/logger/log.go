package logger

import (
	"fmt"
	"os"
)

func Log(arg ...interface{}) {
	if os.Getenv("LOG") == "true" {
		fmt.Println(arg...)
	}
}
