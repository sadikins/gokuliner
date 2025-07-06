package handlers

import (
	"fmt"
	"math"
	"net/http"

	"backend_kalkuliner/database"
	"backend_kalkuliner/models"
	"backend_kalkuliner/utils" // Import utils

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HPPResult struct {
	ResepID     string  `json:"resep_id"`
	ResepNama   string  `json:"resep_nama"`
	HPPPerUnit  float64 `json:"hpp_per_unit"`  // <<< UBAH KE float64
	HPPPerPorsi float64 `json:"hpp_per_porsi"` // <<< UBAH KE float64
}

var (
	ExportedBahanBakuCache = make(map[string]models.BahanBaku)
	ExportedResepCache     = make(map[string]models.Resep)
)

// LoadMasterDataIntoCache memuat semua bahan baku dan resep ke dalam cache yang diekspor
func LoadMasterDataIntoCache() error {
	fmt.Println("Memulai pemuatan master data ke cache...")
	var bahanBakus []models.BahanBaku
	if err := database.DB.Find(&bahanBakus).Error; err != nil {
		fmt.Printf("Error memuat bahan baku ke cache: %v\n", err)
		return fmt.Errorf("gagal memuat bahan baku ke cache: %w", err)
	}
	for _, bb := range bahanBakus {
		ExportedBahanBakuCache[bb.ID] = bb
	}
	fmt.Printf("%d bahan baku dimuat ke cache.\n", len(ExportedBahanBakuCache))

	var reseps []models.Resep
	if err := database.DB.Preload("Komponen").Find(&reseps).Error; err != nil {
		fmt.Printf("Error memuat resep ke cache: %v\n", err)
		return fmt.Errorf("gagal memuat resep ke cache: %w", err)
	}
	for _, r := range reseps {
		ExportedResepCache[r.ID] = r
	}
	fmt.Printf("%d resep dimuat ke cache.\n", len(ExportedResepCache))
	fmt.Println("Pemuatan master data ke cache selesai.")
	return nil
}

// calculateHPPRecursive menghitung HPP secara rekursif
// Sekarang menerima cache yang diekspor sebagai argumen
func calculateHPPRecursive(resepID string, memo map[string]float64) (float64, error) { // <<< UBAH KE float64
	if val, ok := memo[resepID]; ok {
		return val, nil
	}

	resep, ok := ExportedResepCache[resepID]
	if !ok {
		return 0, fmt.Errorf("resep dengan ID %s tidak ditemukan di cache", resepID)
	}

	totalHPP := 0.0 // <<< UBAH KE float64

	for _, komponen := range resep.Komponen {
		var komponenHPP float64 // <<< UBAH KE float64

		if komponen.Kuantitas <= 0 { // <<< UBAH VALIDASI
			return 0, fmt.Errorf("kuantitas komponen '%s' pada resep '%s' harus positif", komponen.KomponenID, resep.Nama)
		}

		if komponen.TipeKomponen == "bahan_baku" {
			bb, ok := ExportedBahanBakuCache[komponen.KomponenID]
			if !ok {
				return 0, fmt.Errorf("bahan baku dengan ID '%s' tidak ditemukan di cache", komponen.KomponenID)
			}

			if bb.NettoPerBeli <= 0 { // <<< UBAH VALIDASI
				return 0, fmt.Errorf("bahan baku '%s' (ID: %s) memiliki netto per beli 0 atau negatif. Tidak dapat menghitung HPP.", bb.Nama, bb.ID)
			}
			hargaPerSatuanPemakaian := bb.HargaBeli / bb.NettoPerBeli // <<< UBAH OPERASI
			komponenHPP = hargaPerSatuanPemakaian

		} else if komponen.TipeKomponen == "resep" {
			subResepHPP, subResepErr := calculateHPPRecursive(komponen.KomponenID, memo)
			if subResepErr != nil {
				return 0, subResepErr
			}

			subResep, ok := ExportedResepCache[komponen.KomponenID]
			if !ok {
				return 0, fmt.Errorf("sub-resep dengan ID %s tidak ditemukan di cache", komponen.KomponenID)
			}

			if subResep.JumlahPorsi <= 0 { // <<< UBAH VALIDASI
				komponenHPP = subResepHPP
				fmt.Printf("Peringatan: Sub-resep '%s' (ID: %s) memiliki JumlahPorsi 0 atau negatif. HPP per unit sub-resep akan sama dengan HPP total sub-resep.\n", subResep.Nama, subResep.ID)
			} else {
				komponenHPP = subResepHPP / subResep.JumlahPorsi // <<< UBAH OPERASI
			}

		} else {
			return 0, fmt.Errorf("tipe komponen tidak valid: %s", komponen.TipeKomponen)
		}

		totalHPP += komponenHPP * komponen.Kuantitas // <<< UBAH OPERASI
	}

	memo[resepID] = totalHPP
	return totalHPP, nil
}

// GetHPPForResep (Diperbarui)
func GetHPPForResep(c *gin.Context) {
	resepID := c.Param("resep_id")

	if err := LoadMasterDataIntoCache(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data master untuk perhitungan HPP: " + err.Error()})
		return
	}

	hpp, err := calculateHPPRecursive(resepID, make(map[string]float64))
	if err != nil {
		fmt.Printf("Error calculating HPP for ResepID %s: %v\n", resepID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung HPP: " + err.Error()})
		return
	}

	resep, ok := ExportedResepCache[resepID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan setelah perhitungan"})
		return
	}

	hppPerPorsi := hpp
	if resep.JumlahPorsi <= 0 {
		fmt.Printf("Peringatan: Resep '%s' (ID: %s) memiliki JumlahPorsi 0 atau negatif. HPP per porsi akan sama dengan HPP per unit.\n", resep.Nama, resep.ID)
	} else {
		hppPerPorsi = hpp / resep.JumlahPorsi
	}

	// >>>>>>> LOGIKA BARU: CEK HPP TERBARU SEBELUM MENYIMPAN <<<<<<<
	var latestHPP models.HPPResult
	// Ambil record HPP terbaru berdasarkan created_at untuk resep ini
	result := database.DB.Where("resep_id = ?", resep.ID).Order("created_at DESC").First(&latestHPP)

	shouldSaveNewHPP := true // Default: simpan record baru

	if result.Error == nil { // Ditemukan record HPP sebelumnya
		// Tentukan epsilon (toleransi perbedaan kecil)
		// Misal, perbedaan kurang dari Rp 0.001 per unit/porsi dianggap tidak signifikan
		epsilon := 0.001

		// Bandingkan HPP per unit dan HPP per porsi
		isHPPPerUnitEffectivelySame := math.Abs(hpp - latestHPP.HPPPerUnit) < epsilon
		isHPPPerPorsiEffectivelySame := math.Abs(hppPerPorsi - latestHPP.HPPPerPorsi) < epsilon

		if isHPPPerUnitEffectivelySame && isHPPPerPorsiEffectivelySame {
			shouldSaveNewHPP = false // Tidak ada perubahan signifikan, jangan simpan
			fmt.Printf("HPP for Resep '%s' (ID: %s) is effectively unchanged (within %.4f tolerance). Not saving new record.\n", resep.Nama, resep.ID, epsilon)
		}
	} else if result.Error != gorm.ErrRecordNotFound {
		// Ini adalah error database selain "record not found"
		// Kita bisa memutuskan untuk return error atau tetap menyimpan (lebih aman jika tidak bisa cek)
		fmt.Printf("Error checking for latest HPP result for ResepID %s: %v. Proceeding to save.\n", resep.ID, result.Error)
		// Jika ingin hentikan: c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal cek HPP terbaru"}) return
	}
	// Jika result.Error == gorm.ErrRecordNotFound, shouldSaveNewHPP tetap true (karena ini HPP pertama)
	// >>>>>>> AKHIR LOGIKA BARU <<<<<<<

	var finalResult models.HPPResult // Variabel untuk menyimpan hasil akhir yang akan dikirim ke frontend

	if shouldSaveNewHPP {
		// Buat objek HPPResult baru dan simpan
		newHPPResult := models.HPPResult{
			ResepID:     resep.ID,
			ResepNama:   resep.Nama,
			HPPPerUnit:  utils.RoundFloat(hpp, 4),      // Bulatkan untuk penyimpanan
			HPPPerPorsi: utils.RoundFloat(hppPerPorsi, 4), // Bulatkan untuk penyimpanan
		}
		if err := database.DB.Create(&newHPPResult).Error; err != nil {
			fmt.Printf("Error saving HPP result for ResepID %s: %v\n", resep.ID, err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan hasil HPP ke database: " + err.Error()})
			return
		}
		fmt.Printf("HPP for Resep '%s' (ID: %s) calculated and SAVED: HPP/Unit=%.4f, HPP/Porsi=%.4f\n",
			resep.Nama, resep.ID, newHPPResult.HPPPerUnit, newHPPResult.HPPPerPorsi)
		finalResult = newHPPResult // Hasil yang baru disimpan
	} else {
		// Jika tidak disimpan (karena tidak ada perubahan signifikan), kembalikan record HPP yang sudah ada
		finalResult = latestHPP
		fmt.Printf("HPP for Resep '%s' (ID: %s) is unchanged. Returning existing record.\n", resep.Nama, resep.ID)
	}

	c.JSON(http.StatusOK, finalResult) // Selalu kirim objek HPPResult yang relevan ke frontend
}