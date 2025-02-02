package config

import "os"

type Config struct {
	ServerAddress           string
	NatsURL                 string
	JWTSecret               string
	ConsumerKey             string
	MaxConcurrentExecutions string
	APIBaseURL              string
	NetworkName             string
}

func LoadConfig() *Config {
	return &Config{
		ServerAddress:           getEnvOrDefault("SERVER_ADDRESS", ":8080"),
		NatsURL:                 getEnvOrDefault("NATS_URL", "nats://localhost:4222"),
		JWTSecret:               getEnvOrDefault("JWT_SECRET", "your-super-secret-key-for-development"),
		ConsumerKey:             getEnvOrDefault("CONSUMER_KEY", "faasapp-key"),
		MaxConcurrentExecutions: getEnvOrDefault("MAX_CONCURRENT_EXECUTIONS", "10"),
		APIBaseURL:              getEnvOrDefault("API_BASE_URL", "http://api:8080/api/function-objects"),
		NetworkName:             getEnvOrDefault("NETWORK_NAME", "apisix"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
