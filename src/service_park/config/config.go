package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string // Database host
	DBPort     string // Database port
	DBUser     string // Database user
	DBPassword string // Database password
	DBName     string // Database name
	ServerPort string // Server port
	GinMode    string // Gin mode release, debug, test
}

// LoadConfig загружает конфигурацию из .env файла
func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Отладка: выводим все переменные
	log.Println("Loading config...")
	log.Printf("POSTGRES_HOST=%s", os.Getenv("POSTGRES_HOST"))
	log.Printf("POSTGRES_USER=%s", os.Getenv("POSTGRES_USER"))
	log.Printf("POSTGRES_PASSWORD=%s", os.Getenv("POSTGRES_PASSWORD"))
	log.Printf("POSTGRES_DB=%s", os.Getenv("POSTGRES_DB"))
	log.Printf("SERVER_PORT=%s", os.Getenv("SERVER_PORT"))
	log.Printf("GIN_MODE=%s", os.Getenv("GIN_MODE"))

	return &Config{
		DBHost:     getEnv("POSTGRES_HOST", "localhost"),
		DBPort:     getEnv("POSTGRES_PORT", "5432"),
		DBUser:     getEnv("POSTGRES_USER", ""),
		DBPassword: getEnv("POSTGRES_PASSWORD", ""),
		DBName:     getEnv("POSTGRES_DB", ""),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		GinMode:    getEnv("GIN_MODE", "debug"),
	}
}

// getEnv возвращает значение переменной или default, если не задана
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
