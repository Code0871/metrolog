package main

import (
	"fmt"
	"log"
	"net/http"

	conf "search_engine/config"
	"search_engine/internal/handlers"
	"search_engine/internal/repository"
	"search_engine/internal/service"
	setup "search_engine/internal/setup"
)

func main() {
	fmt.Println("Load Configs")
	cfg := conf.MustLoadConfig()

	fmt.Println("Init Qdrant client")
	client := setup.MustInitQdrantСlient(
		cfg.QdrantConfigs.QdrantHost,
		cfg.QdrantConfigs.QdrantPort,
	)

	fmt.Println("Init Qdrant collection")
	collection_name := cfg.CollectionConfigs.CollectionName
	vector_size := uint64(cfg.CollectionConfigs.QdrantVectorSize)
	distance_type := cfg.CollectionConfigs.QdrantDistanceType

	setup.MustInitQdrantCollection(client, collection_name, vector_size, distance_type)

	// Инициализация репозитория
	fmt.Println("Init Qdrant repository")
	qdrantRepo := repository.NewQdrantRepository(client)

	// Инициализация сервиса эмбеддингов
	fmt.Println("Init embedding service")
	embeddingService := service.NewEmbeddingService()

	// Инициализация сервиса поиска
	fmt.Println("Init search service")
	searchService := service.NewSearchService(qdrantRepo, embeddingService)

	// Инициализация хендлеров
	fmt.Println("Init handlers")
	searchHandler := handlers.NewSearchHandler(searchService)

	// Регистрация роутов
	fmt.Println("Registering routes")
	http.HandleFunc("/api/v1/upsert/batch", searchHandler.InsertPointsHandler)
	http.HandleFunc("/api/v1/search/hybrid", searchHandler.HybridSearchHandler)
	http.HandleFunc("/api/v1/search/points", searchHandler.GetPointsByIDHandler)
	http.HandleFunc("/health", searchHandler.HealthCheck)

	// Запуск сервера
	port := ":8080"
	fmt.Printf("Starting search service on http://localhost%s\n", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
