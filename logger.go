package logging

import (
	"context"
	"errors"
	"io"
	"os"
	"sync"

	"github.com/rs/zerolog"
)

type getTracingID func(context.Context) string
type loggerOption func(logger *LoggerWrapper) error

var logger LoggerWrapper = LoggerWrapper{log: zerolog.New(os.Stderr).With().Timestamp().Logger()}

type LoggerWrapper struct {
	log zerolog.Logger
}

// WithOutput set output for logger.
func WithOutput(output io.Writer) loggerOption {
	return func(logger *LoggerWrapper) error {
		logger.log = logger.log.Output(output)
		return nil
	}
}

// WithLevel set level for logger.
func WithLevel(level Level) loggerOption {
	return func(logger *LoggerWrapper) error {
		zLogLevel, err := level.ToZerologLevel()
		if err != nil {
			return err
		}

		logger.log = logger.log.Level(zLogLevel)
		return nil
	}
}

// SetUpLogger set up logger, can be and should be called only once.
func SetUpLogger(opts ...loggerOption) (err error) {

	sync.OnceFunc(func() {
		for _, option := range opts {
			if oErr := option(&logger); oErr != nil {
				err = errors.Join(err, oErr)
			}
		}
	})()

	return err
}

type LoggerEventWrapper struct {
	logEvent *zerolog.Event
}

// NOTICE: Should be called only once per event, otherwise unpredictable behavior will be occur.
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

func (l LoggerEventWrapper) Int64(key string, val int64) LoggerEventWrapper {
	l.logEvent.Int64(key, val)
	return l
}

func (l LoggerEventWrapper) Err(err error) LoggerEventWrapper {
	l.logEvent.Err(err)
	return l
}

func (l LoggerEventWrapper) Any(key string, val any) LoggerEventWrapper {
	l.logEvent.Any(key, val)
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

func Fatal() LoggerEventWrapper {
	return LoggerEventWrapper{logEvent: logger.log.Fatal()}
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
