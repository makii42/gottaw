package defaults

import (
	c "github.com/makii42/gottaw/config"
)

func NewNodeYarnDefault(util *defaultsUtil) *NodeYarnDefault {
	return &NodeYarnDefault{
		util: util,
	}
}

type NodeYarnDefault struct {
	util *defaultsUtil
}

func (ny NodeYarnDefault) Name() string {
	return "NodeJS/yarn"
}
func (ny NodeYarnDefault) Test(dir string) bool {
	ny.util.l.Tracef("testing for %s...\n", ny.Name())
	return ny.util.fileExists(dir, "package.json") &&
		ny.util.isExecutable("node") &&
		ny.util.isExecutable("yarn")
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
