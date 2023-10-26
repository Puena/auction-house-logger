package logging

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	globalLog "github.com/rs/zerolog/log"
)

var Logger = globalLog.Output(os.Stdout)

type LogWrapper struct {
	log *zerolog.Logger
}

type LogWrapperEvent struct {
	logEvent *zerolog.Event
}

func (l *LogWrapperEvent) Msg(msg string) *LogWrapperEvent {
	l.logEvent.Msg(msg)
	return l
}

func (l *LogWrapperEvent) Str(key string, val string) *LogWrapperEvent {
	l.logEvent.Str(key, val)
	return l
}

func (l *LogWrapperEvent) Int(key string, val int) *LogWrapperEvent {
	l.logEvent.Int(key, val)
	return l
}

func (l *LogWrapperEvent) Err(err error) *LogWrapperEvent {
	l.logEvent.Err(err)
	return l
}

func (l *LogWrapper) Info() *LogWrapperEvent {
	return &LogWrapperEvent{logEvent: l.log.Info()}
}

func (l *LogWrapper) Debug() *LogWrapperEvent {
	return &LogWrapperEvent{logEvent: l.log.Debug()}
}

func (l *LogWrapper) Error() *LogWrapperEvent {
	return &LogWrapperEvent{logEvent: l.log.Error()}
}

func NewLogWrapper() *LogWrapper {
	return &LogWrapper{log: &Logger}
}

type getTracingID func(context.Context) string

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
