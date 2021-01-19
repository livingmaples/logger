package logger

import (
	"fmt"
	"sync"
	"time"
)

type LogLevel uint32

const (
	// PanicLevel
	PanicLevel LogLevel = iota
	// FatalLevel
	FatalLevel
	// ErrorLevel
	ErrorLevel
	// WarningLevel
	WarningLevel
	// InfoLevel
	InfoLevel
	// DebugLevel
	DebugLevel
	// TraceLevel
	TraceLevel
)

var logInstance *Logger

type Collector interface {
	Ping() (bool, error)
	Write(*[]byte)
}

type OutputFormatter interface {
	Format(string, map[string]string) ([]byte, error)
}

type Event interface {
	Fire(string, map[string]interface{})
}

type Logger struct {
	Level      LogLevel
	Defaults   map[string]interface{}
	Collector  []Collector
	Formatter  OutputFormatter
	Event      map[LogLevel][]Event
	Lock       bool // Not used yet
	Timeout    time.Duration
	SilentMode bool // if true, the verbose mode will disabled
	mu         sync.Mutex
}

func (level LogLevel) Label() string {
	switch level {
	case PanicLevel:
		return "panic"
	case FatalLevel:
		return "fatal"
	case ErrorLevel:
		return "error"
	case WarningLevel:
		return "warning"
	case InfoLevel:
		return "info"
	case DebugLevel:
		return "debug"
	case TraceLevel:
		return "trace"
	default:
		return "unknown"
	}
}

// New initialize logger with given data
func New(logger *Logger) error {
	logInstance = logger
	// Add default fatal and panic events to events list
	if logInstance.Event == nil {
		logInstance.Event = make(map[LogLevel][]Event)
	}

	logInstance.Event[FatalLevel] = append(logInstance.Event[FatalLevel], FatalEvent{})
	logInstance.Event[PanicLevel] = append(logInstance.Event[PanicLevel], PanicEvent{})

	// currently New only returns nil but in the future we may check some conditions
	return nil
}

func (l *Logger) SetLevel(level LogLevel) {
	l.Level = level
}

func (l *Logger) SetDefaults(defaults map[string]interface{}) {
	l.Defaults = defaults
}

func (l *Logger) SetCollector(c Collector) {
	l.Collector = append(l.Collector, c)
}

func (l *Logger) SetFormatter(f OutputFormatter) {
	l.Formatter = f
}

func (l *Logger) SetEvent(level LogLevel, event Event) {
	// If the new event is for FatalLevel or PanicLevel,
	// we must make sure that default events keep at the end of the list
	if level == FatalLevel || level == PanicLevel {
		index := len(l.Event[level]) - 1
		l.Event[level] = append(l.Event[level][:index+1], l.Event[level][index:]...)
		l.Event[level][index] = event
		return
	}

	l.Event[level] = append(l.Event[level], event)

}

func (l *Logger) log(level LogLevel, msg string) error {
	if !l.isLevelEnabled(level) {
		return nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if l.Timeout.Seconds() == 0 {
		return l.handle(msg)
	}

	writeDoneCh := make(chan error, 1)
	go func() {
		writeDoneCh <- l.handle(msg)
	}()

	select {
	case <-time.After(l.Timeout):
		return fmt.Errorf("logging timedout after %f seconds", l.Timeout.Seconds())
	case err := <-writeDoneCh:
		return err
	}
}

// handle don't use Logger ctx directly instead uses context passed from log function
func (l *Logger) handle(msg string) error {
	c, err := l.collector()
	if err != nil {
		fmt.Println(err)
		return err
	}

	r, err := l.formatter().Format(l.Level.Label(), mergeDefaultsToOutput(map[string]string{
		"time":  time.Now().String(),
		"level": l.Level.Label(),
		"msg":   msg,
	}, l.Defaults))
	if err != nil {
		return err
	}

	c.Write(&r)
	l.triggerEvents(msg)

	return nil
}

func (l *Logger) isLevelEnabled(level LogLevel) bool {
	return l.Level >= level
}

func (l *Logger) formatter() OutputFormatter {
	return l.Formatter
}

// collector returns first available log collector
func (l *Logger) collector() (Collector, error) {
	for _, c := range l.Collector {
		ok, err := c.Ping()
		if !ok || err != nil {
			fmt.Println(err)
			l.println(fmt.Sprint("trying another collector..."))
			continue
		}

		l.println(fmt.Sprintf("collector %T selected\n", c))
		return c, nil
	}

	return nil, fmt.Errorf("there's no available log collector")
}

// triggerEvents trigger events. If EventChain enabled, then the next event at the same level will run only
// if the current event returns true.
func (l *Logger) triggerEvents(msg string) {
	if _, ok := l.Event[l.Level]; !ok {
		return
	}

	for _, event := range l.Event[l.Level] {
		event.Fire(msg, l.Defaults)
	}
}

func (l *Logger) println(msg string) {
	if l.SilentMode {
		return
	}

	fmt.Println(msg)
}

func Warning(v ...interface{})                    { logInstance.log(WarningLevel, fmt.Sprint(v...)) }
func WarningF(format string, args ...interface{}) { logInstance.log(WarningLevel, fmt.Sprintf(format, args...)) }

func Trace(v ...interface{})                    { logInstance.log(TraceLevel, fmt.Sprint(v...)) }
func TraceF(format string, args ...interface{}) { logInstance.log(TraceLevel, fmt.Sprintf(format, args...)) }

func Debug(v ...interface{})                    { logInstance.log(DebugLevel, fmt.Sprint(v...)) }
func DebugF(format string, args ...interface{}) { logInstance.log(DebugLevel, fmt.Sprintf(format, args...)) }

func Error(v ...interface{})                    { logInstance.log(ErrorLevel, fmt.Sprint(v...)) }
func ErrorF(format string, args ...interface{}) { logInstance.log(ErrorLevel, fmt.Sprintf(format, args...)) }

// Fatal log message and run os.Exit
func Fatal(v ...interface{}) { logInstance.log(FatalLevel, fmt.Sprint(v...)) }

// FatalF log message and run os.Exit
func FatalF(format string, args ...interface{}) { logInstance.log(FatalLevel, fmt.Sprintf(format, args...)) }

func Panic(v ...interface{})                    { logInstance.log(PanicLevel, fmt.Sprint(v...)) }
func PanicF(format string, args ...interface{}) { logInstance.log(PanicLevel, fmt.Sprintf(format, args...)) }
