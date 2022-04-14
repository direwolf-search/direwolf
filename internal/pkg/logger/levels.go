package logger

type logLevel uint8

func (ll logLevel) Int() int {
	return int(ll)
}

func (ll logLevel) String() string {
	return logLevels[ll]
}

const (
	infoLevel logLevel = iota
	debugLevel
	errorLevel
	warningLevel
	criticalLevel
	fatalLevel
)

var logLevels = []string{
	"INFO: ",
	"DEBUG: ",
	"ERROR: ",
	"WARNING: ",
	"CRITICAL: ",
	"FATAL: ",
}
