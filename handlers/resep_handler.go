package handlers

import (
	"fmt"
	"net/http"

	"backend_kalkuliner/database"
	"backend_kalkuliner/models"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal" // <<< IMPORT INI
	"gorm.io/gorm"
)

// Struktur input untuk membuat resep baru (Diperbarui)
type CreateResepInput struct {
	Nama        string                `json:"nama" binding:"required"`
	IsSubResep  bool                  `json:"is_sub_resep"`
	JumlahPorsi decimal.Decimal       `json:"jumlah_porsi"` // <<< UBAH TIPE INI
	Komponen    []models.ResepKomponen `json:"komponen"`
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
		JumlahPorsi: input.JumlahPorsi,
	}
	if resep.JumlahPorsi.IsZero() || resep.JumlahPorsi.IsNegative() { // Validasi decimal
		resep.JumlahPorsi = decimal.NewFromFloat(1.0)
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

		if compInput.Kuantitas.IsZero() || compInput.Kuantitas.IsNegative() { // Validasi Kuantitas
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kuantitas komponen harus lebih dari 0."})
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

	// 1. Update data resep utama
	existingResep.Nama = input.Nama
	existingResep.IsSubResep = input.IsSubResep
	existingResep.JumlahPorsi = input.JumlahPorsi // ASSIGN NILAI INI
	if existingResep.JumlahPorsi.IsZero() || existingResep.JumlahPorsi.IsNegative() {
		existingResep.JumlahPorsi = decimal.NewFromFloat(1.0)
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

	// 2. Kelola ResepKomponen: Hapus yang lama, tambahkan yang baru (Diperbarui Validasi)
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

		if compInput.Kuantitas.IsZero() || compInput.Kuantitas.IsNegative() {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Kuantitas komponen harus lebih dari 0."})
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

// GetReseps, GetResepByID, DeleteResep (Tidak Berubah pada logika, hanya memastikan import)
func GetReseps(c *gin.Context) {
	var reseps []models.Resep
	if err := database.DB.Preload("Komponen").Find(&reseps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep"})
		return
	}
	c.JSON(http.StatusOK, reseps)
}

// KomponenDetailResponse adalah DTO untuk detail komponen dalam resep
type KomponenDetailResponse struct {
	ID        string          `json:"id"`        // ID BahanBaku atau Resep komponen
	Nama      string          `json:"nama"`      // Nama BahanBaku atau Resep komponen
	Kuantitas decimal.Decimal `json:"kuantitas"` // Kuantitas penggunaan dalam resep ini
	Tipe      string          `json:"tipe"`      // 'bahan_baku' atau 'resep'
	Satuan    string          `json:"satuan,omitempty"` // Satuan Pemakaian untuk bahan baku, atau kosong untuk resep
	HargaUnit decimal.Decimal `json:"harga_unit,omitempty"` // Harga per unit pemakaian (bahan baku) atau HPP per porsi (resep)
}

// ResepDetailResponse adalah DTO untuk detail resep lengkap
type ResepDetailResponse struct {
	ID          string                  `json:"id"`
	Nama        string                  `json:"nama"`
	IsSubResep  bool                    `json:"is_sub_resep"`
	JumlahPorsi decimal.Decimal         `json:"jumlah_porsi"`
	Komponen    []KomponenDetailResponse `json:"komponen"` // Menggunakan DTO KomponenDetailResponse
	CreatedAt   string                  `json:"created_at"` // Format string untuk kemudahan frontend
	UpdatedAt   string                  `json:"updated_at"` // Format string untuk kemudahan frontend
}

// GetResepByID mengambil satu resep berdasarkan ID (Diperbarui)
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
				// Handle case where bahan baku might be deleted but still referenced
				detail.Nama = "[Bahan Baku Tidak Ditemukan]"
				detail.Satuan = "" // Kosongkan satuan
				detail.HargaUnit = decimal.Zero
			} else {
				detail.Nama = bb.Nama
				detail.Satuan = bb.SatuanPemakaian
				// Hitung harga per unit pemakaian untuk display
				if bb.NettoPerBeli.IsZero() || bb.NettoPerBeli.IsNegative() {
					detail.HargaUnit = decimal.Zero // Hindari pembagian nol
				} else {
					detail.HargaUnit = bb.HargaBeli.Div(bb.NettoPerBeli).Round(4) // Bulatkan untuk display
				}
			}
		} else if komp.TipeKomponen == "resep" {
			var subResep models.Resep
			if err := database.DB.First(&subResep, "id = ?", komp.KomponenID).Error; err != nil {
				detail.Nama = "[Resep Tidak Ditemukan]"
				detail.Satuan = ""
				detail.HargaUnit = decimal.Zero
			} else {
				detail.Nama = subResep.Nama
				// Untuk resep, satuan bisa jadi porsi atau unit resep
				detail.Satuan = fmt.Sprintf("per %s porsi", subResep.JumlahPorsi.StringFixed(0)) // Contoh
				// Harga unit resep bisa jadi HPP per porsi atau per unit, tergantung konteks
				// Untuk detail ini, kita hanya tampilkan namanya, HPP dihitung terpisah
				detail.HargaUnit = decimal.Zero // Tidak relevan di sini untuk harga unit komponen
			}
		}
		resepDetail.Komponen = append(resepDetail.Komponen, detail)
	}

	c.JSON(http.StatusOK, resepDetail)
}



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


func DuplicateResep(c *gin.Context) {
    resepID := c.Param("id")

    tx := database.DB.Begin()
    if tx.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memulai transaksi database"})
        return
    }

    // 1. Ambil resep asli beserta komponennya
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

    // 2. Buat objek resep baru (salinan)
    newResep := models.Resep{
        Nama:        originalResep.Nama + " (Copy)", // Ubah nama agar unik
        IsSubResep:  originalResep.IsSubResep,
        JumlahPorsi: originalResep.JumlahPorsi,
        // ID akan otomatis digenerate oleh BeforeCreate
    }
    if err := tx.Create(&newResep).Error; err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat resep duplikat: " + err.Error()})
        return
    }

    // 3. Duplikasi komponen resep
    for _, originalComp := range originalResep.Komponen {
        newResepKomponen := models.ResepKomponen{
            ResepID:      newResep.ID,         // Link ke ID resep baru
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