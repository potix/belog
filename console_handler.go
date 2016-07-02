package belog

import (
	"os"
	"sync"
)

type OutputType int

const (
	OutputTypeStdout OutputType = 1
	OutputTypeStderr            = 2
)

type ConsoleHandler struct {
	outputType OutputType
	mutex      *sync.RWMutex
}

func (h *ConsoleHandler) Open() {
}

func (h *ConsoleHandler) Write(loggerName string, logEvent LogEvent, formattedLog string) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	switch h.outputType {
	case OutputTypeStdout:
		os.Stdout.WriteString(formattedLog)
	case OutputTypeStderr:
		os.Stderr.WriteString(formattedLog)
	}
}

func (h *ConsoleHandler) Flush() {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	switch h.outputType {
	case OutputTypeStdout:
		os.Stdout.Sync()
	case OutputTypeStderr:
		os.Stderr.Sync()
	}
}

func (h *ConsoleHandler) Close() {
}

func (h *ConsoleHandler) SetOutputType(outputType OutputType) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.outputType = outputType
}

func NewConsoleHandler() (consoleHandler *ConsoleHandler) {
	return &ConsoleHandler{
		outputType: OutputTypeStdout,
		mutex:      new(sync.RWMutex),
	}
}

func init() {
	RegisterHandler("ConsoleHandler", func() (handler Handler) {
		return NewConsoleHandler()
	})
}
