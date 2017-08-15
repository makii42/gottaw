package daemon

import (
	"bytes"
	"context"
	"time"

	t "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/makii42/gottaw/docker"
)

var (
	client  docker.Client
	timeout *time.Duration
)

func init() {
	tmpTimeout, err := time.ParseDuration("30s")
	if err != nil {
		panic(err)
	}
	timeout = &tmpTimeout
}

// Sidecar describes a backend service for a build pipeline or server
type (
	Sidecar interface {
		Start(ctx *context.Context) error
		Stop(ctx *context.Context) error
		Reload(ctx *context.Context) error
	}

	Docker interface {
	}

	sidecar struct {
		name        string
		image       string
		environment map[string]string
		logs        *bytes.Buffer
		containerID string
	}
)

func (sc *sidecar) Start(ctx context.Context) error {
	container, err := client.ContainerCreate(
		ctx,
		&container.Config{
			Image:        sc.image,
			AttachStdout: true,
			AttachStderr: true,
		},
		&container.HostConfig{},
		nil,
		sc.name,
	)
	if err != nil {
		return err
	}
	sc.containerID = container.ID
	err = client.ContainerStart(
		ctx,
		container.ID,
		t.ContainerStartOptions{},
	)
	if err != nil {
		return err
	}
	return nil
}

func (sc *sidecar) Stop(ctx context.Context) error {
	if err := client.ContainerStop(ctx, sc.containerID, timeout); err != nil {
		return err
	}
	if err := client.ContainerRemove(ctx, sc.containerID, t.ContainerRemoveOptions{}); err != nil {
		return err
	}
	sc.containerID = ""
	return nil
}
