package handler

import (
	"os"
	"path/filepath"
	"time"
)

const (
	defaultAsyncFlushInterval = 1
)

type RotationFileHandler struct {
	LogFileName        string
	LogDirPath         string
	MaxGen             int
	MaxSize            int
	DailyRotation      bool
	AsyncFlushInterval int
	bufferManager      buffer.BufferManager
	scheduledFlush     bool
	lastModifiedTime   time.Time
	logDileSize        int64
	logFile            *os.File
	mutex              *sync.Mutex
}

func (a *RotationFileHandler) SetBufferManager(bufferManager buffer.BufferManager) (err error) {
	a.bufferManager = bufferManager
	return nil
}

func (a *RotationFileHandler) Open() {
	a.openLogFile()
}

func (a *RotationFileHandler) Write(logString string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if bufferManager != nil {
		logBuffer, needFlush := bufferManager.AddBuffer(logString)
		if needFlush {
			a.writeLog(logBuffer)
		}
		// timer flush
		if !scheduledFlush {
			a.scheduledFlush = true
			go a.timerFlush()
		}
	} else {
		a.writeLog(logString)
	}
}

func (a *RotationFileHandler) Flush() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if bufferManager != nil {
		logBuffer, needFlush := bufferManager.DrainBuffer()
		if needFlush {
			a.writeLog(logBuffer)
		}
	}
	if c.logFile != nil {
		c.logFile.Sync()
	}
}

func (a *RotationFileHandler) Close() {
	if a.logFile == nil {
		return
	}
	a.logFile.Close()
	a.logFile = nil
	a.lastModifiedTime = time.Time{}
	a.logFileSize = 0
}

func (a *RotationFileHandler) timerFlush() {
	<-time.After(time.Second)
	a.mutex.Lock()
	defer a.mutex.Unlock()
	logBuffer, needFlush := bufferManager.DrainBuffer()
	if needFlush {
		a.writeLog(logBuffer)
	}
	a.scheduledFlush = false
}

func (a *RotationFileHandler) writeLog(logBuffer string) {
	a.openLogFile()
	now := time.Now()
	a.rotateLogFile(now)
	if a.logFile == nil {
		// statistics
		return
	}
	wlen, err := a.logFile.writeString(logBuffer)
	if err != nil {
		// sstatistics
		return
	}
	a.lastModifiedTime = now
	a.logFileSize += int64(wlen)
}

func (a *RotationFileHandler) openLogFile() {
	if a.logFile != nil {
		return
	}
	// make directories
	err := os.MkdirAll(a.LogDirPath, 0755)
	if err != nil {
		// statistics
		return
	}
	// open log file
	logFilePath := filepath.Join(c.LogDirPath, c.LogFileName)
	file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		// statisticsa
		return
	}
	a.logFile = file
	// get stat info
	fileInfo, _ := a.logFile.Stat()
	a.lastModifiedTime = fileInfo.ModTime()
	a.logFileSize = fileInfo.Size()
}

func (a *RotationFileHandler) rotateLogFile(now time.Time) {
	if (a.lastModifiedTime.Year() == now.Year() &&
		a.lastModifiedTime.YearDay() == now.YearDay()) &&
		(a.MaxSize <= 0 || a.logFileSize < a.MaxSize) {
		return
	}
	// get rotated file path
	rotatedLogFilePath := a.getRotatedLogFilePath()

	logFilePath := filepath.Join(c.LogDirPath, c.LogFileName)
	// rename
	err := os.Rename(logFilePath, rotatedLogFilePath)
	if err != nil {
		// statistics
		return
	}
	// open new log file
	file, err := os.OpenFile(logFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		// statistics
		return
	}
	a.logFile.Close()
	a.logFile = file
	a.lastModifiedTime = now
	a.logFileSize = 0
}

func (a *RotationFileHandler) getRotatedLogFilePath() {
	rotatedLogFilePath := filepath.Join(c.LogDirPath, rotatedLogFileName)
}

func NewRotationFileHandler() (handler Handler) {
	return &rotationFile{
		LogFileName:        fmt.Sprintf("%v.log", filepath.Base(os.Args[0])),
		LogDirPath:         "/var/log/",
		MaxGen:             7,
		DailyRotation:      true,
		AsyncFlushInterval: defaultAsyncFlushInterval,
	}
}

func init() {
	RegisterHandler("RotationFile", NewRotationFileHandler)
}
