package logging

import (
	"context"
	"io"
	"os"

	"github.com/rs/zerolog"
)

type getTracingID func(context.Context) string

var logger = LoggerWrapper{log: zerolog.New(os.Stderr).With().Timestamp().Logger()}

type LoggerWrapper struct {
	log zerolog.Logger
}

func NewLogger(output io.Writer, level Level) error {
	zLogLevel, err := level.ToZerologLevel()
	if err != nil {
		return err
	}

	logger = LoggerWrapper{log: zerolog.New(output).Level(zLogLevel).With().Timestamp().Logger()}
	return nil
}

type LoggerEventWrapper struct {
	logEvent *zerolog.Event
}

func (l LoggerEventWrapper) Msg(msg string) LoggerEventWrapper {
	l.logEvent.Msg(msg)
	return l
}

func (l LoggerEventWrapper) Str(key string, val string) LoggerEventWrapper {
	l.logEvent.Str(key, val)
	return l
}

func (l LoggerEventWrapper) Int(key string, val int) LoggerEventWrapper {
	l.logEvent.Int(key, val)
	return l
}

func (l LoggerEventWrapper) Err(err error) LoggerEventWrapper {
	l.logEvent.Err(err)
	return l
}

func Info() LoggerEventWrapper {
	return LoggerEventWrapper{logEvent: logger.log.Info()}
}

func Debug() LoggerEventWrapper {
	return LoggerEventWrapper{logEvent: logger.log.Debug()}
}

func Error() LoggerEventWrapper {
	return LoggerEventWrapper{logEvent: logger.log.Error()}
}

type tracingIDHook struct {
	tracingID getTracingID
}

// NewTracingIDHook returns a new logger hook.
// Required function getTracingID returns the tracing_id
// to add to the log event.
func NewTracingIDHook(tracingID getTracingID) zerolog.Hook {
	return tracingIDHook{tracingID: tracingID}
}

// Run implements zerolog.Hook interface.
func (h tracingIDHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	e.Str("tracing_id", h.tracingID(ctx)).Msg(msg)
}
