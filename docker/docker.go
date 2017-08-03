package docker

import (
	t "github.com/docker/docker/api/types"
)

type (
	// Docker is a mockable interface over our glue code towards docker itself
	// to enable testability
	Docker interface {
		EnsureImage(image string) error
		StartContainer() (Container, error)
	}
	// Container is a handle over a docker container to manage the lifecycle of a service.
	Container interface {
		ID()
		Name()
		Start()
		Stop()
		Restart()
	}

	container struct {
		container t.Container
	}
)
