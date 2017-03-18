package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsIgnoreWorksAsExpected(t *testing.T) {
	ignored, err := isIgnored("./foo.go", cfgExcludes("bla.go"))
	assert.Nil(t, err)
	assert.False(t, ignored)
}

func TestIsIgnoreIgnoresPlainFile(t *testing.T) {
	ignored, err := isIgnored("/bla", cfgExcludes("\\/bla"))
	assert.Nil(t, err)
	assert.True(t, ignored)
}
func TestIsIgnoreIgnoresSuffix(t *testing.T) {
	ignored, err := isIgnored("foo.bla", cfgExcludes("*bla"))
	assert.Nil(t, err)
	assert.True(t, ignored)
}

func TestIsIgnoreIgnoresAllInDir(t *testing.T) {
	ignored, err := isIgnored("./.git/hooks", cfgExcludes("./.git/**"))
	assert.Nil(t, err)
	assert.True(t, ignored)
}

func TestIsIgnoreIgnoresAllInGitAsWeDontHaveDoubleStar(t *testing.T) {
	c := cfgExcludes("./.git/*/*")
	ignored, err := isIgnored("./.git/hooks/bas", c)
	assert.Nil(t, err)
	assert.True(t, ignored)
}

func cfgExcludes(excludes ...string) *Config {
	return &Config{
		Excludes: excludes,
	}
}
