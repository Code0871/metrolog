package service

import (
	"context"
	"search_engine/internal/model"

	"github.com/qdrant/go-client/qdrant"
)

type EmbeddingProvider interface {
	GetEmbeddingDense(ctx context.Context, texts []string) ([]model.DenseVector, error)
	GetEmbeddingSparse(ctx context.Context, texts []string) ([]model.SparseVector, error)
	GetEmbeddingLate(ctx context.Context, texts []string) ([]model.MultiVector, error)
}

type VectorSearchService interface {
	UpsertBatch(ctx context.Context, collection_name string, points model.BatchUpsertRequest) error
	GetPointByIDBatch(ctx context.Context, collection_name string, ids []string) ([]*qdrant.RetrievedPoint, error)
	DeletePoint(ctx context.Context, collection_name string, passport string) error
	HybridSearch(ctx context.Context, collection_name string, req model.BatchSearchRequest) ([][]*qdrant.ScoredPoint, error)
}
