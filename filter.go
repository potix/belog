package belog

import (
	"github.com/pkg/errors"
)

var (
	filters map[string]func() Filter
)

type Filter interface {
	Evaluate(loggerName string, log LogEvent) bool
}

func GetFilter(name string) (filter Filter, err error) {
	newFunc, ok := filters[name]
	if !ok {
		return nil, errors.Errorf("not found filter (%v)", name)
	}
	return newFunc(), nil
}

func RegisterFilter(name string, newFunc func() Filter) {
	if filters == nil {
		filters = make(map[string]func() Filter)
	}
	filters[name] = newFunc
}

func init() {
	if filters == nil {
		filters = make(map[string]func() Filter)
	}
}
