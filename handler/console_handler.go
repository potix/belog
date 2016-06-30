package handler

import (
	"os"
	"sync"
)

var OutputType int

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

func (h *ConsoleHandler) Write(loggerName string, logEvent *belog.LogEvent, formattedLog string) {
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
		os.Stdout.WriteString(formattedLog)
	case OutputTypeStderr:
		os.Stderr.WriteString(formattedLog)
	}
}

func (h *ConsoleHandler) Close() {
}

func (h *ConsoleHandler) SetOutputType(outputType OutputType) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.outputType = outputType
}

func NewConsoleHandler() (handler Handler) {
	return &ConsoleHandler{
		outputType: OutputTypeStdout,
		mutex:      new(sync.RWMutex),
	}
}

func init() {
	RegisterHandler("ConsoleHandler", NewConsoleHandler)
}
