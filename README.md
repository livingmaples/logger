## Install

```console
GOPRIVATE=github.com/livingmaples/logger go get github.com/livingmaples/logger
```


## What is packages/logger?

Logger is a log collector solution designed to work within applications.

* setting default parameters to attach in every log entity
* format output as TEXT and JSON
* support Stderr and Graylog2 as a collector
* support custom events to fire after specific log level triggered 

### Initializing in main function

```go
import (
	"livingmaples/logger"
	"livingmaples/logger/collectors"
	"livingmaples/logger/formatters"
	"time"
)

logger.New(&logger.Logger{
		Level: logger.TraceLevel,
		Defaults: map[string]interface{}{
			"foo": "bar",
			"moo": "goo",
		},
		Collector:  []Collector{&collectors.StderrCollector{}},
		Formatter:  formatters.JsonFormatter{},
		SilentMode: false, // Disable writing to console
		Timeout:    time.Second * 3,
	})
```

## Supported Log Levels

```go 
logger.Trace("Something very low level.")
logger.Debug("Useful debugging information.")
logger.Info("Something noteworthy happened!")
logger.Warn("You should probably take a look at this.")
logger.Error("Something failed but I'm not quitting.")
// By default FatalEvent event attached to Fatal and calls os.Exit(1) after logging
logger.Fatal("Bye.")
// By default PanicEvent event attached to Panic and calls panic() after logging
logger.Panic("I'm bailing.")
````
* Examples from logrus

## TODO
* Write more tests
* Implement io.Writer