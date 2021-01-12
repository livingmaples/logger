package logger

import (
	"context"
	"livingmaples/packages/logger/collectors"
	"livingmaples/packages/logger/formatters"
	"testing"
)

func TestJsonFormatter(t *testing.T) {
	New(context.TODO(), &Logger{
		Level: TraceLevel,
		Defaults: map[string]interface{}{
			"foo": "bar",
			"moo": "goo",
		},
		Collector: []Collector{collectors.StderrCollector{}},
		Formatter: formatters.JsonFormatter{},
	})

	Debug("Test logging in json formatter")
}

func TestTextFormatter(t *testing.T) {
	New(context.TODO(), &Logger{
		Level: TraceLevel,
		Defaults: map[string]interface{}{
			"foo": "bar",
			"moo": "goo",
		},
		Collector: []Collector{collectors.StderrCollector{}},
		Formatter: formatters.TextFormatter{},
	})

	Debug("Test logging in text formatter")
}
