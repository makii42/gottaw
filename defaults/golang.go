package defaults

import (
	c "github.com/makii42/gottaw/config"
)

type GolangDefault struct{}

func (g GolangDefault) Name() string {
	return "Golang"
}
func (g GolangDefault) Test(dir string) bool {
	return isExecutable("go") &&
		filesMatch(dir, "*.go")
}
func (g GolangDefault) Config() *c.Config {
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
