package config

import (
	"fmt"
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

type SearchParamsConfig struct {
	QdrantDenseCount  int
	QdrantSparseCount int
	QdrantMultiCount  int
}

type CollectionConfig struct {
	CollectionName     string
	QdrantDistanceType qdrant.Distance
	QdrantVectorSize   int
}

type EmbeddingServiceConfig struct {
	ServiceHost string
	ServicePort int
}

type MiinstanceServiceConfig struct {
	ServiceHost string
	ServicePort int
}

type Config struct {
	QdrantConfigs            QdrantConfig
	CollectionConfigs        CollectionConfig
	EmbeddingServiceConfigs  EmbeddingServiceConfig
	MiinstanceServiceConfigs MiinstanceServiceConfig
	SearchParamsConfigs      SearchParamsConfig
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
		CollectionConfigs: CollectionConfig{
			CollectionName:     getEnv("collection_name", "miinstance_park"),
			QdrantDistanceType: QdrantDistanceConverter(getEnv("qdrant_distance_type", "Cosine")),
			QdrantVectorSize:   getEnvAsInt("qdrant_vector_size", 768),
		},
		EmbeddingServiceConfigs: EmbeddingServiceConfig{
			ServiceHost: getEnv("HOST", "localhost"),
			ServicePort: getEnvAsInt("PORT", 8000),
		},
		MiinstanceServiceConfigs: MiinstanceServiceConfig{
			ServiceHost: getEnv("miinstance_host", "localhost"),
			ServicePort: getEnvAsInt("miinstance_port", 8080),
		},
		SearchParamsConfigs: SearchParamsConfig{
			QdrantDenseCount:  getEnvAsInt("dense_count", 25),
			QdrantSparseCount: getEnvAsInt("sparse_count", 25),
			QdrantMultiCount:  getEnvAsInt("multi_count", 25),
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

func (c *EmbeddingServiceConfig) EmbeddingServiceURL() string {
	return fmt.Sprintf("http://%s:%d", c.ServiceHost, c.ServicePort)
}

func (c *MiinstanceServiceConfig) MiinstanceServiceURL() string {
	return fmt.Sprintf("http://%s:%d", c.ServiceHost, c.ServicePort)
}
