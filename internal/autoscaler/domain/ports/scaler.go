package ports

type Scaler interface {
	ScaleUp(workers int) error
	ScaleDown(workers int) error
	GetCurrentWorkers() (int, error)
}
