package belog

import (
	"path/filepath"
	"strings"
	"sync"
)

//
// %(dateTime)
// %(logLevel)
// %(logLevelNum)
// %(pid)
// %(loggerName)
// %(programCounter)
// %(fileName)
// %(shortFileName)
// %(lineNum)
// %(message)
//

type StandardFormatter struct {
	dateTimeLayout string
	layout         string
	mutex          *sync.RWMutex
}

func (f *StandardFormatter) Format(loggerName string, log LogEvent) (formattedLog string) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	replacer := strings.NewReplacer(
		"%(dateTime)", log.Time().Format(f.dateTimeLayout),
		"%(logLevel)", log.LogLevelString(),
		"%(logLevelNum)", string(log.LogLevel()),
		"%(pid)", string(log.Pid()),
		"%(loggerName)", loggerName,
		"%(programCounter)", string(log.Pc()),
		"%(fileName)", log.FileName(),
		"%(shortFileName)", filepath.Base(log.FileName()),
		"%(lineNum)", string(log.LineNum()),
		"%(message)", log.Message())
	return replacer.Replace(f.layout)
}

func (f *StandardFormatter) SetDateTimeLayout(dateTimeLayout string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.dateTimeLayout = dateTimeLayout
}

func (f *StandardFormatter) SetLayout(layout string) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.layout = layout
}

func NewStandardFormatter() (formatter Formatter) {
	return &StandardFormatter{
		dateTimeLayout: "2006-01-02 15:04:05",
		layout:         "%(dateTime) [%(logLevel)] (%(pid)) %(loggerName) %(fileName) %(lineNum) %(message)",
		mutex:          new(sync.RWMutex),
	}
}

func init() {
	RegisterFormatter("StandardFormatter", NewStandardFormatter)
}
