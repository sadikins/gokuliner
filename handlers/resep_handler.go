package handlers

import (
	"net/http"

	"backend_kalkuliner/database"
	"backend_kalkuliner/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Struktur input untuk membuat resep baru (Diperbarui)
type CreateResepInput struct {
	Nama        string                `json:"nama" binding:"required"`
	IsSubResep  bool                  `json:"is_sub_resep"`
	JumlahPorsi float64               `json:"jumlah_porsi"` // <<< TAMBAHKAN INI
	Komponen    []models.ResepKomponen `json:"komponen"`
}

// GetReseps mengambil semua resep (Tidak Berubah)
func GetReseps(c *gin.Context) {
	var reseps []models.Resep
	if err := database.DB.Preload("Komponen").Find(&reseps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep"})
		return
	}
	c.JSON(http.StatusOK, reseps)
}

// CreateResep membuat resep baru beserta komponen-komponennya (Diperbarui)
func CreateResep(c *gin.Context) {
	var input CreateResepInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi database"})
		return
	}

	resep := models.Resep{
		Nama:       input.Nama,
		IsSubResep: input.IsSubResep,
		JumlahPorsi: input.JumlahPorsi, // <<< ASSIGN NILAI INI
	}
	// Pastikan JumlahPorsi tidak nol atau negatif jika ada validasi tambahan yang diperlukan
	if resep.JumlahPorsi <= 0 {
		resep.JumlahPorsi = 1.0 // Default ke 1 jika tidak valid
	}

	if err := tx.Create(&resep).Error; err != nil {
		tx.Rollback()
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"resep_nama_key\" (SQLSTATE 23505)" {
			c.JSON(http.StatusConflict, gin.H{"error": "Nama resep sudah ada."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat resep: " + err.Error()})
		return
	}

	for _, compInput := range input.Komponen {
		if compInput.TipeKomponen != "bahan_baku" && compInput.TipeKomponen != "resep" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tipe komponen tidak valid: " + compInput.TipeKomponen})
			return
		}

		if compInput.TipeKomponen == "bahan_baku" {
			var bb models.BahanBaku
			if err := tx.First(&bb, "id = ?", compInput.KomponenID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Bahan baku dengan ID " + compInput.KomponenID + " tidak ditemukan"})
				return
			}
		} else { // tipe_komponen == "resep"
			var r models.Resep
			if err := tx.First(&r, "id = ?", compInput.KomponenID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Resep dengan ID " + compInput.KomponenID + " tidak ditemukan"})
				return
			}
		}

		resepKomponen := models.ResepKomponen{
			ResepID:      resep.ID,
			KomponenID:   compInput.KomponenID,
			Kuantitas:    compInput.Kuantitas,
			TipeKomponen: compInput.TipeKomponen,
		}
		if err := tx.Create(&resepKomponen).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan komponen resep: " + err.Error()})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"message": "Resep berhasil dibuat", "resep_id": resep.ID})
}

// GetResepByID mengambil satu resep berdasarkan ID (Tidak Berubah)
func GetResepByID(c *gin.Context) {
	id := c.Param("id")
	var resep models.Resep
	if err := database.DB.Preload("Komponen").First(&resep, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resep not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resep"})
		return
	}
	c.JSON(http.StatusOK, resep)
}

// UpdateResep memperbarui resep beserta komponen-komponennya (Diperbarui)
func UpdateResep(c *gin.Context) {
	id := c.Param("id")
	var existingResep models.Resep
	if err := database.DB.Preload("Komponen").First(&existingResep, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resep not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resep"})
		return
	}

	var input CreateResepInput // Gunakan input yang sama dengan Create
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi database"})
		return
	}

	// 1. Update data resep utama
	existingResep.Nama = input.Nama
	existingResep.IsSubResep = input.IsSubResep
	existingResep.JumlahPorsi = input.JumlahPorsi // <<< ASSIGN NILAI INI
	if existingResep.JumlahPorsi <= 0 {
		existingResep.JumlahPorsi = 1.0 // Default ke 1 jika tidak valid
	}

	if err := tx.Save(&existingResep).Error; err != nil {
		tx.Rollback()
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"resep_nama_key\" (SQLSTATE 23505)" {
			c.JSON(http.StatusConflict, gin.H{"error": "Nama resep sudah ada."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui resep: " + err.Error()})
		return
	}

	// 2. Kelola ResepKomponen: Hapus yang lama, tambahkan yang baru (Tidak Berubah)
	if err := tx.Where("resep_id = ?", id).Delete(&models.ResepKomponen{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus komponen resep lama: " + err.Error()})
		return
	}

	for _, compInput := range input.Komponen {
		if compInput.TipeKomponen != "bahan_baku" && compInput.TipeKomponen != "resep" {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Tipe komponen tidak valid: " + compInput.TipeKomponen})
			return
		}

		if compInput.TipeKomponen == "bahan_baku" {
			var bb models.BahanBaku
			if err := tx.First(&bb, "id = ?", compInput.KomponenID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Bahan baku dengan ID " + compInput.KomponenID + " tidak ditemukan"})
				return
			}
		} else { // tipe_komponen == "resep"
			var r models.Resep
			if err := tx.First(&r, "id = ?", compInput.KomponenID).Error; err != nil {
				tx.Rollback()
				c.JSON(http.StatusBadRequest, gin.H{"error": "Resep dengan ID " + compInput.KomponenID + " tidak ditemukan"})
				return
			}
		}

		resepKomponen := models.ResepKomponen{
			ResepID:      id,
			KomponenID:   compInput.KomponenID,
			Kuantitas:    compInput.Kuantitas,
			TipeKomponen: compInput.TipeKomponen,
		}
		if err := tx.Create(&resepKomponen).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menambahkan komponen resep baru: " + err.Error()})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Resep berhasil diperbarui", "resep_id": id})
}

// DeleteResep menghapus resep beserta semua komponennya (Tidak Berubah)
func DeleteResep(c *gin.Context) {
	id := c.Param("id")

	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi database"})
		return
	}

	var resepKomponen models.ResepKomponen
	if err := tx.Where("komponen_id = ? AND tipe_komponen = ?", id, "resep").First(&resepKomponen).Error; err == nil {
		tx.Rollback()
		c.JSON(http.StatusConflict, gin.H{"error": "Resep ini digunakan sebagai komponen di resep lain dan tidak dapat dihapus."})
		return
	} else if err != gorm.ErrRecordNotFound {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengecek penggunaan resep: " + err.Error()})
		return
	}

	if err := tx.Where("resep_id = ?", id).Delete(&models.ResepKomponen{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus komponen resep terkait: " + err.Error()})
		return
	}

	var resep models.Resep
	if err := tx.First(&resep, "id = ?", id).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resep not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resep for deletion"})
		return
	}
	if err := tx.Delete(&resep).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus resep: " + err.Error()})
		return
	}

	tx.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Resep deleted successfully"})
}