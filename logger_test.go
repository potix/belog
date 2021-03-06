package belog

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestDefaultLoggerChange1(t *testing.T) {
	Emerg("test\n")
	Alert("test\n")
	Crit("test\n")
	Error("test\n")
	Warn("test\n")
	Notice("test\n")
	Info("test\n")
	Debug("test\n")
	Trace("test\n")
	Flush()
	filter := NewLogLevelFilter()
	filter.SetLogLevel(LogLevelTrace)
	formatter := NewStandardFormatter()
	formatter.SetAppendNewLine(false)
	formatter.SetDateTimeLayout("2006-01-02 15:04:05 -0700 MST")
	formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(hostname) %(program) %(programCounter) %(loggerName) %(fileName) %(shortFileName) %(lineNum) %(message)")
	handler1 := NewConsoleHandler()
	handler1.SetOutputType(ConsoleOutputTypeStderr)
	handler1.SetConsoleColor(LogLevelEmerg, ConsoleColorBlack)
	handler2 := NewRotationFileHandler()
	handler2.SetLogFileName("belog-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(true)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	handler3 := NewSyslogHandler()
	handler3.SetNetworkAndAddr("", "")
	handler3.SetTag("belog-test")
	handler3.SetFacility("LOCAL7")
	if err := ChangeFilter(filter); err != nil {
		t.Errorf("%+v", err)
	}
	if err := ChangeFormatter(formatter); err != nil {
		t.Errorf("%+v", err)
	}
	handlers := make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
	handlers = append(handlers, handler3)
	if err := ChangeHandlers(handlers); err != nil {
		t.Errorf("%+v", err)
	}
	Emerg("test\n")
	Alert("test\n")
	Crit("test\n")
	Error("test\n")
	Warn("test\n")
	Notice("test\n")
	Info("test\n")
	Debug("test\n")
	Trace("test\n")
	Flush()
}

func TestDefaultLoggerChange2(t *testing.T) {
	Emerg("test\n")
	Alert("test\n")
	Crit("test\n")
	Error("test\n")
	Warn("test\n")
	Notice("test\n")
	Info("test\n")
	Debug("test\n")
	Trace("test\n")
	Flush()
	filter := NewLogLevelFilter()
	filter.SetLogLevel(LogLevelTrace)
	formatter := NewJSONFormatter()
	formatter.SetDateTimeLayout("2006-01-02 15:04:05 -0700 MST")
	handler1 := NewConsoleHandler()
	handler1.SetOutputType(ConsoleOutputTypeStderr)
	handler1.SetConsoleColor(LogLevelEmerg, ConsoleColorBlack)
	handler2 := NewRotationFileHandler()
	handler2.SetLogFileName("belog-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(true)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	handler3 := NewSyslogHandler()
	handler3.SetNetworkAndAddr("", "")
	handler3.SetTag("belog-test")
	handler3.SetFacility("LOCAL7")
	if err := ChangeFilter(filter); err != nil {
		t.Errorf("%+v", err)
	}
	if err := ChangeFormatter(formatter); err != nil {
		t.Errorf("%+v", err)
	}
	handlers := make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
	handlers = append(handlers, handler3)
	if err := ChangeHandlers(handlers); err != nil {
		t.Errorf("%+v", err)
	}
	Emerg("test\n")
	Alert("test\n")
	Crit("test\n")
	Error("test\n")
	Warn("test\n")
	Notice("test\n")
	Info("test\n")
	Debug("test\n")
	Trace("test\n")
	Flush()
}

func TestDefaultLoggerAsync(t *testing.T) {
	os.RemoveAll("/var/tmp/belog-test")
	filter := NewLogLevelFilter()
	filter.SetLogLevel(LogLevelTrace)
	formatter := NewStandardFormatter()
	formatter.SetAppendNewLine(false)
	formatter.SetDateTimeLayout("2006-01-02 15:04:05 -0700 MST")
	formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(hostname) %(program) %(programCounter) %(loggerName) %(fileName) %(shortFileName) %(lineNum) %(message)")
	handler1 := NewConsoleHandler()
	handler1.SetOutputType(ConsoleOutputTypeStderr)
	handler1.SetConsoleColor(LogLevelEmerg, ConsoleColorBlack)
	handler2 := NewRotationFileHandler()
	handler2.SetLogFileName("belog-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(true)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	handler3 := NewSyslogHandler()
	handler3.SetNetworkAndAddr("", "")
	handler3.SetTag("belog-test")
	handler3.SetFacility("LOCAL7")
	if err := ChangeFilter(filter); err != nil {
		t.Errorf("%+v", err)
	}
	if err := ChangeFormatter(formatter); err != nil {
		t.Errorf("%+v", err)
	}
	handlers := make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
	handlers = append(handlers, handler3)
	if err := ChangeHandlers(handlers); err != nil {
		t.Errorf("%+v", err)
	}
	Emerg("test\n")
	Alert("test\n")
	Crit("test\n")
	Error("test\n")
	Warn("test\n")
	Notice("test\n")
	Info("test\n")
	Debug("test\n")
	Trace("test\n")
	// read
	b, err := ioutil.ReadFile("/var/tmp/belog-test/belog-test.log")
	if err != nil {
		t.Errorf("%+v", err)
	}
	if len(b) != 0 {
		t.Errorf("async error")
	}
	time.Sleep(4 * time.Second)
	b, err = ioutil.ReadFile("/var/tmp/belog-test/belog-test.log")
	if err != nil {
		t.Errorf("%+v", err)
	}
	if len(b) == 0 {
		t.Errorf("flush error")
	}
}

func TestDefaultLoggerFlush(t *testing.T) {
	os.RemoveAll("/var/tmp/belog-test")
	filter := NewLogLevelFilter()
	filter.SetLogLevel(LogLevelTrace)
	formatter := NewStandardFormatter()
	formatter.SetDateTimeLayout("2006-01-02 15:04:05 -0700 MST")
	formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(hostname) %(program) %(programCounter) %(loggerName) %(fileName) %(shortFileName) %(lineNum) %(message)")
	handler1 := NewConsoleHandler()
	handler1.SetOutputType(ConsoleOutputTypeStderr)
	handler1.SetConsoleColor(LogLevelEmerg, ConsoleColorBlack)
	handler2 := NewRotationFileHandler()
	handler2.SetLogFileName("belog-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(true)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	handler3 := NewSyslogHandler()
	handler3.SetNetworkAndAddr("", "")
	handler3.SetTag("belog-test")
	handler3.SetFacility("LOCAL7")
	if err := ChangeFilter(filter); err != nil {
		t.Errorf("%+v", err)
	}
	if err := ChangeFormatter(formatter); err != nil {
		t.Errorf("%+v", err)
	}
	handlers := make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
	handlers = append(handlers, handler3)
	if err := ChangeHandlers(handlers); err != nil {
		t.Errorf("%+v", err)
	}
	Emerg("test\n")
	Alert("test\n")
	Crit("test\n")
	Error("test\n")
	Warn("test\n")
	Notice("test\n")
	Info("test\n")
	Debug("test\n")
	Trace("test\n")
	Flush()
	// read
	b, err := ioutil.ReadFile("/var/tmp/belog-test/belog-test.log")
	if err != nil {
		t.Errorf("%+v", err)
	}
	if len(b) == 0 {
		t.Errorf("flush error")
	}
}

func TestDefaultLoggerContentTrace(t *testing.T) {
	os.RemoveAll("/var/tmp/belog-test")
	filter := NewLogLevelFilter()
	filter.SetLogLevel(LogLevelTrace)
	formatter := NewStandardFormatter()
	formatter.SetDateTimeLayout("datetime")
	formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] %(loggerName) %(shortFileName) %(message)")
	handler1 := NewConsoleHandler()
	handler1.SetOutputType(ConsoleOutputTypeStderr)
	handler1.SetConsoleColor(LogLevelEmerg, ConsoleColorBlack)
	handler2 := NewRotationFileHandler()
	handler2.SetLogFileName("belog-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(false)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	handler3 := NewSyslogHandler()
	handler3.SetNetworkAndAddr("", "")
	handler3.SetTag("belog-test")
	handler3.SetFacility("LOCAL7")
	if err := ChangeFilter(filter); err != nil {
		t.Errorf("%+v", err)
	}
	if err := ChangeFormatter(formatter); err != nil {
		t.Errorf("%+v", err)
	}
	handlers := make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
	handlers = append(handlers, handler3)
	if err := ChangeHandlers(handlers); err != nil {
		t.Errorf("%+v", err)
	}
	Emerg("test\n")
	Alert("test\n")
	Crit("test\n")
	Error("test\n")
	Warn("test\n")
	Notice("test\n")
	Info("test\n")
	Debug("test\n")
	Trace("test\n")
	b, err := ioutil.ReadFile("/var/tmp/belog-test/belog-test.log")
	if err != nil {
		t.Errorf("%+v", err)
	}
	if len(b) == 0 {
		t.Errorf("flush error")
	}
	// write
	exp := `datetime [EMERG:1] default logger_test.go test
datetime [ALERT:2] default logger_test.go test
datetime [CRIT:3] default logger_test.go test
datetime [ERROR:4] default logger_test.go test
datetime [WARN:5] default logger_test.go test
datetime [NOTICE:6] default logger_test.go test
datetime [INFO:7] default logger_test.go test
datetime [DEBUG:8] default logger_test.go test
datetime [TRACE:9] default logger_test.go test
`
	if exp != string(b) {
		t.Errorf("mismatch log (exp %v != act %v)", exp, string(b))
	}
}

func TestDefaultLoggerContentNotice(t *testing.T) {
	os.RemoveAll("/var/tmp/belog-test")
	filter := NewLogLevelFilter()
	filter.SetLogLevel(LogLevelNotice)
	formatter := NewStandardFormatter()
	formatter.SetDateTimeLayout("datetime")
	formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] %(loggerName) %(shortFileName) %(message)")
	handler1 := NewConsoleHandler()
	handler1.SetOutputType(ConsoleOutputTypeStderr)
	handler1.SetConsoleColor(LogLevelEmerg, ConsoleColorBlack)
	handler2 := NewRotationFileHandler()
	handler2.SetLogFileName("belog-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(false)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	handler3 := NewSyslogHandler()
	handler3.SetNetworkAndAddr("", "")
	handler3.SetTag("belog-test")
	handler3.SetFacility("LOCAL7")
	if err := ChangeFilter(filter); err != nil {
		t.Errorf("%+v", err)
	}
	if err := ChangeFormatter(formatter); err != nil {
		t.Errorf("%+v", err)
	}
	handlers := make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
	handlers = append(handlers, handler3)
	if err := ChangeHandlers(handlers); err != nil {
		t.Errorf("%+v", err)
	}
	Emerg("test\n")
	Alert("test\n")
	Crit("test\n")
	Error("test\n")
	Warn("test\n")
	Notice("test\n")
	Info("test\n")
	Debug("test\n")
	Trace("test\n")
	b, err := ioutil.ReadFile("/var/tmp/belog-test/belog-test.log")
	if err != nil {
		t.Errorf("%+v", err)
	}
	if len(b) == 0 {
		t.Errorf("flush error")
	}
	// write
	exp := `datetime [EMERG:1] default logger_test.go test
datetime [ALERT:2] default logger_test.go test
datetime [CRIT:3] default logger_test.go test
datetime [ERROR:4] default logger_test.go test
datetime [WARN:5] default logger_test.go test
datetime [NOTICE:6] default logger_test.go test
`
	if exp != string(b) {
		t.Errorf("mismatch log (exp %v != act %v)", exp, string(b))
	}
}

func TestSetLoggerGetLoogerGroup(t *testing.T) {
	os.RemoveAll("/var/tmp/belog-test")
	filter := NewLogLevelFilter()
	filter.SetLogLevel(LogLevelTrace)
	formatter := NewStandardFormatter()
	formatter.SetDateTimeLayout("2006-01-02 15:04:05 -0700 MST")
	formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(hostname) %(program) %(programCounter) %(loggerName) %(fileName) %(shortFileName) %(lineNum) %(message)")
	handler1 := NewConsoleHandler()
	handler1.SetOutputType(ConsoleOutputTypeStderr)
	handler1.SetConsoleColor(LogLevelEmerg, ConsoleColorBlack)
	handler2 := NewRotationFileHandler()
	handler2.SetLogFileName("belog-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(true)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	handler3 := NewSyslogHandler()
	handler3.SetNetworkAndAddr("", "")
	handler3.SetTag("belog-test")
	handler3.SetFacility("LOCAL7")
	handlers := make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
	handlers = append(handlers, handler3)
	if err := SetLogger("logger1", filter, formatter, handlers); err != nil {
		t.Errorf("%+v", err)
	}
	if err := SetLogger("logger2", filter, formatter, handlers); err != nil {
		t.Errorf("%+v", err)
	}
	loggerGroup := GetLoggerGroup("logger1", "logger2")
	loggerGroup.Emerg("test\n")
	loggerGroup.Alert("test\n")
	loggerGroup.Crit("test\n")
	loggerGroup.Error("test\n")
	loggerGroup.Warn("test\n")
	loggerGroup.Notice("test\n")
	loggerGroup.Info("test\n")
	loggerGroup.Debug("test\n")
	loggerGroup.Trace("test\n")
	loggerGroup.Flush()
}

func TestChangeLoggerGroupByLoggerName(t *testing.T) {
	os.RemoveAll("/var/tmp/belog-test")
	filter := NewLogLevelFilter()
	filter.SetLogLevel(LogLevelTrace)
	formatter := NewStandardFormatter()
	formatter.SetDateTimeLayout("2006-01-02 15:04:05 -0700 MST")
	formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(hostname) %(program) %(programCounter) %(loggerName) %(fileName) %(shortFileName) %(lineNum) %(message)")
	handler1 := NewConsoleHandler()
	handler1.SetOutputType(ConsoleOutputTypeStderr)
	handler1.SetConsoleColor(LogLevelEmerg, ConsoleColorBlack)
	handler2 := NewRotationFileHandler()
	handler2.SetLogFileName("belog-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(true)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	handler3 := NewSyslogHandler()
	handler3.SetNetworkAndAddr("", "")
	handler3.SetTag("belog-test")
	handler3.SetFacility("LOCAL7")
	handlers := make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
	handlers = append(handlers, handler3)
	SetLogger("logger1", filter, formatter, handlers)
	SetLogger("logger2", filter, formatter, handlers)
	loggerGroup := GetLoggerGroup("logger1", "logger2")
	filter = NewLogLevelFilter()
	filter.SetLogLevel(LogLevelInfo)
	if err := loggerGroup.ChangeFilterByLoggerName("logger1", filter); err != nil {
		t.Errorf("%+v", err)
	}
	formatter = NewStandardFormatter()
	formatter.SetDateTimeLayout("2006-01-02 15:04:05")
	formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(hostname) %(program) %(loggerName) %(shortFileName) %(shortFileName) %(lineNum) %(message)")
	if err := loggerGroup.ChangeFormatterByLoggerName("logger1",formatter); err != nil {
		t.Errorf("%+v", err)
	}
	handler1 = NewConsoleHandler()
	handler1.SetOutputType(ConsoleOutputTypeStderr)
	handler1.SetConsoleColor(LogLevelEmerg, ConsoleColorBlack)
	handler2 = NewRotationFileHandler()
	handler2.SetLogFileName("belog-test-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(false)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	handler3 = NewSyslogHandler()
	handler3.SetNetworkAndAddr("", "")
	handler3.SetTag("belog-test")
	handler3.SetFacility("LOCAL7")
	handlers = make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
	handlers = append(handlers, handler3)
	if err := loggerGroup.ChangeHandlersByLoggerName("logger1", handlers); err != nil {
		t.Errorf("%+v", err)
	}
	loggerGroup.Emerg("test\n")
	loggerGroup.Alert("test\n")
	loggerGroup.Crit("test\n")
	loggerGroup.Error("test\n")
	loggerGroup.Warn("test\n")
	loggerGroup.Notice("test\n")
	loggerGroup.Info("test\n")
	loggerGroup.Debug("test\n")
	loggerGroup.Trace("test\n")
	loggerGroup.Flush()
}

func TestChangeLoggerGroup(t *testing.T) {
	os.RemoveAll("/var/tmp/belog-test")
	filter := NewLogLevelFilter()
	filter.SetLogLevel(LogLevelTrace)
	formatter := NewStandardFormatter()
	formatter.SetDateTimeLayout("2006-01-02 15:04:05 -0700 MST")
	formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(hostname) %(program) %(programCounter) %(loggerName) %(fileName) %(shortFileName) %(lineNum) %(message)")
	handler1 := NewConsoleHandler()
	handler1.SetOutputType(ConsoleOutputTypeStderr)
	handler1.SetConsoleColor(LogLevelEmerg, ConsoleColorBlack)
	handler2 := NewRotationFileHandler()
	handler2.SetLogFileName("belog-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(true)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	handler3 := NewSyslogHandler()
	handler3.SetNetworkAndAddr("", "")
	handler3.SetTag("belog-test")
	handler3.SetFacility("LOCAL7")
	handlers := make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
	handlers = append(handlers, handler3)
	SetLogger("logger1", filter, formatter, handlers)
	SetLogger("logger2", filter, formatter, handlers)
	loggerGroup := GetLoggerGroup("logger1", "logger2")
	filter = NewLogLevelFilter()
	filter.SetLogLevel(LogLevelInfo)
	if err := loggerGroup.ChangeFilter(filter); err != nil {
		t.Errorf("%+v", err)
	}
	formatter = NewStandardFormatter()
	formatter.SetDateTimeLayout("2006-01-02 15:04:05")
	formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(hostname) %(program) %(loggerName) %(shortFileName) %(shortFileName) %(lineNum) %(message)")
	if err := loggerGroup.ChangeFormatter(formatter); err != nil {
		t.Errorf("%+v", err)
	}
	handler1 = NewConsoleHandler()
	handler1.SetOutputType(ConsoleOutputTypeStderr)
	handler1.SetConsoleColor(LogLevelEmerg, ConsoleColorBlack)
	handler2 = NewRotationFileHandler()
	handler2.SetLogFileName("belog-test-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(false)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	handler3 = NewSyslogHandler()
	handler3.SetNetworkAndAddr("", "")
	handler3.SetTag("belog-test")
	handler3.SetFacility("LOCAL7")
	handlers = make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
	handlers = append(handlers, handler3)
	if err := loggerGroup.ChangeHandlers(handlers); err != nil {
		t.Errorf("%+v", err)
	}
	loggerGroup.Emerg("test\n")
	loggerGroup.Alert("test\n")
	loggerGroup.Crit("test\n")
	loggerGroup.Error("test\n")
	loggerGroup.Warn("test\n")
	loggerGroup.Notice("test\n")
	loggerGroup.Info("test\n")
	loggerGroup.Debug("test\n")
	loggerGroup.Trace("test\n")
	loggerGroup.Flush()
}
