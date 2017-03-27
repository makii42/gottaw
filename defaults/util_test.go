package defaults

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExistsReturnsTrueFile(t *testing.T) {
	assert.True(t, util.fileExists(".", "/defaults_test.go"))
}

func TestFileExistsReturnsFalseIfFileNotExists(t *testing.T) {
	assert.False(t, util.fileExists(".", "/acme.txt"))
}

func TestFileExistsReturnsFalseForDirs(t *testing.T) {
	assert.False(t, util.fileExists(".", ""))
}

func TestDirExistsReturnsTrueIfThere(t *testing.T) {
	assert.True(t, util.dirExists("."))
}

func TestDirExistsReturnsFalseIfNotThere(t *testing.T) {
	assert.False(t, util.dirExists("./node_modules"))
}

func TestDirExistsReturnsFalseIfFile(t *testing.T) {
	assert.False(t, util.dirExists("./defaults_test.go"))
}

func TestFilesMatchReturnsTrueIfMatches(t *testing.T) {
	assert.True(t, util.filesMatch(".", "*go"))
}

func TestFilesMatchReturnsFalseIfNoMatches(t *testing.T) {
	assert.False(t, util.filesMatch(".", "*.exe"))
}
