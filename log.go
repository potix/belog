package belog

import (
	"time"
)

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

//LogEvent is interface of event of log
type LogEvent interface {
	Pid() int
	Time() time.Time
	LogLevel() string
	LogLevelNum() LogLevel
	Pc() uintptr
	FileName() string
	LineNum() int
	Message() string
}

type logInfo struct {
	pid      int
	time     time.Time
	logLevel LogLevel
	pc       uintptr
	fileName string
	lineNum  int
	message  string
}

func (l *logInfo) Pid() (pid int) {
	return l.pid
}

func (l *logInfo) Time() (time time.Time) {
	return l.time
}

func (l *logInfo) LogLevel() (logLevel string) {
	logLevel, ok := logLevelMap[l.logLevel]
	if !ok {
		return "UNKNOWN"
	}
	return logLevel
}

func (l *logInfo) LogLevelNum() (logLevelNum LogLevel) {
	return l.logLevel
}

func (l *logInfo) Pc() (pc uintptr) {
	return l.pc
}

func (l *logInfo) FileName() (fileName string) {
	return l.fileName
}

func (l *logInfo) LineNum() (lineNo int) {
	return l.lineNum
}

func (l *logInfo) Message() (message string) {
	return l.message
}
