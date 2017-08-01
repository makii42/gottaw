package daemon

import (
	"context"
	"log"

	t "github.com/docker/docker/api/types"
	d "github.com/docker/docker/client"

	c "github.com/makii42/gottaw/config"
	o "github.com/makii42/gottaw/output"
)

var client *d.Client

func ensureDockerClient() (*d.Client, error) {
	cli, err := d.NewEnvClient()
	if err != nil {
		return nil, err
	}
	ping, err := cli.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	log.Printf("Connected to Docker API version: %s", ping.APIVersion)
	return cli, nil
}

// Sidecar describes a backend service for a build pipeline or server
type Sidecar interface {
	Start() error
	Stop() error
}

type SidecarRunner struct {
	dockerClient *d.Client
	log          o.Logger
	config       c.Sidecar
	containers   []t.Container
}

type sidecar struct {
	image string
}

// NewRunner returns a Runner that gives you control over all sidecars.
func NewRunner(l o.Logger, sidecarCfg map[string]c.Sidecar) (*SidecarRunner, error) {
	cli, err := ensureDockerClient()
	if err != nil {
		return nil, err
	}
	scr := SidecarRunner{
		dockerClient: cli,
	}

	for name, scconf := range sidecarCfg {
		l.Tracef("%s: %v", name, scconf)
	}

	return &scr, nil
}

func (sr *SidecarRunner) Start() error {

	return nil
}

func (sr *SidecarRunner) Stop() error {
	return nil
}
