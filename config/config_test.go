package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	tt "testing"

	"os"

	"github.com/stretchr/testify/assert"
)

var tempRoot string

const (
	testFile = `growl: lalala
pipeline: this should break it bigtime::: ===
excludes: this as well
`
)

func TestSetupWorks(t *tt.T) {
	// we presume a certain state in this file,
	// so changes might cause test failures
	cfg := Setup("../.gottaw.yml")
	assert.NotNil(t, cfg)
	assert.True(t, len(cfg.Excludes) > 1)
	assert.True(t, len(cfg.Pipeline) > 1)
	assert.True(t, cfg.Growl)
}

func TestSetupPanicsWhenFileNotPresent(t *tt.T) {
	defer assertPanic(t, "setup with non-existent file")
	Setup("./does-not-exist-should-panic.yml")
}

func TestSetupPanicsWhenEmptyStringIsPassed(t *tt.T) {
	defer assertPanic(t, "setup with crap passed")
	Setup("")
}
func TestSetupPanicsWhenCrapIsPassed(t *tt.T) {
	defer assertPanic(t, "setup with crap passed")
	Setup("###")
}

func TestSetupPanicsWhenBrokenFile(t *tt.T) {
	dir, _ := ioutil.TempDir("", "gottaw-config-test")
	filepath := filepath.Join(dir, "testfile")
	defer os.RemoveAll(dir)
	if err := ioutil.WriteFile(filepath, []byte(testFile), 0666); err != nil {
		t.Fatalf("could not write temp %s\n", filepath)
	}
	defer assertPanic(t, "parse")
	Setup(filepath)
	log.Printf("All is well!")
}

func assertPanic(t *tt.T, what string) {
	if r := recover(); r == nil {
		t.Errorf(fmt.Sprintf("%s did not panic as expected", what))
	}
}
