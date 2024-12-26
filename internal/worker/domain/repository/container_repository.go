package repository

import "faas/internal/worker/domain/entity"

type ContainerRepository interface {
	RunContainer(container *entity.Container) error
	StopContainer(containerID string) error
	GetContainerStatus(containerID string) (string, error)
}
