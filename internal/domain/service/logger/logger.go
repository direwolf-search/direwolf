package logger

import (
	"context"
)

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(err error, msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Warning(msg string, keysAndValues ...interface{})
	Critical(err error, msg string, keysAndValues ...interface{})
	Fatal(err error, msg string, keysAndValues ...interface{})

	InfoWithContext(ctx context.Context, msg string, keysAndValues ...interface{})
	ErrorWithContext(ctx context.Context, err error, msg string, keysAndValues ...interface{})
	DebugWithContext(ctx context.Context, msg string, keysAndValues ...interface{})
	WarningWithContext(ctx context.Context, msg string, keysAndValues ...interface{})
	CriticalWithContext(ctx context.Context, err error, msg string, keysAndValues ...interface{})
	FatalWithContext(ctx context.Context, err error, msg string, keysAndValues ...interface{})
}
