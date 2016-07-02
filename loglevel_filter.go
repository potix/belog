package belog

import (
	"sync"
)

type LogLevelFilter struct {
	logLevel LogLevel
	mutex    *sync.RWMutex
}

func (f *LogLevelFilter) Evaluate(loggerName string, logEvent LogEvent) (ok bool) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	if logEvent.LogLevelNum() > f.logLevel {
		return false
	}
	return true
}

func (f *LogLevelFilter) SetLogLevel(logLevel LogLevel) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.logLevel = logLevel
}

func NewLogLevelFilter() (filter Filter) {
	return &LogLevelFilter{
		logLevel: LogLevelInfo,
		mutex:    new(sync.RWMutex),
	}
}

func init() {
	RegisterFilter("LogLevelFilter", NewLogLevelFilter)
}
