package entity

import "time"

type Metric struct {
	ID          string
	WorkerID    string
	CPU         float64
	Memory      float64
	Executions  int64
	CollectedAt time.Time
}

type ScalingPolicy struct {
	MinWorkers      int
	MaxWorkers      int
	CPUThreshold    float64
	MemoryThreshold float64
	ScaleUpStep     int
	ScaleDownStep   int
	CooldownPeriod  time.Duration
}
