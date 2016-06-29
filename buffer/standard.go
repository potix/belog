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
	stringBuffer *bytes.Buffer
}

func (b *StandardBuffer) AddBuffer(logString string) (stringBuffer string, needFlush bool) {
	b.stringBuffer.WriteString(logString)
	if b.stringBuffer.Len() > c.Threshold {
		stringBuffer := b.stringBuffer.String()
		b.stringBuffer.Truncate(0)
		return stringBuffer, true
	}
	return "", false
}

func (b *StandardBuffer) DrainBuffer() (stringBuffer string, needFlush bool) {
	if len(c.stringBuffer) > 0 {
		stringBuffer := b.stringBuffer.String()
		b.stringBuffer.Truncate(0)
		return stringBuffer, true
	}
	return "", false
}

func NewStandardBuffer() (buffer BufferManager) {
	return &StandardBuffer{
		stringBuffer: bytes.NewBuffer(make([]byte, 0, defaultBufferSize)),
		Threshold:    defaultBufferSize,
	}
}

func init() {
	RegisterBuffer("StandardBuffer")
}
