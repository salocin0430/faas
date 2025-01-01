package config

type AutoscalerConfig struct {
	MinWorkers         int    `env:"MIN_WORKERS" envDefault:"2"`
	MaxWorkers         int    `env:"MAX_WORKERS" envDefault:"10"`
	ScaleUpThreshold   int    `env:"SCALE_UP_THRESHOLD" envDefault:"3"`
	ScaleDownThreshold int    `env:"SCALE_DOWN_THRESHOLD" envDefault:"2"`
	CheckInterval      string `env:"CHECK_INTERVAL" envDefault:"30s"`
	CooldownPeriod     string `env:"COOLDOWN_PERIOD" envDefault:"30s"`
}
