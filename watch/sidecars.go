package watch

import (
	c "github.com/makii42/gottaw/config"
	o "github.com/makii42/gottaw/output"
)

type SidecarRunner struct {
	log o.Logger
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
