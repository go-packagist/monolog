package monolog

import (
	"fmt"
	"github.com/go-packagist/logger"
	"time"
)

type Logger struct {
	channel  string
	timezone *time.Location

	handlers   []Handler
	processors []Processor

	logger.Loggerable
}

type Opt func(*Logger)

func NewLogger(channel string, opts ...Opt) *Logger {
	l := &Logger{
		channel: channel,
	}

	for _, opt := range opts {
		opt(l)
	}

	l.init()

	return l
}

func (l *Logger) init() {
	if nil == l.timezone {
		WithTimezone(time.Local)(l)
	}

	l.setLoggerable()
}

func WithChannel(channel string) Opt {
	return func(l *Logger) {
		l.channel = channel
	}
}

func WithTimezone(tz *time.Location) Opt {
	return func(l *Logger) {
		l.timezone = tz
	}
}

func WithHandler(h Handler) Opt {
	return func(l *Logger) {
		l.handlers = append(l.handlers, h)
	}
}

func WithHandlers(hs ...Handler) Opt {
	return func(l *Logger) {
		l.handlers = append(l.handlers, hs...)
	}
}

func WithProcessor(p Processor) Opt {
	return func(l *Logger) {
		l.processors = append(l.processors, p)
	}
}

func WithProcessors(ps ...Processor) Opt {
	return func(l *Logger) {
		l.processors = append(l.processors, ps...)
	}
}

func (l *Logger) Channel() string {
	return l.channel
}

func (l *Logger) Handlers() []Handler {
	return l.handlers
}

func (l *Logger) Processors() []Processor {
	return l.processors
}

func (l *Logger) setLoggerable() {
	l.Loggerable = func(level logger.Level, s string) {
		record := &Record{
			Channel: l.Channel(),
			Message: s,
			Level:   level,
			Time:    time.Now().In(l.timezone),
		}

		for _, handler := range l.Handlers() {
			if !handler.IsHandling(record) {
				continue // skip
			}

			if true == handler.Handle(record) {
				break // handled
			}
		}
	}
}

func (l *Logger) Close() error {
	var errs []error

	for _, handler := range l.handlers {
		if err := handler.Close(); nil != err {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return fmt.Errorf("errors: %v", errs)
}
