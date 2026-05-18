package repository

import (
	"context"
	"search_engine/config"

	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

type QdrantRepository interface {
	// CRUD операции
	Upsert(ctx context.Context, collection_name string, vector []float32, passport, name, mi_type string) error
	UpsertBatch(ctx context.Context, collection_name string, points []*qdrant.PointStruct) error
	GetPoints(ctx context.Context, collection_name string, ids []string) ([]*qdrant.RetrievedPoint, error)
	Delete(ctx context.Context, collection_name string, passport string) error

	// Поиск
	FindNearest(ctx context.Context, collection_name string, vector []float32) ([]*qdrant.ScoredPoint, error)
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

// функция добавления вектора в коллекцию
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

// Пакетная вставка точек
func (qr *qdrantRepository) UpsertBatch(ctx context.Context, collection_name string, points []*qdrant.PointStruct) error {
	batch_size := 250
	for i := 0; i < len(points); i += batch_size {
		end := i + batch_size
		if end > len(points) {
			end = len(points)
		}
		batch := points[i:end]

		upsert_request := &qdrant.UpsertPoints{
			CollectionName: collection_name,
			Points:         batch,
			Wait:           qdrant.PtrOf(false),
		}

		_, err := qr.client.Upsert(ctx, upsert_request)
		if err != nil {
			log.Printf("Ошибка вставки батча: [%d:%d]: %v", i, end, err)
			continue
		}

		log.Printf("Батч [%d:%d] успешно вставлен", i, end)
		time.Sleep(30 * time.Millisecond)
	}

	return nil
}

// получаем точку(и) по ID
func (qr *qdrantRepository) GetPoints(ctx context.Context, collection_name string, ids []string) ([]*qdrant.RetrievedPoint, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	point_ids := make([]*qdrant.PointId, 0, len(ids))

	for _, id := range ids {
		point_ids = append(point_ids, qdrant.NewID(id))
	}

	retrieved_points, err := qr.client.Get(ctx, &qdrant.GetPoints{
		CollectionName: collection_name,
		Ids:            point_ids,
		WithPayload:    qdrant.NewWithPayload(true),
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка получения точек: %w", err)
	}

	return retrieved_points, nil
}

// функция поиска ближайших векторов
func (qr *qdrantRepository) FindNearest(ctx context.Context, collection_name string, vector []float32) ([]*qdrant.ScoredPoint, error) {
	search_result, err := qr.client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: collection_name,
		Query:          qdrant.NewQuery(vector...),
		WithPayload:    qdrant.NewWithPayload(true),
		WithVectors:    qdrant.NewWithVectors(true),
		Limit:          qdrant.PtrOf(uint64(50)),
	})

	if err != nil {
		log.Printf("Ошибка поиска: %v", err)
		return nil, err
	}

	return search_result, nil
}

// Функция удаления точки по совпадению паспорта СИ в payload'е
func (qr *qdrantRepository) Delete(ctx context.Context, collection_name string, passport string) error {
	_, err := qr.client.Delete(ctx, &qdrant.DeletePoints{
		CollectionName: collection_name,
		Points: qdrant.NewPointsSelectorFilter(
			&qdrant.Filter{
				Must: []*qdrant.Condition{
					qdrant.NewMatch("passport", passport),
				},
			},
		),
	})

	if err != nil {
		log.Printf("Ошибка удаления: %v", err)
		return err
	}
	return nil
}

func (qr *qdrantRepository) SearchWithFilter(ctx context.Context, collection_name string, vector []float32, filter *qdrant.Filter) ([]*qdrant.ScoredPoint, error) {
	search_result, err := qr.client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: collection_name,
		Query:          qdrant.NewQuery(vector...),
		WithPayload:    qdrant.NewWithPayload(true),
		WithVectors:    qdrant.NewWithVectors(true),
		Filter:         filter,
		Limit:          qdrant.PtrOf(uint64(50)),
	})

	if err != nil {
		return nil, err
	}

	return search_result, nil
}

func (qr *qdrantRepository) Close() error {
	return qr.client.Close()
}
