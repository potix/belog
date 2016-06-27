package buffer

type BufferManager interface {
	AddBuffer(logString string) (logString string)
	DrainBuffer() (logString string)
}
