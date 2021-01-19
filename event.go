package logger

import (
	"os"
)

type FatalEvent struct{}

func (e FatalEvent) Fire(msg string, defaults map[string]interface{}) {
	os.Exit(1)
}

type PanicEvent struct{}

func (e PanicEvent) Fire(msg string, defaults map[string]interface{}) {
	panic(msg)
}
