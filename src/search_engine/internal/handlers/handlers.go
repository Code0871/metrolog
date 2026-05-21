package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"search_engine/internal/model"
	"search_engine/internal/service"
)

type SearchHandler struct {
	searchService *service.SearchService
}

func NewSearchHandler(searchService *service.SearchService) *SearchHandler {
	return &SearchHandler{
		searchService: searchService,
	}
}

func (h *SearchHandler) HybridSearchHandler(w http.ResponseWriter, r *http.Request) {
	var req model.BatchSearchRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if len(req.Texts) == 0 {
		http.Error(w, "No texts provided", http.StatusBadRequest)
		return
	}

	results, err := h.searchService.HybridSearch(r.Context(), "miinstance_park", req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Search failed: %v", err), http.StatusInternalServerError)
		return
	}

	// Формируем ответ - results это плоский список []*model.Miinstance
	response := make([]map[string]interface{}, len(results))
	for i, miinstance := range results {
		response[i] = map[string]interface{}{
			"passport": miinstance.Passport, // пример поля
			// добавьте другие нужные поля из model.Miinstance
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    response,
		"count":   len(results),
	})
}

// GetPointsByIDHandler - получение точек по ID
func (h *SearchHandler) GetPointsByIDHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		IDs            []string `json:"ids"`
		CollectionName string   `json:"collection_name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	if len(req.IDs) == 0 {
		http.Error(w, "No IDs provided", http.StatusBadRequest)
		return
	}

	points, err := h.searchService.GetPointByIDBatch(r.Context(), req.CollectionName, req.IDs)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get points: %v", err), http.StatusInternalServerError)
		return
	}

	response := make([]map[string]interface{}, len(points))
	for i, point := range points {
		response[i] = map[string]interface{}{
			"id":      point.Id.GetUuid(),
			"payload": point.Payload,
			"vectors": point.Vectors,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    response,
	})
}

// HealthCheck - проверка работоспособности
func (h *SearchHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"service": "search-service",
	})
}

// InsertPointsHandler - вставляет точки через текст, паспорт, имя и т.д. (генерирует эмбеддинги)
func (h *SearchHandler) InsertPointsHandler(w http.ResponseWriter, r *http.Request) {
	var req model.BatchUpsertRequest

	// Декодируем запрос
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid request body: %v", err), http.StatusBadRequest)
		return
	}

	// Валидация
	if len(req.UpsertRequests) == 0 {
		http.Error(w, "No upsert requests provided", http.StatusBadRequest)
		return
	}

	// Выполняем вставку
	err := h.searchService.UpsertBatch(r.Context(), "miinstance_park", req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to insert points: %v", err), http.StatusInternalServerError)
		return
	}

	// Ответ об успехе
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Successfully inserted %d points", len(req.UpsertRequests)),
	})
}
