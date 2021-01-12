package logger

import (
	"context"
	"os"
)

type FatalEvent struct{}

func (e FatalEvent) Fire(ctx context.Context, msg string, defaults map[string]interface{}) {
	os.Exit(1)
}

type PanicEvent struct{}

func (e PanicEvent) Fire(ctx context.Context, msg string, defaults map[string]interface{}) {
	panic(msg)
}
