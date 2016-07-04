package belog

import (
	"os"
	"sync"
)

//OutputType is output type
type OutputType int

const (
	//OutputTypeStdout is output type of stdout
	OutputTypeStdout OutputType = 1
	//OutputTypeStderr is output type of stderr
	OutputTypeStderr = 2
)

//ConsoleHandler is handler of console
type ConsoleHandler struct {
	outputType OutputType
	mutex      *sync.RWMutex
}

//Open is nothing to do
func (h *ConsoleHandler) Open() {
}

//Write is output to console
func (h *ConsoleHandler) Write(loggerName string, logEvent LogEvent, formattedLog string) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	switch h.outputType {
	case OutputTypeStdout:
		_, err := os.Stdout.WriteString(formattedLog)
		if err != nil {
			// statistics
		}
	case OutputTypeStderr:
		_, err := os.Stderr.WriteString(formattedLog)
		if err != nil {
			// statistics
		}
	}
}

//Flush is nothing to do
func (h *ConsoleHandler) Flush() {
}

//Close is nothing to do
func (h *ConsoleHandler) Close() {
}

//SetOutputType is set output type
func (h *ConsoleHandler) SetOutputType(outputType OutputType) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.outputType = outputType
}

//NewConsoleHandler is create ConsoleHandler
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
