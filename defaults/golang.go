package defaults

import (
	c "github.com/makii42/gottaw/config"
)

func NewGolangDefault(util *defaultsUtil) *GolangDefault {
	return &GolangDefault{
		util: util,
	}
}

type GolangDefault struct {
	util *defaultsUtil
}

func (g GolangDefault) Name() string {
	return "Golang"
}

func (g GolangDefault) Test(dir string) bool {
	g.util.l.Tracef("testing for %s...\n", g.Name())
	return g.util.filesMatch(dir, "*.go") && g.util.isExecutable("go")
}

func (g GolangDefault) Config(dir string) *c.Config {
	return &c.Config{
		Excludes: append(
			defaultExcludes,
			"*-go-tmp-umask",
		),
		Pipeline: []string{
			"go get -v .",
			"go build -v .",
			"go test -v ./...",
		},
	}
}
