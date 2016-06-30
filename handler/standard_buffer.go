package buffer

//   configuration sample
//
//    "loggers" {
//         "logger1": {
//             "filter":
//                .....
//             "formatter"
//                .....
//             "appenders": [
//                 {
//                     "type": ...
//                     "params": ....
//		       "buffer": {
//                         "type": "StandardBuffer",
//                         "params": "Threshold=8192"
//                     }
//                 },
//             ]
//         }
//    }
//

import (
	"sync"
)

const (
	defaultBufferSize = 8192
)

type StandardBuffer struct {
	Threshold    int
	logBuffer    *bytes.Buffer
	lastLogEvent *belog.LogEvent
}

func (b *StandardBuffer) AddBuffer(logEvent *belog.LogEvent, formattedLog string) (lastLogEvent *belog.LogEvent, logBuffer string, full bool) {
	b.logBuffer.WriteString(formattedLog)
	b.lastLogEvent = logEvent
	if b.logBuffer.Len() > c.Threshold {
		logBuffer := b.logBuffer.String()
		b.logBuffer.Truncate(0)
		lastLogEvent := b.lastLogEvent
		b.lastLogEvent = nil
		return lastLogEvent, stringBuffer, true
	}
	return nil, "", false
}

func (b *StandardBuffer) DrainBuffer() (lastLogEvent *belog.LogEvent, logBuffer string, remain bool) {
	if len(c.logBuffer) > 0 {
		logBuffer := b.logBuffer.String()
		b.logBuffer.Truncate(0)
		lastLogEvent := b.lastLogEvent
		b.lastLogEvent = nil
		return lastLogEvent, stringBuffer, true
	}
	return nil, "", false
}

func NewStandardBuffer() (buffer BufferManager) {
	return &StandardBuffer{
		logBuffer: bytes.NewBuffer(make([]byte, 0, defaultBufferSize)),
		Threshold: defaultBufferSize,
	}
}

func init() {
	RegisterBuffer("StandardBuffer", NewStandardBuffer)
}
