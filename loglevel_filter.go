package belog

import (
	"sync"
)

//LogLevelFilter is filter of log level
type LogLevelFilter struct {
	logLevel LogLevel
	mutex    *sync.RWMutex
}

//Evaluate is Evaluate log event
func (f *LogLevelFilter) Evaluate(loggerName string, logEvent LogEvent) (ok bool) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	if logEvent.LogLevelNum() > f.logLevel {
		return false
	}
	return true
}

//SetLogLevel is set logger level. outputs the important than this log level.
func (f *LogLevelFilter) SetLogLevel(logLevel LogLevel) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.logLevel = logLevel
}

//NewLogLevelFilter is create LogLevelFilter
func NewLogLevelFilter() (logLevelFilter *LogLevelFilter) {
	return &LogLevelFilter{
		logLevel: LogLevelInfo,
		mutex:    new(sync.RWMutex),
	}
}

func init() {
	RegisterFilter("LogLevelFilter", func() (filter Filter) {
		return NewLogLevelFilter()
	})
}
