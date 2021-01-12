package logger

import (
	"context"
	"livingmaples/packages/logger/collectors"
	"livingmaples/packages/logger/formatters"
	"testing"
)

func TestWrite(t *testing.T) {
	New(context.TODO(), &Logger{
		Level: FatalLevel,
		Defaults: map[string]interface{}{
			"foo": "bar",
			"moo": "goo",
		},
		Collector: []Collector{collectors.StderrCollector{}},
		Formatter: formatters.JsonFormatter{},
	})

	Fatal("Test logging")
}
