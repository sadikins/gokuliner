package handlers

import (
	"net/http"

	"backend_kalkuliner/database"
	"backend_kalkuliner/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetBahanBakus mengambil semua data bahan baku (Sudah ada)
func GetBahanBakus(c *gin.Context) {
	var bahanBakus []models.BahanBaku
	if err := database.DB.Find(&bahanBakus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bahan baku"})
		return
	}
	c.JSON(http.StatusOK, bahanBakus)
}

// CreateBahanBaku membuat data bahan baku baru (Sudah ada)
func CreateBahanBaku(c *gin.Context) {
	var input models.BahanBaku
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&input).Error; err != nil {
		// Menangkap error duplikat nama
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"bahan_baku_nama_key\" (SQLSTATE 23505)" {
			c.JSON(http.StatusConflict, gin.H{"error": "Nama bahan baku sudah ada."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bahan baku: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

// GetBahanBakuByID mengambil satu bahan baku berdasarkan ID
func GetBahanBakuByID(c *gin.Context) {
	id := c.Param("id")
	var bahanBaku models.BahanBaku
	if err := database.DB.First(&bahanBaku, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bahan baku not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bahan baku"})
		return
	}
	c.JSON(http.StatusOK, bahanBaku)
}

// UpdateBahanBaku memperbarui data bahan baku
func UpdateBahanBaku(c *gin.Context) {
	id := c.Param("id")
	var bahanBaku models.BahanBaku
	if err := database.DB.First(&bahanBaku, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bahan baku not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bahan baku"})
		return
	}

	var input models.BahanBaku
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Jangan update ID dan CreatedAt
	bahanBaku.Nama = input.Nama
	bahanBaku.Satuan = input.Satuan
	bahanBaku.HargaBeli = input.HargaBeli

	if err := database.DB.Save(&bahanBaku).Error; err != nil {
		// Menangkap error duplikat nama saat update
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"bahan_baku_nama_key\" (SQLSTATE 23505)" {
			c.JSON(http.StatusConflict, gin.H{"error": "Nama bahan baku sudah ada."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bahan baku: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, bahanBaku)
}

// DeleteBahanBaku menghapus data bahan baku
func DeleteBahanBaku(c *gin.Context) {
	id := c.Param("id")
	var bahanBaku models.BahanBaku
	if err := database.DB.First(&bahanBaku, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Bahan baku not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bahan baku"})
		return
	}

	// Cek apakah bahan baku ini digunakan di resep mana pun sebelum dihapus
	// Ini adalah validasi penting untuk menjaga integritas data.
	// Jika ingin lebih canggih, bisa gunakan JOIN. Untuk sederhana, First.
	var resepKomponen models.ResepKomponen
	if err := database.DB.Where("komponen_id = ? AND tipe_komponen = ?", id, "bahan_baku").First(&resepKomponen).Error; err == nil {
		// Ditemukan penggunaan
		c.JSON(http.StatusConflict, gin.H{"error": "Bahan baku ini digunakan dalam setidaknya satu resep dan tidak dapat dihapus."})
		return
	} else if err != gorm.ErrRecordNotFound {
		// Error selain not found
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check bahan baku usage: " + err.Error()})
		return
	}


	if err := database.DB.Delete(&bahanBaku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bahan baku"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bahan baku deleted successfully"})
}