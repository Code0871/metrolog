package model

// псевдонимы для удобства
type DenseVector []float32   // плотный вектор
type MultiVector [][]float32 // мульти-вектор

// Структура sparse-вектора
type SparseVector struct {
	Indices []uint32  `json:"indices"`
	Values  []float32 `json:"values"`
}

type QdrantPoint struct {
	Id      string                 `json:"id"`
	Vectors []float32              `json:"vectors"`
	Payload map[string]interface{} `json:"payload"`
}
