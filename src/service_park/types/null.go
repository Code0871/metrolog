// ============================================================================
// Пакет types предоставляет кастомные типы для обработки NULL значений
//
// Назначение:
//   - NullInt32 - обработка NULL для int32
//   - NullBool - обработка NULL для bool
//   - NullString - обработка NULL для string
//   - NullTime - обработка NULL для time.Time
//
// Использование:
//   - В моделях вместо sql.NullInt32 используйте types.NullInt32
//   - JSON сериализация: NULL → null, число → число
//
// ============================================================================
package types

import (
	"database/sql"
	"encoding/json"
	"time"
)

// NullInt32 - тип для обработки NULL значений int32 в JSON
type NullInt32 struct {
	sql.NullInt32
}

// MarshalJSON - кастомная сериализация в JSON
func (v NullInt32) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Int32)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON - кастомная десериализация из JSON
func (v *NullInt32) UnmarshalJSON(data []byte) error {
	var num int32
	if err := json.Unmarshal(data, &num); err == nil {
		v.Int32 = num
		v.Valid = true
		return nil
	}

	var null *int32
	if err := json.Unmarshal(data, &null); err == nil && null == nil {
		v.Valid = false
		return nil
	}

	return nil
}

// NewNullInt32 - конструктор для создания NullInt32 с значением
func NewNullInt32(value int32) NullInt32 {
	return NullInt32{
		NullInt32: sql.NullInt32{
			Int32: value,
			Valid: true,
		},
	}
}

// NewNullInt32Null - конструктор для создания NULL значения
func NewNullInt32Null() NullInt32 {
	return NullInt32{
		NullInt32: sql.NullInt32{
			Int32: 0,
			Valid: false,
		},
	}
}

// NullBool - тип для обработки NULL значений bool в JSON
type NullBool struct {
	sql.NullBool
}

// MarshalJSON - кастомная сериализация в JSON
func (v NullBool) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Bool)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON - кастомная десериализация из JSON
func (v *NullBool) UnmarshalJSON(data []byte) error {
	var b bool
	if err := json.Unmarshal(data, &b); err == nil {
		v.Bool = b
		v.Valid = true
		return nil
	}

	var null *bool
	if err := json.Unmarshal(data, &null); err == nil && null == nil {
		v.Valid = false
		return nil
	}

	return nil
}

// NewNullBool - конструктор для создания NullBool с значением
func NewNullBool(value bool) NullBool {
	return NullBool{
		NullBool: sql.NullBool{
			Bool:  value,
			Valid: true,
		},
	}
}

// NullString - тип для обработки NULL значений string в JSON
type NullString struct {
	sql.NullString
}

// MarshalJSON - кастомная сериализация в JSON
func (v NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON - кастомная десериализация из JSON
func (v *NullString) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		v.String = s
		v.Valid = true
		return nil
	}

	var null *string
	if err := json.Unmarshal(data, &null); err == nil && null == nil {
		v.Valid = false
		return nil
	}

	return nil
}

// NewNullString - конструктор для создания NullString с значением
func NewNullString(value string) NullString {
	return NullString{
		NullString: sql.NullString{
			String: value,
			Valid:  true,
		},
	}
}

// NullTime - тип для обработки NULL значений time.Time в JSON
type NullTime struct {
	sql.NullTime
}

// MarshalJSON - кастомная сериализация в JSON
func (v NullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	}
	return json.Marshal(nil)
}

// UnmarshalJSON - кастомная десериализация из JSON
func (v *NullTime) UnmarshalJSON(data []byte) error {
	var t time.Time
	if err := json.Unmarshal(data, &t); err == nil {
		v.Time = t
		v.Valid = true
		return nil
	}

	var null *time.Time
	if err := json.Unmarshal(data, &null); err == nil && null == nil {
		v.Valid = false
		return nil
	}

	return nil
}

// NewNullTime - конструктор для создания NullTime с значением
func NewNullTime(value time.Time) NullTime {
	return NullTime{
		NullTime: sql.NullTime{
			Time:  value,
			Valid: true,
		},
	}
}
