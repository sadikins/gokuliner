package handlers

import (
	"fmt"
	"net/http"

	"backend_kalkuliner/database"
	"backend_kalkuliner/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HPPResult struct {
	ResepID     string  `json:"resep_id"`
	ResepNama   string  `json:"resep_nama"`
	HPPPerUnit  float64 `json:"hpp_per_unit"` // HPP per unit resep (sebelum dibagi porsi)
	HPPPerPorsi float64 `json:"hpp_per_porsi"` // <<< TAMBAHKAN INI
}

// GetHPPForResep menghitung HPP untuk resep tertentu
func GetHPPForResep(c *gin.Context) {
	resepID := c.Param("resep_id")

	var resep models.Resep
	// Preload komponen untuk perhitungan rekursif
	if err := database.DB.Preload("Komponen").First(&resep, "id = ?", resepID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep: " + err.Error()})
		return
	}

	// Peta untuk menyimpan HPP yang sudah dihitung agar tidak menghitung ulang
	memo := make(map[string]float64)

	hpp, err := calculateHPPRecursive(resep.ID, memo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung HPP: " + err.Error()})
		return
	}

	hppPerPorsi := hpp
	if resep.JumlahPorsi > 0 { // Hindari pembagian dengan nol
		hppPerPorsi = hpp / resep.JumlahPorsi
	} else {
		// Opsional: berikan peringatan atau set default jika JumlahPorsi 0
		fmt.Printf("Peringatan: Resep '%s' (ID: %s) memiliki JumlahPorsi 0. HPP per porsi akan sama dengan HPP per unit.\n", resep.Nama, resep.ID)
	}


	c.JSON(http.StatusOK, HPPResult{
		ResepID:     resep.ID,
		ResepNama:   resep.Nama,
		HPPPerUnit:  hpp,
		HPPPerPorsi: hppPerPorsi, // <<< TAMBAHKAN INI
	})
}

// calculateHPPRecursive menghitung HPP secara rekursif (Tidak Berubah, tapi akan menggunakan JumlahPorsi dari Resep)
func calculateHPPRecursive(resepID string, memo map[string]float64) (float64, error) {
	// Cek memoization
	if val, ok := memo[resepID]; ok {
		return val, nil
	}

	var resep models.Resep
	if err := database.DB.Preload("Komponen").First(&resep, "id = ?", resepID).Error; err != nil {
		return 0, fmt.Errorf("resep dengan ID %s tidak ditemukan: %w", resepID, err)
	}

	totalHPP := 0.0

	for _, komponen := range resep.Komponen {
		var komponenHPP float64
		var err error

		if komponen.TipeKomponen == "bahan_baku" {
			var bahanBaku models.BahanBaku
			if err = database.DB.First(&bahanBaku, "id = ?", komponen.KomponenID).Error; err != nil {
				return 0, fmt.Errorf("bahan baku dengan ID %s tidak ditemukan: %w", komponen.KomponenID, err)
			}
			komponenHPP = bahanBaku.HargaBeli
		} else if komponen.TipeKomponen == "resep" {
			// Rekursif panggil untuk sub-resep
			subResepHPP, subResepErr := calculateHPPRecursive(komponen.KomponenID, memo)
			if subResepErr != nil {
				return 0, subResepErr
			}
			// Penting: HPP sub-resep harus dibagi dengan jumlah porsi sub-resep itu sendiri
			// untuk mendapatkan HPP per unit sub-resep yang akan digunakan sebagai "harga beli"
			// saat sub-resep tersebut menjadi komponen di resep lain.
			var subResep models.Resep
			if err = database.DB.First(&subResep, "id = ?", komponen.KomponenID).Error; err != nil {
				return 0, fmt.Errorf("sub-resep dengan ID %s tidak ditemukan: %w", komponen.KomponenID, err)
			}

			if subResep.JumlahPorsi > 0 {
				komponenHPP = subResepHPP / subResep.JumlahPorsi
			} else {
				// Jika sub-resep tidak punya porsi, anggap 1 unit = 1 porsi untuk HPP
				komponenHPP = subResepHPP
				fmt.Printf("Peringatan: Sub-resep '%s' (ID: %s) memiliki JumlahPorsi 0. HPP per unit sub-resep akan sama dengan HPP total sub-resep.\n", subResep.Nama, subResep.ID)
			}

		} else {
			return 0, fmt.Errorf("tipe komponen tidak valid: %s", komponen.TipeKomponen)
		}

		totalHPP += komponenHPP * komponen.Kuantitas
	}

	// Simpan hasil ke memo sebelum kembali
	memo[resepID] = totalHPP
	return totalHPP, nil
}