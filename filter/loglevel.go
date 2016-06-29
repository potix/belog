package filter

import sync

type LogLevelFilter struct {
	logLevel belog.LogLevel
}

func (f *LogLevelFilter) Evaluate(log belog.Log) (ok bool) {
	if log.LogLevel() > f.logLevel {
		return false
	}
	return true
}

func NewLogLevelFilter(logLevel belog.LogLevel) (filter Filter) {
	return &LogLevelFilter{
		logLevel: logLevel,
	}
}
