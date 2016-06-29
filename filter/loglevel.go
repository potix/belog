package filter

type LogLevelFilter struct {
	LogLevel belog.LogLevel
}

func (f *LogLevelFilter) Evaluate(logEvent belog.LogEvent) (ok bool) {
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
	RegisterFilter("LogLevel", NewLogLevelFilter)
}
