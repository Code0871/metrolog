package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"search_engine/config"
	"search_engine/internal/model"
)

type EmbeddingService struct {
	client *http.Client
	apiURL string
}

func NewEmbeddingService() *EmbeddingService {
	cfg := config.MustLoadConfig()
	return &EmbeddingService{
		client: &http.Client{},
		apiURL: cfg.EmbeddingServiceConfigs.EmbeddingServiceURL(),
	}
}

func (es *EmbeddingService) GetEmbeddingDense(ctx context.Context, texts []string) ([]model.DenseVector, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	var result struct {
		Embeddings []model.DenseVector `json:"embeddings"`
	}

	err := es.doRequest(ctx, "/embed/dense", texts, &result)
	if err != nil {
		return nil, err
	}

	return result.Embeddings, nil
}

func (es *EmbeddingService) GetEmbeddingSparse(ctx context.Context, texts []string) ([]model.SparseVector, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	var result struct {
		Embeddings []model.SparseVector `json:"sparse_embeddings"`
	}

	err := es.doRequest(ctx, "/embed/sparse", texts, &result)
	if err != nil {
		return nil, err
	}

	return result.Embeddings, nil
}

func (es *EmbeddingService) GetEmbeddingLate(ctx context.Context, texts []string) ([]model.MultiVector, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	var result struct {
		MultiVectors []model.MultiVector `json:"multi_vectors"`
	}

	err := es.doRequest(ctx, "/embed/late", texts, &result)
	if err != nil {
		return nil, err
	}

	return result.MultiVectors, nil
}

func (es *EmbeddingService) doRequest(ctx context.Context, endpoint string, texts []string, result interface{}) error {
	requestBody := map[string]interface{}{
		"texts": texts,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	url := es.apiURL + endpoint
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := es.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	return nil
}
