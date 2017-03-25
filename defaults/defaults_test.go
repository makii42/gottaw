package defaults

import (
	"path/filepath"
	"testing"

	"io/ioutil"
	"os"

	"fmt"

	"github.com/stretchr/testify/assert"
)

var temp_root string
var original_path string

func TestMain(m *testing.M) {
	dir, err := ioutil.TempDir("", "gottaw-test")
	original_path = os.Getenv("PATH")
	if err != nil {
		panic(err)
	}
	temp_root = dir
	defer os.RemoveAll(dir)
	result := m.Run()
	os.Setenv("PATH", original_path)
	os.Exit(result)
}

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
	tmpdir, err := ioutil.TempDir(temp_root, "nodeyarn")
	if err != nil {
		t.Fatal("could not create tempdir")
	}
	defer os.RemoveAll(tmpdir)
	addFile(t, tmpdir, "package.json", []byte("{name: \"nodepkg\"}"), 0666)
	binFolder := addBinFolder(t, tmpdir)
	addBin(t, binFolder, "node")
	addBin(t, binFolder, "yarn")
	subject := NodeYarnDefault{}
	assert.True(t, subject.Test(tmpdir))
}

func addFile(t *testing.T, dir string, filename string, contents []byte, perm os.FileMode) {
	filepath := filepath.Join(dir, filename)
	if err := ioutil.WriteFile(filepath, contents, perm); err != nil {
		t.Fatalf("could not write temp %s\n", filename)
	}
}

func addBin(t *testing.T, binFolder string, binName string) {
	addFile(
		t,
		binFolder,
		binName,
		[]byte(fmt.Sprintf("#!/bin/sh\necho \"temp bin file %s\"\n", binName)),
		0755,
	)
}

func addBinFolder(t *testing.T, dir string) string {
	binDir := filepath.Join(dir, "bin")
	err := os.Mkdir(binDir, 0777)
	if err != nil {
		t.Fatalf("could not create bin dir %s", binDir)
	}
	os.Setenv(
		"PATH",
		fmt.Sprintf("%s%s%s", binDir, string(os.PathListSeparator), original_path),
	)
	return binDir
}
