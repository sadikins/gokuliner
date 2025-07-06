package handlers

import (
	"net/http"

	"backend_kalkuliner/database"
	"backend_kalkuliner/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateProgramPromoInput untuk input Create/Update
type CreateProgramPromoInput struct {
	NamaPromo           string  `json:"nama_promo" binding:"required"`
	Channel             string  `json:"channel" binding:"required"`
	JenisDiskon         string  `json:"jenis_diskon" binding:"required,oneof=persentase nominal"`
	BesarDiskon         float64 `json:"besar_diskon" binding:"required"` // <<< UBAH KE float64
	MinBelanja          float64 `json:"min_belanja"`  // <<< UBAH KE float64
	MaksimalPotongan    float64 `json:"maksimal_potongan"` // <<< UBAH KE float64
	DitanggungMerchantPersen float64 `json:"ditanggung_merchant_persen"` // <<< UBAH KE float64
	Catatan             string  `json:"catatan"`
}

// CreateProgramPromo membuat program promo baru
func CreateProgramPromo(c *gin.Context) {
	var input CreateProgramPromoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi tambahan
	if input.BesarDiskon <= 0 && input.JenisDiskon == "persentase" { // <<< UBAH VALIDASI
		c.JSON(http.StatusBadRequest, gin.H{"error": "Besar diskon (persentase) harus lebih dari 0."})
		return
	}
	if input.MinBelanja < 0 || input.MaksimalPotongan < 0 || input.DitanggungMerchantPersen < 0 { // <<< UBAH VALIDASI
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nilai nominal tidak boleh negatif."})
		return
	}


	promo := models.ProgramPromo{
		NamaPromo:           input.NamaPromo,
		Channel:             input.Channel,
		JenisDiskon:         input.JenisDiskon,
		BesarDiskon:         input.BesarDiskon,
		MinBelanja:          input.MinBelanja,
		MaksimalPotongan:    input.MaksimalPotongan,
		DitanggungMerchantPersen: input.DitanggungMerchantPersen,
		Catatan:             input.Catatan,
	}

	if err := database.DB.Create(&promo).Error; err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"program_promos_nama_promo_key\" (SQLSTATE 23505)" {
			c.JSON(http.StatusConflict, gin.H{"error": "Nama promo sudah ada."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat program promo: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, promo)
}

// GetProgramPromos mengambil semua program promo
func GetProgramPromos(c *gin.Context) {
	var promos []models.ProgramPromo
	if err := database.DB.Find(&promos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil program promo"})
		return
	}
	c.JSON(http.StatusOK, promos)
}

// GetProgramPromoByID mengambil program promo berdasarkan ID
func GetProgramPromoByID(c *gin.Context) {
	id := c.Param("id")
	var promo models.ProgramPromo
	if err := database.DB.First(&promo, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Program promo tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil program promo"})
		return
	}
	c.JSON(http.StatusOK, promo)
}

// UpdateProgramPromo memperbarui program promo
func UpdateProgramPromo(c *gin.Context) {
	id := c.Param("id")
	var promo models.ProgramPromo
	if err := database.DB.First(&promo, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Program promo tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil program promo"})
		return
	}

	var input CreateProgramPromoInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi tambahan
	if input.BesarDiskon <= 0 && input.JenisDiskon == "persentase" { // <<< UBAH VALIDASI
		c.JSON(http.StatusBadRequest, gin.H{"error": "Besar diskon (persentase) harus lebih dari 0."})
		return
	}
	if input.MinBelanja < 0 || input.MaksimalPotongan < 0 || input.DitanggungMerchantPersen < 0 { // <<< UBAH VALIDASI
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nilai nominal tidak boleh negatif."})
		return
	}


	promo.NamaPromo =           input.NamaPromo
	promo.Channel =             input.Channel
	promo.JenisDiskon =         input.JenisDiskon
	promo.BesarDiskon =         input.BesarDiskon
	promo.MinBelanja =          input.MinBelanja
	promo.MaksimalPotongan =    input.MaksimalPotongan
	promo.DitanggungMerchantPersen = input.DitanggungMerchantPersen
	promo.Catatan =             input.Catatan

	if err := database.DB.Save(&promo).Error; err != nil {
		if err.Error() == "ERROR: duplicate key value violates unique constraint \"program_promos_nama_promo_key\" (SQLSTATE 23505)" {
			c.JSON(http.StatusConflict, gin.H{"error": "Nama promo sudah ada."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui program promo: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, promo)
}

// DeleteProgramPromo menghapus program promo
func DeleteProgramPromo(c *gin.Context) {
	id := c.Param("id")
	var promo models.ProgramPromo
	if err := database.DB.First(&promo, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Program promo tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil program promo"})
		return
	}

	// TODO: Cek apakah promo ini sedang aktif atau terkait dengan data historis
	// Untuk saat ini, kita langsung hapus

	if err := database.DB.Delete(&promo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus program promo"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Program promo berhasil dihapus"})
}