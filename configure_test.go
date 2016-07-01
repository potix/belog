package belog

import (
	"testing"
)

func TestLoadConfiagJson(t *testing.T) {
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
