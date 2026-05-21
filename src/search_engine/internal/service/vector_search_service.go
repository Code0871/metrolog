package service

import (
	"context"
	"fmt"
	"search_engine/internal/model"
	"search_engine/internal/repository"


	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

type SearchService struct {
	repo       repository.QdrantRepository
	embedded   EmbeddingProvider
	miinstance MiinstanceProvider
}

func NewSearchService(repo repository.QdrantRepository, embedded EmbeddingProvider, miinstance MiinstanceProvider) *SearchService {
	return &SearchService{
		repo:       repo,
		embedded:   embedded,
		miinstance: miinstance,
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

func (s *SearchService) HybridSearch(ctx context.Context, collection_name string, req model.BatchSearchRequest) ([]*model.Miinstance, error) {
	if len(req.Texts) == 0 {
		return []*model.Miinstance{}, nil
	}

	var all_miinstances []*model.Miinstance
	// search_results := make([][]*qdrant.ScoredPoint, len(req.Texts))

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

		var passports []string
		// извлекаем паспорта СИ из Payload
		for _, point := range points {
			if payload_map, ok := point.Payload["passport"]; ok {
				passports = append(passports, payload_map.GetStringValue())
			}
		}

		if len(passports) == 0 {
			continue
		}

		miinstances, err := s.miinstance.GetMiinstances(ctx, passports)
		if err != nil {
			return nil, fmt.Errorf("failed to get miinstances for query %d: %w", i, err)
		}

		all_miinstances = append(all_miinstances, miinstances...)
	}

	return all_miinstances, nil
}
