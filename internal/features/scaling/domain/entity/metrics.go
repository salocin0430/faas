package entity

// Metrics representa las métricas del sistema
type Metrics struct {
	QueueLength    int     `json:"queue_length"`
	ActiveWorkers  int     `json:"active_workers"`
	ProcessingTime float64 `json:"processing_time"`
	Timestamp      int64   `json:"timestamp"`
}

// MetricsCollector define la interfaz para recolectar métricas
type MetricsCollector interface {
	Collect() Metrics
	PublishMetrics(metrics Metrics) error
}
