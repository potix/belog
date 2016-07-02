package belog

import (
	"github.com/pkg/errors"
)

var (
	filters map[string]func() Filter
)

//Filter is interface of fileter
type Filter interface {
	Evaluate(loggerName string, log LogEvent) bool
}

func getFilter(name string) (filter Filter, err error) {
	newFunc, ok := filters[name]
	if !ok {
		return nil, errors.Errorf("not found filter (%v)", name)
	}
	return newFunc(), nil
}

//RegisterFilter is register filter
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
