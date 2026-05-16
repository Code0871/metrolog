package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/qdrant/go-client/qdrant"
)

type QdrantConfig struct {
	QdrantHost string
	QdrantPort int
}

type CollectionConfig struct {
	CollectionName     string
	QdrantDistanceType qdrant.Distance
	QdrantVectorSize   int
}

type ModelConfig struct {
	ModelName string
}

type Config struct {
	QdrantConfigs    QdrantConfig
	CollectionConfig CollectionConfig
	ModelConfig      ModelConfig
}

func MustLoadConfig() *Config {

	if err := godotenv.Load("config/config.env"); err != nil {
		if !os.IsNotExist(err) {
			log.Printf("Warning: error loading .env file: %v", err)
		}
		panic("No .env file found, using system environment variables")
	}

	return &Config{
		QdrantConfigs: QdrantConfig{
			QdrantHost: getEnv("qdrant_host", "localhost"),
			QdrantPort: getEnvAsInt("qdrant_port_grpc", 6334),
		},
		CollectionConfig: CollectionConfig{
			CollectionName:     getEnv("collection_name", "miinstance_park"),
			QdrantDistanceType: QdrantDistanceConverter(getEnv("qdrant_distance_type", "Cosine")),
			QdrantVectorSize:   getEnvAsInt("qdrant_vector_size", 768),
		},
		ModelConfig: ModelConfig{
			ModelName: getEnv("model_from_hugging_face", "sentence-transformers/paraphrase-multilingual-MiniLM-L12-v2"),
		},
	}
}

// получение значений из .env
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// конвертация значений в int
func getEnvAsInt(key string, defaultValue int) int {
	if value, err := strconv.Atoi(os.Getenv(key)); err == nil {
		return value
	}
	return defaultValue
}

// Конвертер для типов дистанций Qdrant
func QdrantDistanceConverter(distance_type string) qdrant.Distance {
	switch strings.ToLower(distance_type) {
	case "cosine":
		return qdrant.Distance_Cosine
	case "dot":
		return qdrant.Distance_Dot
	case "manhattan":
		return qdrant.Distance_Manhattan
	case "euclid":
		return qdrant.Distance_Euclid
	default:
		panic("unknown distance type")
	}
}
