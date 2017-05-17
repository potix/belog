package belog

import (
	"fmt"
	"github.com/pkg/errors"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

//LogLevel is log level
type LogLevel int

const (
	//LogLevelEmerg is log level of emergency
	LogLevelEmerg LogLevel = iota + 1
	//LogLevelAlert is log level of alert
	LogLevelAlert
	//LogLevelCrit is log level of critical
	LogLevelCrit
	//LogLevelError is log level of error
	LogLevelError
	//LogLevelWarn is log level of warning
	LogLevelWarn
	//LogLevelNotice is log level of notice
	LogLevelNotice
	//LogLevelInfo is log level of info
	LogLevelInfo
	//LogLevelDebug is log level of debug
	LogLevelDebug
	//LogLevelTrace is log level of trace
	LogLevelTrace
)

var (
	program       string
	pid           int
	hostname      string
	defaultLogger *logger
	loggers       map[string]*logger
	loggersMutex  *sync.RWMutex
)

//
// logger group wrapper
//

//LoggerGroup is logger group
type LoggerGroup struct {
	loggers map[string]*logger
}

func (l *LoggerGroup) logBase(logLevel LogLevel, message string) {
	logInfo := &logInfo{
		program:  program,
		pid:      pid,
		hostname: hostname,
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

//Emerg is output log of emergency level with logger group
func (l *LoggerGroup) Emerg(format string, args ...interface{}) {
	l.logBase(LogLevelEmerg, fmt.Sprintf(format, args...))
}

//Alert is output log of alert level with logger group
func (l *LoggerGroup) Alert(format string, args ...interface{}) {
	l.logBase(LogLevelAlert, fmt.Sprintf(format, args...))
}

//Crit is output log of critical level with logger group
func (l *LoggerGroup) Crit(format string, args ...interface{}) {
	l.logBase(LogLevelCrit, fmt.Sprintf(format, args...))
}

//Error is output log of error level with logger group
func (l *LoggerGroup) Error(format string, args ...interface{}) {
	l.logBase(LogLevelError, fmt.Sprintf(format, args...))
}

//Warn is output log of warn level with logger group
func (l *LoggerGroup) Warn(format string, args ...interface{}) {
	l.logBase(LogLevelWarn, fmt.Sprintf(format, args...))
}

//Notice is output log of notice level with logger group
func (l *LoggerGroup) Notice(format string, args ...interface{}) {
	l.logBase(LogLevelNotice, fmt.Sprintf(format, args...))
}

//Info is output log of info level with logger group
func (l *LoggerGroup) Info(format string, args ...interface{}) {
	l.logBase(LogLevelInfo, fmt.Sprintf(format, args...))
}

//Debug is output log of debug level with logger group
func (l *LoggerGroup) Debug(format string, args ...interface{}) {
	l.logBase(LogLevelDebug, fmt.Sprintf(format, args...))
}

//Trace is output log of trace level with logger group
func (l *LoggerGroup) Trace(format string, args ...interface{}) {
	l.logBase(LogLevelTrace, fmt.Sprintf(format, args...))
}

//Flush is flush log with logger group
func (l *LoggerGroup) Flush() {
	for _, logger := range l.loggers {
		logger.flush()
	}
}

//ChangeFilterByLoggerName is change fileter by logger name of logger group
func (l *LoggerGroup) ChangeFilterByLoggerName(name string, filter Filter) (error) {
	logger, ok := l.loggers[name]
	if !ok {
		return errors.Errorf("not found name")
	}
	return logger.changeFilter(filter)
}

//ChangeFormatterByLoggerName is change formatter by logger name of logger group
func (l *LoggerGroup) ChangeFormatterByLoggerName(name string, formatter Formatter) (error) {
	logger, ok := l.loggers[name]
	if !ok {
		return errors.Errorf("not found name")
	}
	return logger.changeFormatter(formatter)
}

//ChangeHandlersByLoggerName is change handlers by logger name of logger group
func (l *LoggerGroup) ChangeHandlersByLoggerName(name string, handlers []Handler) (error) {
	logger, ok := l.loggers[name]
	if !ok {
		return errors.Errorf("not found name")
	}
	return logger.changeHandlers(handlers)
}

//ChangeFilter is change fileter of logger group
func (l *LoggerGroup) ChangeFilter(filter Filter) (err error) {
	for _, logger := range l.loggers {
		err = logger.changeFilter(filter)
		if err != nil {
			return err
		}
	}
	return nil
}

//ChangeFormatter is change formatter of logger group
func (l *LoggerGroup) ChangeFormatter(formatter Formatter) (err error) {
	for _, logger := range l.loggers {
		err = logger.changeFormatter(formatter)
		if err != nil {
			return err
		}
	}
	return nil
}

//ChangeHandlers is change handlers of logger group
func (l *LoggerGroup) ChangeHandlers(handlers []Handler) (err error) {
	for _, logger := range l.loggers {
		err = logger.changeHandlers(handlers)
		if err != nil {
			return err
		}
	}
	return nil
}

//GetLoggerGroup is get logger group
func GetLoggerGroup(names ...string) (loggerGroup *LoggerGroup) {
	loggersMutex.RLock()
	defer loggersMutex.RUnlock()
	loggerGroup = &LoggerGroup{
		loggers: make(map[string]*logger),
	}
	for _, name := range names {
		logger, ok := loggers[name]
		if !ok {
			loggerGroup.loggers[name] = defaultLogger
		} else {
			loggerGroup.loggers[name] = logger
		}
	}
	return loggerGroup
}

//SetLogger is set logger
func SetLogger(name string, filter Filter, formatter Formatter, handlers []Handler) (err error) {
	if  filter == nil || formatter == nil || handlers == nil || len(handlers) == 0 {
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
		mutex:     new(sync.RWMutex),
	}
	for _, handler := range handlers {
		if !handler.IsOpened() {
			handler.Open()
		}
	}
	return nil
}

//
// default logger wrapper
//

func logBase(logLevel LogLevel, message string) {
	logInfo := &logInfo{
		program:  program,
		pid:      pid,
		hostname: hostname,
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
	defaultLogger.log("default", logInfo)
}

//Emerg is output log of emergency level with default logger
func Emerg(format string, args ...interface{}) {
	logBase(LogLevelEmerg, fmt.Sprintf(format, args...))
}

//Alert is output log of alert level with default logger
func Alert(format string, args ...interface{}) {
	logBase(LogLevelAlert, fmt.Sprintf(format, args...))
}

//Crit is output log of critical level with default logger
func Crit(format string, args ...interface{}) {
	logBase(LogLevelCrit, fmt.Sprintf(format, args...))
}

//Error is output log of error level with default logger
func Error(format string, args ...interface{}) {
	logBase(LogLevelError, fmt.Sprintf(format, args...))
}

//Warn is output log of warning level with default logger
func Warn(format string, args ...interface{}) {
	logBase(LogLevelWarn, fmt.Sprintf(format, args...))
}

//Notice is output log of notice level with default logger
func Notice(format string, args ...interface{}) {
	logBase(LogLevelNotice, fmt.Sprintf(format, args...))
}

//Info is output log of info level with default logger
func Info(format string, args ...interface{}) {
	logBase(LogLevelInfo, fmt.Sprintf(format, args...))
}

//Debug is output log of debug level with default logger
func Debug(format string, args ...interface{}) {
	logBase(LogLevelDebug, fmt.Sprintf(format, args...))
}

//Trace is output log of trace level with default logger
func Trace(format string, args ...interface{}) {
	logBase(LogLevelTrace, fmt.Sprintf(format, args...))
}

//Flush is flush log of default logger
func Flush() {
	defaultLogger.flush()
}

//ChangeFilter is change filter of default logger
func ChangeFilter(filter Filter) (err error) {
	return defaultLogger.changeFilter(filter)
}

//ChangeFormatter is change formatter of default logger
func ChangeFormatter(formatter Formatter) (err error) {
	return defaultLogger.changeFormatter(formatter)
}

//ChangeHandlers is change handler of default logger
func ChangeHandlers(handlers []Handler) (err error) {
	return defaultLogger.changeHandlers(handlers)
}

//
// logger
//

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
	formattedLog, err := l.formatter.Format(loggerName, logEvent)
	if err != nil {
		// statistics
		return
	}
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
		if handler.IsOpened() {
			handler.Close()
		}
	}
	l.handlers = handlers
	for _, handler := range l.handlers {
		if !handler.IsOpened() {
			handler.Open()
		}
	}
	return nil
}

func init() {
	program = filepath.Base(os.Args[0])
	pid = os.Getpid()
	name, err := os.Hostname()
	if err == nil {
		hostname = name
	}
	loggers = make(map[string]*logger)
	loggersMutex = new(sync.RWMutex)
	h := NewConsoleHandler()
	defaultLogger = &logger{
		filter:    NewLogLevelFilter(),
		formatter: NewStandardFormatter(),
		handlers:  []Handler{h},
		mutex:     new(sync.RWMutex),
	}
	if !h.IsOpened() {
		h.Open()
	}
}
