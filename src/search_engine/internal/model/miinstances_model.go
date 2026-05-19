package model

// Структура запроса на поиск
type SearchRequest struct {
	Text string `json:"text" validate:"required"`
}

type UpsertRequest struct {
	Text     string `json:"text" validate:"required"`
	Passport string `json:"passport" validate:"required"`
	Name     string `json:"name" validate:"required"`
	MiType   string `json:"mi_type" validate:"required"`
}

// Структура результата поиска
type SearchResult struct {
	Id      string                 `json:"id"`
	Score   float32                `json:"score"`
	Payload map[string]interface{} `json:"payload"`
}

// Струкутра ответа на запрос поиска
type SearchResponses struct {
	Responses []SearchResult `json:"responses"`
}

// Структура запроса на пакетный поиск
type BatchSearchRequest struct {
	Texts []SearchRequest `json:"texts" validate:"required"`
}

// Структура запроса на пакетную вставку
type BatchUpsertRequest struct {
	UpsertRequests []UpsertRequest `json:"upsert_requests" validate:"required"`
}

type PointRequest struct {
	ID string `json:"id"`
}

type PointResponse struct {
	ID      string                 `json:"id"`
	Payload map[string]interface{} `json:"payload"`
	Vector  []float32              `json:"vector,omitempty"` // если нужно
}

type BatchPointRequest struct {
	Points []PointRequest `json:"points"`
}

type BatchPointResponse struct {
	Points []PointResponse `json:"points"`
}
