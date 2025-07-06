package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Resep struct {
	ID          string          `gorm:"primaryKey;type:uuid" json:"id"`
	Nama        string          `gorm:"unique;not null;type:varchar(255)" json:"nama"`
	IsSubResep  bool            `gorm:"not null;default:false" json:"is_sub_resep"`
	JumlahPorsi float64 `gorm:"type:decimal(10,4);default:1.0" json:"jumlah_porsi"` // <<< UBAH TIPE INI
	Komponen    []ResepKomponen `gorm:"foreignKey:ResepID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"komponen,omitempty"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

func (resep *Resep) BeforeCreate(tx *gorm.DB) (err error) {
	if resep.ID == "" {
		resep.ID = uuid.New().String()
	}
	return
}