package docker

import (
	"context"
	"faas/internal/worker/domain/entity"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

type DockerContainerManager struct {
	client *client.Client
}

func NewDockerContainerManager() (*DockerContainerManager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return nil, err
	}
	return &DockerContainerManager{client: cli}, nil
}

func (m *DockerContainerManager) RunContainer(container *entity.Container) error {
	ctx := context.Background()

	resp, err := m.client.ContainerCreate(ctx,
		&container.Config{
			Image:     container.ImageURL,
			Resources: container.Resources,
		},
		nil,
		nil,
		nil,
		container.ID,
	)
	if err != nil {
		return err
	}

	return m.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
}
