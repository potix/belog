package belog

import (
	"encoding/json"
	"sync"
)

type jsonLogInfo struct {
	LoggerName string
	Program    string
	Pid        int
	Hostname   string
	Time       string
	LogLevel   string
	Pc         uintptr
	FileName   string
	LineNum    int
	Message    string
	Attrs      map[string]interface{}
}

//JsonFormatter is format json string
type JsonFormatter struct {
	dateTimeLayout string
	mutex          *sync.RWMutex
}

//Format is format log event to json string
func (f *JsonFormatter) Format(loggerName string, log LogEvent) (formattedLog string, err error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	jsonLogInfo := &jsonLogInfo{
		LoggerName: loggerName,
		Program:    log.Program(),
		Pid:        log.Pid(),
		Hostname:   log.Hostname(),
		Time:       log.Time().Format(f.dateTimeLayout),
		LogLevel:   log.LogLevel(),
		Pc:         log.Pc(),
		FileName:   log.FileName(),
		LineNum:    log.LineNum(),
		Message:    log.Message(),
		Attrs:      log.GetAttrs(),
	}
	serialized, err := json.Marshal(jsonLogInfo)
	if err != nil {
		return "", err
	}
	return string(append(serialized, '\n')), nil
}

//SetDateTimeLayout is set layout of date and time. See Time.Format.
func (f *JsonFormatter) SetDateTimeLayout(dateTimeLayout string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.dateTimeLayout = dateTimeLayout
}

//NewJsonFormatter is create JsonFormatter
func NewJsonFormatter() (jsonFormatter *JsonFormatter) {
	return &JsonFormatter{
		dateTimeLayout: "2006-01-02 15:04:05 -0700 MST",
		mutex:          new(sync.RWMutex),
	}
}

func init() {
	RegisterFormatter("JsonFormatter", func() (formatter Formatter) {
		return NewJsonFormatter()
	})
}
