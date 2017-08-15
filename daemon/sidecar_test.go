package daemon

import (
	"os"
	tt "testing"

	"github.com/stretchr/testify/assert"

	c "github.com/makii42/gottaw/config"
	"github.com/makii42/gottaw/output"
)

var (
	logger output.Logger
)

func TestMain(m *tt.M) {
	cfg := &c.Config{}
	l, _ := output.NewLog(cfg)
	logger = l
	os.Exit(m.Run())
}

func TestSidecarRunnerCreationWithEmptyList(t *tt.T) {
	sidecarCfg := make(map[string]c.Sidecar)
	r, err := NewRunner(logger, sidecarCfg)
	assert.Nil(t, err)
	assert.NotNil(t, r)
}
func TestSidecarRunnerCreationWithSideCar(t *tt.T) {
	sidecarCfg := map[string]c.Sidecar{
		"foo": c.Sidecar{},
	}
	r, err := NewRunner(logger, sidecarCfg)
	assert.Nil(t, err)
	assert.NotNil(t, r)
}
