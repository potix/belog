package appender

import (
	"os"
	"path/filepath"
	"time"
)

const (
	defaultAsyncFlushInterval = 1
)

type RotationFileAppender struct {
	LogFileName        string
	LogDirPath         string
	MaxGen             int
	MaxSize            int
	DailyRotation      bool
	AsyncFlushInterval int
	bufferManager      buffer.BufferManager
	scheduledFlush     bool
	mutex              *sync.Mutex
}

func (a *RotationFileAppender) SetBufferManager(bufferManager buffer.BufferManager) (err error) {
	a.bufferManager = bufferManager
	return nil
}

func (a *RotationFileAppender) Write(logString string) {
	now := time.Now()
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if bufferManager != nil {
		logBuffer, needFlush := bufferManager.AddBuffer(logString)
		if needFlush {
			a.writeLog(now, logBuffer)
		}
		if !scheduledFlush {
			a.scheduledFlush = true
			go a.waitAndFlush(now)
		}
	} else {
		a.writeLog(now, logString)
	}
}

func (a *RotationFileAppender) Flush() {
	now := time.Now()
	a.mutex.Lock()
	defer a.mutex.Unlock()
	if bufferManager != nil {
		logBuffer, needFlush := bufferManager.DrainBuffer()
		if needFlush {
			a.writeLog(now, logBuffer)
		}
		// file flush
	} else {
		// file flush
	}
}

func (a *RotationFileAppender) waitAndFlush(now time.Time) {
	<-time.After(time.Second)
	a.mutex.Lock()
	defer a.mutex.Unlock()
	logBuffer, needFlush := bufferManager.DrainBuffer()
	if needFlush {
		a.writeLog(now, logBuffer)
	}
	a.scheduledFlush = false
}

func (a *RotationFileAppender) writeLog(now time.Time, logBuffer string) {
	sync.Once.Do(createLogDirPath)
	file, err := getFile(time.Now())
	if err != nil {
		// statistics
	}
	_, err := file.writeString(logBuffer)
	if err != nil {
		// sstatistics
	}
}

func (a *RotationFileAppender) createLogDirPath() (err error) {
	err := os.MkdirAll(a.LogDirPath, 0755)
	if err != nil {
		// statistics
	}
}

func (a *RotationFileAppender) getFile(now time.Time) (file *File, err error) {
	// check rotation
	newFileName, needRotate := checkRotation(now)
	if a.oldFileName != newFile
	

	if needRotation {
		// open new file 
		open new file	
		// close old file
		close old file
		oldFile
	}

	// XXXX
	file, ok := a.newFile
	if !ok {
		return nil, errors.Errorf("not found file")
	}

	return file, nil
}

func NewRotationFileAppender() (appender Appender) {
	return &rotationFile{
		LogFileName:        fmt.Sprintf("%v.log", filepath.Base(os.Args[0])),
		LogDirPath:         "/var/log/",
		MaxGen:             7,
		DailyRotation:      true,
		AsyncFlushInterval: defaultAsyncFlushInterval,
	}
}

func init() {
	appenders["RotationFile"] = NewRotationFile
}
