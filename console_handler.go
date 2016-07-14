package belog

import (
	"fmt"
	"os"
	"sync"
)

type ConsoleColor int

const (
	ConsoleNoColor           ConsoleColor = 0
	ConsoleColorBlack                     = 30
	ConsoleColorRed                       = 31
	ConsoleColorGreen                     = 32
	ConsoleColorYellow                    = 33
	ConsoleColorBlue                      = 34
	ConsoleColorMagenta                   = 35
	ConsoleColorCyan                      = 36
	ConsoleColorLightGray                 = 37
	ConsoleColorDarkGray                  = 90
	ConsoleColorLightRed                  = 91
	ConsoleColorLightGreen                = 92
	ConsoleColorLightYellow               = 93
	ConsoleColorLightBlue                 = 94
	ConsoleColorLightMagenta              = 95
	ConsoleColorLightCyan                 = 96
	ConsoleColorWhite                     = 97
)

var (
	colorMap = map[LogLevel]ConsoleColor{
		LogLevelEmerg:  ConsoleNoColor,
		LogLevelAlert:  ConsoleNoColor,
		LogLevelCrit:   ConsoleNoColor,
		LogLevelError:  ConsoleNoColor,
		LogLevelWarn:   ConsoleNoColor,
		LogLevelNotice: ConsoleNoColor,
		LogLevelInfo:   ConsoleNoColor,
		LogLevelDebug:  ConsoleNoColor,
		LogLevelTrace:  ConsoleNoColor,
	}
)

//ConsoleOutputType is output type
type ConsoleOutputType int

const (
	//ConsoleOutputTypeStdout is output type of stdout
	ConsoleOutputTypeStdout ConsoleOutputType = 1
	//ConsoleOutputTypeStderr is output type of stderr
	ConsoleOutputTypeStderr = 2
)

//ConsoleHandler is handler of console
type ConsoleHandler struct {
	outputType ConsoleOutputType
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
	case ConsoleOutputTypeStdout:
		color, ok := colorMap[logEvent.LogLevelNum()]
		if !ok {
			color = ConsoleColorBlue
		}
		_, err := os.Stdout.WriteString(fmt.Sprintf("\x1b[%dm", color))
		if err != nil {
			// statistics
		}
		_, err = os.Stdout.WriteString(formattedLog)
		if err != nil {
			// statistics
		}
		_, err = os.Stdout.WriteString("\x1b[0m")
		if err != nil {
			// statistics
		}
	case ConsoleOutputTypeStderr:
		color, ok := colorMap[logEvent.LogLevelNum()]
		if !ok {
			color = ConsoleColorBlue
		}
		_, err := os.Stderr.WriteString(fmt.Sprintf("\x1b[%dm", color))
		if err != nil {
			// statistics
		}
		_, err = os.Stderr.WriteString(formattedLog)
		if err != nil {
			// statistics
		}
		_, err = os.Stderr.WriteString("\033[0m")
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
func (h *ConsoleHandler) SetOutputType(outputType ConsoleOutputType) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.outputType = outputType
}

//SetConsoleColor is set color by log level
func (h *ConsoleHandler) SetConsoleColor(loglevel LogLevel, color ConsoleColor) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	colorMap[loglevel] = color
}

//NewConsoleHandler is create ConsoleHandler
func NewConsoleHandler() (consoleHandler *ConsoleHandler) {
	return &ConsoleHandler{
		outputType: ConsoleOutputTypeStdout,
		mutex:      new(sync.RWMutex),
	}
}

func init() {
	RegisterHandler("ConsoleHandler", func() (handler Handler) {
		return NewConsoleHandler()
	})
}
