package service

import (
	"context"
	"fmt"
	_ "fmt"
	_ "log"
	_ "search_engine/config"
	"search_engine/internal/model"
	"search_engine/internal/repository"
	_ "time"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

type SearchService struct {
	repo     repository.QdrantRepository
	embedded EmbeddingProvider
}

func NewSearchService(repo repository.QdrantRepository, embedded EmbeddingProvider) *SearchService {
	return &SearchService{
		repo:     repo,
		embedded: embedded,
	}
}

func (s *SearchService) UpsertBatch(ctx context.Context, collection_name string, req model.BatchUpsertRequest) error {
	if len(req.UpsertRequests) == 0 {
		return nil
	}

	texts := make([]string, len(req.UpsertRequests))
	for i, r := range req.UpsertRequests {
		texts[i] = r.Text
	}

	// получаем эмбеддинги для записи
	dense, err := s.embedded.GetEmbeddingDense(ctx, texts)
	if err != nil {
		return err
	}
	sparse, err := s.embedded.GetEmbeddingSparse(ctx, texts)
	if err != nil {
		return err
	}
	multi, err := s.embedded.GetEmbeddingLate(ctx, texts)
	if err != nil {
		return err
	}

	points := make([]*qdrant.PointStruct, len(req.UpsertRequests))
	for i, r := range req.UpsertRequests {
		// Создаем срез DenseVector для multi-vector
		multi_vectors := make([]*qdrant.DenseVector, len(multi[i]))
		for j, vec := range multi[i] {
			multi_vectors[j] = &qdrant.DenseVector{Data: vec}
		}

		points[i] = &qdrant.PointStruct{
			Id: qdrant.NewID(uuid.New().String()),
			Vectors: &qdrant.Vectors{
				VectorsOptions: &qdrant.Vectors_Vectors{
					Vectors: &qdrant.NamedVectors{
						Vectors: map[string]*qdrant.Vector{
							"dense": {
								Vector: &qdrant.Vector_Dense{
									Dense: &qdrant.DenseVector{Data: dense[i]},
								},
							},
							"sparse": {
								Vector: &qdrant.Vector_Sparse{
									Sparse: &qdrant.SparseVector{
										Values:  sparse[i].Values,
										Indices: sparse[i].Indices,
									},
								},
							},
							"multi": {
								Vector: &qdrant.Vector_MultiDense{
									MultiDense: &qdrant.MultiDenseVector{
										Vectors: multi_vectors,
									},
								},
							},
						},
					},
				},
			},
			Payload: qdrant.NewValueMap(map[string]any{
				"passport": r.Passport,
				"name":     r.Name,
				"mi_type":  r.MiType,
			}),
		}
	}

	// Запись в Qdrant
	return s.repo.UpsertBatch(ctx, collection_name, points)
}

func (s *SearchService) GetPointByIDBatch(ctx context.Context, collection_name string, ids []string) ([]*qdrant.RetrievedPoint, error) {
	if len(ids) == 0 {
		return []*qdrant.RetrievedPoint{}, nil
	}

	return s.repo.GetPoints(ctx, collection_name, ids)
}

func (s *SearchService) DeletePoint(ctx context.Context, collection_name string, passport string) error {
	if passport == "" {
		return nil
	}

	return s.repo.Delete(ctx, collection_name, passport)
}

func (s *SearchService) HybridSearch(ctx context.Context, collection_name string, req model.BatchSearchRequest) ([][]*qdrant.ScoredPoint, error) {
	if len(req.Texts) == 0 {
		return [][]*qdrant.ScoredPoint{}, nil
	}

	results := make([][]*qdrant.ScoredPoint, len(req.Texts))

	for i, searchReq := range req.Texts {
		// Получаем эмбеддинги для одного текста
		dense, err := s.embedded.GetEmbeddingDense(ctx, []string{searchReq.Text})
		if err != nil {
			return nil, fmt.Errorf("dense embedding failed for query %d: %w", i, err)
		}

		sparse, err := s.embedded.GetEmbeddingSparse(ctx, []string{searchReq.Text})
		if err != nil {
			return nil, fmt.Errorf("sparse embedding failed for query %d: %w", i, err)
		}

		multi, err := s.embedded.GetEmbeddingLate(ctx, []string{searchReq.Text})
		if err != nil {
			return nil, fmt.Errorf("late embedding failed for query %d: %w", i, err)
		}

		// Извлекаем данные из sparse вектора
		indices := make([]uint32, len(sparse[0].Indices))
		for j, idx := range sparse[0].Indices {
			indices[j] = uint32(idx)
		}

		values := sparse[0].Values

		// Выполняем гибридный поиск для этого запроса
		points, err := s.repo.GetNearestPointsHybrid(
			ctx,
			collection_name,
			dense[0],
			indices,
			values,
			multi[0],
		)
		if err != nil {
			return nil, fmt.Errorf("hybrid search failed for query %d: %w", i, err)
		}

		results[i] = points
	}

	return results, nil
}
