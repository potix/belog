package formatter

import (
	"github.com/pkg/errors"
)

var (
	dateTimeFormatters map[string]func() dateTimeFormatter
	formatters         map[string]func() Formatter
)

type Formatter interface {
	SetDateTimeFormatter(dateTimeFormatter DateTimeFormatter)
	Format(loggerName string, log belog.Log) (logString string)
}

type DateTimeFormatter interface {
	DateTimeFormat(log *belog.Log) (dateTime string)
}

var (
	buffers map[string]func() BufferManager
)

func GetFormatter(name string) (formatter Formatter, err error) {
	newFunc, ok := formatters[name]
	if !ok {
		return nil, errors.Errorf("not found formatter (%v)", name)
	}
	return newFunc(), nil
}

func RegisterFormatter(name string, newFunc func() Formatter) {
	formaters[name] = newFunc
}

func GetDateTimeFormatter(name string) (dateTimeFormatter DateTimeFormatter, err error) {
	newFunc, ok := dateTimeFormatters[name]
	if !ok {
		return nil, errors.Errorf("not found date time formatter (%v)", name)
	}
	return newFunc(), nil
}

func RegisterDateTimeFormatter(name string, newFunc func() DateTimeFormatter) {
	dateTimeFormaters[name] = newFunc
}
