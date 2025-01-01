package docker

import (
	"context"
	"fmt"
	"log"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

type DockerContainerManager struct {
	client  *client.Client
	network string
}

func NewContainerManager() (*DockerContainerManager, error) {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithVersion("1.46"),
	)
	if err != nil {
		return nil, err
	}
	return &DockerContainerManager{
		client:  cli,
		network: "faas_apisix", // Red del docker-compose
	}, nil
}

//const WORKER_START_INDEX = 100 // Mismo índice que usa el scaler

func (m *DockerContainerManager) RunContainer(name string, replicas int) error {
	// Configuración del contenedor
	config := &container.Config{
		Image: "faas-worker:latest",
		Env: []string{
			"NATS_URL=nats://nats:4222",
		},
	}

	// Configuración del host
	hostConfig := &container.HostConfig{
		NetworkMode: container.NetworkMode(m.network),
		// Agregar configuración para acceso a Docker socket si es necesario
		Mounts: []mount.Mount{
			{
				Type:   mount.TypeBind,
				Source: "/var/run/docker.sock",
				Target: "/var/run/docker.sock",
			},
		},
	}

	// Crear y arrancar contenedores
	currentCount, err := m.GetContainerCount(name)
	if err != nil {
		return err
	}

	for i := 0; i < replicas; i++ {
		containerName := fmt.Sprintf("%s-%d", name, currentCount+i+WORKER_START_INDEX)

		// Crear contenedor
		resp, err := m.client.ContainerCreate(
			context.Background(),
			config,
			hostConfig,
			nil,
			nil,
			containerName,
		)
		if err != nil {
			return fmt.Errorf("error creating container: %v", err)
		}

		// Iniciar contenedor
		if err := m.client.ContainerStart(context.Background(), resp.ID, container.StartOptions{}); err != nil {
			return fmt.Errorf("error starting container: %v", err)
		}
	}
	return nil
}

func (m *DockerContainerManager) StopContainer(name string) error {
	// Primero detener el contenedor
	timeoutSeconds := int(10)
	if err := m.client.ContainerStop(context.Background(), name, container.StopOptions{Timeout: &timeoutSeconds}); err != nil {
		return fmt.Errorf("failed to stop container %s: %v", name, err)
	}

	// Luego eliminar el contenedor
	if err := m.client.ContainerRemove(context.Background(), name, container.RemoveOptions{
		Force: true, // Forzar eliminación si es necesario
	}); err != nil {
		return fmt.Errorf("failed to remove container %s: %v", name, err)
	}

	log.Printf("Container %s stopped and removed", name)
	return nil
}

func (m *DockerContainerManager) GetContainerCount(name string) (int, error) {
	containers, err := m.client.ContainerList(context.Background(), container.ListOptions{
		All:     true,
		Filters: filters.NewArgs(filters.Arg("name", name)),
	})
	if err != nil {
		return 0, err
	}
	return len(containers), nil
}

// Método helper para obtener el cliente (necesario para el scaler)
func (m *DockerContainerManager) Client() *client.Client {
	return m.client
}
