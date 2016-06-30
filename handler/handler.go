package Handler

import (
	"github.com/pkg/errors"
)

var (
	handlers map[string]func() Handler
	buffers  map[string]func() BufferManager
)

type Handler interface {
	SetBufferManager(bufferManager buffer.BufferManager) (err error)
	Open()
	Write(loggerName string, logEvent *belog.LogEvent, formattedLog string)
	Flush()
	Close()
}

type BufferManager interface {
	AddBuffer(logEvent *belog.LogEvent, formattedLog string) (lastLogEvent *belog.LogEvent, logBuffer string, full bool)
	DrainBuffer() (lastLogEvent *belog.LogEvent, logbuffer string, remain bool)
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

func GetBuffer(name string) (buffer BufferManager, err error) {
	newFunc, ok := buffers[name]
	if !ok {
		return nil, errors.Errorf("not found buffer (%v)", name)
	}
	return newFunc(), nil
}

func RegisterBuffer(name string, newFunc func() BufferManager) {
	buffers[name] = newFunc
}
