package defaults

import (
	c "github.com/makii42/gottaw/config"
)

func NewWerckerDefault(util *defaultsUtil) *WerckerDefault {
	return &WerckerDefault{
		util: util,
	}
}

type WerckerDefault struct {
	util *defaultsUtil
}

func (w WerckerDefault) Name() string {
	return "Wercker"
}
func (w WerckerDefault) Test(dir string) bool {
	w.util.l.Tracef("testing for %s...\n", w.Name())
	return w.util.fileExists(dir, "wercker.yml")
}
func (w WerckerDefault) Config(dir string) *c.Config {
	if w.util.isExecutable("wercker") {

	}
	return nil
}
