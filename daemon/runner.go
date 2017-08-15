package daemon

import (
	"bytes"
	"context"

	c "github.com/makii42/gottaw/config"
	"github.com/makii42/gottaw/docker"
	o "github.com/makii42/gottaw/output"
)

type (
	sidecarRunner struct {
		dockerClient docker.Client
		log          o.Logger
		config       map[string]c.Sidecar
		sides        []*sidecar
	}
	Runner interface {
		Start(ctx context.Context) error
		Stop(ctx context.Context) error
	}
)

// NewRunner returns a Runner that gives you control over all sidecars.
func NewRunner(l o.Logger, sidecarCfg map[string]c.Sidecar) (Runner, error) {
	cli, err := docker.NewClient()
	if err != nil {
		return nil, err
	}
	scr := sidecarRunner{
		dockerClient: cli,
	}
	var sidecars []*sidecar
	for name, scconf := range sidecarCfg {
		l.Tracef("%s: %v", name, scconf)
		sidecar := &sidecar{
			name:        name,
			image:       scconf.Image,
			environment: scconf.Environment,
			logs:        &bytes.Buffer{},
		}
		sidecars = append(sidecars, sidecar)
	}
	scr.sides = sidecars

	return &scr, nil
}

func (sr *sidecarRunner) Start(ctx context.Context) error {
	for _, side := range sr.sides {
		if err := side.Start(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (sr *sidecarRunner) Stop(ctx context.Context) error {
	for _, side := range sr.sides {
		if err := side.Stop(ctx); err != nil {
			return err
		}
	}
	return nil
}
