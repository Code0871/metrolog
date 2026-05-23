package model

import (
	"time"
)

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

// Структура запроса на пакетный поиск
type BatchSearchRequest struct {
	Texts []SearchRequest `json:"texts" validate:"required"`
}

// Структура запроса на пакетную вставку
type BatchUpsertRequest struct {
	UpsertRequests []UpsertRequest `json:"upsert_requests" validate:"required"`
}

type Miinstance struct {
	Passport          *string    `json:"passport"`
	Name              *string    `json:"name"`
	Type              *string    `json:"type"`
	StateCondition    *string    `json:"state_condition"`
	TechCondition     *string    `json:"tech_condition"`
	IssueDate         *time.Time `json:"issue_date,omitempty"`
	CommissioningDate *time.Time `json:"commissioning_date,omitempty"`
	IsFit             *bool      `json:"is_fit"`
	MPI               *int32     `json:"mpi"`
}
