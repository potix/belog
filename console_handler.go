package belog

import (
	"fmt"
	"os"
	"sync"
)

//ConsoleColor is console color
type ConsoleColor int

const (
	//ConsoleNoColor is no color
	ConsoleNoColor           ConsoleColor = 0
	//ConsoleColorBlack is black
	ConsoleColorBlack                     = 30
	//ConsoleColorRed is red
	ConsoleColorRed                       = 31
	//ConsoleColorGreen is green
	ConsoleColorGreen                     = 32
	//ConsoleColorYellow is yellow
	ConsoleColorYellow                    = 33
	//ConsoleColorBlue is blue
	ConsoleColorBlue                      = 34
	//ConsoleColorMagenta is magenta
	ConsoleColorMagenta                   = 35
	//ConsoleColorCyan is cyan
	ConsoleColorCyan                      = 36
	//ConsoleColorLightGray is light gray
	ConsoleColorLightGray                 = 37
	//ConsoleColorDarkGray is dark gray 
	ConsoleColorDarkGray                  = 90
	//ConsoleColorLightRed light red
	ConsoleColorLightRed                  = 91
	//ConsoleColorLightGreen light freen
	ConsoleColorLightGreen                = 92
	//ConsoleColorLightYellow is light yellow
	ConsoleColorLightYellow               = 93
	//ConsoleColorLightBlue is light blue
	ConsoleColorLightBlue                 = 94
	//ConsoleColorLightMagenta is light magenta
	ConsoleColorLightMagenta              = 95
	//ConsoleColorLightCyan is light cyan
	ConsoleColorLightCyan                 = 96
	//ConsoleColorWhite is white
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

//IsOpened is nothing to do
func (h *ConsoleHandler) IsOpened() (bool) {
	return true
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
