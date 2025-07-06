package database

import (
	"fmt"
	"log"

	"backend_kalkuliner/config"
	"backend_kalkuliner/models" // Penting: import semua model yang akan dimigrasi

	"gorm.io/driver/postgres" // Atau "gorm.io/driver/mysql" jika Anda menggunakan MySQL
	"gorm.io/gorm"
	"gorm.io/gorm/logger" // Untuk logging GORM yang lebih baik (opsional)
)

// DB adalah instance GORM database global
var DB *gorm.DB

// InitDB menginisialisasi koneksi ke database dan melakukan migrasi skema tabel.
// Program akan berhenti (log.Fatalf) jika ada error kritis saat koneksi atau migrasi.
func InitDB(cfg config.Config) {
	var err error

	// Data Source Name (DSN) untuk PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)


	// Membuka koneksi database dengan GORM
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Konfigurasi logger GORM (opsional, tapi bagus untuk debugging)
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	log.Println("Koneksi database berhasil dibuat!")

	// Melakukan migrasi otomatis untuk semua model
	// Penting: Setiap kali Anda menambah atau mengubah struct model, pastikan
	// untuk menambahkan atau memperbarui entry di sini.
	err = DB.AutoMigrate(
		&models.BahanBaku{},
		&models.Resep{},
		&models.ResepKomponen{},
		&models.HPPResult{},    // Hasil perhitungan HPP
		&models.HargaJual{},    // Data harga jual yang tersimpan
		&models.ProgramPromo{}, // Data program promo
	)
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi skema database: %v", err)
	}
	log.Println("Migrasi skema database selesai!")
}