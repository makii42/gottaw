package defaults

import (
	"path/filepath"
	"testing"

	"io/ioutil"
	"os"

	"fmt"

	"github.com/makii42/gottaw/config"
	"github.com/makii42/gottaw/output"
	"github.com/stretchr/testify/assert"
)

var packageJsonContents = []byte("{name: \"nodepkg\"}")
var tempRoot string
var util *defaultsUtil
var logger *output.Logger

var golang, nodeNpm, nodeYarn, javaMaven Default

func TestMain(m *testing.M) {
	// deps in trace - YES thats not quiet by default
	logger = output.NewLogger(output.NOTICE, &config.Config{})
	util = newDefaultsUtil(logger)

	// test default objects
	golang = NewGolangDefault(util)
	nodeNpm = NewNodeNpmDefault(util)
	nodeYarn = NewNodeYarnDefault(util)
	javaMaven = NewJavaMavenDefault(util)

	// create and rollback test root directory
	dir, err := ioutil.TempDir("", "gottaw-test")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(dir)
	tempRoot = dir

	// create and rollback path
	originalPath := os.Getenv("PATH")
	os.Setenv("PATH", "")
	fmt.Printf("FIXED PATH '%s'", os.Getenv("PATH"))
	defer os.Setenv("PATH", originalPath)

	result := m.Run()
	os.Exit(result)
}

// These tests do ensure the defaults are recognized properly.
// TODO: Take windows into the fold by providing "exe"/"bat" suffixes
// for tests on windows.

func TestGolangDefault(t *testing.T) {
	tmpDir := createGolangEnv(t, tempRoot)
	// positive
	assert.True(t, golang.Test(tmpDir))

	// negative tests
	assert.False(t, nodeYarn.Test(tmpDir))
	assert.False(t, nodeNpm.Test(tmpDir))
	assert.False(t, javaMaven.Test(tmpDir))
}

func TestNodeYarnDefault(t *testing.T) {
	tmpDir := createNodeYarnEnv(t, tempRoot)
	// positive
	assert.True(t, nodeYarn.Test(tmpDir))

	// negative tests
	assert.False(t, golang.Test(tmpDir))
	assert.False(t, nodeNpm.Test(tmpDir))
	assert.False(t, javaMaven.Test(tmpDir))
}

func TestNodeNpmDefault(t *testing.T) {
	tmpDir := createNodeNpmEnv(t, tempRoot)
	// positive
	assert.True(t, nodeNpm.Test(tmpDir))

	// negative tests
	assert.False(t, golang.Test(tmpDir))
	assert.False(t, nodeYarn.Test(tmpDir))
	assert.False(t, javaMaven.Test(tmpDir))
}

func createNodeYarnEnv(t *testing.T, tempRoot string) string {
	tmpDir, err := ioutil.TempDir(tempRoot, "nodeyarn-")
	if err != nil {
		t.Fatal("could not create tempdir")
	}
	addFile(t, tmpDir, "package.json", packageJsonContents, 0666)
	binFolder := addBinFolder(t, tmpDir)
	addBin(t, binFolder, "node")
	addBin(t, binFolder, "yarn")
	return tmpDir
}

func createNodeNpmEnv(t *testing.T, tempRoot string) string {
	tmpDir, err := ioutil.TempDir(tempRoot, "nodenpm-")
	if err != nil {
		t.Fatal("could not create temp dir")
	}
	addFile(t, tmpDir, "package.json", packageJsonContents, 0666)
	binFolder := addBinFolder(t, tmpDir)
	addBin(t, binFolder, "node")
	addBin(t, binFolder, "npm")
	return tmpDir
}

func createGolangEnv(t *testing.T, tempRoot string) string {
	tmpDir, err := ioutil.TempDir(tempRoot, "golang-")
	if err != nil {
		t.Fatal("could not create temp dir")
	}
	addFile(t, tmpDir, "main.go", packageJsonContents, 0666)
	addFile(t, tmpDir, "foobar.go", packageJsonContents, 0666)
	binFolder := addBinFolder(t, tmpDir)
	addBin(t, binFolder, "go")
	return tmpDir
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
	os.Setenv("PATH", binDir)
	return binDir
}
