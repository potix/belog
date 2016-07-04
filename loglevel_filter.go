package belog

import (
	"sync"
)

//LogLevelFilter is filter of log level
type LogLevelFilter struct {
	logLevel    LogLevel
	chainFilter Filter
	mutex       *sync.RWMutex
}

//Evaluate is Evaluate log event
func (f *LogLevelFilter) Evaluate(loggerName string, logEvent LogEvent) (ok bool) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	if logEvent.LogLevelNum() > f.logLevel {
		return false
	}
	if f.chainFilter == nil {
		return true
	}
	return f.chainFilter.Evaluate(loggerName, logEvent)
}

//SetLogLevel is set logger level. outputs the important than this log level.
func (f *LogLevelFilter) SetLogLevel(logLevel LogLevel) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.logLevel = logLevel
}

//SetChainFilter is set chain filter.
func (f *LogLevelFilter) SetChainFilter(chainFilter Filter) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.chainFilter = chainFilter
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
