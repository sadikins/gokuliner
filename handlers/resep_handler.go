package handlers

import (
	"fmt" // Diperlukan untuk math.Round
	"net/http"

	"backend_kalkuliner/database"
	"backend_kalkuliner/models"
	"backend_kalkuliner/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)



type KomponenDetailResponse struct {
	ID        string  `json:"id"`        // ID BahanBaku atau Resep komponen
	Nama      string  `json:"nama"`      // Nama BahanBaku atau Resep komponen
	Kuantitas float64 `json:"kuantitas"` // Kuantitas penggunaan dalam resep ini
	Tipe      string  `json:"tipe"`      // 'bahan_baku' atau 'resep'
	Satuan    string  `json:"satuan,omitempty"` // Satuan Pemakaian untuk bahan baku, atau kosong untuk resep
	HargaUnit float64 `json:"harga_unit,omitempty"` // Harga per unit pemakaian (bahan baku) atau HPP per porsi (resep)
}

// ResepDetailResponse adalah DTO untuk detail resep lengkap
// Ini adalah respons API, bukan model database
type ResepDetailResponse struct {
	ID          string                  `json:"id"`
	Nama        string                  `json:"nama"`
	IsSubResep  bool                    `json:"is_sub_resep"`
	JumlahPorsi float64                 `json:"jumlah_porsi"`
	Komponen    []KomponenDetailResponse `json:"komponen"` // Menggunakan DTO KomponenDetailResponse
	CreatedAt   string                  `json:"created_at"` // Format string untuk kemudahan frontend
	UpdatedAt   string                  `json:"updated_at"` // Format string untuk kemudahan frontend
}


type CreateResepInput struct {
	Nama        string  `json:"nama" binding:"required"`
	IsSubResep  bool    `json:"is_sub_resep"`
	JumlahPorsi float64 `json:"jumlah_porsi"`
	Komponen    []models.ResepKomponen `json:"komponen"`
}

// GetReseps mengambil semua resep
func GetReseps(c *gin.Context) {
	var reseps []models.Resep
	// Preload komponen agar bisa dikirim ke frontend
	if err := database.DB.Preload("Komponen").Find(&reseps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep"})
		return
	}
	c.JSON(http.StatusOK, reseps)
}

// CreateResep membuat resep baru beserta komponen-komponennya
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
		JumlahPorsi: input.JumlahPorsi,
	}
	if resep.JumlahPorsi <= 0 { // Validasi
		resep.JumlahPorsi = 1.0
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

		if compInput.Kuantitas <= 0 { // Validasi
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kuantitas komponen harus lebih dari 0."})
			return
		}

		// Validasi apakah komponenID benar-benar ada
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

// GetResepByID mengambil satu resep berdasarkan ID
func GetResepByID(c *gin.Context) {
	id := c.Param("id")
	var resep models.Resep
	// Preload komponen
	if err := database.DB.Preload("Komponen").First(&resep, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resep not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep"})
		return
	}

	// Membuat response DTO
	resepDetail := ResepDetailResponse{
		ID:          resep.ID,
		Nama:        resep.Nama,
		IsSubResep:  resep.IsSubResep,
		JumlahPorsi: resep.JumlahPorsi,
		Komponen:    []KomponenDetailResponse{},
		CreatedAt:   resep.CreatedAt.Format("2006-01-02 15:04:05"), // Format tanggal
		UpdatedAt:   resep.UpdatedAt.Format("2006-01-02 15:04:05"), // Format tanggal
	}

	// Mengisi detail komponen
	for _, komp := range resep.Komponen {
		detail := KomponenDetailResponse{
			ID:        komp.KomponenID,
			Kuantitas: komp.Kuantitas,
			Tipe:      komp.TipeKomponen,
		}

		if komp.TipeKomponen == "bahan_baku" {
			var bb models.BahanBaku
			if err := database.DB.First(&bb, "id = ?", komp.KomponenID).Error; err != nil {
				detail.Nama = "[Bahan Baku Tidak Ditemukan]"
				detail.Satuan = ""
				detail.HargaUnit = 0.0 // Default
			} else {
				detail.Nama = bb.Nama
				detail.Satuan = bb.SatuanPemakaian
				// Hitung harga per unit pemakaian untuk display
				if bb.NettoPerBeli <= 0 {
					detail.HargaUnit = 0.0
				} else {
					detail.HargaUnit = utils.RoundFloat(bb.HargaBeli / bb.NettoPerBeli, 4) // Bulatkan untuk display
				}
			}
		} else if komp.TipeKomponen == "resep" {
			var subResep models.Resep
			if err := database.DB.First(&subResep, "id = ?", komp.KomponenID).Error; err != nil {
				detail.Nama = "[Resep Tidak Ditemukan]"
				detail.Satuan = ""
				detail.HargaUnit = 0.0
			} else {
				detail.Nama = subResep.Nama
				detail.Satuan = fmt.Sprintf("per %.0f porsi", subResep.JumlahPorsi) // Contoh
				detail.HargaUnit = 0.0
			}
		}
		resepDetail.Komponen = append(resepDetail.Komponen, detail)
	}

	c.JSON(http.StatusOK, resepDetail)
}

// UpdateResep memperbarui resep beserta komponen-komponennya
func UpdateResep(c *gin.Context) {
	id := c.Param("id")
	var existingResep models.Resep
	if err := database.DB.Preload("Komponen").First(&existingResep, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resep not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep"})
		return
	}

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

	existingResep.Nama = input.Nama
	existingResep.IsSubResep = input.IsSubResep
	existingResep.JumlahPorsi = input.JumlahPorsi
	if existingResep.JumlahPorsi <= 0 { // Validasi
		existingResep.JumlahPorsi = 1.0
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

		if compInput.Kuantitas <= 0 { // Validasi
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kuantitas komponen harus lebih dari 0."})
			return
		}

		// Validasi apakah komponenID benar-benar ada
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

// GetReseps mengambil semua resep (Tidak Berubah)
// func GetReseps(c *gin.Context) {
// 	var reseps []models.Resep
// 	if err := database.DB.Preload("Komponen").Find(&reseps).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, reseps)
// }

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

// DuplicateResep membuat salinan dari resep yang sudah ada (Tidak Berubah)
func DuplicateResep(c *gin.Context) {
	resepID := c.Param("id")

	tx := database.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi database"})
		return
	}

	var originalResep models.Resep
	if err := tx.Preload("Komponen").First(&originalResep, "id = ?", resepID).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resep asli tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep asli: " + err.Error()})
		return
	}

	newResep := models.Resep{
		Nama:        originalResep.Nama + " (Copy)",
		IsSubResep:  originalResep.IsSubResep,
		JumlahPorsi: originalResep.JumlahPorsi,
	}
	if err := tx.Create(&newResep).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat resep duplikat: " + err.Error()})
		return
	}

	for _, originalComp := range originalResep.Komponen {
		newResepKomponen := models.ResepKomponen{
			ResepID:      newResep.ID,
			KomponenID:   originalComp.KomponenID,
			Kuantitas:    originalComp.Kuantitas,
			TipeKomponen: originalComp.TipeKomponen,
		}
		if err := tx.Create(&newResepKomponen).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menduplikasi komponen resep: " + err.Error()})
			return
		}
	}

	tx.Commit()
	c.JSON(http.StatusCreated, gin.H{"message": "Resep berhasil diduplikasi", "resep_id_baru": newResep.ID, "nama_resep_baru": newResep.Nama})
}