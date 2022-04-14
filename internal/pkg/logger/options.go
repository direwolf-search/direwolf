package logger

import (
	"io"
	"log"
)

type loggerOption func(l *logger)

func setFields(level logLevel, options int, writer io.Writer) *log.Logger {
	var (
		num int
		w   io.Writer
	)

	if options != 0 {
		num = options
	} else {
		num = defaultOptions
	}

	if writer != nil {
		w = writer
	} else {
		w = defaultErrorWriter
	}

	return log.New(w, level.String(), num)
}

func InfoOptions(options int, writer io.Writer) loggerOption {
	stdLogger := setFields(infoLevel, options, writer)
	return func(l *logger) {
		l.infoLogger = stdLogger
	}
}
func ErrorOptions(options int, writer io.Writer) loggerOption {
	stdLogger := setFields(errorLevel, options, writer)
	return func(l *logger) {
		l.errorLogger = stdLogger
	}
}
func DebugOptions(options int, writer io.Writer) loggerOption {
	stdLogger := setFields(debugLevel, options, writer)
	return func(l *logger) {
		l.debugLogger = stdLogger
	}
}
func WarningOptions(options int, writer io.Writer) loggerOption {
	stdLogger := setFields(warningLevel, options, writer)
	return func(l *logger) {
		l.warningLogger = stdLogger
	}
}
func CriticalOptions(options int, writer io.Writer) loggerOption {
	stdLogger := setFields(criticalLevel, options, writer)
	return func(l *logger) {
		l.criticalLogger = stdLogger
	}
}
func FatalOptions(options int, writer io.Writer) loggerOption {
	stdLogger := setFields(fatalLevel, options, writer)
	return func(l *logger) {
		l.fatalLogger = stdLogger
	}
}
