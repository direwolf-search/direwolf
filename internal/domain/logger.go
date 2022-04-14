package domain

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(err error, msg string, keysAndValues ...interface{})
	Printf(format string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Warning(msg string, keysAndValues ...interface{})
	Critical(err error, msg string, keysAndValues ...interface{})
	Fatal(err error, msg string, keysAndValues ...interface{})
}
