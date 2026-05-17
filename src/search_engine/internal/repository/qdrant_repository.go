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
	Delete(ctx context.Context, collection_name string, ids []uint64) error

	// Поиск
	FindNearest(ctx context.Context, collection_name string, vector []float32) ([]*qdrant.ScoredPoint, error)
	//Search(ctx context.Context, collection_name string, vector []float32) ([]*qdrant.ScoredPoint, error)
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

func (qr *qdrantRepository) FindNearest(ctx context.Context, collection_name string, vector []float32) ([]*qdrant.ScoredPoint, error) {
	search_result, err := qr.client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: collection_name,
		Query:          qdrant.NewQuery(vector...),
		WithPayload:    qdrant.NewWithPayload(true),
		WithVectors:    qdrant.NewWithVectors(true),
	})

	if err != nil {
		return nil, err
	}

	return search_result, nil
}

func (qr *qdrantRepository) Delete(ctx context.Context, collection_name string, ids []uint64) error {
	return nil

}

func (qr *qdrantRepository) SearchWithFilter(ctx context.Context, collection_name string, vector []float32, filter *qdrant.Filter) ([]*qdrant.ScoredPoint, error) {
	search_result, err := qr.client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: collection_name,
		Query:          qdrant.NewQuery(vector...),
	})

	if err != nil {
		return nil, err
	}

	return search_result, nil
}

func (qr *qdrantRepository) Close() error {
	return qr.client.Close()
}
