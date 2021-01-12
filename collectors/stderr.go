package collectors

import (
	"context"
	"log"
	"os"
)

type StderrCollector struct {
}

func (c StderrCollector) Ping(ctx context.Context) bool {
	return true
}

func (c StderrCollector) Write(ctx context.Context, data []byte) {
	log.SetOutput(os.Stderr)
	log.Print(string(data))
}
