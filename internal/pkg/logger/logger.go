// Package logger implements logger interface with logging levels and custom configs
package logger

import (
	"log"
	"os"

	"direwolf/internal/domain"
)

var (
	defaultOptions     = log.Lshortfile | log.Ldate | log.Ltime
	defaultWriter      = os.Stdout
	defaultErrorWriter = os.Stderr
)

func stdLogger(level logLevel) *log.Logger {
	if level == errorLevel {
		return log.New(defaultErrorWriter, level.String(), defaultOptions)
	}
	return log.New(defaultWriter, level.String(), defaultOptions)
}

// logger implements domain.Logger and github.com/robfig/cron/v3 Logger
type logger struct {
	infoLogger, errorLogger, debugLogger, warningLogger, criticalLogger, fatalLogger, printfLogger *log.Logger
}

func NewLogger(opts ...loggerOption) domain.Logger {
	l := &logger{
		infoLogger:     stdLogger(infoLevel),
		errorLogger:    stdLogger(errorLevel),
		debugLogger:    stdLogger(debugLevel),
		warningLogger:  stdLogger(warningLevel),
		criticalLogger: stdLogger(criticalLevel),
		fatalLogger:    stdLogger(fatalLevel),
		printfLogger:   log.New(defaultWriter, "", defaultOptions),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

// Info prints to the underlying standard logger with prefix "INFO: ".
// Call of Info equivalent to log.Printf() with msg as first formatted value
func (l *logger) Info(msg string, keysAndValues ...interface{}) {
	var (
		format = "%v "
		values = make([]interface{}, 0)
	)

	values = append(values, msg)

	if len(keysAndValues) > 0 {
		for _, val := range keysAndValues {
			if v, ok := val.(map[string]interface{}); ok {
				for key, value := range v {
					format += key + ": %v "
					values = append(values, value)
				}
			} else {
				format += " %v, "
				values = append(values, val)
			}
		}

		l.infoLogger.Printf(format, values...)
	} else {
		l.infoLogger.Println(msg)
	}
}

// Error prints to the underlying standard logger with prefix "ERROR: ".
// Call of Error equivalent to log.Printf() with msg as first and err as second formatted values
func (l *logger) Error(err error, msg string, keysAndValues ...interface{}) {
	var (
		format = "%v error: %s "
		values = make([]interface{}, 0)
	)

	values = append(values, msg)
	values = append(values, err)

	if len(keysAndValues) > 0 {
		for _, val := range keysAndValues {
			if v, ok := val.(map[string]interface{}); ok {
				for key, value := range v {
					format += key + ": %v "
					values = append(values, value)
				}
			} else {
				format += " %v, "
				values = append(values, val)
			}
		}

		l.errorLogger.Printf(format, values...)
	} else {
		l.errorLogger.Println(msg)
	}
}

// Debug prints to the underlying standard logger with prefix "DEBUG: ".
// Call of Debug equivalent to log.Printf() with msg as first formatted value
func (l *logger) Debug(msg string, keysAndValues ...interface{}) {
	var (
		format = "%v "
		values = make([]interface{}, 0)
	)

	values = append(values, msg)

	if len(keysAndValues) > 0 {
		for _, val := range keysAndValues {
			if v, ok := val.(map[string]interface{}); ok {
				for key, value := range v {
					format += key + ": %v "
					values = append(values, value)
				}
			} else {
				format += " %v, "
				values = append(values, val)
			}
		}

		l.debugLogger.Printf(format, values...)
	} else {
		l.debugLogger.Println(msg)
	}
}

// Warning prints to the underlying standard logger with prefix "WARNING: ".
// Call of Warning equivalent to log.Printf() with msg as first formatted value
func (l *logger) Warning(msg string, keysAndValues ...interface{}) {
	var (
		format = "%v "
		values = make([]interface{}, 0)
	)

	values = append(values, msg)

	if len(keysAndValues) > 0 {
		for _, val := range keysAndValues {
			if v, ok := val.(map[string]interface{}); ok {
				for key, value := range v {
					format += key + ": %v "
					values = append(values, value)
				}
			} else {
				format += " %v, "
				values = append(values, val)
			}
		}

		l.warningLogger.Printf(format, values...)
	} else {
		l.warningLogger.Println(msg)
	}
}

// Critical prints to the underlying standard logger with prefix "CRITICAL: ".
// Call of Critical equivalent to log.Printf() with msg as first and err as second formatted values
func (l *logger) Critical(err error, msg string, keysAndValues ...interface{}) {
	var (
		format = "%v error: %s "
		values = make([]interface{}, 0)
	)

	values = append(values, msg)
	values = append(values, err)

	if len(keysAndValues) > 0 {
		for _, val := range keysAndValues {
			if v, ok := val.(map[string]interface{}); ok {
				for key, value := range v {
					format += key + ": %v "
					values = append(values, value)
				}
			} else {
				format += " %v, "
				values = append(values, val)
			}
		}

		l.criticalLogger.Printf(format, values...)
	} else {
		l.criticalLogger.Println(msg)
	}
}

// Fatal prints to the underlying standard logger with prefix "FATAL: ".
// Call of Fatal equivalent to log.Printf() with msg as first and err as second formatted values
func (l *logger) Fatal(err error, msg string, keysAndValues ...interface{}) {
	var (
		format = "%v error: %s "
		values = make([]interface{}, 0)
	)

	values = append(values, msg)
	values = append(values, err)

	if len(keysAndValues) > 0 {
		for _, val := range keysAndValues {
			if v, ok := val.(map[string]interface{}); ok {
				for key, value := range v {
					format += key + ": %v "
					values = append(values, value)
				}
			} else {
				format += " %v, "
				values = append(values, val)
			}
		}

		l.fatalLogger.Printf(format, values...)
	} else {
		l.fatalLogger.Println(msg)
	}
}

func (l *logger) Printf(format string, ii ...interface{}) {
	l.printfLogger.Printf(format, ii...)
}
