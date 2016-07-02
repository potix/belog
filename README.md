# belog

simple logger package

## belog

- logging process
```
 logging   --------        -----------        ---------
     ---> | filter | ---> | formatter | -+-> | handler | ---> console
           --------        -----------   |    ---------
                                         |    ---------
                                         +-> | handler | ---> file
                                              ---------
       |----------------- logger -----------------------|
```

## use default logger
* use default logger

```
	belog.Info("test")
```

## change filter

* change filter of default logger

```
        filter := NewLogLevelFilter()
        filter.SetLogLevel(LogLevelTrace)
        err := ChangeFilter(filter)
        if err != nil {
                fmt.Println(err)
        }
```

## change formatter

* change formatter of default logger

```
        formatter := NewStandardFormatter()
        formatter.SetDateTimeLayout("2006-01-02 15:04:05 -0700 MST")
        formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(programCounter) %(loggerName) %(fileName) %(shortFileName) %(lineNum) %(message)")
        err := ChangeFormatter(formatter) 
        if err != nil {
               fmt.Println(err)
        }
```

## change handler

* change filter of default logger

```
        handler := NewRotationFileHandler()
        handler.SetLogFileName("belog-test.log")
        handler.SetLogDirPath("/var/tmp/belog-test")
        handler.SetMaxAge(2)
        handler.SetMaxSize(65535)
        handler.SetAsync(true)
        handler.SetAsyncFlushInterval(3)
        handler.SetBufferSize(2048)
        handlers := make([]Handler, 0)
        handlers = append(handlers, handler1)
        err := ChangeHandlers(handlers)
        if err != nil {
                fmt.Println(err)
        }
```
 
## setup loggers (no default logger)

* setup loggers
  - if more than one of the logger is required

### use SetLogger
* call new function and call setter and call SetLogger 

```
func init() {
        filter := NewLogLevelFilter()
        filter.SetLogLevel(LogLevelTrace)
        formatter := NewStandardFormatter()
        formatter.SetDateTimeLayout("2006-01-02 15:04:05 -0700 MST")
        formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(programCounter) %(loggerName) %(fileName) %(shortFileName) %(lineNum) %(message)")
        handler := NewRotationFileHandler()
        handler.SetLogFileName("belog-test.log")
        handler.SetLogDirPath("/var/tmp/belog-test")
        handler.SetMaxAge(10)
        handler.SetMaxSize(1024 * 1024 * 1024)
        handlers := make([]Handler, 0)
        handlers = append(handlers, handler1)
        handlers = append(handlers, handler2)
        SetLogger("mylogger1", filter, formatter, handlers)
}
```

### use LoadConfig
* create config and call LoadConfig
  - yaml or json or toml
  - see test directory sample

```
--- sample.yaml ---
loggers:
  mylogger:
    filter:
      structname: LogLevelFilter
      structsetters: []
    formatter:
      structname: StandardFormatter
      structsetters:
      - settername: SetDateTimeLayout
        setterparams:
        - 2006-01-02 15:04:05 -0700 MST
      - settername: SetLayout
        setterparams:
        - '%(dateTime) [%(logLevel)] (%(pid)) %(programCounter) %(loggerName) %(fileName) %(lineNum) %(message)'
    handlers:
    - structname: ConsoleHandler
      structsetters: []
    - structname: RotationFileHandler
      structsetters: []
```

```
func init() {
        if err := LoadConfig("sample.yaml"); err != nil {
               fmt.Println(err)
        }
}
```

## use logger 
* get logger
  - can able to get mutiple logger

```
func init() {
	logger := GetLogger("mylogger1", "mylogger2")
	logger.Info("test")
}
```

## custom fileter
*  create stuct have filter interface

```
type Filter interface {
        Evaluate(loggerName string, log LogEvent) bool
}
```

* create new function and register filter

```
func NewLogLevelFilter() (logLevelFilter *LogLevelFilter) {
        return &LogLevelFilter{
                logLevel: LogLevelInfo,
                mutex:    new(sync.RWMutex),
        }
}

func init() {
        RegisterFilter("LogLevelFilter", func() (filter Filter) {
                return NewLogLevelFilter()
        })
}
```

## custom formatter
* create struct have formatter interface

```
type Formatter interface {
        Format(loggerName string, log LogEvent) (logString string)
}
```

* create new function and register formatter

```
func NewStandardFormatter() (standardFormatter *StandardFormatter) {
        return &StandardFormatter{
                dateTimeLayout: "2006-01-02 15:04:05",
                layout:         "%(dateTime) [%(logLevel)] (%(pid)) %(loggerName) %(fileName) %(lineNum) %(message)",
                mutex:          new(sync.RWMutex),
        }
}

func init() {
        RegisterFormatter("StandardFormatter", func() (formatter Formatter) {
                return NewStandardFormatter()
        })
}
```

## custom handler
*  create struct have handler interface

```
type Handler interface {
        Open()
        Write(loggerName string, logEvent LogEvent, formattedLog string)
        Flush()
        Close()
}
```

* create new function and register handler

```
func NewConsoleHandler() (consoleHandler *ConsoleHandler) {
        return &ConsoleHandler{
                outputType: OutputTypeStdout,
                mutex:      new(sync.RWMutex),
        }
}

func init() {
        RegisterHandler("ConsoleHandler", func() (handler Handler) {
                return NewConsoleHandler()
        })
}
```
