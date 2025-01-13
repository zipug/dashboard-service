package logger

import (
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

type Logger struct {
	log *zerolog.Logger
}

type LoggerEvent func(*zerolog.Event)

type LoggerAction string

const (
	Info  LoggerAction = "info"
	Debug LoggerAction = "debug"
	Warn  LoggerAction = "warn"
	Error LoggerAction = "error"
)

func WithStrAttr(key, value string) LoggerEvent {
	return func(e *zerolog.Event) {
		e.Str(key, value)
	}
}

func WithInt64Attr(key string, value int64) LoggerEvent {
	return func(e *zerolog.Event) {
		e.Int64(key, value)
	}
}

func WithUInt64Attr(key string, value uint64) LoggerEvent {
	return func(e *zerolog.Event) {
		e.Uint64(key, value)
	}
}

func WithStringArrayAttr(key string, value []string) LoggerEvent {
	return func(e *zerolog.Event) {
		arr := zerolog.Arr()
		for _, el := range value {
			arr.Str(el)
		}
		e.Array(key, arr)
	}
}

func WithInt64ArrayAttr(key string, value []int64) LoggerEvent {
	return func(e *zerolog.Event) {
		arr := zerolog.Arr()
		for _, el := range value {
			arr.Int64(int64(el))
		}
		e.Array(key, arr)
	}
}

func WithErrAttr(value error) LoggerEvent {
	return func(e *zerolog.Event) {
		e.Err(value)
	}
}

func New() *Logger {
	output := zerolog.ConsoleWriter{Out: os.Stdout}
	log := zerolog.New(output).
		With().
		Timestamp().
		Logger().
		Level(zerolog.DebugLevel)
	return &Logger{&log}
}

func (log *Logger) Log(action LoggerAction, message string, opts ...LoggerEvent) {
	l := log.log
	var event *zerolog.Event
	switch action {
	case Info:
		event = l.Info()
	case Debug:
		event = l.Debug()
	case Warn:
		event = l.Warn()
	case Error:
		event = l.Error()
	default:
		fmt.Println("invalid log type")
		return
	}

	for _, opt := range opts {
		opt(event)
	}

	event.Msg(message)
}
