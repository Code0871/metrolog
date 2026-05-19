package model

// псевдонимы для удобства
type DenseVector []float32   // плотный вектор
type MultiVector [][]float32 // мульти-вектор вектор

// Структура sparse-вектора
type SparseVector struct {
	Indices []uint32  `json:"indices"`
	Values  []float32 `json:"values"`
}
