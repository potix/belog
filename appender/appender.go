package appender

import (
	"github.com/pkg/errors"
)

var (
	appenders map[string]func() Appender
)

type Appender interface {
	SetBufferManager(bufferManager buffer.BufferManager) (err error)
	Start()
	Write(logString string)
	Flush()
	Stop()
}

func GetAppender(name string) (appender Appender, err error) {
	appender, ok := appenders[name]
	if !ok {
		return nil, errors.Errorf("not found appender (%v)", name)
	}
	return appender, nil
}

func RegisterAppender(name string, newFunc func() AppenderManager) {
	appenders[name] = newFunc
}
