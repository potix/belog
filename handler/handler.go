package Handler

import (
	"github.com/pkg/errors"
)

var (
	handlers map[string]func() Handler
)

type Handler interface {
	SetBufferManager(bufferManager buffer.BufferManager) (err error)
	Open()
	Write(logString string)
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
	handlers[name] = newFunc
}
