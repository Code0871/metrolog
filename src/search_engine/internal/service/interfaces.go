package service

import (
	"context"
	"search_engine/internal/model"
)

type EmbeddingProvider interface {
	GetEmbeddingDense(ctx context.Context, texts []string) ([]model.DenseVector, error)
	GetEmbeddingSparse(ctx context.Context, texts []string) ([]model.SparseVector, error)
	GetEmbeddingLate(ctx context.Context, texts []string) ([]model.MultiVector, error)
}

type VectorSearchService interface {
	Upsert(ctx context.Context, collection_name string, point model.UpsertRequest) error
	UpsertBatch(ctx context.Context, collection_name string, points model.BatchUpsertRequest) error
	SearchNearestPoints(ctx context.Context, collection_name string, req model.SearchRequest) (model.SearchResponses, error)
	GetPointByIDBatch(ctx context.Context, collection_name string, ids model.BatchPointRequest) (model.BatchPointResponse, error)
}
