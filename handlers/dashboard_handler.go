package handlers

import (
	"fmt"
	"net/http"

	"backend_kalkuliner/database"
	"backend_kalkuliner/models"

	// Untuk RoundFloat
	"github.com/gin-gonic/gin"
)

// TopResepHPPResult DTO untuk top resep HPP tertinggi
type TopResepHPPResult struct {
	ResepID       string  `json:"resep_id"`
	ResepNama     string  `json:"resep_nama"`
	HPPPerPorsi   float64 `json:"hpp_per_porsi"`
	Channel       string  `json:"channel,omitempty"` // Jika ingin menampilkan channel dari harga_jual terkait
	HargaJualKotor float64 `json:"harga_jual_kotor,omitempty"` // Jika ingin menampilkan harga jual dari harga_jual terkait
}

// GetDashboardSummary mengambil data ringkasan dashboard
// Untuk saat ini hanya mengembalikan dummy/hitung cepat
func GetDashboardSummary(c *gin.Context) {
	// Total Bahan Baku
	var totalBahanBaku int64
	if err := database.DB.Model(&models.BahanBaku{}).Count(&totalBahanBaku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil total bahan baku"})
		return
	}

	// Total Resep
	var totalResep int64
	if err := database.DB.Model(&models.Resep{}).Count(&totalResep).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil total resep"})
		return
	}

	// Total Biaya Operasional (placeholder)
	totalBiayaOperasional := 0.0

	// Top 5 Resep dengan HPP per Porsi Tertinggi

	var topResepsHPPFormatted []TopResepHPPResult

	// Query menggunakan ROW_NUMBER() untuk mendapatkan HPP terbaru per resep_name
	// Kemudian, urutkan dan limit 5
	// PARTITION BY resep_name ensures uniqueness by resep_name, selecting the latest one
	subQuery := database.DB.Table("hpp_results").
		Select("resep_id, resep_nama, hpp_per_porsi, created_at, ROW_NUMBER() OVER (PARTITION BY resep_nama ORDER BY created_at DESC) as rn").
		Where("hpp_per_porsi > ?", 0) // Hanya ambil resep dengan HPP > 0 (opsional)

	err := database.DB.Table("(?) as ranked_hpp", subQuery). // Menggunakan subquery sebagai tabel virtual
		Where("rn = ?", 1). // Hanya ambil baris terbaru for each unique resep_name
		Order("hpp_per_porsi DESC"). // Urutkan berdasarkan HPP per porsi (tertinggi)
		Limit(5). // Ambil 5 teratas
		Find(&topResepsHPPFormatted).Error // Masukkan hasilnya ke struct DTO

	if err != nil {
		fmt.Printf("Error mengambil top resep HPP: %v\n", err) // Debug log error
		// Jangan kembalikan error fatal, cukup kirim data kosong jika ada masalah DB
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil top resep HPP", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_bahan_baku":      totalBahanBaku,
		"total_resep":           totalResep,
		"total_biaya_operasional": totalBiayaOperasional,
		"top_reseps_hpp":        topResepsHPPFormatted,
	})
}