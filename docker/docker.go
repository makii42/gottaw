package docker

//go:generate mockgen -destination ./docker_mocks.go -package docker -source=docker.go

import (
	"io"
	"time"

	"context"

	"github.com/docker/docker/api/types"
	t "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
)

type (
	// Docker is a interface over our abstraction over docker
	Docker interface {
		EnsureImage(image string) error
		StartContainer() (Container, error)
	}
	// Client is our interface for the docker client to define what we use and to mock it.
	Client interface {
		ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error)
		ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error)
		ContainerStart(ctx context.Context, containerID string, options types.ContainerStartOptions) error
		ContainerStop(ctx context.Context, container string, timeout *time.Duration) error
		ContainerRemove(ctx context.Context, container string, options types.ContainerRemoveOptions) error
	}
	// Container is a handle over a docker container to manage the lifecycle of a service.
	Container interface {
		ID() string
		Name() string
		Start() error
		Stop() error
		Restart() error
	}

	dockerProxy struct {
		cli Client
	}

	cntnr struct {
		container t.Container
	}
)

func NewClient() (Client, error) {
	return client.NewEnvClient()
}
