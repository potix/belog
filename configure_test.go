package belog

import (
	"testing"
)

func TestLoadConfiagJSON1(t *testing.T) {
	if err := LoadConfig("./test/sample1.jsn"); err != nil {
		t.Errorf("%+v", err)
	}
}

func TestLoadConfiagJSON2(t *testing.T) {
	if err := LoadConfig("./test/sample1.json"); err != nil {
		t.Errorf("%+v", err)
	}
}

func TestLoadConfiagToml1(t *testing.T) {
	if err := LoadConfig("./test/sample1.tml"); err != nil {
		t.Errorf("%+v", err)
	}
}

func TestLoadConfiagToml2(t *testing.T) {
	if err := LoadConfig("./test/sample1.toml"); err != nil {
		t.Errorf("%+v", err)
	}
}

func TestLoadConfiagYaml1(t *testing.T) {
	if err := LoadConfig("./test/sample1.yml"); err != nil {
		t.Errorf("%+v", err)
	}
}

func TestLoadConfiagYaml2(t *testing.T) {
	if err := LoadConfig("./test/sample1.yaml"); err != nil {
		t.Errorf("%+v", err)
	}
}

func TestGetLogger(t *testing.T) {
	if err := LoadConfig("./test/sample1.jsn"); err != nil {
		t.Errorf("%+v", err)
	}
	l1, ok := loggers["test1"]
	if !ok {
		t.Errorf("not found test1 logger")
	}
	l2, ok := loggers["test2"]
	if !ok {
		t.Errorf("not found test2 logger")
	}
	loggerGroup := GetLoggerGroup("test1", "test2")
	ll1, ok := loggerGroup.loggers["test1"]
	if !ok {
		t.Errorf("not found test1 logger")
	}
	ll2, ok := loggerGroup.loggers["test2"]
	if !ok {
		t.Errorf("not found test2 logger")
	}
	if l1 != ll1 {
		t.Errorf("l1 != ll1")
	}
	if l2 != ll2 {
		t.Errorf("l2 != ll2")
	}
}
