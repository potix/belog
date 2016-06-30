package formatter

import (
	"github.com/pkg/errors"
)

var (
	formatters map[string]func() Formatter
)

type Formatter interface {
	Format(loggerName string, log belog.Log) (logString string)
}

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
