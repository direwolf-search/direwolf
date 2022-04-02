package logger

import (
	"context"
)

type Logger interface {
	Info(ctx context.Context, message string, errs ...error)
	Debug(ctx context.Context, message string, errs ...error)
	Warning(ctx context.Context, message string, errs ...error)
	Error(ctx context.Context, message string, errs ...error)
	Critical(ctx context.Context, message string, errs ...error)
	Fatal(ctx context.Context, message string, errs ...error)
}
