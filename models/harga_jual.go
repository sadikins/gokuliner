package models

import (
	"time"

	"gorm.io/gorm"
)

// HargaJual merepresentasikan data harga jual yang dihitung untuk suatu resep/produk
type HargaJual struct {
	ID                 string          `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	ResepID            string          `gorm:"type:uuid;not null" json:"resep_id"`
	Resep              Resep           `gorm:"foreignKey:ResepID" json:"resep,omitempty"` // Relasi ke model Resep
	NamaProduk         string          `gorm:"type:varchar(255);not null" json:"nama_produk"`
	Channel            string          `gorm:"type:varchar(50);not null" json:"channel"` // <<< TAMBAHKAN INI: Channel penjualan (GoFood, GrabFood, Internal, etc.)
	HPP                float64 `gorm:"type:decimal(18,4);not null" json:"hpp"` // HPP dari resep terkait
	JumlahPorsiProduk  float64 `gorm:"type:decimal(18,4);not null" json:"jumlah_porsi_produk"` // Jumlah porsi yang dihasilkan produk ini

	// Kriteria Perhitungan
	MetodePerhitungan  string          `json:"metode_perhitungan"`
	NilaiKriteria      float64 `json:"nilai_kriteria"`

	// Biaya Tambahan (ini adalah persentase yang digunakan dalam perhitungan dasar HargaJualKotor)
	PajakPersen        float64 `json:"pajak_persen"`
	KomisiChannelPersen float64 `json:"komisi_channel_persen"`

	// Hasil Perhitungan
	HargaJualKotor     float64 `json:"harga_jual_kotor"`
	HargaJualBersih    float64 `json:"harga_jual_bersih"`
	TotalPajak         float64 `json:"total_pajak"`
	TotalKomisi        float64 `json:"total_komisi"`
	Profit             float64 `json:"profit"`
	ProfitPersen       float64 `json:"profit_persen"`

	CreatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}