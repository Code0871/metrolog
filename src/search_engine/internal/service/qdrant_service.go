package service

import (
	_ "context"
	_ "fmt"
	_ "log"
	_ "search_engine/config"
	"search_engine/internal/repository"
	_ "time"

	_ "github.com/qdrant/go-client/qdrant"
)

type SearchService struct {
	repo     repository.QdrantRepository
	embedded EmbeddingProvider
}
