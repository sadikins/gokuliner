package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BahanBaku struct {
	ID        string    `gorm:"primaryKey;type:varchar(36)" json:"id"`
	Nama      string    `gorm:"unique;not null;type:varchar(255)" json:"nama"`
	Satuan    string    `gorm:"not null;type:varchar(50)" json:"satuan"`
	HargaBeli float64   `gorm:"not null;type:decimal(18,2)" json:"harga_beli"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// BeforeCreate is a GORM hook to set UUID before creating a record
func (b *BahanBaku) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return
}