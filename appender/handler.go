package Handler

import (
	"github.com/pkg/errors"
)

var (
	Handlers map[string]func() Handler
)

type Handler interface {
	SetBufferManager(bufferManager buffer.BufferManager) (err error)
	Write(logString string)
	Flush()
}

func GetHandler(name string) (Handler Handler, err error) {
	Handler, ok := Handlers[name]
	if !ok {
		return nil, errors.Errorf("not found Handler (%v)", name)
	}
	return Handler, nil
}

func RegisterHandler(name string, newFunc func() HandlerManager) {
	Handlers[name] = newFunc
}
