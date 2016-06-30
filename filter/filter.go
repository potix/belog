package filter

import (
	"github.com/pkg/errors"
)

var (
	filters map[string]func() Filter
)

type Filter interface {
	Evaluate(loggerName string, log belog.Log) bool
}

func GetFilter(name string) (filter Filter, err error) {
	newFunc, ok := filters[name]
	if !ok {
		return nil, errors.Errorf("not found filter (%v)", name)
	}
	return newFunc(), nil
}

func RegisterFilter(name string, newFunc func() Filter) {
	filters[name] = newFunc
}
