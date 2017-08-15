package docker

//go:generate mockgen -destination ./docker_mocks.go -package docker -source=docker.go

import (
	"io"

	"golang.org/x/net/context"

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
		ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, containerName string) (container.ContainerCreateCreatedBody, error)
		ImagePull(ctx context.Context, refStr string, options types.ImagePullOptions) (io.ReadCloser, error)
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

func newDockerClient() (Client, error) {
	return client.NewEnvClient()
}
