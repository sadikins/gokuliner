package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type BahanBaku struct {
	ID              string          `gorm:"primaryKey;type:uuid" json:"id"`
	Nama            string          `gorm:"unique;not null;type:varchar(255)" json:"nama"`
	Kategori        string          `gorm:"type:varchar(100);not null" json:"kategori"`
	HargaBeli       decimal.Decimal `gorm:"type:decimal(18,4);not null" json:"harga_beli"`   // <<< UBAH TIPE INI
	SatuanBeli      string          `gorm:"type:varchar(50);not null" json:"satuan_beli"`
	NettoPerBeli    decimal.Decimal `gorm:"type:decimal(10,4);not null" json:"netto_per_beli"` // <<< UBAH TIPE INI
	SatuanPemakaian string          `gorm:"type:varchar(50);not null" json:"satuan_pemakaian"`
	Catatan         string          `gorm:"type:text" json:"catatan"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

// BeforeCreate is a GORM hook to set UUID before creating a record
func (b *BahanBaku) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return
}

