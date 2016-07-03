# belog

logger package

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

## includes
* filter
  - LogLevelFilter: filter by log level
* formatter
  - StandardFormatter: standard formatter
* handler
  - consoleHadnler: output to console
  - syslogHadnler: output to syslog
  - RotationFileHandler: output to file and rotate file

## use default logger
* use default logger

```
        belog.Emerg("test\n")
        belog.Alert("test\n")
        belog.Crit("test\n")
        belog.Error("test\n")
        belog.Warn("test\n")
        belog.Notice("test\n")
        belog.Info("test\n")
        belog.Debug("test\n")
        belog.Trace("test\n")
```

## change filter

* change filter of default logger

```
        filter := belog.NewLogLevelFilter()
        filter.SetLogLevel(belog.LogLevelTrace)
        err := belog.ChangeFilter(filter)
        if err != nil {
                fmt.Println(err)
        }
```

## change formatter

* change formatter of default logger

```
        formatter := belog.NewStandardFormatter()
        formatter.SetDateTimeLayout("2006-01-02 15:04:05 -0700 MST")
        formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(programCounter) %(loggerName) %(fileName) %(shortFileName) %(lineNum) %(message)")
        err := belog.ChangeFormatter(formatter) 
        if err != nil {
               fmt.Println(err)
        }
```

## change handler

* change filter of default logger

```
        handler := belog.NewRotationFileHandler()
        handler.SetLogFileName("belog-test.log")
        handler.SetLogDirPath("/var/tmp/belog-test")
        handler.SetMaxAge(2)
        handler.SetMaxSize(65535)
        handler.SetAsync(true)
        handler.SetAsyncFlushInterval(3)
        handler.SetBufferSize(2048)
        handlers := make([]belog.Handler, 0)
        handlers = append(handlers, handler1)
        err := belog.ChangeHandlers(handlers)
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
        filter := belog.NewLogLevelFilter()
        filter.SetLogLevel(belog.LogLevelTrace)
        formatter := belog.NewStandardFormatter()
        formatter.SetDateTimeLayout("2006-01-02 15:04:05 -0700 MST")
        formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(programCounter) %(loggerName) %(fileName) %(shortFileName) %(lineNum) %(message)")
        handler := belog.NewRotationFileHandler()
        handler.SetLogFileName("belog-test.log")
        handler.SetLogDirPath("/var/tmp/belog-test")
        handler.SetMaxAge(10)
        handler.SetMaxSize(1024 * 1024 * 1024)
        handlers := make([]belog.Handler, 0)
        handlers = append(handlers, handler1)
        handlers = append(handlers, handler2)
        belog.SetLogger("mylogger1", filter, formatter, handlers)
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
        if err := belog.LoadConfig("sample.yaml"); err != nil {
               fmt.Println(err)
        }
}
```

## use logger (no default logger)
* get logger
  - can get mutiple logger

```
func init() {
	logger := belog.GetLogger("mylogger1", "mylogger2")
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
  - name is struct name of custom filter

```
func NewLogLevelFilter() (logLevelFilter *LogLevelFilter) {
        return &LogLevelFilter{
                logLevel: LogLevelInfo,
                mutex:    new(sync.RWMutex),
        }
}

func init() {
        belog.RegisterFilter("LogLevelFilter", func() (filter belog.Filter) {
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
  - name is struct name of custom formatter

```
func NewStandardFormatter() (standardFormatter *StandardFormatter) {
        return &StandardFormatter{
                dateTimeLayout: "2006-01-02 15:04:05",
                layout:         "%(dateTime) [%(logLevel)] (%(pid)) %(loggerName) %(fileName) %(lineNum) %(message)",
                mutex:          new(sync.RWMutex),
        }
}

func init() {
        belog.RegisterFormatter("StandardFormatter", func() (formatter belog.Formatter) {
                return NewStandardFormatter()
        })
}
```

## custom handler
*  create struct have handler interface
  - name is struct name of custom handler

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
        belog.RegisterHandler("ConsoleHandler", func() (handler belog.Handler) {
                return NewConsoleHandler()
        })
}
```
