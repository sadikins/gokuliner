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

type HPPResult struct {
	ResepID     string          `json:"resep_id"`
	ResepNama   string          `json:"resep_nama"`
	HPPPerUnit  decimal.Decimal `json:"hpp_per_unit"` // <<< UBAH TIPE INI
	HPPPerPorsi decimal.Decimal `json:"hpp_per_porsi"` // <<< UBAH TIPE INI
}

// loadMasterDataIntoCache (Diperbarui untuk memuat decimal.Decimal)
var (
	bahanBakuCache = make(map[string]models.BahanBaku)
	resepCache     = make(map[string]models.Resep)
)

func loadMasterDataIntoCache() error {
	var bahanBakus []models.BahanBaku
	if err := database.DB.Find(&bahanBakus).Error; err != nil {
		return fmt.Errorf("gagal memuat bahan baku ke cache: %w", err)
	}
	for _, bb := range bahanBakus {
		bahanBakuCache[bb.ID] = bb
	}

	var reseps []models.Resep
	if err := database.DB.Preload("Komponen").Find(&reseps).Error; err != nil {
		return fmt.Errorf("gagal memuat resep ke cache: %w", err)
	}
	for _, r := range reseps {
		resepCache[r.ID] = r
	}
	return nil
}

// calculateHPPRecursive (Diperbarui untuk menggunakan decimal.Decimal)
func calculateHPPRecursive(resepID string, memo map[string]decimal.Decimal) (decimal.Decimal, error) { // <<< UBAH TIPE KEMBALIAN DAN MAP
	if val, ok := memo[resepID]; ok {
		return val, nil
	}

	resep, ok := resepCache[resepID]
	if !ok {
		return decimal.Decimal{}, fmt.Errorf("resep dengan ID %s tidak ditemukan di cache", resepID)
	}

	totalHPP := decimal.Zero // <<< INISIALISASI DENGAN decimal.Zero

	for _, komponen := range resep.Komponen {
		var komponenHPP decimal.Decimal // <<< UBAH TIPE
		// var err error

		if komponen.Kuantitas.IsZero() || komponen.Kuantitas.IsNegative() { // Validasi Kuantitas (pakai decimal method)
			return decimal.Decimal{}, fmt.Errorf("kuantitas komponen '%s' pada resep '%s' harus positif", komponen.KomponenID, resep.Nama)
		}

		if komponen.TipeKomponen == "bahan_baku" {
			bb, ok := bahanBakuCache[komponen.KomponenID]
			if !ok {
				return decimal.Decimal{}, fmt.Errorf("bahan baku dengan ID '%s' tidak ditemukan di cache", komponen.KomponenID)
			}

			if bb.NettoPerBeli.IsZero() || bb.NettoPerBeli.IsNegative() { // Validasi NettoPerBeli
				return decimal.Decimal{}, fmt.Errorf("bahan baku '%s' (ID: %s) memiliki netto per beli 0 atau negatif. Tidak dapat menghitung HPP.", bb.Nama, bb.ID)
			}
			// Harga per unit pemakaian = Harga Beli per Satuan Beli / Netto per Beli
			hargaPerSatuanPemakaian := bb.HargaBeli.Div(bb.NettoPerBeli) // <<< GUNAKAN METHOD .Div()
			komponenHPP = hargaPerSatuanPemakaian

		} else if komponen.TipeKomponen == "resep" {
			subResepHPP, subResepErr := calculateHPPRecursive(komponen.KomponenID, memo)
			if subResepErr != nil {
				return decimal.Decimal{}, subResepErr
			}

			subResep, ok := resepCache[komponen.KomponenID]
			if !ok {
				return decimal.Decimal{}, fmt.Errorf("sub-resep dengan ID %s tidak ditemukan di cache", komponen.KomponenID)
			}

			if subResep.JumlahPorsi.IsZero() || subResep.JumlahPorsi.IsNegative() { // Validasi JumlahPorsi
				komponenHPP = subResepHPP
				fmt.Printf("Peringatan: Sub-resep '%s' (ID: %s) memiliki JumlahPorsi 0 atau negatif. HPP per unit sub-resep akan sama dengan HPP total sub-resep.\n", subResep.Nama, subResep.ID)
			} else {
				komponenHPP = subResepHPP.Div(subResep.JumlahPorsi) // <<< GUNAKAN METHOD .Div()
			}

		} else {
			return decimal.Decimal{}, fmt.Errorf("tipe komponen tidak valid: %s", komponen.TipeKomponen)
		}

		totalHPP = totalHPP.Add(komponenHPP.Mul(komponen.Kuantitas)) // <<< GUNAKAN METHOD .Add() dan .Mul()
	}

	memo[resepID] = totalHPP
	return totalHPP, nil
}

// GetHPPForResep (Diperbarui untuk menggunakan decimal.Decimal)
func GetHPPForResep(c *gin.Context) {
	resepID := c.Param("resep_id")

	var resep models.Resep
	if err := database.DB.Preload("Komponen").First(&resep, "id = ?", resepID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Resep tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil resep: " + err.Error()})
		return
	}

	// Memuat data master ke cache.
	if err := loadMasterDataIntoCache(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memuat data master untuk perhitungan HPP: " + err.Error()})
		return
	}

	hpp, err := calculateHPPRecursive(resep.ID, make(map[string]decimal.Decimal)) // <<< UBAH TIPE MAP
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghitung HPP: " + err.Error()})
		return
	}

	hppPerPorsi := hpp
	if resep.JumlahPorsi.IsZero() || resep.JumlahPorsi.IsNegative() {
		fmt.Printf("Peringatan: Resep '%s' (ID: %s) memiliki JumlahPorsi 0 atau negatif. HPP per porsi akan sama dengan HPP per unit.\n", resep.Nama, resep.ID)
	} else {
		hppPerPorsi = hpp.Div(resep.JumlahPorsi) // <<< GUNAKAN METHOD .Div()
	}

	c.JSON(http.StatusOK, HPPResult{
		ResepID:     resep.ID,
		ResepNama:   resep.Nama,
		HPPPerUnit:  hpp,
		HPPPerPorsi: hppPerPorsi,
	})
}