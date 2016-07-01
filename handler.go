package belog

import (
	"github.com/pkg/errors"
)

var (
	handlers map[string]func() Handler
)

type Handler interface {
	Open()
	Write(loggerName string, logEvent LogEvent, formattedLog string)
	Flush()
	Close()
}

func GetHandler(name string) (Handler Handler, err error) {
	newFunc, ok := handlers[name]
	if !ok {
		return nil, errors.Errorf("not found Handler (%v)", name)
	}
	return newFunc(), nil
}

func RegisterHandler(name string, newFunc func() Handler) {
	if handlers == nil {
		handlers = make(map[string]func() Handler)
	}
	handlers[name] = newFunc
}

func init() {
	if handlers == nil {
		handlers = make(map[string]func() Handler)
	}
}
