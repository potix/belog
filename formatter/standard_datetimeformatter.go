package formatter

import (
	"sync"
)

type StandardDateTimeFormatter struct {
}

func (f *StandardDateTimeFormatter) DateTimeFormat(log belog.Log) {
	// format
	// XXXXXXX
}

type StandardFormatter struct {
	dateTimeFormatter DateTimeFormatter
	mutex             *sync.RWMutex
}

func (f *StandardFormatter) Format(log belog.Log) {
	f.mutex.RLock()
	defer f.mutex.RUnlock()
	// format
	// XXXXXX
}

func (f *StandardFormatter) SetDateTimeFormatter(dateTimeFormatter DateTimeFormatter) {
	f.mutex.Lock()
	defer f.mutex.Unlock()
	f.dateTimeFormatter = dateTimeFormatter

}
