package filter

type LogLevelFilter struct {
	LogLevel belog.LogLevel
}

func (f *LogLevelFilter) Evaluate(loggerName string, logEvent belog.LogEvent) (ok bool) {
	if logEvent.LogLevel() > f.LogLevel {
		return false
	}
	return true
}

func NewLogLevelFilter() (filter Filter) {
	return &LogLevelFilter{
		LogLevel: belog.LogLevelInfo,
	}
}

func init() {
	RegisterFilter("LogLevelFilter", NewLogLevelFilter)
}
