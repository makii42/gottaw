package defaults

import (
	"fmt"
	"os"
	"os/exec"
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
func (def *defaultsUtil) isExecutable(names ...string) bool {
	for _, name := range names {
		binPath, err := exec.LookPath(name)
		if err == nil && binPath != "" {
			return true
		}
	}
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
