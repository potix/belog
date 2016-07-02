package belog

import (
	"github.com/pkg/errors"
)

var (
	formatters map[string]func() Formatter
)

//Formatter is interface of formatter
type Formatter interface {
	Format(loggerName string, log LogEvent) (logString string)
}

func getFormatter(name string) (formatter Formatter, err error) {
	newFunc, ok := formatters[name]
	if !ok {
		return nil, errors.Errorf("not found formatter (%v)", name)
	}
	return newFunc(), nil
}

//RegisterFormatter is register formatter
func RegisterFormatter(name string, newFunc func() Formatter) {
	if formatters == nil {
		formatters = make(map[string]func() Formatter)
	}
	formatters[name] = newFunc
}

func init() {
	if formatters == nil {
		formatters = make(map[string]func() Formatter)
	}
}
