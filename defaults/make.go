package defaults

import (
	c "github.com/makii42/gottaw/config"
)

func NewMakeDefault(util *defaultsUtil) *MakeDefault {
	return &MakeDefault{
		util: util,
	}
}

type MakeDefault struct {
	util *defaultsUtil
}

func (g MakeDefault) Name() string {
	return "make"
}
func (g MakeDefault) Test(dir string) bool {
	g.util.l.Tracef("testing for %s...\n", g.Name())
	return g.util.filesMatch(dir, "Makefile") && g.util.isExecutable("make")

}
func (g MakeDefault) Config(dir string) *c.Config {
	return &c.Config{
		Excludes: append(
			defaultExcludes,
			"*-go-tmp-umask",
		),
		Pipeline: []string{
			"make",
		},
	}
}
