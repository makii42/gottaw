package defaults

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/makii42/gottaw/output"
)

var defaultExcludes = []string{
	".git",
	".hg",
	".vscode",
	".idea",
	".gitrecinfo",
}

const foundMsg = "✅  found!"
const notFoundMsg = "❌  not found!"

func newDefaultsUtil(l *output.Logger) *defaultsUtil {
	return &defaultsUtil{l: l}
}

type defaultsUtil struct {
	l *output.Logger
}

func (def *defaultsUtil) fileExists(dir string, file string) bool {
	abs := path.Join(dir, file)
	def.l.Tracef("file '%s'? ", abs)
	if info, err := os.Stat(abs); err == nil {
		def.l.Traceln(foundMsg)
		return info.Mode().IsRegular()
	}
	def.l.Traceln(notFoundMsg)
	return false
}
func (def *defaultsUtil) dirExists(path string) bool {
	def.l.Tracef("dir '%s'? ", path)
	if info, err := os.Stat(path); err == nil {
		def.l.Traceln(foundMsg)
		return info.Mode().IsDir()
	}
	def.l.Traceln(notFoundMsg)
	return false
}
func (def *defaultsUtil) isExecutable(name string) bool {
	path := os.Getenv("PATH")
	def.l.Tracef("executable '%s' on PATH? ", name)
	// find binary in path and ensure it has some x-es
	for _, dir := range strings.Split(path, string(os.PathListSeparator)) {
		if file, err := os.Stat(filepath.Join(dir, name)); err == nil &&
			file.Mode().IsRegular() &&
			file.Mode().Perm()&0111 != 0 {
			def.l.Traceln(foundMsg)
			return true
		}
	}
	def.l.Traceln(notFoundMsg)

	return false
}

func (def *defaultsUtil) filesMatch(dir string, pattern string) bool {
	def.l.Tracef("testing for files matching '%s' in '%s' ... ", pattern, dir)
	var sep string
	if !strings.HasSuffix(dir, string(os.PathSeparator)) {
		sep = string(os.PathSeparator)
	}
	p := fmt.Sprintf("%s%s%s", dir, sep, pattern)
	matches, err := filepath.Glob(p)
	if err != nil {
		panic(err)
	}
	if len(matches) > 0 {
		def.l.Traceln(foundMsg)
		return true
	}
	def.l.Traceln(notFoundMsg)
	return false
}
