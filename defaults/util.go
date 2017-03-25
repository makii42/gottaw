package defaults

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var defaultExcludes = []string{
	".git",
	".hg",
	".vscode",
	".idea",
	".gitrecinfo",
}

func fileExists(dir string, file string) bool {
	if info, err := os.Stat(path.Join(dir, file)); err == nil {
		return info.Mode().IsRegular()
	}
	return false
}
func dirExists(path string) bool {
	if info, err := os.Stat(path); err == nil {
		return info.Mode().IsDir()
	}
	return false
}
func isExecutable(name string) bool {
	path := os.Getenv("PATH")
	// find binary in path and ensure it has some x-es
	for _, dir := range strings.Split(path, string(os.PathListSeparator)) {
		if file, err := os.Stat(filepath.Join(dir, name)); err == nil &&
			file.Mode().IsRegular() &&
			file.Mode().Perm()&0111 != 0 {
			return true
		}
	}
	return false
}

func filesMatch(dir string, pattern string) bool {
	var sep string
	if !strings.HasSuffix(dir, string(os.PathSeparator)) {
		dir = string(os.PathSeparator)
	}
	p := fmt.Sprintf("%s%s%s", dir, sep, pattern)
	matches, err := filepath.Glob(p)
	if err != nil {
		panic(err)
	}
	return len(matches) > 0
}
