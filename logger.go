package logger

import (
	"context"
	"fmt"
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
	Ping(context.Context) bool
	Write(context.Context, []byte)
}

type OutputFormatter interface {
	Format(context.Context, string, map[string]string) ([]byte, error)
}

type Event interface {
	Fire(context.Context, string, map[string]interface{})
}

type Logger struct {
	Level      LogLevel
	Defaults   map[string]interface{}
	Collector  []Collector
	Formatter  OutputFormatter
	Event      map[LogLevel][]Event
	EventChain bool // If true, events at the same level will run only if the previous event returns true.
	Timeout    time.Duration
	Lock       bool // Not used yet
	ctx        context.Context
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
// Please note that, never pass a timeout context instead set Timeout.
func New(ctx context.Context, logger *Logger) error {
	logInstance = logger
	logInstance.ctx = ctx

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

	// If Timeout passed to the logger initializer we can make a timeout context and pass it to the `handler` function.
	if l.Timeout.Seconds() > 0 {
		ctx, cnl := context.WithTimeout(l.ctx, l.Timeout)
		defer cnl()

		select {
		case <-ctx.Done():
			return fmt.Errorf("log preparing canceled because: %s", l.ctx.Err())
		default:
			return l.handle(ctx, msg)
		}
	}

	select {
	case <-l.ctx.Done():
		return fmt.Errorf("log preparing canceled because: %s", l.ctx.Err())
	default:
		return l.handle(l.ctx, msg)
	}
}

// handle don't use Logger ctx directly instead uses context passed from log function
func (l *Logger) handle(ctx context.Context, msg string) error {
	c, err := l.collector(ctx)
	if err != nil {
		return err
	}

	r, err := l.formatter().Format(ctx, l.Level.Label(), mergeDefaultsToOutput(map[string]string{
		"time":  time.Now().String(),
		"level": l.Level.Label(),
		"msg":   msg,
	}, l.Defaults))
	if err != nil {
		return err
	}

	c.Write(ctx, r)
	l.triggerEvents(ctx, msg)

	return nil
}

func (l *Logger) isLevelEnabled(level LogLevel) bool {
	return l.Level >= level
}

func (l *Logger) formatter() OutputFormatter {
	return l.Formatter
}

// collector returns first available log collector
func (l *Logger) collector(ctx context.Context) (Collector, error) {
	for _, c := range l.Collector {
		if c.Ping(ctx) {
			return c, nil
		}
	}

	return nil, fmt.Errorf("there's no available log collector")
}

// triggerEvents trigger events. If EventChain enabled, then the next event at the same level will run only
// if the current event returns true.
func (l *Logger) triggerEvents(ctx context.Context, msg string) {
	if _, ok := l.Event[l.Level]; !ok {
		return
	}

	select {
	case <-ctx.Done():
		return
	default:
		for _, event := range l.Event[l.Level] {
			event.Fire(ctx, msg, l.Defaults)
		}
	}
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
