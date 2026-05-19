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
	UpsertBatch(ctx context.Context, collection_name string, points model.BatchUpsertRequest) error
	GetPointByIDBatch(ctx context.Context, collection_name string, ids model.BatchPointRequest) (model.BatchPointResponse, error)
	DeletePoint(ctx context.Context, collection_name string, passport string) error
	HybridSearch(ctx context.Context, collecction_name string, req model.BatchSearchRequest) (model.SearchResponses, error)
}
