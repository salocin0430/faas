package docker

import (
	"faas/internal/autoscaler/domain/ports"
	"fmt"
	"log"
)

type DockerScaler struct {
	containerManager *DockerContainerManager
	serviceName      string
}

const WORKER_START_INDEX = 100 // Initial index for workers

func NewDockerScaler(serviceName string) (ports.Scaler, error) {
	manager, err := NewContainerManager()
	if err != nil {
		return nil, err
	}
	return &DockerScaler{
		containerManager: manager,
		serviceName:      serviceName,
	}, nil
}

func (s *DockerScaler) GetCurrentWorkers() (int, error) {
	count, err := s.containerManager.GetContainerCount(s.serviceName)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (s *DockerScaler) ScaleUp(workers int) error {
	return s.scale(workers, true)
}

func (s *DockerScaler) ScaleDown(workers int) error {
	return s.scale(workers, false)
}

func (s *DockerScaler) scale(delta int, up bool) error {
	current, err := s.GetCurrentWorkers()
	if err != nil {
		return err
	}

	target := current
	if up {
		target += delta
	} else {
		target -= delta
		if target < 2 {
			target = 2
		}
	}

	log.Printf("Scaling %s: Current=%d, Target=%d",
		map[bool]string{true: "UP", false: "DOWN"}[up],
		current, target)

	// Stop existing containers if needed
	if !up {
		for i := current - 1; i >= target; i-- {
			containerName := s.getContainerName(i)
			log.Printf("Stopping container: %s", containerName)
			if err := s.containerManager.StopContainer(containerName); err != nil {
				return fmt.Errorf("failed to stop container %s: %v", containerName, err)
			}
		}
	}

	// Create new containers if needed
	if target > current {
		return s.containerManager.RunContainer(s.serviceName, target-current)
	}

	return nil
}

func (s *DockerScaler) getContainerName(index int) string {
	return fmt.Sprintf("%s-%d", s.serviceName, index+WORKER_START_INDEX)
}
