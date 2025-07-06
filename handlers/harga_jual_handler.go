package handlers

import (
	"fmt"
	"math"
	"net/http"

	"backend_kalkuliner/database"
	"backend_kalkuliner/models"
	"backend_kalkuliner/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CalculateHargaJualInput adalah struktur input komprehensif untuk perhitungan dan penyimpanan harga jual
// MENGGABUNGKAN semua field kriteria perhitungan optimal
type CalculateHargaJualInput struct {
	ResepID             string  `json:"resep_id" binding:"required"`
	NamaProduk          string  `json:"nama_produk" binding:"required"`
	Channel             string  `json:"channel" binding:"required"` // Channel penjualan
	JumlahPorsiProduk   float64 `json:"jumlah_porsi_produk" binding:"required,gt=0"`

	// --- Kriteria Perhitungan Optimal (hanya satu yang akan diisi) ---
	SelectedCriteria string `json:"selectedCriteria" binding:"required"` // Kriteria mana yang dipilih

	MinProfitNetSalesPersen    *float64 `json:"min_profit_net_sales_persen"`    // % dari Net Sales
	MinProfitRpHPP             *float64 `json:"min_profit_rp_hpp"`              // Rp dari HPP
	MinProfitPersenHPP         *float64 `json:"min_profit_persen_hpp"`          // % dari HPP
	MinProfitXLipatHPP         *float64 `json:"min_profit_x_lipat_hpp"`         // x Lipat dari HPP
	MaxHPPNetSalesPersen       *float64 `json:"max_hpp_net_sales_persen"`       // % dari Net Sales
	TargetNetSalesXLipatHPP    *float64 `json:"target_net_sales_x_lipat_hpp"`    // x Lipat dari HPP
	TargetNetSalesRp           *float64 `json:"target_net_sales_rp"`            // Rp
	TargetHargaJualRp          *float64 `json:"target_harga_jual_rp"`           // Rp (Ini adalah HJK langsung)
	ConsumerPaysIncludingTaxRp *float64 `json:"consumer_pays_including_tax_rp"` // Rp (Ini adalah HJK, diasumsikan sudah termasuk pajak)
	TargetHargaJualExclTaxRp   *float64 `json:"target_harga_jual_excl_tax_rp"`  // Rp (Ini adalah HJK, diasumsikan belum termasuk pajak)

	// Biaya Operasional (juga digunakan untuk kalkulasi optimal)
	PajakPersen         float64 `json:"pajak_persen" binding:"gte=0,lte=100"`
	KomisiChannelPersen float64 `json:"komisi_channel_persen" binding:"gte=0,lte=100"`

	// Ini akan menjadi hasil kalkulator optimal yang disimpan
	// `MetodePerhitungan` dan `NilaiKriteria` di struct models.HargaJual
	// akan diisi berdasarkan hasil kalkulasi di backend.
	HargaBulat bool `json:"harga_bulat"` // Untuk pembulatan hasil optimal
}

// HargaJualResponse adalah DTO untuk hasil perhitungan harga jual (dikirim ke frontend)
// Ini juga akan menjadi response untuk perhitungan optimal
type HargaJualResponse struct {
	ID                 string  `json:"id,omitempty"`
	ResepID            string  `json:"resep_id"`
	NamaProduk         string  `json:"nama_produk"`
	Channel            string  `json:"channel"`
	HPP                float64 `json:"hpp"`
	JumlahPorsiProduk  float64 `json:"jumlah_porsi_produk"`
	MetodePerhitungan  string  `json:"metode_perhitungan"` // Ini akan menyimpan kriteria optimal yang dipilih
	NilaiKriteria      float64 `json:"nilai_kriteria"`     // Ini akan menyimpan calculatedHargaJualKotor
	PajakPersen        float64 `json:"pajak_persen"`
	KomisiChannelPersen float64 `json:"komisi_channel_persen"`

	HargaJualKotor     float64 `json:"harga_jual_kotor"`
	HargaJualBersih    float64 `json:"harga_jual_bersih"`
	TotalPajak         float64 `json:"total_pajak"`
	TotalKomisi        float64 `json:"total_komisi"`
	Profit             float64 `json:"profit"`
	ProfitPersen       float64 `json:"profit_persen"`

	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
	ResepNama string `json:"resep_nama,omitempty"`

	// Field tambahan untuk hasil perhitungan optimal (digunakan saat response)
	MetodeTerkalkulasi string `json:"metode_terkalkulasi,omitempty"` // Kriteria mana yang akhirnya digunakan (untuk display)
}


// calculateHargaJualOptimalLogic adalah fungsi inti yang menghitung HJK dari berbagai kriteria.
// Ini adalah fungsionalitas inti dari "Kalkulator Harga Jual Optimal" yang kini terintegrasi.
func calculateHargaJualOptimalLogic(
	hppProdukTotal float64,
	pajakPersen float64,
	komisiChannelPersen float64,
	inputCriteria CalculateHargaJualInput, // Menerima seluruh input untuk mengecek kriteria
) (
	calculatedHargaJualKotor float64,
	metodeTerkalkulasi string,
	err error,
) {
	calculatedHargaJualKotor = 0.0
	metodeTerkalkulasi = ""
	err = nil

	// Validasi dasar biaya operasional
	totalBiayaOperasionalPersen := (komisiChannelPersen + pajakPersen) / 100.0
	pembagiBiayaOperasional := 1.0 - totalBiayaOperasionalPersen
	if pembagiBiayaOperasional <= 0 {
		err = fmt.Errorf("Total Komisi dan Pajak tidak boleh 100%% atau lebih.")
		return
	}

	// === Logika Perhitungan untuk Setiap Kriteria ===
	// Perhitungan ini akan bekerja mundur atau maju untuk menemukan HargaJualKotor

	if inputCriteria.SelectedCriteria == "min_profit_net_sales_persen" && inputCriteria.MinProfitNetSalesPersen != nil {
		metodeTerkalkulasi = "MinProfitNetSalesPersen"
		profitPersenDariNetSales := *inputCriteria.MinProfitNetSalesPersen / 100.0
		if profitPersenDariNetSales >= 1.0 { // Profit 100% atau lebih dari net sales
			err = fmt.Errorf("Profit margin dari net sales tidak boleh 100%% atau lebih.")
			return
		}
		netSalesTarget := hppProdukTotal / (1.0 - profitPersenDariNetSales)
		calculatedHargaJualKotor = netSalesTarget / pembagiBiayaOperasional

	} else if inputCriteria.SelectedCriteria == "min_profit_rp_hpp" && inputCriteria.MinProfitRpHPP != nil {
		metodeTerkalkulasi = "MinProfitRpHPP"
		profitNominal := *inputCriteria.MinProfitRpHPP
		hargaJualBersihTarget := hppProdukTotal + profitNominal
		calculatedHargaJualKotor = hargaJualBersihTarget / pembagiBiayaOperasional

	} else if inputCriteria.SelectedCriteria == "min_profit_persen_hpp" && inputCriteria.MinProfitPersenHPP != nil {
		metodeTerkalkulasi = "MinProfitPersenHPP"
		profitPersenDariHPP := *inputCriteria.MinProfitPersenHPP / 100.0
		hargaJualBersihTarget := hppProdukTotal * (1.0 + profitPersenDariHPP)
		calculatedHargaJualKotor = hargaJualBersihTarget / pembagiBiayaOperasional

	} else if inputCriteria.SelectedCriteria == "min_profit_x_lipat_hpp" && inputCriteria.MinProfitXLipatHPP != nil {
		metodeTerkalkulasi = "MinProfitXLipatHPP"
		profitXLipat := *inputCriteria.MinProfitXLipatHPP
		hargaJualBersihTarget := hppProdukTotal * (1.0 + profitXLipat)
		calculatedHargaJualKotor = hargaJualBersihTarget / pembagiBiayaOperasional

	} else if inputCriteria.SelectedCriteria == "max_hpp_net_sales_persen" && inputCriteria.MaxHPPNetSalesPersen != nil {
		metodeTerkalkulasi = "MaxHPPNetSalesPersen"
		maxHPPNetSalesPersen := *inputCriteria.MaxHPPNetSalesPersen / 100.0
		if maxHPPNetSalesPersen <= 0 {
			err = fmt.Errorf("Persentase HPP maksimal dari net sales harus lebih dari 0.")
			return
		}
		netSalesTarget := hppProdukTotal / maxHPPNetSalesPersen
		calculatedHargaJualKotor = netSalesTarget / pembagiBiayaOperasional

	} else if inputCriteria.SelectedCriteria == "target_net_sales_x_lipat_hpp" && inputCriteria.TargetNetSalesXLipatHPP != nil {
		metodeTerkalkulasi = "TargetNetSalesXLipatHPP"
		netSalesXLipat := *inputCriteria.TargetNetSalesXLipatHPP
		netSalesTarget := hppProdukTotal * netSalesXLipat
		calculatedHargaJualKotor = netSalesTarget / pembagiBiayaOperasional

	} else if inputCriteria.SelectedCriteria == "target_net_sales_rp" && inputCriteria.TargetNetSalesRp != nil {
		metodeTerkalkulasi = "TargetNetSalesRp"
		netSalesTarget := *inputCriteria.TargetNetSalesRp
		calculatedHargaJualKotor = netSalesTarget / pembagiBiayaOperasional

	} else if inputCriteria.SelectedCriteria == "target_harga_jual_rp" && inputCriteria.TargetHargaJualRp != nil {
		metodeTerkalkulasi = "TargetHargaJualRp"
		calculatedHargaJualKotor = *inputCriteria.TargetHargaJualRp

	} else if inputCriteria.SelectedCriteria == "consumer_pays_including_tax_rp" && inputCriteria.ConsumerPaysIncludingTaxRp != nil {
		metodeTerkalkulasi = "ConsumerPaysIncludingTaxRp"
		calculatedHargaJualKotor = *inputCriteria.ConsumerPaysIncludingTaxRp // Diasumsikan ini adalah HJK yang sudah termasuk pajak

	} else if inputCriteria.SelectedCriteria == "target_harga_jual_excl_tax_rp" && inputCriteria.TargetHargaJualExclTaxRp != nil {
		metodeTerkalkulasi = "TargetHargaJualExclTaxRp"
		hargaJualBersihTarget := *inputCriteria.TargetHargaJualExclTaxRp
		calculatedHargaJualKotor = hargaJualBersihTarget / pembagiBiayaOperasional
	} else {
		err = fmt.Errorf("Kriteria perhitungan harga jual tidak valid atau tidak dipilih.")
		return
	}

	// Validasi dasar agar tidak ada hasil negatif atau sangat besar tak terhingga
	if calculatedHargaJualKotor <= 0 || math.IsInf(calculatedHargaJualKotor, 0) || math.IsNaN(calculatedHargaJualKotor) {
		err = fmt.Errorf("Hasil perhitungan harga jual tidak valid. Periksa input dan kriteria.")
		return
	}

	return
}


// CalculateAndSaveHargaJual menghitung dan menyimpan harga jual baru (Kini Menggabungkan Kalkulasi Optimal)
func CalculateAndSaveHargaJual(c *gin.Context) {
	var input CalculateHargaJualInput // Input ini sekarang komprehensif
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid input: %v", err.Error())})
		return
	}

	// Validasi awal input yang esensial (sebelumnya ada di main.go, sekarang di sini)
	if input.JumlahPorsiProduk <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jumlah porsi produk harus lebih dari 0."})
		return
	}
	if input.ResepID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resep harus dipilih."})
		return
	}
	if input.NamaProduk == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama produk tidak boleh kosong."})
		return
	}
	if input.Channel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Channel penjualan tidak boleh kosong."})
		return
	}


	// Ambil HPP resep (terbaru)
	var hppResult models.HPPResult
	if err := database.DB.Where("resep_id = ?", input.ResepID).Order("created_at DESC").First(&hppResult).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "HPP untuk resep ini belum dihitung. Harap hitung HPP terlebih dahulu."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil HPP resep: " + err.Error()})
		return
	}

	// Dapatkan detail resep dari cache yang diekspor
	resepDetail, ok := ExportedResepCache[input.ResepID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Detail resep tidak ditemukan di cache. Coba restart server."})
		return
	}

	// Panggil logika perhitungan harga jual optimal yang sekarang terintegrasi
	// HPPProdukTotal dari hppResult.HPPPerUnit (HPP Total resep)
	calculatedHargaJualKotor, metodeTerkalkulasi, errCalc := calculateHargaJualOptimalLogic(
		hppResult.HPPPerPorsi,
		input.PajakPersen,
		input.KomisiChannelPersen,
		input, // Teruskan seluruh input untuk kriteria
	)
	if errCalc != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kesalahan perhitungan optimal: " + errCalc.Error()})
		return
	}

	// Lanjutkan dengan perhitungan breakdown berdasarkan calculatedHargaJualKotor ini
	totalKomisi := calculatedHargaJualKotor * (input.KomisiChannelPersen / 100.0)
	totalPajak := calculatedHargaJualKotor * (input.PajakPersen / 100.0)
	hargaJualBersih := calculatedHargaJualKotor - totalKomisi - totalPajak
	netSales := hargaJualBersih // Jika tidak ada promo/ongkir tambahan di sini

	profit := netSales - hppResult.HPPPerPorsi // Profit dari NetSales dikurangi HPP Total
	profitPersen := 0.0
	if hppResult.HPPPerPorsi > 0 {
		profitPersen = (profit / hppResult.HPPPerPorsi) * 100.0
	}


	// Simpan Hasil Perhitungan ke database
	hargaJual := models.HargaJual{
		ResepID:             input.ResepID,
		NamaProduk:          input.NamaProduk,
		Channel:             input.Channel,
		HPP:                 utils.RoundFloat(hppResult.HPPPerPorsi, 4), // Simpan HPP Total resep sebagai HPP dasar produk
		JumlahPorsiProduk:   utils.RoundFloat(input.JumlahPorsiProduk, 4),
		// Metode dan Nilai Kriteria akan mencerminkan hasil optimal
		MetodePerhitungan:   metodeTerkalkulasi, // <<< Simpan metode kriteria optimal
		NilaiKriteria:       utils.RoundFloat(calculatedHargaJualKotor, 4), // <<< Simpan HJKotor sebagai nilai kriteria
		PajakPersen:         utils.RoundFloat(input.PajakPersen, 2),
		KomisiChannelPersen: utils.RoundFloat(input.KomisiChannelPersen, 2),
		HargaJualKotor:      utils.RoundFloat(calculatedHargaJualKotor, 4),
		HargaJualBersih:     utils.RoundFloat(hargaJualBersih, 4),
		TotalPajak:          utils.RoundFloat(totalPajak, 4),
		TotalKomisi:         utils.RoundFloat(totalKomisi, 4),
		Profit:              utils.RoundFloat(profit, 4),
		ProfitPersen:        utils.RoundFloat(profitPersen, 2),
	}

	if err := database.DB.Create(&hargaJual).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan harga jual: " + err.Error()})
		return
	}

	// Kirim Response ke frontend
	c.JSON(http.StatusCreated, HargaJualResponse{
		ID:                  hargaJual.ID,
		ResepID:             hargaJual.ResepID,
		NamaProduk:          hargaJual.NamaProduk,
		Channel:             hargaJual.Channel,
		HPP:                 hargaJual.HPP,
		JumlahPorsiProduk:   hargaJual.JumlahPorsiProduk,
		MetodePerhitungan:   hargaJual.MetodePerhitungan, // Metode dari optimal
		NilaiKriteria:       hargaJual.NilaiKriteria,     // Nilai HJKotor hasil optimal
		PajakPersen:         hargaJual.PajakPersen,
		KomisiChannelPersen: hargaJual.KomisiChannelPersen,
		HargaJualKotor:      hargaJual.HargaJualKotor,
		HargaJualBersih:     hargaJual.HargaJualBersih,
		TotalPajak:          hargaJual.TotalPajak,
		TotalKomisi:         hargaJual.TotalKomisi,
		Profit:              hargaJual.Profit,
		ProfitPersen:        hargaJual.ProfitPersen,
		CreatedAt:           hargaJual.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:           hargaJual.UpdatedAt.Format("2006-01-02 15:04:05"),
		ResepNama:           resepDetail.Nama,
		MetodeTerkalkulasi:  metodeTerkalkulasi, // Ini adalah kriteria yang benar-benar digunakan untuk kalkulasi
	})
}

// GetHargaJuals mengambil semua harga jual yang tersimpan (tidak berubah)
func GetHargaJuals(c *gin.Context) {
	var hargaJuals []models.HargaJual
	if err := database.DB.Preload("Resep").Find(&hargaJuals).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil daftar harga jual"})
		return
	}

	var responses []HargaJualResponse
	for _, hj := range hargaJuals {
		resepNama := "N/A"
		if hj.Resep.Nama != "" {
			resepNama = hj.Resep.Nama
		}
		responses = append(responses, HargaJualResponse{
			ID:                  hj.ID,
			ResepID:             hj.ResepID,
			NamaProduk:          hj.NamaProduk,
			Channel:             hj.Channel,
			HPP:                 hj.HPP,
			JumlahPorsiProduk:   hj.JumlahPorsiProduk,
			MetodePerhitungan:   hj.MetodePerhitungan,
			NilaiKriteria:       hj.NilaiKriteria,
			PajakPersen:         hj.PajakPersen,
			KomisiChannelPersen: hj.KomisiChannelPersen,
			HargaJualKotor:      hj.HargaJualKotor,
			HargaJualBersih:     hj.HargaJualBersih,
			TotalPajak:          hj.TotalPajak,
			TotalKomisi:         hj.TotalKomisi,
			Profit:              hj.Profit,
			ProfitPersen:        hj.ProfitPersen,
			CreatedAt:           hj.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:           hj.UpdatedAt.Format("2006-01-02 15:04:05"),
			ResepNama:           resepNama,
		})
	}
	c.JSON(http.StatusOK, responses)
}

// GetHargaJualByID mengambil satu harga jual berdasarkan ID (Diperbarui untuk response)
func GetHargaJualByID(c *gin.Context) {
	id := c.Param("id")
	var hargaJual models.HargaJual
	if err := database.DB.Preload("Resep").First(&hargaJual, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Harga jual tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil harga jual"})
		return
	}

	resepNama := "N/A"
	if hargaJual.Resep.Nama != "" {
		resepNama = hargaJual.Resep.Nama
	}

	c.JSON(http.StatusOK, HargaJualResponse{
		ID:                  hargaJual.ID,
		ResepID:             hargaJual.ResepID,
		NamaProduk:          hargaJual.NamaProduk,
		Channel:             hargaJual.Channel,
		HPP:                 hargaJual.HPP,
		JumlahPorsiProduk:   hargaJual.JumlahPorsiProduk,
		MetodePerhitungan:   hargaJual.MetodePerhitungan,
		NilaiKriteria:       hargaJual.NilaiKriteria,
		PajakPersen:         hargaJual.PajakPersen,
		KomisiChannelPersen: hargaJual.KomisiChannelPersen,
		HargaJualKotor:      hargaJual.HargaJualKotor,
		HargaJualBersih:     hargaJual.HargaJualBersih,
		TotalPajak:          hargaJual.TotalPajak,
		TotalKomisi:         hargaJual.TotalKomisi,
		Profit:              hargaJual.Profit,
		ProfitPersen:        hargaJual.ProfitPersen,
		CreatedAt:           hargaJual.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:           hargaJual.UpdatedAt.Format("2006-01-02 15:04:05"),
		ResepNama:           resepNama,
		MetodeTerkalkulasi:  hargaJual.MetodePerhitungan, // Metode yang digunakan saat disimpan
	})
}

// UpdateHargaJual memperbarui data harga jual yang sudah ada (Diperbarui untuk menggunakan logika optimal)
func UpdateHargaJual(c *gin.Context) {
	id := c.Param("id")
	var existingHargaJual models.HargaJual
	if err := database.DB.First(&existingHargaJual, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Harga jual tidak ditemukan"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil harga jual: " + err.Error()})
		return
	}

	var input CalculateHargaJualInput // Menggunakan input komprehensif
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid input: %v", err.Error())})
		return
	}

	// Validasi awal input yang esensial
	if input.JumlahPorsiProduk <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Jumlah porsi produk harus lebih dari 0."})
		return
	}
	if input.ResepID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resep harus dipilih."})
		return
	}
	if input.NamaProduk == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nama produk tidak boleh kosong."})
		return
	}
	if input.Channel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Channel penjualan tidak boleh kosong."})
		return
	}

	// Ambil HPP resep
	var hppResult models.HPPResult
	if err := database.DB.Where("resep_id = ?", input.ResepID).Order("created_at DESC").First(&hppResult).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "HPP untuk resep ini belum dihitung. Harap hitung HPP terlebih dahulu."})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil HPP resep: " + err.Error()})
		return
	}

	resepDetail, ok := ExportedResepCache[input.ResepID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Detail resep tidak ditemukan di cache. Coba restart server."})
		return
	}

	// Panggil logika perhitungan harga jual optimal yang sekarang terintegrasi
	calculatedHargaJualKotor, metodeTerkalkulasi, errCalc := calculateHargaJualOptimalLogic(
		hppResult.HPPPerPorsi,
		input.PajakPersen,
		input.KomisiChannelPersen,
		input, // Teruskan seluruh input untuk kriteria
	)
	if errCalc != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Kesalahan perhitungan optimal: " + errCalc.Error()})
		return
	}

	// Lanjutkan dengan perhitungan breakdown berdasarkan calculatedHargaJualKotor ini
	totalKomisi := calculatedHargaJualKotor * (input.KomisiChannelPersen / 100.0)
	totalPajak := calculatedHargaJualKotor * (input.PajakPersen / 100.0)
	hargaJualBersih := calculatedHargaJualKotor - totalKomisi - totalPajak
	netSales := hargaJualBersih

	profit := netSales - hppResult.HPPPerPorsi
	profitPersen := 0.0
	if hppResult.HPPPerPorsi > 0 {
		profitPersen = (profit / hppResult.HPPPerPorsi) * 100.0
	}

	// Perbarui objek existingHargaJual dengan nilai-nilai baru
	existingHargaJual.ResepID = input.ResepID
	existingHargaJual.NamaProduk = input.NamaProduk
	existingHargaJual.Channel = input.Channel
	existingHargaJual.HPP = utils.RoundFloat(hppResult.HPPPerPorsi, 4)
	existingHargaJual.JumlahPorsiProduk = utils.RoundFloat(input.JumlahPorsiProduk, 4)
	existingHargaJual.MetodePerhitungan = metodeTerkalkulasi // <<< Simpan metode kriteria optimal
	existingHargaJual.NilaiKriteria = utils.RoundFloat(calculatedHargaJualKotor, 4) // <<< Simpan HJKotor sebagai nilai kriteria
	existingHargaJual.PajakPersen = utils.RoundFloat(input.PajakPersen, 2)
	existingHargaJual.KomisiChannelPersen = utils.RoundFloat(input.KomisiChannelPersen, 2)
	existingHargaJual.HargaJualKotor = utils.RoundFloat(calculatedHargaJualKotor, 4)
	existingHargaJual.HargaJualBersih = utils.RoundFloat(hargaJualBersih, 4)
	existingHargaJual.TotalPajak = utils.RoundFloat(totalPajak, 4)
	existingHargaJual.TotalKomisi = utils.RoundFloat(totalKomisi, 4)
	existingHargaJual.Profit = utils.RoundFloat(profit, 4)
	existingHargaJual.ProfitPersen = utils.RoundFloat(profitPersen, 2)

	// Simpan ke database
	if err := database.DB.Save(&existingHargaJual).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memperbarui harga jual: " + err.Error()})
		return
	}

	// Kirim Response yang diperbarui
	c.JSON(http.StatusOK, HargaJualResponse{
		ID:                  existingHargaJual.ID,
		ResepID:             existingHargaJual.ResepID,
		NamaProduk:          existingHargaJual.NamaProduk,
		Channel:             existingHargaJual.Channel,
		HPP:                 existingHargaJual.HPP,
		JumlahPorsiProduk:   existingHargaJual.JumlahPorsiProduk,
		MetodePerhitungan:   existingHargaJual.MetodePerhitungan,
		NilaiKriteria:       existingHargaJual.NilaiKriteria,
		PajakPersen:         existingHargaJual.PajakPersen,
		KomisiChannelPersen: existingHargaJual.KomisiChannelPersen,
		HargaJualKotor:      existingHargaJual.HargaJualKotor,
		HargaJualBersih:     existingHargaJual.HargaJualBersih,
		TotalPajak:          existingHargaJual.TotalPajak,
		TotalKomisi:         existingHargaJual.TotalKomisi,
		Profit:              existingHargaJual.Profit,
		ProfitPersen:        existingHargaJual.ProfitPersen,
		CreatedAt:           existingHargaJual.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:           existingHargaJual.UpdatedAt.Format("2006-01-02 15:04:05"),
		ResepNama:           resepDetail.Nama,
		MetodeTerkalkulasi:  metodeTerkalkulasi, // Gunakan hasil dari calculateHargaJualOptimalLogic
	})
}

// DeleteHargaJual menghapus harga jual berdasarkan ID (Tidak Berubah)
func DeleteHargaJual(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Delete(&models.HargaJual{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus harga jual"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Harga jual berhasil dihapus"})
}