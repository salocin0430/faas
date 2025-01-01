package ports

type ContainerManager interface {
	RunContainer(name string, replicas int) error
	StopContainer(name string) error
	GetContainerCount(name string) (int, error)
}
