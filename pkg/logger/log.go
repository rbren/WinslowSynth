package logger

import (
	"os"
	"runtime/debug"

	"github.com/sirupsen/logrus"
)

func init() {
	level := os.Getenv("LOG")
	if level == "debug" {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func Recover(msg string, r interface{}) {
	logrus.Errorf("PANIC at %s %v", msg, r)
	debug.PrintStack()
}
