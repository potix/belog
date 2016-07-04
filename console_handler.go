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
		LogLevelEmerg:  ConsoleColorLightMagenta,
		LogLevelAlert:  ConsoleColorLightRed,
		LogLevelCrit:   ConsoleColorMagenta,
		LogLevelError:  ConsoleColorRed,
		LogLevelWarn:   ConsoleColorYellow,
		LogLevelNotice: ConsoleColorGreen,
		LogLevelInfo:   ConsoleColorBlue,
		LogLevelDebug:  ConsoleColorCyan,
		LogLevelTrace:  ConsoleColorLightGray,
	}
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
	case OutputTypeStderr:
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
func (h *ConsoleHandler) SetOutputType(outputType OutputType) {
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
		outputType: OutputTypeStdout,
		mutex:      new(sync.RWMutex),
	}
}

func init() {
	RegisterHandler("ConsoleHandler", func() (handler Handler) {
		return NewConsoleHandler()
	})
}
