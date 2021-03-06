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
	Program() (program string)
	Pid() (pid int)
	Hostname() (hostname string)
	Time() (time time.Time)
	LogLevel() (logLevel string)
	LogLevelNum() (logLevelNum LogLevel)
	Pc() (pc uintptr)
	FileName() (fileName string)
	LineNum() (lineNum int)
	Message() (message string)
	SetAttr(key string, value interface{})
	GetAttr(key string) (value interface{})
	GetAttrs() map[string]interface{}
}

type logInfo struct {
	program  string
	pid      int
	hostname string
	time     time.Time
	logLevel LogLevel
	pc       uintptr
	fileName string
	lineNum  int
	message  string
	attrs    map[string]interface{}
}

//Program is return program
func (l *logInfo) Program() (program string) {
	return l.program
}

//Pid is return process id
func (l *logInfo) Pid() (pid int) {
	return l.pid
}

//Hostname is return hostname
func (l *logInfo) Hostname() (hostname string) {
	return l.hostname
}

//Time is return time
func (l *logInfo) Time() (time time.Time) {
	return l.time
}

//LogLevel is return log level
func (l *logInfo) LogLevel() (logLevel string) {
	logLevel, ok := logLevelMap[l.logLevel]
	if !ok {
		return "UNKNOWN"
	}
	return logLevel
}

//LogLevelNum is return log level number
func (l *logInfo) LogLevelNum() (logLevelNum LogLevel) {
	return l.logLevel
}

//Pc is return program counter
func (l *logInfo) Pc() (pc uintptr) {
	return l.pc
}

//FileName is return file name
func (l *logInfo) FileName() (fileName string) {
	return l.fileName
}

//LineNum is line number
func (l *logInfo) LineNum() (lineNo int) {
	return l.lineNum
}

//Message is return message
func (l *logInfo) Message() (message string) {
	return l.message
}

//SetAttr is set attribute
func (l *logInfo) SetAttr(key string, value interface{}) {
	if l.attrs == nil {
		l.attrs = make(map[string]interface{})
	}
	l.attrs[key] = value
}

//GetAttr is get attribute
func (l *logInfo) GetAttr(key string) (value interface{}) {
	if l.attrs == nil {
		return nil
	}
	value, ok := l.attrs[key]
	if !ok {
		return nil
	}
	return value
}

//GetAttrs is get attributes
func (l *logInfo) GetAttrs() map[string]interface{} {
	return l.attrs
}
