package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsIgnoreWorksAsExpected(t *testing.T) {
	ignored := isIgnored("./foo.go", cfgExcludes("bla.go"))
	assert.False(t, ignored)
}

func TestIsIgnoreIgnoresPlainFile(t *testing.T) {
	ignored := isIgnored("/bla", cfgExcludes("/bla"))
	assert.True(t, ignored)
}
func TestIsIgnoreIgnoresSuffix(t *testing.T) {
	ignored := isIgnored("foo.bla", cfgExcludes("*bla"))
	assert.True(t, ignored)
}

func TestIsIgnoreIgnoresAllInDir(t *testing.T) {
	ignored := isIgnored("./.git/hooks", cfgExcludes("./.git/**"))
	assert.True(t, ignored)
}

func TestIsIgnoreIgnoresAllInGitAsWeDontHaveDoubleStar(t *testing.T) {
	c := cfgExcludes("./.git/*/*")
	ignored := isIgnored("./.git/hooks/bas", c)
	assert.True(t, ignored)
}

func cfgExcludes(excludes ...string) *Config {
	return &Config{
		Excludes: excludes,
	}
}
