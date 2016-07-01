package belog

import (
	"bytes"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	rotationFileDefaultAsyncFlushInterval = 1
	rotationFileDefaultBufferSize         = 8192
)

type RotationFileHandler struct {
	logFileName        string
	logDirPath         string
	maxAge             int
	maxSize            int
	async              bool
	asyncFlushInterval int
	logBuffer          *rotationFileLogBuffer
	scheduledFlush     bool
	logFileSize        int64
	lastModifiedTime   time.Time
	logFile            *os.File
	mutex              *sync.Mutex
}

func (h *RotationFileHandler) Open() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.openLogFile()
}

func (h *RotationFileHandler) Write(loggerName string, logEvent LogEvent, formattedLog string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.async {
		lastLogEvent, formattedLog, full := h.logBuffer.AddBuffer(logEvent, formattedLog)
		if full {
			h.writeLog(lastLogEvent.Time(), formattedLog)
		} else {
			// timer flush
			if !scheduledFlush {
				h.scheduledFlush = true
				go h.logBufferFlushTimer()
			}
		}
	} else {
		h.writeLog(logEvent.Time(), formattedLog)
	}
}

func (h *RotationFileHandler) Flush() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.async {
		h.logBufferFlush()
	}
	if h.logFile != nil {
		h.logFile.Sync()
	}
}

func (h *RotationFileHandler) Close() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.logFile == nil {
		return
	}
	h.flushBase()
	h.logFile.Close()
	h.logFile = nil
	h.lastModifiedTime = time.Time{}
	h.logFileSize = 0
}

func (h *RotationFileHandler) SetLogFileName(logFileName string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.logFileName = logFileName
}

func (h *RotationFileHandler) SetLogDirPath(logDirPath string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.logDirPath = logDirPath
}

func (h *RotationFileHandler) SetMaxAge(maxAge int) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.maxAge = maxAge
}

func (h *RotationFileHandler) SetMaxSize(maxSize int) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.maxSize = maxSize
}

func (h *RotationFileHandler) SetAsync(async bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.async == async {
		return
	}
	if h.async == ture && async == false {
		h.logBufferFlush()
	}
	h.async = async
}

func (h *RotationFileHandler) SetAsyncFlushInterval(asyncFlushInterval int) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.asyncFlushInterval = asyncFlushInterval
}

func (h *RotationFileHandler) logBufferFlushTimer() {
	<-time.After(time.Second)
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.logBufferFlush()
	h.scheduledFlush = false
}

func (h *RotationFileHandler) logBufferFlush() {
	lastLogEvent, logBuffer, remain := h.logBuffer.DrainBuffer()
	if remain {
		h.writeLog(lastLogEvent.Time(), logBuffer)
	}
}

func (h *RotationFileHandler) writeLog(logTime time.Time, logBuffer string) {
	h.openLogFile()
	h.rotateLogFile(lastLogTime)
	if h.logFile == nil {
		// statistics
		return
	}
	wlen, err := h.logFile.writeString(logBuffer)
	if err != nil {
		// sstatistics
		return
	}
	h.lastModifiedTime = lastLogTime
	h.logFileSize += int64(wlen)
}

func (h *RotationFileHandler) openLogFile() {
	if h.logFile != nil {
		return
	}
	// make directories
	err := os.MkdirAll(h.LogDirPath, 0755)
	if err != nil {
		// statistics
		return
	}
	// open log file
	logFilePath := filepath.Join(h.LogDirPath, h.LogFileName)
	file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		// statisticsa
		return
	}
	h.logFile = file
	// get stat info
	fileInfo, _ := h.logFile.Stat()
	h.lastModifiedTime = fileInfo.ModTime()
	h.logFileSize = fileInfo.Size()
}

func (h *RotationFileHandler) rotateLogFile(lastLogTime time.Time) {
	if (h.lastModifiedTime.Year() == lastLogTime.Year() &&
		h.lastModifiedTime.YearDay() == lastLogTime.YearDay()) &&
		(h.MaxSize <= 0 || h.logFileSize < h.MaxSize) {
		return
	}
	logFilePath := filepath.Join(h.LogDirPath, h.LogFileName)
	// get rotated file path
	rotatedLogDirPath, rotatedLogFilePath := h.getRotatedLogFilePath()
	if err := os.MkdirAll(rotatedLogDirPath, 0755); err != nil {
		// statistics
		return
	}
	// rename
	if err := os.Rename(logFilePath, rotatedLogFilePath); err != nil {
		// statistics
		return
	}
	// open new log file
	file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		// statistics
		return
	}
	h.logFile.Close()
	h.logFile = file
	h.lastModifiedTime = lastLogTime
	h.logFileSize = 0
	h.deleteOldLogFiles()
}

func (h *RotationFileHandler) getRotatedLogFilePath() {
	idx := 1
	date := h.lastModifiedTime.strftime("%Y-%m-%d")
	for {
		rotatedLogFileName := fmt.Sprintf("%v.%v.%v", h.LogFileName, date, idx)
		rotatedLogDirPath := filepath.Join(h.LogDirPath, date)
		rotatedLogFilePath := filepath.Join(rotatedLogDirPath, rotatedLogFileName)
		_, err := os.Stat(rotatedLogFilePath)
		if err != nil {
			return rotatedLogDirPath, rotatedLogFilePath
		}
		idx += 1
	}
}

func (h *RotationFileHandler) deleteOldLogFiles() {
	files, err := ioutil.ReadDir(h.LogDirPath)
	if err != nil {
		return
	}
	oldTime := h.lastModifiedTime.AddDate(0, 0, -1*h.MaxAge)
	_, offset := oldTime.Zone()
	oldAdjustTime := time.Unix((oldTime.Unix()/86400)*86400-int(offset), 0)
	const dirLayout = "2006-01-02"
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		dirTime, err := time.Parse(dirLayout, file.Name())
		if err != nil {
			// statistics
			continue
		}
		if dirTime.Before(oldAdjustTime) {
			os.RemoveAll(filepath.Join(h.LogDirPath, file.Name()))
		}
	}
}

func NewRotationFileHandler() (handler Handler) {
	return &RotationFileHandler{
		logFileName:        fmt.Sprintf("%v.log", filepath.Base(os.Args[0])),
		logDirPath:         fmt.Sprintf("/var/log/%v", filepath.Base(os.Args[0])),
		maxGen:             7,
		async:              false,
		asyncFlushInterval: rotationFileDefaultAsyncFlushInterval,
		logBuffer:          newLogBuffer(),
		mutex:              new(sync.Mutex),
	}
}

type rotationFileLogBuffer struct {
	Threshold    int
	buffer       *bytes.Buffer
	lastLogEvent *LogEvent
}

func (b *rotationFileLogBuffer) addBuffer(logEvent *LogEvent, formattedLog string) (lastLogEvent *LogEvent, logBuffer string, full bool) {
	b.buffer.WriteString(formattedLog)
	b.lastLogEvent = logEvent
	if b.buffer.Len() > c.Threshold {
		logBuffer := b.buffer.String()
		b.buffer.Truncate(0)
		lastLogEvent := b.lastLogEvent
		b.lastLogEvent = nil
		return lastLogEvent, logBuffer, true
	}
	return nil, "", false
}

func (b *rotationFileLogBuffer) drainBuffer() (lastLogEvent *LogEvent, logBuffer string, remain bool) {
	if b.buffer.Len() > 0 {
		logBuffer := b.buffer.String()
		b.buffer.Truncate(0)
		lastLogEvent := b.lastLogEvent
		b.lastLogEvent = nil
		return lastLogEvent, logBuffer, true
	}
	return nil, "", false
}

func newRotationFileLogBuffer() (rotationFileLogBuffer *rotationFileLogBuffer) {
	return &rotationFileLogBuffer{
		buffer:    bytes.NewBuffer(make([]byte, 0, rotationFileDefaultBufferSize)),
		Threshold: rotationFileDefaultBufferSize,
	}
}

func init() {
	RegisterHandler("RotationFileHandler", NewRotationFileHandler)
}
