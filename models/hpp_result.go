// backend/models/hpp_result.go
package models

import (
	"time" // <<< Pastikan time diimport
)

type HPPResult struct {
    // GORM akan otomatis menambahkan ID sebagai primary key jika tidak ada
    // ID string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"` // Opsional, GORM bisa mengurus ini
    ResepID     string          `json:"resep_id"`
    ResepNama   string          `json:"resep_nama"`
    HPPPerUnit  float64 `json:"hpp_per_unit"`
    HPPPerPorsi float64 `json:"hpp_per_porsi"`
    CreatedAt   time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"` // <<< PASTIKAN INI ADA
    UpdatedAt   time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"` // <<< PASTIKAN INI ADA
}