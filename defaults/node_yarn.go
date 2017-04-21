package defaults

import (
	"path/filepath"

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

func (g NodeYarnDefault) Config(dir string) *c.Config {
	config := c.Config{
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
	if nodePkg, err := LoadNodePackage(filepath.Join(dir, "package.json")); err == nil {
		nodePkg.FillPipeline(&config, "yarn")
	}
	return &config
}
