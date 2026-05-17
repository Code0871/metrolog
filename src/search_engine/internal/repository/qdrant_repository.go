package repository

import (
	"context"
	"search_engine/config"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

// TODO: Сделать функцию для обновления коллекции

type QdrantRepository interface {
	// CRUD операции
	Upsert(ctx context.Context, collection_name string, vector []float32, passport, name, mi_type string) error
	Get(ctx context.Context, collection_name string, ids []uint64) ([]*qdrant.PointStruct, error)
	Delete(ctx context.Context, collection_name string, ids []uint64) error

	// Поиск
	Search(ctx context.Context, collection_name string, vector []float32) ([]*qdrant.ScoredPoint, error)
	SearchWithFilter(ctx context.Context, collection_name string, vector []float32, filter *qdrant.Filter) ([]*qdrant.ScoredPoint, error)

	Close() error
}

// qdrantRepository реализация интерфейса
type qdrantRepository struct {
	client *qdrant.Client
}

func NewQdrantRepository(client *qdrant.Client, config *config.Config) QdrantRepository {
	return &qdrantRepository{
		client: client,
	}
}

func (qr *qdrantRepository) Upsert(ctx context.Context, collection_name string, vector []float32, passport, name, mi_type string) error {
	_, err := qr.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collection_name,
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewID(uuid.New().String()),
				Vectors: &qdrant.Vectors{VectorsOptions: &qdrant.Vectors_Vector{Vector: &qdrant.Vector{Data: vector}}},
				Payload: qdrant.NewValueMap(map[string]any{
					"passport": passport,
					"name":     name,
					"mi_type":  mi_type,
				}),
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (qr *qdrantRepository) Get(ctx context.Context, collection_name string, ids []uint64) ([]*qdrant.PointStruct, error) {
	return nil, nil
}

func (qr *qdrantRepository) Delete(ctx context.Context, collection_name string, ids []uint64) error {
	return nil

}

func (qr *qdrantRepository) Search(ctx context.Context, collection_name string, vector []float32) ([]*qdrant.ScoredPoint, error) {
	return nil, nil
}

func (qr *qdrantRepository) SearchWithFilter(ctx context.Context, collection_name string, vector []float32, filter *qdrant.Filter) ([]*qdrant.ScoredPoint, error) {
	return nil, nil
}

func (qr *qdrantRepository) Close() error {
	return qr.client.Close()
}
