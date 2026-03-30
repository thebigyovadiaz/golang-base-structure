// Package logger provides support for initializing the log system using the gcp logger format.
package logger

import (
	"context"
	"io"
	"log"
	"runtime"
	"time"

	"log/slog"
)

// Level represents different logging levels.
type Level slog.Level

// A set of possible logging levels.
const (
	LevelDebug = Level(slog.LevelDebug)
	LevelInfo  = Level(slog.LevelInfo)
	LevelWarn  = Level(slog.LevelWarn)
	LevelError = Level(slog.LevelError)
)

// RequiredFieldsFunc represents a function that can return required fields to be logged at root level, like the trace id from
// the specified context.
type RequiredFieldsFunc func(ctx context.Context) []any

// Logger represents a logger for logging information.
type Logger struct {
	handler            slog.Handler
	requiredFieldsFunc RequiredFieldsFunc
}

// New constructs a new log for application use.
func New(w io.Writer, minLevel Level, serviceName string, requiredFieldsFunc RequiredFieldsFunc) *Logger {
	// Replace msg, level, source, and time keys to message, severity, timestamp, and file respectively.
	f := func(groups []string, a slog.Attr) slog.Attr {
		switch a.Key {
		case slog.MessageKey:
			return slog.Attr{Key: "message", Value: a.Value}

		case slog.LevelKey:
			return slog.Attr{Key: "severity", Value: a.Value}

		case slog.TimeKey:
			return slog.Attr{Key: "timestamp", Value: a.Value}

		case slog.SourceKey:
			return slog.Attr{Key: "source", Value: a.Value}
		}

		return a
	}

	// Construct the slog JSON handler for use.
	handler := slog.Handler(slog.NewJSONHandler(w, &slog.HandlerOptions{AddSource: true, Level: slog.Level(minLevel), ReplaceAttr: f}))

	// Attributes to add to every log.
	attrs := []slog.Attr{
		{Key: "serviceID", Value: slog.StringValue(serviceName)},
	}

	// Add those attributes and capture the final handler.
	handler = handler.WithAttrs(attrs)

	return &Logger{
		handler:            handler,
		requiredFieldsFunc: requiredFieldsFunc,
	}
}

// NewStdLogger returns a standard library Logger that wraps the slog Logger.
func NewStdLogger(logger *Logger, level Level) *log.Logger {
	return slog.NewLogLogger(logger.handler, slog.Level(level))
}

// Debug logs at LevelDebug with the given context.
func (log *Logger) Debug(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelDebug, 3, msg, args...)
}

// Debugc logs the information at the specified call stack position.
func (log *Logger) Debugc(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, LevelDebug, caller, msg, args...)
}

// Info logs at LevelInfo with the given context.
func (log *Logger) Info(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelInfo, 3, msg, args...)
}

// Infoc logs the information at the specified call stack position.
func (log *Logger) Infoc(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, LevelInfo, caller, msg, args...)
}

// Warn logs at LevelWarn with the given context.
func (log *Logger) Warn(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelWarn, 3, msg, args...)
}

// Warnc logs the information at the specified call stack position.
func (log *Logger) Warnc(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, LevelWarn, caller, msg, args...)
}

// Error logs at LevelError with the given context.
func (log *Logger) Error(ctx context.Context, msg string, args ...any) {
	log.write(ctx, LevelError, 3, msg, args...)
}

// Errorc logs the information at the specified call stack position.
func (log *Logger) Errorc(ctx context.Context, caller int, msg string, args ...any) {
	log.write(ctx, LevelError, caller, msg, args...)
}

func (log *Logger) write(ctx context.Context, level Level, caller int, msg string, args ...any) {
	slogLevel := slog.Level(level)

	if !log.handler.Enabled(ctx, slogLevel) {
		return
	}

	var pcs [1]uintptr
	runtime.Callers(caller, pcs[:])

	r := slog.NewRecord(time.Now(), slogLevel, msg, pcs[0])
	//r := slog.Record{Level: slogLevel, PC: pcs[0]}

	if log.requiredFieldsFunc != nil {
		r.Add(log.requiredFieldsFunc(ctx)...)
	}
	r.AddAttrs(slog.Group("customFields", args...))
	//r.Add(args...)

	log.handler.Handle(ctx, r)
}
