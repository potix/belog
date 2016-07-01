package filter

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
	if logEvent.LogLevel() > f.LogLevel {
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
		LogLevel: LogLevelInfo,
		mutex:    new(sync.RWMutex),
	}
}

func init() {
	RegisterFilter("LogLevelFilter", NewLogLevelFilter)
}