package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type QdrantConfig struct {
	QdrantHost string
	QdrantPort string
}

type Config struct {
	QdrantConfigs QdrantConfig
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	return &Config{
		QdrantConfigs: QdrantConfig{
			QdrantHost: getEnv("QDRANT_HOST", "localhost"),
			QdrantPort: getEnv("QDRANT_PORT", "6334"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
