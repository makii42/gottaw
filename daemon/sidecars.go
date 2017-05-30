package daemon

import (
	"fmt"

	d "github.com/docker/docker/client"

	c "github.com/makii42/gottaw/config"
	o "github.com/makii42/gottaw/output"
)

var client *d.Client

func init() {
	c, err := d.NewEnvClient()
	if err != nil {
		panic(fmt.Errorf("cannot create docker client - please install docker before using sidecars"))
	}
	client = c
}

// Sidecar describes a backend service for a build pipeline or server
type Sidecar interface {
	Start() error
	Stop() error
}

type SidecarRunner struct {
	dockerClient d.Client
	log          o.Logger
	config       c.Sidecar
}

type sidecar struct {
	image string
}

func NewSidecarRunner(l o.Logger, sidecarCfg map[string]c.Sidecar) (*SidecarRunner, error) {
	scr := SidecarRunner{}

	for name, scconf := range sidecarCfg {
		l.Tracef("%s: %v", name, scconf)
	}

	return &scr, nil
}

func (scr *SidecarRunner) Reload(sidcarCfg map[string]c.Sidecar) error {
	return nil
}
