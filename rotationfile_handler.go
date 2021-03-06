package belog

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	rotationFileDefaultAsyncFlushInterval = 1
	rotationFileDefaultBufferSize         = 8192
	rotationFileDirLayout                 = "2006-01-02"
)

//RotationFileHandler is handler of file with rotation
type RotationFileHandler struct {
	logFileName        string
	logDirPath         string
	maxAge             int
	maxSize            int64
	async              bool
	asyncFlushInterval int
	bufferSize         int
	buffer             *bytes.Buffer
	lastLogEvent       LogEvent
	scheduledFlush     bool
	logFileSize        int64
	lastModifiedTime   time.Time
	logFile            *os.File
	mutex              *sync.Mutex
}

//IsOpened is opened
func (h *RotationFileHandler) IsOpened() (bool) {
	h.mutex.Lock()
	f := h.logFile
	h.mutex.Unlock()
	if f == nil  {
		return false
	}
	return true
}

//Open is open file
func (h *RotationFileHandler) Open() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.openLogFile()
}

//Write is write formatted log
func (h *RotationFileHandler) Write(loggerName string, logEvent LogEvent, formattedLog string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.async {
		lastLogEvent, logBuffer, full := h.pushBuffer(logEvent, formattedLog)
		if full {
			h.writeLog(lastLogEvent.Time(), logBuffer)
		} else {
			// timer flush
			if !h.scheduledFlush {
				h.scheduledFlush = true
				go h.logBufferFlushTimer()
			}
		}
	} else {
		h.writeLog(logEvent.Time(), formattedLog)
	}
}

//Flush is flush buffer and sync
func (h *RotationFileHandler) Flush() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.async {
		h.logBufferFlush()
	}
	if h.logFile != nil {
		err := h.logFile.Sync()
		if err != nil {
			// statistics
		}
	}
}

//Close is close File
func (h *RotationFileHandler) Close() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.logFile == nil {
		return
	}
	if h.async {
		h.logBufferFlush()
	}
	err := h.logFile.Close()
	if err != nil {
		// statistics
	}
	h.logFile = nil
	h.lastModifiedTime = time.Time{}
	h.logFileSize = 0
}

//SetLogFileName is set log file name
func (h *RotationFileHandler) SetLogFileName(logFileName string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.logFileName = logFileName
}

//SetLogDirPath is set log dir path
func (h *RotationFileHandler) SetLogDirPath(logDirPath string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.logDirPath = logDirPath
}

//SetMaxAge is set max age
func (h *RotationFileHandler) SetMaxAge(maxAge int) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.maxAge = maxAge
}

//SetMaxSize is set max log file size. if wrote size over this size, file is rotated.
func (h *RotationFileHandler) SetMaxSize(maxSize int64) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.maxSize = maxSize
}

//SetAsync is set async mode
func (h *RotationFileHandler) SetAsync(async bool) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.async == async {
		return
	}
	if h.async == true && async == false {
		h.logBufferFlush()
	}
	h.async = async
}

//SetAsyncFlushInterval is set flush timer interfal
func (h *RotationFileHandler) SetAsyncFlushInterval(asyncFlushInterval int) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.asyncFlushInterval = asyncFlushInterval
}

//SetBufferSize is set buffer size. if buffer of async size over this size, write buffer to file.
func (h *RotationFileHandler) SetBufferSize(bufferSize int) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.bufferSize = bufferSize
}

func (h *RotationFileHandler) logBufferFlushTimer() {
	time.Sleep(time.Second)
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.logBufferFlush()
	h.scheduledFlush = false
}

func (h *RotationFileHandler) logBufferFlush() {
	lastLogEvent, logBuffer, remain := h.popBuffer()
	if remain {
		h.writeLog(lastLogEvent.Time(), logBuffer)
	}
}

func (h *RotationFileHandler) writeLog(logTime time.Time, logBuffer string) {
	h.openLogFile()
	h.rotateLogFile(logTime)
	if h.logFile == nil {
		// statistics
		return
	}
	wlen, err := h.logFile.WriteString(logBuffer)
	if err != nil {
		// sstatistics
		return
	}
	h.lastModifiedTime = logTime
	h.logFileSize += int64(wlen)
}

func (h *RotationFileHandler) openLogFile() {
	if h.logFile != nil {
		return
	}
	// make directories
	err := os.MkdirAll(h.logDirPath, os.FileMode(0755))
	if err != nil {
		// statistics
		return
	}
	// open log file
	logFilePath := filepath.Join(h.logDirPath, h.logFileName)
	file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.FileMode(0644))
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
		(h.maxSize <= 0 || h.logFileSize < h.maxSize) {
		return
	}
	logFilePath := filepath.Join(h.logDirPath, h.logFileName)
	// get rotated file path
	rotatedLogDirPath, rotatedLogFilePath := h.getRotatedLogFilePath()
	if err := os.MkdirAll(rotatedLogDirPath, os.FileMode(0755)); err != nil {
		// statistics
		return
	}
	// rename
	if err := os.Rename(logFilePath, rotatedLogFilePath); err != nil {
		// statistics
		return
	}
	// open new log file
	file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.FileMode(0644))
	if err != nil {
		// statistics
		return
	}
	if err := h.logFile.Close(); err != nil {
		// statistics
	}
	h.logFile = file
	h.lastModifiedTime = lastLogTime
	h.logFileSize = 0
	h.deleteOldLogFiles()
}

func (h *RotationFileHandler) getRotatedLogFilePath() (rotatedLogDirPath string, rotatedLogFilePath string) {
	idx := 1
	date := h.lastModifiedTime.Format(rotationFileDirLayout)
	for {
		rotatedLogFileName := fmt.Sprintf("%v.%v.%v", h.logFileName, date, idx)
		rotatedLogDirPath := filepath.Join(h.logDirPath, date)
		rotatedLogFilePath := filepath.Join(rotatedLogDirPath, rotatedLogFileName)
		_, err := os.Stat(rotatedLogFilePath)
		if err != nil {
			return rotatedLogDirPath, rotatedLogFilePath
		}
		idx++
	}
}

func (h *RotationFileHandler) deleteOldLogFiles() {
	files, err := ioutil.ReadDir(h.logDirPath)
	if err != nil {
		return
	}
	oldTime := h.lastModifiedTime.AddDate(0, 0, -1*h.maxAge)
	_, offset := oldTime.Zone()
	oldAdjustTime := time.Unix((oldTime.Unix()/86400)*86400-int64(offset), 0)
	for _, file := range files {
		if !file.IsDir() {
			continue
		}
		dirTime, err := time.Parse(rotationFileDirLayout, file.Name())
		if err != nil {
			// statistics
			continue
		}
		if dirTime.Before(oldAdjustTime) {
			err := os.RemoveAll(filepath.Join(h.logDirPath, file.Name()))
			if err != nil {
				// statistics
			}
		}
	}
}

func (h *RotationFileHandler) pushBuffer(logEvent LogEvent, formattedLog string) (lastLogEvent LogEvent, logBuffer string, full bool) {
	_, err := h.buffer.WriteString(formattedLog)
	if err != nil {
		// statistics
	}
	h.lastLogEvent = logEvent
	if h.buffer.Len() > h.bufferSize {
		return h.popBufferBase()
	}
	return nil, "", false
}

func (h *RotationFileHandler) popBuffer() (lastLogEvent LogEvent, logBuffer string, remain bool) {
	if h.buffer.Len() > 0 {
		return h.popBufferBase()
	}
	return nil, "", false
}

func (h *RotationFileHandler) popBufferBase() (lastLogEvent LogEvent, logBuffer string, remain bool) {
	logBuffer = h.buffer.String()
	h.buffer.Truncate(0)
	lastLogEvent = h.lastLogEvent
	h.lastLogEvent = nil
	return lastLogEvent, logBuffer, true
}

//NewRotationFileHandler is create RotationFileHandler
func NewRotationFileHandler() (rotationFileHandler *RotationFileHandler) {
	return &RotationFileHandler{
		logFileName:        fmt.Sprintf("%v.log", filepath.Base(os.Args[0])),
		logDirPath:         fmt.Sprintf("/var/log/%v", filepath.Base(os.Args[0])),
		maxAge:             7,
		async:              false,
		asyncFlushInterval: rotationFileDefaultAsyncFlushInterval,
		buffer:             bytes.NewBuffer(make([]byte, 0, rotationFileDefaultBufferSize)),
		bufferSize:         rotationFileDefaultBufferSize,
		mutex:              new(sync.Mutex),
	}
}

func init() {
	RegisterHandler("RotationFileHandler", func() (handler Handler) {
		return NewRotationFileHandler()
	})
}
