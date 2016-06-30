package belog

import (
	"log"
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
	LogLevelMap = map[LogLevel]string{
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
	pid           int
	defaultLogger *logger
	loggers       map[string]*logger
	loggersMutex  *sync.RWMutex
)

type LoggerHandler struct {
	loggers map[name]*logger
}

func (l *LoggerHandler) logBase(logLevel LogLevel, message string) {
	log := &log{
		pid:      pid,
		time:     time.Now(),
		logLevel: logLevel,
		message:  message,
	}
	pc, fileName, lineNum, ok := runtime.Caller(2)
	if ok {
		log.pc = pc
		log.fileName = fileNamw
		log.lineNum = lineNum
	}
	for name, logger := range l.loggers {
		logger.log(name, log)
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
	return logger.changeFormatter(filter)
}

func (l *LoggerHandler) ChangeHandlers(name string, handlers []handlers) (err error) {
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
		loggers: make(map[string]lLogger),
	}
	for name := range names {
		logger, ok := loggers[name]
		if !ok {
			loggerHandler.loggers[name] = defaultLogger
		} else {
			loggerHandler.loggers[name] = loggers[name]
		}
	}
	return loggerHandler
}

// setup logger
func SetLogger(name string, filter filter.Filter, formatter formatter.Formatter, handlers []handler.Handler) (err error) {
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
		formatter: formater,
		handlers:  handlers,
	}
	for handler := range handlers {
		handler.Open()
	}
}

func logBase(logLevel LogLevel, message string) {
	log := &log{
		pid:      pid,
		time:     time.Now(),
		logLevel: logLevel,
		message:  message,
	}
	pc, fileName, lineNum, ok := runtime.Caller(2)
	if ok {
		log.pc = pc
		log.fileName = fileNamw
		log.lineNum = lineNum
	}
	defaultLogger.log("", log)
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

func (l *Logger) log(name string, log LogEvent) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	if ok := logger.filter.Evalute(log); !ok {
		return
	}
	logString := logger.formatter.format(name, log)
	for handler := range handlers {
		handler.Write(logString)
	}
}

func (l *Logger) changeFilter(filter Filter) (err error) {
	if filter == nil {
		return errors.Errorf("invalid argument")
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.filter = filter
}

func (l *Logger) changeFormatter(formatter Formatter) (err error) {
	if formatter == nil {
		return errors.Errorf("invalid argument")
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	l.formatter = formatter

}

func (l *Logger) changeHandlers(handlers []Handler) (err error) {
	if handlers == nil || len(handlers) == 0 {
		return errors.Errorf("invalid argument")
	}
	l.mutex.Lock()
	defer l.mutex.Unlock()
	for handler := range l.handlers {
		handler.Close()
	}
	l.handlers = handlers
	for handler := range l.handlers {
		handler.Open()
	}
}

func init() {
	pid = os.Getpid()
	loggers = make(map[string]*Logger)
	loggersMutex = new(sync.RWMutex)
	h := handler.NewConsoleHandler()
	defaultLogger = &logger{
		filter:    filter.NewLogLevelFilter(),
		formatter: formatter.NewStandardFormatter(),
		handlers:  []Handler{h},
	}
	h.Open()
}
