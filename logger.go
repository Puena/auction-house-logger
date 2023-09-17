package logging

import (
	"context"
	"os"

	"github.com/rs/zerolog"
	globalLog "github.com/rs/zerolog/log"
)

var Logger = globalLog.Output(os.Stdout)

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
