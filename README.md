# belog

logger package

## belog

- logging process
```
                                                                        ---
 logging   --------        -----------        ---------                  |
     ---> | filter1| ---> | formatter1| -+-> | handler1| ---> console    |
           --------        -----------   |    ---------                  |
                                         |    ---------                  |
                                         +-> | handler2| ---> file       |
                                              ---------                  |
       |----------------- logger -----------------------|                |
                                                                         | logger group
 logging   --------        -----------        ---------                  |
     ---> | filter2| ---> | formatter2| -+-> | handler1| ---> console    |
           --------        -----------   |    ---------                  |
                                         |    ---------                  |
                                         +-> | handler2| ---> file       |
                                              ---------                  |
       |----------------- logger -----------------------|                |
                                                                        ---
```

## included componets
* filter
  * LogLevelFilter
    - filter by log level.
* formatter
  * StandardFormatter
    - standard formatter
  * JSONFormatter
    - json formatter
* handler
  * ConsoleHadnler
    - output to console.
    - color is supported.
  * SyslogHadnler
    - output to syslog
  * RotationFileHandler
    - output to file.
    - rotation is supported.

## logging with default logger

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

## change filter of default logger

```
        filter := belog.NewLogLevelFilter()
        filter.SetLogLevel(belog.LogLevelTrace)
        err := belog.ChangeFilter(filter)
        if err != nil {
                fmt.Println(err)
        }
```

## change formatter of default logger

```
        formatter := belog.NewStandardFormatter()
        formatter.SetDateTimeLayout("2006-01-02 15:04:05 -0700 MST")
        formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(programCounter) %(loggerName) %(fileName) %(shortFileName) %(lineNum) %(message)")
        err := belog.ChangeFormatter(formatter) 
        if err != nil {
               fmt.Println(err)
        }
```

## change handlers of default logger

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
 
## setup custom loggers 

- It is requred, if you need multiple logger.

### setup Logger

- logger name 'default' is not used, because reserved by defaultLogger.

```
func init() {
	// create filter
        filter := belog.NewLogLevelFilter()
        filter.SetLogLevel(belog.LogLevelTrace)

	// create formatter
        formatter := belog.NewStandardFormatter()
        formatter.SetDateTimeLayout("2006-01-02 15:04:05 -0700 MST")
        formatter.SetLayout("%(dateTime) [%(logLevel):%(logLevelNum)] (%(pid)) %(programCounter) %(loggerName) %(fileName) %(shortFileName) %(lineNum) %(message)")

	// create handler
        handler := belog.NewRotationFileHandler()
        handler.SetLogFileName("belog-test.log")
        handler.SetLogDirPath("/var/tmp/belog-test")
        handler.SetMaxAge(10)
        handler.SetMaxSize(1024 * 1024 * 1024)

	// add handlers
        handlers := make([]belog.Handler, 0)
        handlers = append(handlers, handler1)
        handlers = append(handlers, handler2)

	// set logger
        belog.SetLogger("mylogger1", filter, formatter, handlers)
}
```

## get logger

- You can get mutiple logger.

```
func init() {
	logger := belog.GetLoggerGroup("mylogger1", "mylogger2")
	logger.Info("test")
}
```

### setup logger from config file

- Loadable config format are toml or yaml of json.
  - See test directory samples.
- logger name 'default' is not used, because reserved by defaultLogger.

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

### setup logger from ConfigLoggers object

```
type my struct {
	Alice *Person
        Logger *belog.ConfigLoggers
}

func init() {

	...

        if err := SetupLoggers(my.Logger); err != nil {
               fmt.Println(err)
        }
}
```

## create custom fileter

- Your filter struct have to method of filter interface.

```
type Filter interface {
        Evaluate(loggerName string, log LogEvent) bool
}
```

- You must create function, it must create filter struct and must return pointer of filter struct.

```
func NewLogLevelFilter() (logLevelFilter *LogLevelFilter) {
        return &LogLevelFilter{
                logLevel: LogLevelInfo,
                mutex:    new(sync.RWMutex),
        }
}
```

- Finally, you must register filter in init function.
  - First argument of RegisterFilter function is your filter name.
  - This process is requied,if you setup logger from config or from configLogger object. 

```
func init() {
        belog.RegisterFilter("LogLevelFilter", func() (filter belog.Filter) {
                return NewLogLevelFilter()
        })
}
```

## create custom formatter

- Your formatter struct have to method of formatter interface.

```
type Formatter interface {
        Format(loggerName string, log LogEvent) (logString string)
}
```

- You must create function, it must create formatter struct and must return pointer of formatter struct.

```
func NewStandardFormatter() (standardFormatter *StandardFormatter) {
        return &StandardFormatter{
                dateTimeLayout: "2006-01-02 15:04:05",
                layout:         "%(dateTime) [%(logLevel)] (%(pid)) %(loggerName) %(fileName) %(lineNum) %(message)",
                mutex:          new(sync.RWMutex),
        }
}
```

- Finally, you must register formatter in init function.
  - First argument of RegisterFormatter function is your formatter name.
  - This process is requied, if you setup logger from config or from configLogger object. 

```
func init() {
        belog.RegisterFormatter("StandardFormatter", func() (formatter belog.Formatter) {
                return NewStandardFormatter()
        })
}
```

## create custom handler

- Your handler struct have to method of handler interface.

```
type Handler interface {
        IsOpened()
        Open()
        Write(loggerName string, logEvent LogEvent, formattedLog string)
        Flush()
        Close()
}
```

- You must create function, it must create handler struct and must return pointer of handler struct.

```
func NewConsoleHandler() (consoleHandler *ConsoleHandler) {
        return &ConsoleHandler{
                outputType: OutputTypeStdout,
                mutex:      new(sync.RWMutex),
        }
}
```

- Finally, you must register handler in init function.
  - First argument of RegisterHandler function is your handler name.
  - This process is requied,if you setup logger from config or from configLogger object. 

```
func init() {
        belog.RegisterHandler("ConsoleHandler", func() (handler belog.Handler) {
                return NewConsoleHandler()
        })
}
```

## log event

* LogEvent interface

```
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
```
