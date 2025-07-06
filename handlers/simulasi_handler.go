package handlers

import (
	"fmt"
	"math" // Pastikan import ini ada
	"net/http"

	"backend_kalkuliner/database"
	"backend_kalkuliner/models"
	"backend_kalkuliner/utils" // Untuk utils.RoundFloat

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SimulasiInput adalah struktur input untuk simulasi promo dari frontend
type SimulasiInput struct {
	// Bagian 1: Pilih Menu (Data Awal untuk Simulasi)
	HargaJualKotorProduk          float64 `json:"harga_jual_kotor_produk" binding:"required,gt=0"`
	HPPProduk                     float64 `json:"hpp_produk" binding:"required,gt=0"`
	NamaMenu                      string  `json:"nama_menu"`
	ChannelMenu                   string  `json:"channel_menu"`

	// Bagian 2: Set Ketentuan & Pilih Promo
	JumlahPorsiPembelian          float64 `json:"jumlah_porsi_pembelian" binding:"required,gt=0"`

	// Ongkir
	IsPromoOngkir                 bool    `json:"is_promo_ongkir"`
	SimulatedOngkirDitanggungMerchant float64 `json:"simulated_ongkir_ditanggung_merchant" binding:"gte=0"`

	// Promo Channel
	IsPakaiPromoChannel           bool    `json:"is_pakai_promo_channel"`
	SelectedPromoID               string  `json:"selected_promo_id" binding:"required_if=IsPakaiPromoChannel true"`

	// Komisi & Pajak (untuk simulasi)
	SimulatedKomisiChannelPersen float64 `json:"simulated_komisi_channel_persen" binding:"gte=0,lte=100"`
	SimulatedPajakPersen         float64 `json:"simulated_pajak_persen" binding:"gte=0,lte=100"`
}

// SimulasiResult adalah struktur output dari perhitungan simulasi ke frontend
type SimulasiResult struct {
	// =======================================================
	// Kategori 1: Data Awal (sesuai gambar)
	// =======================================================
	NamaMenu                    string  `json:"nama_menu"`
	ChannelMenu                 string  `json:"channel_menu"`
	JumlahPorsiPembelian        float64 `json:"jumlah_porsi_pembelian"`
	HPPProdukTotal              float64 `json:"hpp_produk_total"`
	HargaJualKotorProduk        float64 `json:"harga_jual_kotor_produk"`
	HargaJualTotalKotor         float64 `json:"harga_jual_total_kotor"`

	// =======================================================
	// Kategori 2: Detail Promo Channel Terpilih (sesuai gambar)
	// =======================================================
	NamaPromoTerpilih           string  `json:"nama_promo_terpilih"`
	JenisDiskonPromo            string  `json:"jenis_diskon_promo"`
	BesarDiskonPromo            float64 `json:"besar_diskon_promo"`
	MinBelanjaPromo             float64 `json:"min_belanja_promo"`
	MaksimalPotonganPromo       float64 `json:"maksimal_potongan_promo"`
	DitanggungMerchantPromoPersen float64 `json:"ditanggung_merchant_promo_persen"`
	CatatanPromo                string  `json:"catatan_promo"`
	PromoApplied                bool    `json:"promo_applied"` // Apakah promo benar-benar diterapkan

	// =======================================================
	// Kategori 3: Hasil Perhitungan (sesuai gambar)
	// =======================================================

	// Sub-kategori: Bagi Konsumen
	HargaJualUntukKonsumen      float64 `json:"harga_jual_untuk_konsumen"`
	DiskonPromoKonsumen         float64 `json:"diskon_promo_konsumen"`
	HargaAkhirKonsumen          float64 `json:"harga_akhir_konsumen"`

	// Sub-kategori: Biaya Promo Channel
	PotonganPromoDitanggungChannel float64 `json:"potongan_promo_ditanggung_channel"`
	PotonganPromoDitanggungMerchant float64 `json:"potongan_promo_ditanggung_merchant"`
	BiayaKomisiChannel            float64 `json:"biaya_komisi_channel"`
	BiayaPajak                    float64 `json:"biaya_pajak"`
	BiayaSubsidiOngkir            float64 `json:"biaya_subsidi_ongkir"`

	// Sub-kategori: Perhitungan Net Sales
	SalesSebelumKomisiPajakOngkir float64 `json:"sales_sebelum_komisi_pajak_ongkir"`
	NetSales                      float64 `json:"net_sales"`

	// Sub-kategori: Hasil Akhir
	GrossProfit                   float64 `json:"gross_profit"`
	HPPTerhadapNetSalesPersen     float64 `json:"hpp_terhadap_net_sales_persen"`
	GrossProfitTerhadapNetSalesPersen float64 `json:"gross_profit_terhadap_net_sales_persen"`
}

// SimulatePromoAndCommission menghitung simulasi promo, komisi, pajak, dan ongkir
func SimulatePromoAndCommission(c *gin.Context) {
	fmt.Println(">>> SimulatePromoAndCommission handler dipanggil <<<") // Debug log

	var input SimulasiInput
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Printf("Error binding JSON: %v\n", err) // Debug log error binding
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid input: %v", err.Error())})
		return
	}
	fmt.Printf("Input diterima: %+v\n", input) // Debug log input yang ter-bind

	// Inisialisasi hasil simulasi dengan nilai-nilai awal dari input
	simulasiResult := SimulasiResult{
		NamaMenu:                  input.NamaMenu,
		ChannelMenu:               input.ChannelMenu,
		JumlahPorsiPembelian:      utils.RoundFloat(input.JumlahPorsiPembelian, 0), // Bulatkan porsi ke integer terdekat
		HargaJualKotorProduk:      utils.RoundFloat(input.HargaJualKotorProduk, 2), // Bulatkan harga per unit
	}

	// Hitungan Awal (berdasarkan input)
	simulasiResult.HPPProdukTotal = input.HPPProduk * input.JumlahPorsiPembelian
	simulasiResult.HPPProdukTotal = utils.RoundFloat(simulasiResult.HPPProdukTotal, 4) // Pembulatan 4 desimal
	simulasiResult.HargaJualTotalKotor = input.HargaJualKotorProduk * input.JumlahPorsiPembelian
	simulasiResult.HargaJualTotalKotor = utils.RoundFloat(simulasiResult.HargaJualTotalKotor, 4) // Pembulatan 4 desimal
	simulasiResult.HargaJualUntukKonsumen = simulasiResult.HargaJualTotalKotor // Harga awal yang dilihat konsumen

	fmt.Printf("Data awal terhitung. HPP Total: %.4f, HJK Total: %.4f\n", simulasiResult.HPPProdukTotal, simulasiResult.HargaJualTotalKotor) // Debug log


	// 1. Ambil dan set Detail Promo Channel Terpilih
	if input.IsPakaiPromoChannel && input.SelectedPromoID != "" {
		fmt.Println("Mencoba mengambil detail promo...") // Debug log
		var promoProgram models.ProgramPromo
		if err := database.DB.First(&promoProgram, "id = ?", input.SelectedPromoID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Program promo tidak ditemukan. Harap refresh halaman dan pilih promo yang valid."})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil program promo: " + err.Error()})
			}
			return
		}
		simulasiResult.NamaPromoTerpilih = promoProgram.NamaPromo
		simulasiResult.JenisDiskonPromo = promoProgram.JenisDiskon
		simulasiResult.BesarDiskonPromo = utils.RoundFloat(promoProgram.BesarDiskon, 2)
		simulasiResult.MinBelanjaPromo = utils.RoundFloat(promoProgram.MinBelanja, 2)
		simulasiResult.MaksimalPotonganPromo = utils.RoundFloat(promoProgram.MaksimalPotongan, 2)
		simulasiResult.DitanggungMerchantPromoPersen = utils.RoundFloat(promoProgram.DitanggungMerchantPersen, 2)
		simulasiResult.CatatanPromo = promoProgram.Catatan
		fmt.Printf("Detail promo ditemukan: %+v\n", promoProgram) // Debug log
	} else {
		fmt.Println("Tidak pakai promo channel atau promo ID kosong.") // Debug log
	}

	// 2. Hitung Diskon Promo Konsumen (Total Potongan Promo)
	fmt.Println("Menghitung diskon promo konsumen...") // Debug log
	diskonPromoKonsumen := 0.0
	potonganDitanggungMerchant := 0.0
	potonganDitanggungChannel := 0.0

	if simulasiResult.NamaPromoTerpilih != "" && simulasiResult.HargaJualTotalKotor >= simulasiResult.MinBelanjaPromo {
		simulasiResult.PromoApplied = true
		if simulasiResult.JenisDiskonPromo == "persentase" {
			potonganAwal := simulasiResult.HargaJualTotalKotor * (simulasiResult.BesarDiskonPromo / 100.0)
			diskonPromoKonsumen = math.Min(potonganAwal, simulasiResult.MaksimalPotonganPromo)
		} else if simulasiResult.JenisDiskonPromo == "nominal" {
			diskonPromoKonsumen = simulasiResult.BesarDiskonPromo
		}

		potonganDitanggungMerchant = diskonPromoKonsumen * (simulasiResult.DitanggungMerchantPromoPersen / 100.0)
		potonganDitanggungChannel = diskonPromoKonsumen - potonganDitanggungMerchant
	}

	simulasiResult.DiskonPromoKonsumen = utils.RoundFloat(diskonPromoKonsumen, 2)
	simulasiResult.PotonganPromoDitanggungMerchant = utils.RoundFloat(potonganDitanggungMerchant, 2)
	simulasiResult.PotonganPromoDitanggungChannel = utils.RoundFloat(potonganDitanggungChannel, 2)
	fmt.Printf("Diskon promo konsumen terhitung: %.2f\n", simulasiResult.DiskonPromoKonsumen)

	// 3. Hitung Harga Akhir Konsumen
	fmt.Println("Menghitung harga akhir konsumen...") // Debug log
	simulasiResult.HargaAkhirKonsumen = simulasiResult.HargaJualUntukKonsumen - simulasiResult.DiskonPromoKonsumen
	simulasiResult.HargaAkhirKonsumen = utils.RoundFloat(simulasiResult.HargaAkhirKonsumen, 2)
	if simulasiResult.HargaAkhirKonsumen < 0 {
		simulasiResult.HargaAkhirKonsumen = 0.0
	}
	fmt.Printf("Harga akhir konsumen: %.2f\n", simulasiResult.HargaAkhirKonsumen)

	// 4. Hitung Biaya Komisi Channel (berdasarkan HargaJualTotalKotor)
	fmt.Println("Menghitung biaya komisi channel...") // Debug log
	simulasiResult.BiayaKomisiChannel = simulasiResult.HargaJualTotalKotor * (input.SimulatedKomisiChannelPersen / 100.0)
	simulasiResult.BiayaKomisiChannel = utils.RoundFloat(simulasiResult.BiayaKomisiChannel, 2)
	fmt.Printf("Biaya komisi channel: %.2f\n", simulasiResult.BiayaKomisiChannel)


	// 5. Hitung Biaya Pajak (berdasarkan HargaJualTotalKotor)
	fmt.Println("Menghitung biaya pajak...") // Debug log
	simulasiResult.BiayaPajak = simulasiResult.HargaJualTotalKotor * (input.SimulatedPajakPersen / 100.0)
	simulasiResult.BiayaPajak = utils.RoundFloat(simulasiResult.BiayaPajak, 2)
	fmt.Printf("Biaya pajak: %.2f\n", simulasiResult.BiayaPajak)


	// 6. Hitung Biaya Subsidi Ongkir
	fmt.Println("Menghitung biaya subsidi ongkir...") // Debug log
	if input.IsPromoOngkir {
		simulasiResult.BiayaSubsidiOngkir = utils.RoundFloat(input.SimulatedOngkirDitanggungMerchant, 2)
	} else {
		simulasiResult.BiayaSubsidiOngkir = 0.0
	}
	fmt.Printf("Biaya subsidi ongkir: %.2f\n", simulasiResult.BiayaSubsidiOngkir)


	// 7. Hitung Sales Sebelum Komisi/Pajak/Ongkir
	fmt.Println("Menghitung sales sebelum komisi/pajak/ongkir...") // Debug log
	simulasiResult.SalesSebelumKomisiPajakOngkir = simulasiResult.HargaJualTotalKotor - simulasiResult.PotonganPromoDitanggungMerchant
	simulasiResult.SalesSebelumKomisiPajakOngkir = utils.RoundFloat(simulasiResult.SalesSebelumKomisiPajakOngkir, 2)
	if simulasiResult.SalesSebelumKomisiPajakOngkir < 0 { // Pastikan tidak negatif
		simulasiResult.SalesSebelumKomisiPajakOngkir = 0.0
	}
	fmt.Printf("Sales sebelum komisi/pajak/ongkir: %.2f\n", simulasiResult.SalesSebelumKomisiPajakOngkir)


	// 8. Hitung Net Sales
	fmt.Println("Menghitung net sales...") // Debug log
	simulasiResult.NetSales = simulasiResult.SalesSebelumKomisiPajakOngkir - simulasiResult.BiayaKomisiChannel - simulasiResult.BiayaPajak - simulasiResult.BiayaSubsidiOngkir
	simulasiResult.NetSales = utils.RoundFloat(simulasiResult.NetSales, 2)
	fmt.Printf("Net sales: %.2f\n", simulasiResult.NetSales)


	// 9. Hitung Gross Profit
	fmt.Println("Menghitung gross profit...") // Debug log
	simulasiResult.GrossProfit = simulasiResult.NetSales - simulasiResult.HPPProdukTotal
	simulasiResult.GrossProfit = utils.RoundFloat(simulasiResult.GrossProfit, 2)
	fmt.Printf("Gross profit: %.2f\n", simulasiResult.GrossProfit)


	// 10. Hitung HPP terhadap Net Sales (%)
	fmt.Println("Menghitung persentase HPP terhadap net sales...") // Debug log
	if simulasiResult.NetSales == 0 {
		simulasiResult.HPPTerhadapNetSalesPersen = 0.0
	} else {
		simulasiResult.HPPTerhadapNetSalesPersen = (simulasiResult.HPPProdukTotal / simulasiResult.NetSales) * 100.0
	}
	simulasiResult.HPPTerhadapNetSalesPersen = utils.RoundFloat(simulasiResult.HPPTerhadapNetSalesPersen, 2)
	fmt.Printf("HPP terhadap net sales: %.2f%%\n", simulasiResult.HPPTerhadapNetSalesPersen)


	// 11. Hitung Gross Profit terhadap Net Sales (%)
	fmt.Println("Menghitung persentase gross profit terhadap net sales...") // Debug log
	if simulasiResult.NetSales == 0 {
		simulasiResult.GrossProfitTerhadapNetSalesPersen = 0.0
	} else {
		simulasiResult.GrossProfitTerhadapNetSalesPersen = (simulasiResult.GrossProfit / simulasiResult.NetSales) * 100.0
	}
	simulasiResult.GrossProfitTerhadapNetSalesPersen = utils.RoundFloat(simulasiResult.GrossProfitTerhadapNetSalesPersen, 2)
	fmt.Printf("Gross profit terhadap net sales: %.2f%%\n", simulasiResult.GrossProfitTerhadapNetSalesPersen)


	fmt.Println("Perhitungan simulasi selesai. Mengirim response.")
	c.JSON(http.StatusOK, simulasiResult)
}