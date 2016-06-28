package buffer

import (
	"github.com/pkg/errors"
)

var (
	buffers map[string]func() BufferManager
)

type BufferManager interface {
	AddBuffer(logString string) (stringBuffer string, needFlush bool)
	DrainBuffer() (stringBuffer string, needFlush bool)
}

func GetBuffer(name string) (buffer BufferManager, err error) {
	buffer, ok := buffers[name]
	if !ok {
		return nil, errors.Errorf("not found buffer (%v)", name)
	}
	return buffer, nil
}

func RegisterBuffer(name string, newFunc func() BufferManager) {
	buffers[name] = newFunc
}
