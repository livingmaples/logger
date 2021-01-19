package collectors

import (
	"log"
	"os"
)

type StderrCollector struct {
}

func (c StderrCollector) Ping() (bool, error) {
	return true, nil
}

func (c StderrCollector) Write(data *[]byte) {
	log.SetOutput(os.Stderr)
	log.SetFlags(0)
	log.Print(string(*data))
}
