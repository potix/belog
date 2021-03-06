package belog

import (
	"path/filepath"
	"strconv"
	"strings"
	"sync"
)

//StandardFormatter is standard formatter
//this formatter is replace particular tags.
type StandardFormatter struct {
	appendNewLine   bool
	dateTimeLayout  string
	layout          string
	mutex           *sync.RWMutex
}

//Format is format log event
func (f *StandardFormatter) Format(loggerName string, log LogEvent) (formattedLog string, err error) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	logMessage := log.Message()
	if f.appendNewLine {
		if len(logMessage) == 0 || logMessage[len(logMessage) - 1:] != "\n" {
			logMessage = logMessage + "\n"
		}
	}
	replacer := strings.NewReplacer(
		"%(dateTime)", log.Time().Format(f.dateTimeLayout),
		"%(logLevel)", log.LogLevel(),
		"%(logLevelNum)", strconv.Itoa(int(log.LogLevelNum())),
		"%(program)", log.Program(),
		"%(pid)", strconv.Itoa(log.Pid()),
		"%(hostname)", log.Hostname(),
		"%(loggerName)", loggerName,
		"%(programCounter)", strconv.FormatUint(uint64(log.Pc()), 16),
		"%(fileName)", log.FileName(),
		"%(shortFileName)", filepath.Base(log.FileName()),
		"%(lineNum)", strconv.Itoa(log.LineNum()),
		"%(message)", logMessage)
	return replacer.Replace(f.layout), nil
}

//SetAppendNewLine is set append new line.
func (f *StandardFormatter) SetAppendNewLine(appendNewLine bool) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.appendNewLine = appendNewLine
}

//SetDateTimeLayout is set layout of date and time. See Time.Format.
func (f *StandardFormatter) SetDateTimeLayout(dateTimeLayout string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.dateTimeLayout = dateTimeLayout
}

//SetLayout is set layout.
//usable tags is follow:
//   %(dateTime)       : date adn time
//   %(logLevel)       : log level
//   %(logLevelNum)    : log level number
//   %(program)        : program name
//   %(pid)            : process id
//   %(hostname)       : hostname
//   %(loggerName)     : loggername
//   %(programCounter) : program counter
//   %(fileName)       : filename (full path)
//   %(shortFileName)  : short file name (basename only)
//   %(lineNum)        : line number
//   %(message)        : message
func (f *StandardFormatter) SetLayout(layout string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.layout = layout
}

//NewStandardFormatter is create StandardFormatter
func NewStandardFormatter() (standardFormatter *StandardFormatter) {
	return &StandardFormatter{
		appendNewLine:  true,
		dateTimeLayout: "2006-01-02 15:04:05",
		layout:         "%(dateTime) [%(logLevel)] (%(pid)) %(program) %(loggerName) %(fileName) %(lineNum) %(message)",
		mutex:          new(sync.RWMutex),
	}
}

func init() {
	RegisterFormatter("StandardFormatter", func() (formatter Formatter) {
		return NewStandardFormatter()
	})
}
