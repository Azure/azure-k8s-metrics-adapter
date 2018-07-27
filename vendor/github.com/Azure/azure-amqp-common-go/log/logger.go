package log

import (
	"context"

	"github.com/opentracing/opentracing-go"
	tag "github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

type (
	// Logger is the interface for opentracing logging
	Logger interface {
		Info(msg string, fields ...log.Field)
		Error(err error, fields ...log.Field)
		Fatal(msg string, fields ...log.Field)
		Debug(msg string, fields ...log.Field)
	}

	spanLogger struct {
		span opentracing.Span
	}

	nopLogger struct{}
)

// For will return a logger for a given context
func For(ctx context.Context) Logger {
	if span := opentracing.SpanFromContext(ctx); span != nil {
		return &spanLogger{
			span: span,
		}
	}
	return new(nopLogger)
}

func (sl spanLogger) Info(msg string, fields ...log.Field) {
	sl.logToSpan("info", msg, fields...)
}

func (sl spanLogger) Error(err error, fields ...log.Field) {
	tag.Error.Set(sl.span, true)
	sl.logToSpan("error", err.Error(), fields...)
}

func (sl spanLogger) Fatal(msg string, fields ...log.Field) {
	tag.Error.Set(sl.span, true)
	sl.logToSpan("fatal", msg, fields...)
}

func (sl spanLogger) Debug(msg string, fields ...log.Field) {
	sl.logToSpan("debug", msg, fields...)
}

func (sl spanLogger) logToSpan(level string, msg string, fields ...log.Field) {
	fields = append(fields, log.String("event", msg), log.String("level", level))
	sl.span.LogFields(fields...)
}

func (sl nopLogger) Info(msg string, fields ...log.Field)  {}
func (sl nopLogger) Error(err error, fields ...log.Field)  {}
func (sl nopLogger) Fatal(msg string, fields ...log.Field) {}
func (sl nopLogger) Debug(msg string, fields ...log.Field) {}
