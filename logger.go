package belog

import (
	"log"
	"sync"
	"time"
)

type logLevel int

const (
	logLevelEmergency logLevel = 1
	logLevelAlert     logLevel = 2
	logLevelCritical  logLevel = 3
	logLevelError     logLevel = 4
	logLevelWarning   logLevel = 5
	logLevelInfo      logLevel = 6
	logLevelDebug     logLevel = 7
	logLevelTrace     logLevel = 8
)

var (
	logLevelMap = map[logLevel]string{
		logLevelEmergency: "EMERG",
		logLevelAlert:     "ALERT",
		logLevelCritical:  "CRIT",
		logLevelError:     "ERROR",
		logLevelWarning:   "WARN",
		logLevelInfo:      "INFO",
		logLevelDebug:     "DEBUG",
		logLevelTrace:     "TRACE",
	}
	loggers map[string]*Logger{}
	mutex   sync.RWMutex
)

type LogEvent interface {
	Time() time
	logLevel() int
	logLevelString() int
	FileName() string
	FuncName() string
	LineNo() int
	Message()
}

type log struct {
	time time
	logLevel int
	fileName string
	funcName string
	lineNo int
	msg string
}

func (l *log) Time() (time time) {
	return l.time
}

func (l *log) XXXGETTER() (XXX XXX) {
	return XXXX
}


type Loggers struct {
	loggers []*Logger
}

func (l *Loggers )printf() {
	XXXXXX
}

func (l *Loggers )Info() {

}

func (l *Loggers )Warn() {
}


type Logger struct {
	name string
	// filter interface
	filter Filter
	// formatter interface
	formatter Formatter
	// appender interfaces
	appenders []Appender
	mutex     sync.RWMutex
}

func (l *Logger) printf() {
	
}

func GetLoggers(name ...string) (loggers *Loggers) {

}




