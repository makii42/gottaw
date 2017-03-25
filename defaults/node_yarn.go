package defaults

import (
	c "github.com/makii42/gottaw/config"
)

type NodeYarnDefault struct{}

func (d NodeYarnDefault) Name() string {
	return "NodeJS/yarn"
}
func (g NodeYarnDefault) Test(dir string) bool {
	return fileExists(dir, "package.json") &&
		isExecutable("node") &&
		isExecutable("yarn")
}
func (g NodeYarnDefault) Config() *c.Config {
	return &c.Config{
		Excludes: append(
			defaultExcludes,
			"node_modules",
			"lib",
			"yarn.lock",
		),
		Pipeline: []string{
			"yarn",
			"yarn test",
		},
	}
}
