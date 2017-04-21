package defaults

import (
	"path/filepath"

	c "github.com/makii42/gottaw/config"
)

func NewNodeNpmDefault(defUtil *defaultsUtil) *NodeNpmDefault {
	return &NodeNpmDefault{util: defUtil}
}

type NodeNpmDefault struct {
	util *defaultsUtil
}

func (nn NodeNpmDefault) Name() string {
	return "NodeJS/npm"
}

func (nn NodeNpmDefault) Test(dir string) bool {
	nn.util.l.Tracef("testing for %s...\n", nn.Name())
	return nn.util.fileExists(dir, "package.json") &&
		nn.util.isExecutable("node") &&
		nn.util.isExecutable("npm")
}

func (nn NodeNpmDefault) Config(dir string) *c.Config {
	config := c.Config{
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
	if nodePkg, err := LoadNodePackage(filepath.Join(dir, "package.json")); err == nil {
		nodePkg.FillPipeline(&config, "npm")
	}
	return &config
}
