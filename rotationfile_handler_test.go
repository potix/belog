package belog

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestRotationFile(t *testing.T) {
	now := time.Now()
	os.RemoveAll("/var/tmp/belog-test")
	if err := os.MkdirAll("/var/tmp/belog-test", os.FileMode(0755)); err != nil {
		t.Errorf("%+v", err)
	}
	old := now.AddDate(0, 0, -1)
	oldstr := old.Format("200601021504.05")
	cmd := exec.Command("touch", "-t", oldstr, "/var/tmp/belog-test/belog-test.log")
	if err := cmd.Run(); err != nil {
		t.Errorf("%+v", err)
	}
	logInfo := &logInfo{
		time: now,
	}
	rotationFilehandler := NewRotationFileHandler()
	rotationFilehandler.SetLogFileName("belog-test.log")
	rotationFilehandler.SetLogDirPath("/var/tmp/belog-test")
	rotationFilehandler.SetMaxAge(3)
	rotationFilehandler.SetMaxSize(500)
	for i := 0; i < 100; i++ {
		rotationFilehandler.Write("test", logInfo, "01234567890")
	}
	files, err := ioutil.ReadDir("/var/tmp/belog-test")
	if err != nil {
		t.Errorf("%+v", err)
	}
	if len(files) != 3 {
		t.Errorf("files count mismatch /var/tmp/belog-test")
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
	dirName := old.Format(rotationFileDirLayout)
	files, err = ioutil.ReadDir(fmt.Sprintf("/var/tmp/belog-test/%v", dirName))
	if err != nil {
		t.Errorf("%+v", err)
	}
	if len(files) != 1 {
		t.Errorf("files count mismatch /var/tmp/belog-test/%v", dirName)
	}
	dirName = now.Format(rotationFileDirLayout)
	files, err = ioutil.ReadDir(fmt.Sprintf("/var/tmp/belog-test/%v", dirName))
	if err != nil {
		t.Errorf("%+v", err)
	}
	if len(files) != 2 {
		t.Errorf("files count mismatch /var/tmp/belog-test/%v", dirName)
	}
}

func TestRotationFileCleanup(t *testing.T) {
	now := time.Now()
	os.RemoveAll("/var/tmp/belog-test")
	if err := os.MkdirAll("/var/tmp/belog-test", os.FileMode(0755)); err != nil {
		t.Errorf("%+v", err)
	}
	cmd := exec.Command("touch", "-t", "201601010101.00", "/var/tmp/belog-test/belog-test.log")
	if err := cmd.Run(); err != nil {
		t.Errorf("%+v", err)
	}
	logInfo := &logInfo{
		time: now,
	}
	rotationFilehandler := NewRotationFileHandler()
	rotationFilehandler.SetLogFileName("belog-test.log")
	rotationFilehandler.SetLogDirPath("/var/tmp/belog-test")
	rotationFilehandler.SetMaxAge(3)
	rotationFilehandler.SetMaxSize(500)
	for i := 0; i < 100; i++ {
		rotationFilehandler.Write("test", logInfo, "01234567890")
	}
	files, err := ioutil.ReadDir("/var/tmp/belog-test")
	if err != nil {
		t.Errorf("%+v", err)
	}
	if len(files) != 2 {
		t.Errorf("files count mismatch /var/tmp/belog-test")
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
	dirName := now.Format(rotationFileDirLayout)
	files, err = ioutil.ReadDir(fmt.Sprintf("/var/tmp/belog-test/%v", dirName))
	if err != nil {
		t.Errorf("%+v", err)
	}
	if len(files) != 2 {
		t.Errorf("files count mismatch /var/tmp/belog-test/%v", dirName)
	}
}
