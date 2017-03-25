package defaults

import (
	c "github.com/makii42/gottaw/config"
)

type NodeNpmDefault struct{}

func (g NodeNpmDefault) Name() string {
	return "NodeJS/npm"
}
func (g NodeNpmDefault) Test(dir string) bool {
	return fileExists(dir, "package.json") &&
		isExecutable("node") &&
		isExecutable("npm")
}
func (g NodeNpmDefault) Config() *c.Config {
	return &c.Config{
		Excludes: append(
			defaultExcludes,
			"node_modules",
			"lib",
		),
		Pipeline: []string{
			"npm install",
			"npm test",
		},
	}
}
