package logger

import (
	"livingmaples/packages/logger/collectors"
	"livingmaples/packages/logger/formatters"
	"testing"
	"time"
)

func TestJsonFormatter(t *testing.T) {
	New(&Logger{
		Level: TraceLevel,
		Defaults: map[string]interface{}{
			"foo": "bar",
			"moo": "goo",
		},
		Collector:  []Collector{&collectors.StderrCollector{}},
		Formatter:  formatters.JsonFormatter{},
		SilentMode: false,
		Timeout:    time.Second * 3,
	})

	//l := func(i int) {
	//	Debug("Test logging in json formatter", i)
	//}
	Debug("Test logging in json formatter")
	//for i := 0; i < 200; i++ {
	//	go l(i)
	//}

	//time.Sleep(time.Second * 5)
}

//func TestTextFormatter(t *testing.T) {
//	New(&Logger{
//		Level: TraceLevel,
//		Defaults: map[string]interface{}{
//			"foo": "bar",
//			"moo": "goo",
//		},
//		Collector: []Collector{&collectors.StderrCollector{}},
//		Formatter: formatters.TextFormatter{},
//	})
//
//	Debug("Test logging in text formatter")
//}
