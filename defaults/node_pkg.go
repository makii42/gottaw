package defaults

import (
	"fmt"
	"io/ioutil"

	"github.com/makii42/gottaw/config"

	yaml "gopkg.in/yaml.v2"
)

type NodePackage struct {
	Name    string            `json:"name"`
	Version string            `json:"version"`
	Scripts map[string]string `json:"scripts"`
}

func LoadNodePackage(pkgFile string) (*NodePackage, error) {
	source, err := ioutil.ReadFile(pkgFile)
	if err != nil {
		return nil, err
	}
	var pkg NodePackage
	if err := yaml.Unmarshal(source, &pkg); err != nil {
		return nil, err
	}
	return &pkg, nil
}

func (n *NodePackage) filterScripts() []string {
	res := []string{"install"}
	for _, script := range scripts {
		if _, ok := n.Scripts[script]; ok {
			res = append(res, script)
		}
	}
	return res
}

func (n *NodePackage) FillPipeline(c *config.Config, cmd string) {
	scripts := n.filterScripts()
	p := make([]string, len(scripts))
	for i, script := range scripts {
		p[i] = fmt.Sprintf("npm %s", script)
	}
	c.Pipeline = p
}

var scripts []string = []string{"test", "start"}
