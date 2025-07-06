package models

import (
	"time"

	"github.com/google/uuid" // Pastikan package ini diimpor
	"gorm.io/gorm"
)

type ResepKomponen struct {
	ID           string          `gorm:"primaryKey;type:uuid" json:"id"`
	ResepID      string          `gorm:"type:uuid;not null" json:"resep_id"`
	KomponenID   string          `gorm:"type:uuid;not null" json:"komponen_id"`
	Kuantitas    float64 		 `gorm:"type:decimal(10,4);not null" json:"kuantitas"` // <<< UBAH TIPE INI
	TipeKomponen string          `gorm:"type:varchar(50);not null" json:"tipe_komponen"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// BeforeCreate is a GORM hook to set UUID before creating a record for ResepKomponen
func (r *ResepKomponen) BeforeCreate(tx *gorm.DB) (err error) {
	if r.ID == "" {
		r.ID = uuid.New().String()
	}
	return
}