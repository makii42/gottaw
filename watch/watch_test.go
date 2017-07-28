package watch

import (
	"io/ioutil"
	"testing"

	"gopkg.in/fsnotify.v1"

	"os"

	"path/filepath"

	"github.com/makii42/gottaw/config"
	"github.com/stretchr/testify/assert"
)

func TestIsIgnoreWorksAsExpected(t *testing.T) {
	ignored := isIgnored("./foo.go", cfgExcludes("bla.go"))
	assert.False(t, ignored)
}

func TestIsIgnoreIgnoresPlainFile(t *testing.T) {
	ignored := isIgnored("bla", cfgExcludes("bla"))
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

func TestIsIgnoredIgnoresDeepFileSpec(t *testing.T) {
	ignored := isIgnored("foo/bar/snae.cfg", cfgExcludes("foo/bar/snae.cfg"))
	assert.True(t, ignored)
}

func TestIsIgnoreIgnoresAllInGitAsWeDontHaveDoubleStar(t *testing.T) {
	c := cfgExcludes("./.git/*/*")
	ignored := isIgnored("./.git/hooks/bas", c)
	assert.True(t, ignored)
}

func TestWatchDirRecursivelyWorksAsExpected(t *testing.T) {
	folders := []string{
		"qux",
		"foo/baem",
		"bar/ignored/snae",
		"node_modules/dep1",
		"target/classes",
	}
	cfg := cfgExcludes("bar/ignored", "node_modules", "target")
	tmpdir, err := ioutil.TempDir("", "gottaw-test-watch-dir-recursive")
	if err != nil {
		panic(err)
	}
	defer os.RemoveAll(tmpdir)
	cfg.WorkingDirectory = tmpdir

	for _, d := range folders {
		if err := os.MkdirAll(filepath.Join(tmpdir, d), 0755); err != nil {
			panic(err)
		}
	}
	trkr := NewTestTrkr()
	err = watchDirRecursive(tmpdir, trkr, cfg)
	assert.NoError(t, err, "Unexpected error while adding tracked folders")
	//fmt.Printf("%#v\n", trkr)
	exp := map[string]bool{
		filepath.Join(tmpdir, "qux"):               true,
		filepath.Join(tmpdir, "foo"):               true,
		filepath.Join(tmpdir, "foo/baem"):          true,
		filepath.Join(tmpdir, "bar"):               true,
		filepath.Join(tmpdir, "bar/ignored"):       false,
		filepath.Join(tmpdir, "bar/ignored/snae"):  false,
		filepath.Join(tmpdir, "node_modules"):      false,
		filepath.Join(tmpdir, "node_modules/dep1"): false,
		filepath.Join(tmpdir, "target"):            false,
		filepath.Join(tmpdir, "target/classes"):    false,
	}

	for dir, shouldBeTracked := range exp {
		assert.Equal(t, shouldBeTracked, trkr.IsTracked(dir), "wrong tracking status: %s", dir)
	}
}

type TestTrkr struct {
	w map[string]bool
}

func NewTestTrkr() Tracker {
	return &TestTrkr{
		w: make(map[string]bool),
	}
}

func (tt *TestTrkr) Tracked() []string {
	var keys []string
	for k := range tt.w {
		keys = append(keys, k)
	}
	return keys
}

func (tt *TestTrkr) Add(p string) error {
	tt.w[p] = true
	return nil
}

func (tt *TestTrkr) IsTracked(p string) bool {
	_, ok := tt.w[p]
	return ok
}

func (tt *TestTrkr) Remove(p string) {
	if _, ok := tt.w[p]; ok {
		delete(tt.w, p)
	}
}

func (tt *TestTrkr) Events() chan fsnotify.Event {
	return nil
}
func (tt *TestTrkr) Errors() chan error {
	return nil
}
func (tt *TestTrkr) Close() error {
	return nil
}

func cfgExcludes(excludes ...string) *config.Config {
	return &config.Config{
		Excludes: excludes,
	}
}
