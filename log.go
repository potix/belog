package belog

var (
	logLevelMap = map[LogLevel]string{
		LogLevelEmerg:  "EMERG",
		LogLevelAlert:  "ALERT",
		LogLevelCrit:   "CRIT",
		LogLevelError:  "ERROR",
		LogLevelWarn:   "WARN",
		LogLevelNotice: "NOTICE",
		LogLevelInfo:   "INFO",
		LogLevelDebug:  "DEBUG",
		LogLevelTrace:  "TRACE",
	}
)

type LogEvent interface {
	Pid() int
	Time() time.Time
	LogLevel() LogLevel
	LogLevelString() string
	Pc() uintptr
	FileName() string
	LineNum() int
	Message() string
}

type log struct {
	pid      int
	time     time.Time
	logLevel LogLevel
	pc       uintptr
	fileName string
	lineNum  int
	message  string
}

func (l *log) Pid() (pid int) {
	return l.pid
}

func (l *log) Time() (time time.Time) {
	return l.time
}

func (l *log) LogLevel() (logLevel LogLevel) {
	return l.logLevel
}

func (l *log) LogLevelString() (logLevelString string) {
	return logLevelMap[l.logLevel]
}

func (l *log) Pc() (pc uintptr) {
	return l.pc
}

func (l *log) FileName() (fileName string) {
	return l.fileName
}

func (l *log) LineNum() (lineNo int) {
	return l.lineNum
}

func (l *log) Message() (message string) {
	return l.message
}
