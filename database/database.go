package database

import (
	"fmt"
	"log"

	"backend_kalkuliner/config" // Sesuaikan path package Anda
	"backend_kalkuliner/models" // Sesuaikan path package Anda

	"gorm.io/driver/postgres" // ganti dengan "gorm.io/driver/mysql" jika pakai MySQL
	"gorm.io/gorm"
)

var DB *gorm.DB

// ConnectDatabase menginisialisasi koneksi ke database
func ConnectDatabase(cfg config.Config) {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)




	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	// DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
	// 	Logger: logger.Default.LogMode(logger.Info),
	// })


	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established!", dsn)

	// Migrasi skema database
	err = DB.AutoMigrate(
		&models.BahanBaku{},
		&models.Resep{},
		&models.ResepKomponen{},
	)
	if err != nil {
		log.Fatalf("Failed to auto migrate database: %v", err)
	}
	log.Println("Database migration completed!")
}