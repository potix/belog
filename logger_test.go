package belog

import (
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestDefaultLoggerChange(t *testing.T) {
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
	formatter.SetDateTimeLayout("2006-01-02 15:04:05 -0700 MST")
	formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(programCounter) %(loggerName) %(fileName) %(shortFileName) %(lineNum) %(message)")
	handler1 := NewConsoleHandler()
	handler1.SetOutputType(OutputTypeStderr)
	handler2 := NewRotationFileHandler()
	handler2.SetLogFileName("belog-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(true)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	if err := ChangeFilter(filter); err != nil {
		t.Errorf("%+v", err)
	}
	if err := ChangeFormatter(formatter); err != nil {
		t.Errorf("%+v", err)
	}
	handlers := make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
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
	formatter.SetDateTimeLayout("2006-01-02 15:04:05 -0700 MST")
	formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(programCounter) %(loggerName) %(fileName) %(shortFileName) %(lineNum) %(message)")
	handler1 := NewConsoleHandler()
	handler1.SetOutputType(OutputTypeStderr)
	handler2 := NewRotationFileHandler()
	handler2.SetLogFileName("belog-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(true)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	if err := ChangeFilter(filter); err != nil {
		t.Errorf("%+v", err)
	}
	if err := ChangeFormatter(formatter); err != nil {
		t.Errorf("%+v", err)
	}
	handlers := make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
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
	formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(programCounter) %(loggerName) %(fileName) %(shortFileName) %(lineNum) %(message)")
	handler1 := NewConsoleHandler()
	handler1.SetOutputType(OutputTypeStderr)
	handler2 := NewRotationFileHandler()
	handler2.SetLogFileName("belog-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(true)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	if err := ChangeFilter(filter); err != nil {
		t.Errorf("%+v", err)
	}
	if err := ChangeFormatter(formatter); err != nil {
		t.Errorf("%+v", err)
	}
	handlers := make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
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
	handler1.SetOutputType(OutputTypeStderr)
	handler2 := NewRotationFileHandler()
	handler2.SetLogFileName("belog-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(false)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	if err := ChangeFilter(filter); err != nil {
		t.Errorf("%+v", err)
	}
	if err := ChangeFormatter(formatter); err != nil {
		t.Errorf("%+v", err)
	}
	handlers := make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
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
	exp := `datetime [EMERG:1]  logger_test.go test
datetime [ALERT:2]  logger_test.go test
datetime [CRIT:3]  logger_test.go test
datetime [ERROR:4]  logger_test.go test
datetime [WARN:5]  logger_test.go test
datetime [NOTICE:6]  logger_test.go test
datetime [INFO:7]  logger_test.go test
datetime [DEBUG:8]  logger_test.go test
datetime [TRACE:9]  logger_test.go test
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
	handler1.SetOutputType(OutputTypeStderr)
	handler2 := NewRotationFileHandler()
	handler2.SetLogFileName("belog-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(false)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	if err := ChangeFilter(filter); err != nil {
		r.Errorf("%+v", err)
	}
	if err := ChangeFormatter(formatter); err != nil {
		t.Errorf("%+v", err)
	}
	handlers := make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
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
	exp := `datetime [EMERG:1]  logger_test.go test
datetime [ALERT:2]  logger_test.go test
datetime [CRIT:3]  logger_test.go test
datetime [ERROR:4]  logger_test.go test
datetime [WARN:5]  logger_test.go test
datetime [NOTICE:6]  logger_test.go test
`
	if exp != string(b) {
		t.Errorf("mismatch log (exp %v != act %v)", exp, string(b))
	}
}

func TestSetLoggerGetLooger(t *testing.T) {
	os.RemoveAll("/var/tmp/belog-test")
	filter := NewLogLevelFilter()
	filter.SetLogLevel(LogLevelTrace)
	formatter := NewStandardFormatter()
	formatter.SetDateTimeLayout("2006-01-02 15:04:05 -0700 MST")
	formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(programCounter) %(loggerName) %(fileName) %(shortFileName) %(lineNum) %(message)")
	handler1 := NewConsoleHandler()
	handler1.SetOutputType(OutputTypeStderr)
	handler2 := NewRotationFileHandler()
	handler2.SetLogFileName("belog-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(true)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	handlers := make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
	if err := SetLogger("logger1", filter, formatter, handlers); err != nil {
		t.Errorf("%+v", err)
	}
	if err := SetLogger("logger2", filter, formatter, handlers); err != nil {
		t.Errorf("%+v", err)
	}
	logger := GetLogger("logger1", "logger2")
	logger.Emerg("test\n")
	logger.Alert("test\n")
	logger.Crit("test\n")
	logger.Error("test\n")
	logger.Warn("test\n")
	logger.Notice("test\n")
	logger.Info("test\n")
	logger.Debug("test\n")
	logger.Trace("test\n")
	logger.Flush()
}

func TestChangeLogger(t *testing.T) {
	os.RemoveAll("/var/tmp/belog-test")
	filter := NewLogLevelFilter()
	filter.SetLogLevel(LogLevelTrace)
	formatter := NewStandardFormatter()
	formatter.SetDateTimeLayout("2006-01-02 15:04:05 -0700 MST")
	formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(programCounter) %(loggerName) %(fileName) %(shortFileName) %(lineNum) %(message)")
	handler1 := NewConsoleHandler()
	handler1.SetOutputType(OutputTypeStderr)
	handler2 := NewRotationFileHandler()
	handler2.SetLogFileName("belog-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(true)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	handlers := make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
	SetLogger("logger1", filter, formatter, handlers)
	SetLogger("logger2", filter, formatter, handlers)
	logger := GetLogger("logger1", "logger2")
	filter = NewLogLevelFilter()
	filter.SetLogLevel(LogLevelInfo)
	if err := logger.ChangeFilter("logger1", filter); err != nil {
		t.Errorf("%+v", err)
	}
	formatter = NewStandardFormatter()
	formatter.SetDateTimeLayout("2006-01-02 15:04:05")
	formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(loggerName) %(shortFileName) %(shortFileName) %(lineNum) %(message)")
	if err := logger.ChangeFormatter("logger1", formatter); err != nil {
		t.Errorf("%+v", err)
	}
	handler1 = NewConsoleHandler()
	handler1.SetOutputType(OutputTypeStderr)
	handler2 = NewRotationFileHandler()
	handler2.SetLogFileName("belog-test-test.log")
	handler2.SetLogDirPath("/var/tmp/belog-test")
	handler2.SetMaxAge(2)
	handler2.SetMaxSize(65535)
	handler2.SetAsync(false)
	handler2.SetAsyncFlushInterval(3)
	handler2.SetBufferSize(2048)
	handlers = make([]Handler, 0, 0)
	handlers = append(handlers, handler1)
	handlers = append(handlers, handler2)
	if err := logger.ChangeHandlers("logger1", handlers); err != nil {
		t.Errorf("%+v", err)
	}
	logger.Emerg("test\n")
	logger.Alert("test\n")
	logger.Crit("test\n")
	logger.Error("test\n")
	logger.Warn("test\n")
	logger.Notice("test\n")
	logger.Info("test\n")
	logger.Debug("test\n")
	logger.Trace("test\n")
	logger.Flush()
}
