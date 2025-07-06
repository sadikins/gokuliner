package models

import (
	"time"
	// "github.com/shopspring/decimal" // <<< PASTIKAN INI DIHAPUS
	"gorm.io/gorm"
)

type ProgramPromo struct {
    ID                  string          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
    NamaPromo           string          `gorm:"unique;not null;type:varchar(255)" json:"nama_promo"`
    Channel             string          `json:"channel"`
    JenisDiskon         string          `json:"jenis_diskon"`
    BesarDiskon         float64         `json:"besar_diskon"` // Menggunakan float64
    MinBelanja          float64         `json:"min_belanja"`  // Menggunakan float64
    MaksimalPotongan    float64         `json:"maksimal_potongan"` // Menggunakan float64
    DitanggungMerchantPersen float64         `json:"ditanggung_merchant_persen"` // Menggunakan float64
    Catatan             string          `json:"catatan"`

    CreatedAt           time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt           time.Time       `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// GORM hook: BeforeCreate untuk memastikan nilai default
func (p *ProgramPromo) BeforeCreate(tx *gorm.DB) (err error) {
    // ID sudah default oleh gorm uuid
    // Perbarui perbandingan dan assignment untuk float64
    if p.MinBelanja == 0.0 { // <<< UBAH PERBANDINGAN DAN NILAI
        p.MinBelanja = 0.0
    }
    if p.MaksimalPotongan == 0.0 { // <<< UBAH PERBANDINGAN DAN NILAI
        p.MaksimalPotongan = 0.0
    }
    if p.DitanggungMerchantPersen == 0.0 { // <<< UBAH PERBANDINGAN DAN NILAI
        p.DitanggungMerchantPersen = 0.0
    }
    return
}