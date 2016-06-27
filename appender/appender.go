package appender

var (
	appenders map[string]func() Appender
)

type Appender interface {
	SetBufferManager(bufferManager buffer.BufferManager)
	Write(logString string)
	Flush()
}
