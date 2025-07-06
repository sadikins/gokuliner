package main

import (
	"log" // Pastikan package log diimport

	"backend_kalkuliner/config"
	"backend_kalkuliner/database"
	"backend_kalkuliner/handlers"

	// Import models package (diperlukan untuk AutoMigrate)
	// Import utils package (diperlukan untuk RoundFloat)
	"github.com/gin-contrib/cors" // Untuk konfigurasi CORS
	"github.com/gin-gonic/gin"    // Framework web Gin
)

func main() {
	// 1. Muat konfigurasi aplikasi dari variabel lingkungan (termasuk dari file .env)
	cfg, err := config.LoadConfig() // config.LoadConfig mengembalikan (Config, error)
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err) // Hentikan program jika konfigurasi gagal dimuat
	}
	log.Println("Konfigurasi aplikasi berhasil dimuat.")

	// 2. Inisialisasi koneksi database GORM dan migrasi skema tabel
	// Fungsi InitDB akan menangani koneksi database dan menjalankan AutoMigrate untuk semua model.
	// Program akan berhenti (log.Fatalf) jika ada error saat koneksi atau migrasi.
	database.InitDB(cfg)
	log.Println("Database berhasil diinisialisasi dan skema dimigrasi.")

	// 3. Panggil fungsi utilitas untuk mendaftarkan validator kustom (untuk shopspring/decimal)
	// Kita sudah kembali ke float64, jadi validator custom ini tidak lagi digunakan untuk decimal.
	// Jika ada validator custom lain di utils yang relevan, biarkan pemanggilan ini.
	// utils.RegisterDecimalValidator() // Baris ini tidak lagi relevan jika decimal dihapus.
	// Jika file utils/validator.go Anda kosong atau sudah dihapus, Anda bisa menghapus baris import utils dan pemanggilan ini.
	log.Println("Validator kustom berhasil didaftarkan (jika ada).")


	// 4. Muat Master Data ke Cache
	// Data seperti Bahan Baku dan Resep dimuat ke cache di memori untuk akses cepat oleh handler.
	// Ini krusial untuk performa perhitungan HPP dan Simulasi Promo.
	if err := handlers.LoadMasterDataIntoCache(); err != nil {
		log.Fatalf("Gagal memuat master data ke cache: %v", err) // Hentikan program jika pemuatan cache gagal
	}
	log.Println("Master data berhasil dimuat ke cache aplikasi.")

	// 5. Inisialisasi Gin Router
	// Menggunakan gin.Default() yang sudah menyertakan logger dan recovery middleware.
	router := gin.Default()

	// 6. Konfigurasi CORS (Cross-Origin Resource Sharing)
	// Penting untuk mengizinkan frontend Vue.js (yang berjalan di origin berbeda) berkomunikasi dengan backend.
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Izinkan origin frontend Vue Anda
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}, // Metode HTTP yang diizinkan
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},        // Header yang diizinkan di expose ke browser
		ExposeHeaders:    []string{"Content-Length"},                        // Header yang diizinkan di expose ke browser
		AllowCredentials: true,                                              // Izinkan pengiriman kredensial (misal: cookies)
		MaxAge:           86400,                                             // Durasi cache preflight request (24 jam)
	}))
	log.Println("Konfigurasi CORS berhasil diterapkan.")

	// 7. Daftarkan Routes API
	// Mengelompokkan semua rute di bawah prefix "/api".
	api := router.Group("/api")
	{


		// Routes untuk Modul Bahan Baku (CRUD)
		api.GET("/bahan-bakus", handlers.GetBahanBakus)
		api.POST("/bahan-bakus", handlers.CreateBahanBaku)
		api.GET("/bahan-bakus/:id", handlers.GetBahanBakuByID)
		api.PUT("/bahan-bakus/:id", handlers.UpdateBahanBaku)
		api.DELETE("/bahan-bakus/:id", handlers.DeleteBahanBaku)
		log.Println("Routes Modul Bahan Baku terdaftar.")

		// Routes untuk Modul Resep (CRUD & Duplikasi)
		api.GET("/reseps", handlers.GetReseps)
		api.POST("/reseps", handlers.CreateResep)
		api.GET("/reseps/:id", handlers.GetResepByID)
		api.PUT("/reseps/:id", handlers.UpdateResep)
		api.DELETE("/reseps/:id", handlers.DeleteResep)
		api.POST("/reseps/:id/duplicate", handlers.DuplicateResep) // Endpoint duplikasi resep
		log.Println("Routes Modul Resep terdaftar.")

		// Routes untuk Perhitungan HPP
		api.GET("/hpp/:resep_id", handlers.GetHPPForResep) // Menghitung dan menyimpan HPP per resep
		log.Println("Routes Perhitungan HPP terdaftar.")

		// Routes untuk Modul Harga Jual (Perhitungan & CRUD Data Tersimpan)
		api.POST("/harga-juals/calculate", handlers.CalculateAndSaveHargaJual) // Menghitung dan menyimpan harga jual baru
		api.GET("/harga-juals", handlers.GetHargaJuals)                      // Mengambil daftar harga jual tersimpan
		api.GET("/harga-juals/:id", handlers.GetHargaJualByID)               // Mengambil detail harga jual tersimpan
		api.PUT("/harga-juals/:id", handlers.UpdateHargaJual)                // Memperbarui harga jual tersimpan
		api.DELETE("/harga-juals/:id", handlers.DeleteHargaJual)             // Menghapus harga jual tersimpan
		log.Println("Routes Modul Harga Jual terdaftar.")



		// Routes untuk Modul Program Promo (CRUD)
		api.POST("/program-promos", handlers.CreateProgramPromo)
		api.GET("/program-promos", handlers.GetProgramPromos)
		api.GET("/program-promos/:id", handlers.GetProgramPromoByID)
		api.PUT("/program-promos/:id", handlers.UpdateProgramPromo)
		api.DELETE("/program-promos/:id", handlers.DeleteProgramPromo)
		log.Println("Routes Modul Program Promo terdaftar.")

		// Route Simulasi Promo
		api.POST("/simulasi-promo", handlers.SimulatePromoAndCommission) // Endpoint untuk menjalankan simulasi promo
		log.Println("Route Modul Simulasi Promo terdaftar.")

		//Routes untuk Dashboard
		api.GET("/dashboard", handlers.GetDashboardSummary)
		log.Println("Routes Dashboard terdaftar.")
	}

	// 8. Jalankan Server
	appPort := cfg.AppPort // Ambil port dari konfigurasi
	log.Printf("Server berjalan di http://localhost:%s", appPort)
	router.Run(":" + appPort) // Memulai server Gin pada port yang ditentukan
}