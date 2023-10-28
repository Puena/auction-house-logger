package logging

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

type Level int8

const (
	// DebugLevel defines debug log level.
	DebugLevel Level = iota
	// InfoLevel defines info log level.
	InfoLevel
	// WarnLevel defines warn log level.
	WarnLevel
	// ErrorLevel defines error log level.
	ErrorLevel
	// FatalLevel defines fatal log level.
	FatalLevel
	// PanicLevel defines panic log level.
	PanicLevel
	// NoLevel defines an absent log level.
	NoLevel
	// Disabled disables the logger.
	Disabled

	// TraceLevel defines trace log level.
	TraceLevel Level = -1
	// Values less than TraceLevel are handled as numbers.
)

var (
	LevelFieldMarshalFunc = func(l Level) string {
		return l.String()
	}
)

func (l Level) String() string {
	switch l {
	case TraceLevel:
		return "trace"
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	case FatalLevel:
		return "fatal"
	case PanicLevel:
		return "panic"
	case Disabled:
		return "disabled"
	case NoLevel:
		return ""
	}
	return strconv.Itoa(int(l))
}

func (l Level) ToZerologLevel() (zerolog.Level, error) {
	switch l {
	case TraceLevel:
		return zerolog.TraceLevel, nil
	case DebugLevel:
		return zerolog.DebugLevel, nil
	case InfoLevel:
		return zerolog.InfoLevel, nil
	case WarnLevel:
		return zerolog.WarnLevel, nil
	case ErrorLevel:
		return zerolog.ErrorLevel, nil
	case FatalLevel:
		return zerolog.FatalLevel, nil
	case PanicLevel:
		return zerolog.PanicLevel, nil
	case Disabled:
		return zerolog.Disabled, nil
	case NoLevel:
		return zerolog.NoLevel, nil
	}
	return zerolog.NoLevel, fmt.Errorf("Unknown Level: '%d'", l)
}

// ParseLevel converts a level string into a zerolog Level value.
// returns an error if the input string does not match known values.
func ParseLevel(levelStr string) (Level, error) {
	switch {
	case strings.EqualFold(levelStr, LevelFieldMarshalFunc(TraceLevel)):
		return TraceLevel, nil
	case strings.EqualFold(levelStr, LevelFieldMarshalFunc(DebugLevel)):
		return DebugLevel, nil
	case strings.EqualFold(levelStr, LevelFieldMarshalFunc(InfoLevel)):
		return InfoLevel, nil
	case strings.EqualFold(levelStr, LevelFieldMarshalFunc(WarnLevel)):
		return WarnLevel, nil
	case strings.EqualFold(levelStr, LevelFieldMarshalFunc(ErrorLevel)):
		return ErrorLevel, nil
	case strings.EqualFold(levelStr, LevelFieldMarshalFunc(FatalLevel)):
		return FatalLevel, nil
	case strings.EqualFold(levelStr, LevelFieldMarshalFunc(PanicLevel)):
		return PanicLevel, nil
	case strings.EqualFold(levelStr, LevelFieldMarshalFunc(Disabled)):
		return Disabled, nil
	case strings.EqualFold(levelStr, LevelFieldMarshalFunc(NoLevel)):
		return NoLevel, nil
	}
	i, err := strconv.Atoi(levelStr)
	if err != nil {
		return NoLevel, fmt.Errorf("Unknown Level String: '%s', defaulting to NoLevel", levelStr)
	}
	if i > 127 || i < -128 {
		return NoLevel, fmt.Errorf("Out-Of-Bounds Level: '%d', defaulting to NoLevel", i)
	}
	return Level(i), nil
}

// UnmarshalText implements encoding.TextUnmarshaler to allow for easy reading from toml/yaml/json formats
func (l *Level) UnmarshalText(text []byte) error {
	if l == nil {
		return errors.New("can't unmarshal a nil *Level")
	}
	var err error
	*l, err = ParseLevel(string(text))
	return err
}

// MarshalText implements encoding.TextMarshaler to allow for easy writing into toml/yaml/json formats
func (l Level) MarshalText() ([]byte, error) {
	return []byte(LevelFieldMarshalFunc(l)), nil
}
