package handlers

import (
	"net/http"

	"backend_kalkuliner/database"
	"backend_kalkuliner/models"

	"github.com/gin-gonic/gin" // <<< IMPORT INI
	"gorm.io/gorm"
)

// CreateBahanBaku (Diperbarui)
func CreateBahanBaku(c *gin.Context) {
	var input models.BahanBaku
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi tambahan untuk field baru (menggunakan metode decimal)
	if input.HargaBeli.IsZero() || input.HargaBeli.IsNegative() || input.NettoPerBeli.IsZero() || input.NettoPerBeli.IsNegative() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Harga beli dan netto per beli harus lebih dari 0."})
		return
	}

	if err := database.DB.Create(&input).Error; err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"bahan_baku_nama_key\" (SQLSTATE 23505)" {
			c.JSON(http.StatusConflict, gin.H{"error": "Nama bahan baku sudah ada."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create bahan baku: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

// UpdateBahanBaku (Diperbarui)
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

	// Validasi tambahan untuk field baru
	if input.HargaBeli.IsZero() || input.HargaBeli.IsNegative() || input.NettoPerBeli.IsZero() || input.NettoPerBeli.IsNegative() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Harga beli dan netto per beli harus lebih dari 0."})
		return
	}

	bahanBaku.Nama = input.Nama
	bahanBaku.Kategori = input.Kategori
	bahanBaku.HargaBeli = input.HargaBeli
	bahanBaku.SatuanBeli = input.SatuanBeli
	bahanBaku.NettoPerBeli = input.NettoPerBeli
	bahanBaku.SatuanPemakaian = input.SatuanPemakaian
	bahanBaku.Catatan = input.Catatan

	if err := database.DB.Save(&bahanBaku).Error; err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"bahan_baku_nama_key\" (SQLSTATE 23505)" {
			c.JSON(http.StatusConflict, gin.H{"error": "Nama bahan baku sudah ada. Silakan gunakan nama lain."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update bahan baku: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, bahanBaku)
}

// GetBahanBakus, GetBahanBakuByID, DeleteBahanBaku (Tidak Berubah pada logika, hanya memastikan import)
func GetBahanBakus(c *gin.Context) {
	var bahanBakus []models.BahanBaku
	if err := database.DB.Find(&bahanBakus).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bahan baku"})
		return
	}
	c.JSON(http.StatusOK, bahanBakus)
}

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

	var resepKomponen models.ResepKomponen
	if err := database.DB.Where("komponen_id = ? AND tipe_komponen = ?", id, "bahan_baku").First(&resepKomponen).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Bahan baku ini digunakan dalam setidaknya satu resep dan tidak dapat dihapus."})
		return
	} else if err != gorm.ErrRecordNotFound {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check bahan baku usage: " + err.Error()})
		return
	}


	if err := database.DB.Delete(&bahanBaku).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete bahan baku"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Bahan baku deleted successfully"})
}