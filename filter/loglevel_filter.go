package filter

type LogLevelFilter struct {
	logLevel belog.LogLevel
	mutex    *sync.RWMutex
}

func (f *LogLevelFilter) Evaluate(loggerName string, logEvent belog.LogEvent) (ok bool) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	if logEvent.LogLevel() > f.LogLevel {
		return false
	}
	return true
}

func (f *LogLevelFilter) SetLogLevel(logLevel belog.LogLevel) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.logLevel = logLevel
}

func NewLogLevelFilter() (filter Filter) {
	return &LogLevelFilter{
		LogLevel: belog.LogLevelInfo,
		mutex:    new(sync.RWMutex),
	}
}

func init() {
	RegisterFilter("LogLevelFilter", NewLogLevelFilter)
}
