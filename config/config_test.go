package config

import (
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
	File = "../.gottaw.yml"
	cfg := Load()
	assert.NotNil(t, cfg)
	assert.True(t, len(cfg.Excludes) > 1)
	assert.True(t, len(cfg.Pipeline) > 1)
	assert.True(t, cfg.Growl)
}

func TestSetupPanicsWhenFileNotPresent(t *tt.T) {
	File = "./does-not-exist-should-panic.yml"
	Load()
}

func TestSetupPanicsWhenEmptyStringIsPassed(t *tt.T) {
	File = ""
	Load()
}
func TestSetupPanicsWhenCrapIsPassed(t *tt.T) {
	File = "###"
	Load()
}

func TestSetupPanicsWhenBrokenFile(t *tt.T) {
	dir, _ := ioutil.TempDir("", "gottaw-config-test")
	File = filepath.Join(dir, "testfile")
	defer os.RemoveAll(dir)
	if err := ioutil.WriteFile(File, []byte(testFile), 0666); err != nil {
		t.Fatalf("could not write temp %s\n", File)
	}
	Load()
	log.Printf("All is well!")
}

func TestSerializeConfigDoesTheJob(t *tt.T) {
	cfg := &Config{
		Growl:    true,
		Pipeline: []string{"echo \"Hello, World\""},
		Excludes: []string{".git", ".hg"},
	}
	data, err := SerializeConfig(cfg)
	assert.Nil(t, err)
	assert.Equal(t, `excludes:
- .git
- .hg
pipeline:
- echo "Hello, World"
growl: true
`, string(data))
}
