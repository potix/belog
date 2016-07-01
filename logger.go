package belog

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"runtime"
	"sync"
	"time"
)

type LogLevel int

const (
	LogLevelEmerg  LogLevel = 1
	LogLevelAlert           = 2
	LogLevelCrit            = 3
	LogLevelError           = 4
	LogLevelWarn            = 5
	LogLevelNotice          = 6
	LogLevelInfo            = 7
	LogLevelDebug           = 8
	LogLevelTrace           = 9
)

var (
	pid           int
	defaultLogger *logger
	loggers       map[string]*logger
	loggersMutex  *sync.RWMutex
)

type LoggerHandler struct {
	loggers map[string]*logger
}

func (l *LoggerHandler) logBase(logLevel LogLevel, message string) {
	logInfo := &logInfo{
		pid:      pid,
		time:     time.Now(),
		logLevel: logLevel,
		message:  message,
	}
	pc, fileName, lineNum, ok := runtime.Caller(2)
	if ok {
		logInfo.pc = pc
		logInfo.fileName = fileName
		logInfo.lineNum = lineNum
	}
	for name, logger := range l.loggers {
		logger.log(name, logInfo)
	}
}

func (l *LoggerHandler) Emerg(format string, args ...interface{}) {
	l.logBase(LogLevelEmerg, fmt.Sprintf(format, args...))
}

func (l *LoggerHandler) Alert(format string, args ...interface{}) {
	l.logBase(LogLevelAlert, fmt.Sprintf(format, args...))
}

func (l *LoggerHandler) Crit(format string, args ...interface{}) {
	l.logBase(LogLevelCrit, fmt.Sprintf(format, args...))
}

func (l *LoggerHandler) Error(format string, args ...interface{}) {
	l.logBase(LogLevelError, fmt.Sprintf(format, args...))
}

func (l *LoggerHandler) Warn(format string, args ...interface{}) {
	l.logBase(LogLevelWarn, fmt.Sprintf(format, args...))
}

func (l *LoggerHandler) Notice(format string, args ...interface{}) {
	l.logBase(LogLevelNotice, fmt.Sprintf(format, args...))
}

func (l *LoggerHandler) Info(format string, args ...interface{}) {
	l.logBase(LogLevelInfo, fmt.Sprintf(format, args...))
}

func (l *LoggerHandler) Debug(format string, args ...interface{}) {
	l.logBase(LogLevelDebug, fmt.Sprintf(format, args...))
}

func (l *LoggerHandler) Trace(format string, args ...interface{}) {
	l.logBase(LogLevelTrace, fmt.Sprintf(format, args...))
}

func (l *LoggerHandler) Flush() {
	for _, logger := range l.loggers {
		logger.flush()
	}
}

func (l *LoggerHandler) ChangeFilter(name string, filter Filter) (err error) {
	logger, ok := l.loggers[name]
	if !ok {
		return errors.Errorf("not found name")
	}
	return logger.changeFilter(filter)
}

func (l *LoggerHandler) ChangeFormatter(name string, formatter Formatter) (err error) {
	logger, ok := l.loggers[name]
	if !ok {
		return errors.Errorf("not found name")
	}
	return logger.changeFormatter(formatter)
}

func (l *LoggerHandler) ChangeHandlers(name string, handlers []Handler) (err error) {
	logger, ok := l.loggers[name]
	if !ok {
		return errors.Errorf("not found name")
	}
	return logger.changeHandlers(handlers)
}

func GetLogger(names ...string) (loggerHandle *LoggerHandler) {
	loggersMutex.RLock()
	defer loggersMutex.RUnlock()
	loggerHandler := &LoggerHandler{
		loggers: make(map[string]*logger),
	}
	for _, name := range names {
		logger, ok := loggers[name]
		if !ok {
			loggerHandler.loggers[name] = defaultLogger
		} else {
			loggerHandler.loggers[name] = logger
		}
	}
	return loggerHandler
}

// set logger
func SetLogger(name string, filter Filter, formatter Formatter, handlers []Handler) (err error) {
	if name == "" || filter == nil || formatter == nil || handlers == nil || len(handlers) == 0 {
		return errors.Errorf("invalid argument")
	}
	loggersMutex.Lock()
	defer loggersMutex.Unlock()
	if _, ok := loggers[name]; ok {
		return errors.Errorf("already esixts logger")
	}
	loggers[name] = &logger{
		filter:    filter,
		formatter: formatter,
		handlers:  handlers,
	}
	for _, handler := range handlers {
		handler.Open()
	}
	return nil
}

func logBase(logLevel LogLevel, message string) {
	logInfo := &logInfo{
		pid:      pid,
		time:     time.Now(),
		logLevel: logLevel,
		message:  message,
	}
	pc, fileName, lineNum, ok := runtime.Caller(2)
	if ok {
		logInfo.pc = pc
		logInfo.fileName = fileName
		logInfo.lineNum = lineNum
	}
	defaultLogger.log("", logInfo)
}

func Emerg(format string, args ...interface{}) {
	logBase(LogLevelEmerg, fmt.Sprintf(format, args...))
}

func Alert(format string, args ...interface{}) {
	logBase(LogLevelAlert, fmt.Sprintf(format, args...))
}

func Crit(format string, args ...interface{}) {
	logBase(LogLevelCrit, fmt.Sprintf(format, args...))
}

func Error(format string, args ...interface{}) {
	logBase(LogLevelError, fmt.Sprintf(format, args...))
}

func Warn(format string, args ...interface{}) {
	logBase(LogLevelWarn, fmt.Sprintf(format, args...))
}

func Notice(format string, args ...interface{}) {
	logBase(LogLevelNotice, fmt.Sprintf(format, args...))
}

func Info(format string, args ...interface{}) {
	logBase(LogLevelInfo, fmt.Sprintf(format, args...))
}

func Debug(format string, args ...interface{}) {
	logBase(LogLevelDebug, fmt.Sprintf(format, args...))
}

func Trace(format string, args ...interface{}) {
	logBase(LogLevelTrace, fmt.Sprintf(format, args...))
}

func Flush() {
	defaultLogger.flush()
}

// change default logger filter
func ChangeFilter(filter Filter) (err error) {
	return defaultLogger.changeFilter(filter)
}

// change default logger formatter
func ChangeFormatter(formatter Formatter) (err error) {
	return defaultLogger.changeFormatter(formatter)
}

// change default logger handler
func ChangeHandlers(handlers []Handler) (err error) {
	return defaultLogger.changeHandlers(handlers)
}

type logger struct {
	filter    Filter
	formatter Formatter
	handlers  []Handler
	mutex     *sync.RWMutex
}

func (l *logger) log(loggerName string, logEvent LogEvent) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	if ok := l.filter.Evaluate(loggerName, logEvent); !ok {
		return
	}
	formattedLog := l.formatter.Format(loggerName, logEvent)
	for _, handler := range l.handlers {
		handler.Write(loggerName, logEvent, formattedLog)
	}
}

func (l *logger) flush() {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	for _, handler := range l.handlers {
		handler.Flush()
	}
}

func (l *logger) changeFilter(filter Filter) (err error) {
	if filter == nil {
		return errors.Errorf("invalid argument")
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.filter = filter
	return nil
}

func (l *logger) changeFormatter(formatter Formatter) (err error) {
	if formatter == nil {
		return errors.Errorf("invalid argument")
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.formatter = formatter
	return nil
}

func (l *logger) changeHandlers(handlers []Handler) (err error) {
	if handlers == nil || len(handlers) == 0 {
		return errors.Errorf("invalid argument")
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	for _, handler := range l.handlers {
		handler.Close()
	}
	l.handlers = handlers
	for _, handler := range l.handlers {
		handler.Open()
	}
	return nil
}

func init() {
	pid = os.Getpid()
	loggers = make(map[string]*logger)
	loggersMutex = new(sync.RWMutex)
	h := NewConsoleHandler()
	defaultLogger = &logger{
		filter:    NewLogLevelFilter(),
		formatter: NewStandardFormatter(),
		handlers:  []Handler{h},
		mutex:     new(sync.RWMutex),
	}
	h.Open()
}
