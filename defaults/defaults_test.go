package defaults

import (
	"path/filepath"
	"testing"

	"io/ioutil"
	"os"

	"github.com/stretchr/testify/assert"
)

func TestFileExistsReturnsTrueFile(t *testing.T) {
	assert.True(t, fileExists(".", "/defaults_test.go"))
}
func TestFileExistsReturnsFalseIfFileNotExists(t *testing.T) {
	assert.False(t, fileExists(".", "/acme.txt"))
}
func TestFileExistsReturnsFalseForDirs(t *testing.T) {
	assert.False(t, fileExists(".", ""))
}

func TestDirExistsReturnsTrueIfThere(t *testing.T) {
	assert.True(t, dirExists("."))
}
func TestDirExistsReturnsFalseIfNotThere(t *testing.T) {
	assert.False(t, dirExists("./node_modules"))
}
func TestDirExistsReturnsFalseIfFile(t *testing.T) {
	assert.False(t, dirExists("./defaults_test.go"))
}

// These tests do ensure the defaults are recognized properly.
// TODO: I need a way to properly fix up the test environment, as
// it requires the single binaries to be installed locally.

func TestNodeYarnDefault(t *testing.T) {
	tmpdir, err := ioutil.TempDir("/tmp", "nodeyarn")
	if err != nil {
		t.Fatal("could not create tempdir")
	}
	defer os.RemoveAll(tmpdir)
	pkgJson := []byte("{name: \"nodepkg\"}")
	pkgName := filepath.Join(tmpdir, "package.json")
	if err := ioutil.WriteFile(pkgName, pkgJson, 0666); err != nil {
		t.Fatal("could not write package.json")
	}

	subject := NodeYarnDefault{}
	assert.True(t, subject.Test(tmpdir))
}
