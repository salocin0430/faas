package ports

type StreamMonitor interface {
	GetPendingMessages() (int, error)
}
