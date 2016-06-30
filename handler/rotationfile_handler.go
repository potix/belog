package handler

import (
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	defaultAsyncFlushInterval = 1
)

type RotationFileHandler struct {
	LogFileName        string
	LogDirPath         string
	MaxAge             int
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
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.bufferManager = bufferManager
	return nil
}

func (a *RotationFileHandler) Open() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	a.openLogFile()
}

func (a *RotationFileHandler) Write(loggerName string, logEvent *belog.LogEvent, formattedLog string) {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if bufferManager != nil {
		lastLogEvent, formattedLog, needFlush := bufferManager.AddBuffer(logEvent, formattedLog)
		if needFlush {
			a.writeLog(lastLogEvent.Time(), formattedLog)
		}
		// timer flush
		if !scheduledFlush {
			a.scheduledFlush = true
			go a.timerFlush()
		}
	} else {
		a.writeLog(logEvent.Time(), formattedLog)
	}
}

func (a *RotationFileHandler) Flush() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if bufferManager != nil {
		lastLogEvent, logBuffer, needFlush := bufferManager.DrainBuffer()
		if needFlush {
			a.writeLog(lastLogEvent.Time(), logBuffer)
		}
	}
	if a.logFile != nil {
		a.logFile.Sync()
	}
}

func (a *RotationFileHandler) Close() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
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
	lastLogEvent, logBuffer, needFlush := bufferManager.DrainBuffer()
	if needFlush {
		a.writeLog(lastLogEvent.Time(), logBuffer)
	}
	a.scheduledFlush = false
}

func (a *RotationFileHandler) writeLog(logTime time.Time, logBuffer string) {
	a.openLogFile()
	a.rotateLogFile(lastLogTime)
	if a.logFile == nil {
		// statistics
		return
	}
	wlen, err := a.logFile.writeString(logBuffer)
	if err != nil {
		// sstatistics
		return
	}
	a.lastModifiedTime = lastLogTime
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
	logFilePath := filepath.Join(a.LogDirPath, a.LogFileName)
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

func (a *RotationFileHandler) rotateLogFile(lastLogTime time.Time) {
	if (a.lastModifiedTime.Year() == lastLogTime.Year() &&
		a.lastModifiedTime.YearDay() == lastLogTime.YearDay()) &&
		(a.MaxSize <= 0 || a.logFileSize < a.MaxSize) {
		return
	}
	logFilePath := filepath.Join(a.LogDirPath, a.LogFileName)
	// get rotated file path
	rotatedLogDirPath, rotatedLogFilePath := a.getRotatedLogFilePath()
	err := os.MkdirAll(rotatedLogDirPath, 0755)
	if err != nil {
		// statistics
		return
	}
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
	a.lastModifiedTime = lastLogTime
	a.logFileSize = 0
	a.deleteOldLog()
}

func (a *RotationFileHandler) getRotatedLogFilePath() {
	idx := 1
	date := a.lastModifiedTime.strftime("%Y-%m-%d")
	for {
		rotatedLogFileName := fmt.Sprintf("%v.%v.%v", a.LogFileName, date, idx)
		rotatedLogDirPath = filepath.Join(a.LogDirPath, date)
		rotatedLogFilePath = filepath.Join(rotatedLogDirPath, rotatedLogFileName)
		_, err := os.Stat(rotatedLogFilePath)
		if err != nil {
			return rotatedLogDirPath, rotatedLogFilePath
		}
		idx += 1
	}
}

func (a *RotationFileHandler) deleteOldLog() {
	oldDate = a.lastModifiedTime.AddDate(0, -1*(a.MaxAge+1), 0)
	oldRotatedLogDirPath = filepath.Join(a.LogDirPath, oldDate)
	os.RemoveAll(oldRotatedLogDirPath)
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
	RegisterHandler("RotationFileHandler", NewRotationFileHandler)
}
