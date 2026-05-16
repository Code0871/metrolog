package repository

import (
	"context"
	"search_engine/config"

	"github.com/qdrant/go-client/qdrant"
)

type QdrantRepository interface {
	// CRUD операции
	Upsert(ctx context.Context, collection_name string, points []*qdrant.PointStruct) error
	Get(ctx context.Context, collection_name string, ids []uint64) ([]*qdrant.PointStruct, error)
	Delete(ctx context.Context, collection_name string, ids []uint64) error

	// Поиск
	Search(ctx context.Context, collection_name string, vector []float32, limit int) ([]*qdrant.ScoredPoint, error)
	SearchWithFilter(ctx context.Context, collection_name string, vector []float32, filter *qdrant.Filter) ([]*qdrant.ScoredPoint, error)

	Close() error
}

// qdrantRepository реализация интерфейса
type qdrantRepository struct {
	client *qdrant.Client
	config *config.Config
}

func NewQdrantRepository(client *qdrant.Client, config *config.Config) QdrantRepository {
	return &qdrantRepository{
		client: client,
		config: config,
	}
}

func (qr *qdrantRepository) Upsert(ctx context.Context, collection_name string, points []*qdrant.PointStruct) error {

}

func (qr *qdrantRepository) Get(ctx context.Context, collection_name string, ids []uint64) ([]*qdrant.PointStruct, error) {
	
}

func (qr *qdrantRepository) Delete(ctx context.Context, collection_name string, ids []uint64) error {

}

func (qr *qdrantRepository) Search(ctx context.Context, collection_name string, vector []float64) ([]*qdrant.ScoredPoint, error) {

}

func (qr *qdrantRepository) SearchWithFilter(ctx context.Context, collection_name string, vector []float32, filter *qdrant.Filter) ([]*qdrant.ScoredPoint, error) {

}

func (qr *qdrantRepository) Close() error {
	return qr.client.Close()
}
