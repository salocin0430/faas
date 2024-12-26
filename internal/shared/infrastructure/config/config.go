package config

import "os"

type Config struct {
	ServerAddress string
	NatsURL       string
	JWTSecret     string
	ConsumerKey   string
}

func LoadConfig() *Config {
	return &Config{
		ServerAddress: getEnvOrDefault("SERVER_ADDRESS", ":8080"),
		NatsURL:       getEnvOrDefault("NATS_URL", "nats://localhost:4222"),
		JWTSecret:     getEnvOrDefault("JWT_SECRET", "your-super-secret-key-for-development"),
		ConsumerKey:   getEnvOrDefault("CONSUMER_KEY", "faasapp-key"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
