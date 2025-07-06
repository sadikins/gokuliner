<script setup>
import { ref, onMounted, watch } from 'vue'
  import axios from 'axios'
import { formatCurrency, formatPercentage } from '@/utils/formatters'

const API_BASE_URL = 'http://localhost:8080/api'

// --- State Data ---
const hargaJuals = ref([]) // Daftar Harga Jual tersimpan (sebagai menu dasar untuk simulasi)
const programPromos = ref([]) // Daftar Program Promo yang tersedia

// --- State Input Simulasi ---
const simulasiInputForm = ref({
  // Bagian 1: Pilih Menu
  selectedHargaJualId: '', // ID Harga Jual yang dipilih dari dropdown
  nama_menu: '',          // Nama menu yang dipilih (diambil dari hargaJual)
  channel_menu: '',       // Channel menu yang dipilih (diambil dari hargaJual)
  hpp_produk: 0,          // HPP produk satuan (diambil dari hargaJual)
  harga_jual_kotor_produk: 0, // Harga jual produk satuan (diambil dari hargaJual)

  // Bagian 2: Set Ketentuan & Pilih Promo
  jumlah_porsi_pembelian: 1, // Jumlah porsi yang ingin disimulasikan

  // Ongkir
  is_promo_ongkir: false, // Checkbox untuk mengaktifkan/menonaktifkan promo ongkir
  simulated_ongkir_ditanggung_merchant: 0, // Nominal ongkir yang ditanggung merchant

  // Promo Channel
  is_pakai_promo_channel: false, // Checkbox untuk mengaktifkan/menonaktifkan promo channel
  selected_promo_id: '',         // ID Program Promo yang dipilih

  // Komisi & Pajak (untuk simulasi)
  simulated_komisi_channel_persen: 0, // Persentase komisi channel
  simulated_pajak_persen: 0           // Persentase pajak
})

// --- State Hasil Simulasi ---
const simulasiResult = ref(null) // Objek untuk menyimpan dan menampilkan hasil perhitungan simulasi

// --- Fetch Data Awal ---
// Mengambil daftar harga jual yang sudah tersimpan untuk dropdown "Pilih Menu"
const fetchHargaJuals = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/harga-juals`)
    hargaJuals.value = response.data
    console.log('Daftar Harga Jual berhasil diambil:', response.data);
  } catch (error) {
    console.error('Error fetching harga juals:', error)
    alert('Gagal mengambil daftar harga jual untuk pilihan menu. Pastikan ada data di modul Harga Jual.')
  }
}

// Mengambil daftar program promo yang tersedia untuk dropdown "Pilih Promo"
const fetchProgramPromos = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/program-promos`)
    programPromos.value = response.data
    console.log('Daftar Program Promo berhasil diambil:', response.data);
  } catch (error) {
    console.error('Error fetching program promos:', error)
    alert('Gagal mengambil daftar program promo.')
  }
}

// --- Watcher untuk Pilihan Menu ---
// Akan otomatis mengisi field HPP dan Harga Jual saat menu dipilih
watch(() => simulasiInputForm.value.selectedHargaJualId, (newId) => {
  if (newId) {
    const selectedHJ = hargaJuals.value.find(hj => hj.id === newId)
    if (selectedHJ) {
      simulasiInputForm.value.nama_menu = selectedHJ.nama_produk
      simulasiInputForm.value.channel_menu = selectedHJ.channel || 'N/A'
      simulasiInputForm.value.hpp_produk = parseFloat(selectedHJ.hpp)
      simulasiInputForm.value.harga_jual_kotor_produk = parseFloat(selectedHJ.harga_jual_kotor)
      // >>> AMBIL DARI selectedHJ <<<
      simulasiInputForm.value.simulated_komisi_channel_persen = parseFloat(selectedHJ.komisi_channel_persen)
      simulasiInputForm.value.simulated_pajak_persen = parseFloat(selectedHJ.pajak_persen)

      // Reset pilihan promo jika menu berubah
      simulasiInputForm.value.is_pakai_promo_channel = false
      simulasiInputForm.value.selected_promo_id = ''
      simulasiInputForm.value.is_promo_ongkir = false
      simulasiInputForm.value.simulated_ongkir_ditanggung_merchant = 0
    }
  } else {
    // Reset data jika tidak ada menu terpilih
    simulasiInputForm.value.nama_menu = ''
    simulasiInputForm.value.channel_menu = ''
    simulasiInputForm.value.hpp_produk = 0
    simulasiInputForm.value.harga_jual_kotor_produk = 0
    simulasiInputForm.value.simulated_komisi_channel_persen = 0
    simulasiInputForm.value.simulated_pajak_persen = 0
  }
}, {immediate: true})

// --- Handler Submit Simulasi ---
// Mengirim data input ke backend untuk perhitungan simulasi
const handleSubmitSimulasi = async () => {
  // Validasi dasar di frontend sebelum mengirim payload
  if (!simulasiInputForm.value.selectedHargaJualId && simulasiInputForm.value.nama_menu === '') {
    alert('Pilih Menu atau isi "Nama Menu" dan "HPP Produk" secara manual terlebih dahulu.')
    return
  }
  if (simulasiInputForm.value.hpp_produk <= 0 || simulasiInputForm.value.harga_jual_kotor_produk <= 0) {
    alert('HPP Produk dan Harga Jual Produk harus lebih dari 0.')
    return
  }
  if (simulasiInputForm.value.is_promo_ongkir && simulasiInputForm.value.simulated_ongkir_ditanggung_merchant <= 0) {
    alert('Nominal ongkir ditanggung merchant harus lebih dari 0 jika promo ongkir aktif.')
    return
  }
  if (simulasiInputForm.value.is_pakai_promo_channel && !simulasiInputForm.value.selected_promo_id) {
    alert('Pilih Promo jika "Pakai Promo Channel" aktif.')
    return
  }
  if (simulasiInputForm.value.simulated_komisi_channel_persen < 0 || simulasiInputForm.value.simulated_komisi_channel_persen > 100 ||
      simulasiInputForm.value.simulated_pajak_persen < 0 || simulasiInputForm.value.simulated_pajak_persen > 100) {
    alert('Persentase Komisi dan Pajak harus antara 0 dan 100.')
    return
  }


  try {
    // Membuat payload sesuai struct SimulasiInput di Go
    const payload = {
      harga_jual_kotor_produk: parseFloat(simulasiInputForm.value.harga_jual_kotor_produk),
      hpp_produk: parseFloat(simulasiInputForm.value.hpp_produk),
      nama_menu: simulasiInputForm.value.nama_menu,
      channel_menu: simulasiInputForm.value.channel_menu,
      jumlah_porsi_pembelian: parseFloat(simulasiInputForm.value.jumlah_porsi_pembelian),
      is_promo_ongkir: simulasiInputForm.value.is_promo_ongkir,
      simulated_ongkir_ditanggung_merchant: parseFloat(simulasiInputForm.value.simulated_ongkir_ditanggung_merchant),
      is_pakai_promo_channel: simulasiInputForm.value.is_pakai_promo_channel,
      selected_promo_id: simulasiInputForm.value.selected_promo_id,
      simulated_komisi_channel_persen: parseFloat(simulasiInputForm.value.simulated_komisi_channel_persen),
      simulated_pajak_persen: parseFloat(simulasiInputForm.value.simulated_pajak_persen)
    }

    const response = await axios.post(`${API_BASE_URL}/simulasi-promo`, payload)
    simulasiResult.value = response.data
    console.log('Hasil Simulasi:', simulasiResult.value);
  } catch (error) {
    console.error('Error simulating promo:', error.response ? error.response.data : error)
    const errorMessage = error.response && error.response.data && error.response.data.error
                         ? error.response.data.error
                         : 'Gagal melakukan simulasi. Periksa input Anda dan log backend.';
    alert(errorMessage);
  }
}

// --- Fungsi Reset Form Simulasi ---
const resetSimulasiForm = () => {
  simulasiInputForm.value = {
    selectedHargaJualId: '',
    nama_menu: '',
    channel_menu: '',
    hpp_produk: 0,
    harga_jual_kotor_produk: 0,
    jumlah_porsi_pembelian: 1,
    is_promo_ongkir: false,
    simulated_ongkir_ditanggung_merchant: 0,
    is_pakai_promo_channel: false,
    selected_promo_id: '',
    simulated_komisi_channel_persen: 0,
    simulated_pajak_persen: 0
  }
  simulasiResult.value = null // Bersihkan hasil simulasi
}

// --- On Mounted ---
// Panggil fungsi fetch saat komponen dimuat
onMounted(() => {
  fetchHargaJuals()
  fetchProgramPromos()
})

// --- Helper untuk Formatting ---
</script>
<template>
  <div class="">
    <h1 class="text-3xl font-bold text-gray-800 mb-6">Simulasi Promo Channel & Event Normal</h1>

    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <div class="bg-gray-50 rounded-lg p-6 shadow-md">
        <h2 class="text-xl font-semibold text-gray-700 mb-4">1. Pilih Menu</h2>
        <div class="mb-4">
          <label for="selectMenu" class="block text-gray-700 text-sm font-bold mb-2">Pilih Menu:</label>
          <select id="selectMenu" v-model="simulasiInputForm.selectedHargaJualId"
                  class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
            <option value="">-- Pilih Menu Tersimpan --</option>
            <option v-for="hj in hargaJuals" :key="hj.id" :value="hj.id">{{ hj.nama_produk }} (Rp {{ formatCurrency(hj.harga_jual_kotor) }})</option>
          </select>
        </div>
        <div class="mb-4">
          <label for="namaMenu" class="block text-gray-700 text-sm font-bold mb-2">Nama Menu:</label>
          <input type="text" id="namaMenu" v-model="simulasiInputForm.nama_menu"
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 bg-gray-100 leading-tight focus:outline-none" readonly>
        </div>
        <div class="mb-4">
          <label for="channelMenu" class="block text-gray-700 text-sm font-bold mb-2">Channel Menu:</label>
          <input type="text" id="channelMenu" v-model="simulasiInputForm.channel_menu"
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 bg-gray-100 leading-tight focus:outline-none" readonly>
        </div>
        <div class="mb-4">
          <label for="hppProduk" class="block text-gray-700 text-sm font-bold mb-2">HPP Satuan (Rp):</label>
          <input type="number" id="hppProduk" v-model.number="simulasiInputForm.hpp_produk" min="0" step="0.01"
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
        </div>
        <div class="mb-4">
          <label for="hargaJualKotorProduk" class="block text-gray-700 text-sm font-bold mb-2">Harga Jual Satuan (Rp):</label>
          <input type="number" id="hargaJualKotorProduk" v-model.number="simulasiInputForm.harga_jual_kotor_produk"  min="0" step="0.01"
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline hidden">
          <input type="text" :value="simulasiInputForm.harga_jual_kotor_produk.toLocaleString('id-ID', { minimumFractionDigits: 0, maximumFractionDigits: 2 })"
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">

        </div>
      </div>

      <div class="bg-gray-50 rounded-lg p-6 shadow-md">
        <h2 class="text-xl font-semibold text-gray-700 mb-4">2. Set Ketentuan & Pilih Promo</h2>
        <form @submit.prevent="handleSubmitSimulasi">
          <div class="mb-4">
            <label for="jumlahPembelian" class="block text-gray-700 text-sm font-bold mb-2">Jumlah Porsi Pembelian:</label>
            <input type="number" id="jumlahPembelian" v-model.number="simulasiInputForm.jumlah_porsi_pembelian" min="0" max="100" step="0.01"
            class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
          </div>
          <div class="mb-4 flex items-center">
            <input type="checkbox" id="isPromoOngkir" v-model="simulasiInputForm.is_promo_ongkir"
                   class="mr-2 h-4 w-4 text-lime-600 focus:ring-lime-500 border-lime-300 rounded">
            <label for="isPromoOngkir" class="text-gray-700 text-sm font-bold">Promo Ongkir</label>
          </div>
          <div class="mb-4" v-if="simulasiInputForm.is_promo_ongkir">
            <label for="ongkirDitanggung" class="block text-gray-700 text-sm font-bold mb-2">Ongkir Ditanggung Merchant (Rp):</label>
            <input type="number" id="ongkirDitanggung" v-model.number="simulasiInputForm.simulated_ongkir_ditanggung_merchant" min="0" step="100"
                   class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
          </div>

          <div class="mb-4 flex items-center">
            <input type="checkbox" id="isPakaiPromoChannel" v-model="simulasiInputForm.is_pakai_promo_channel"
                   class="mr-2 h-4 w-4 text-lime-600 focus:ring-lime-500 border-gray-300 rounded">
            <label for="isPakaiPromoChannel" class="text-gray-700 text-sm font-bold">Pakai Promo Channel</label>
          </div>
          <div class="mb-4" v-if="simulasiInputForm.is_pakai_promo_channel">
            <label for="selectPromo" class="block text-gray-700 text-sm font-bold mb-2">Pilih Promo:</label>
            <select id="selectPromo" v-model="simulasiInputForm.selected_promo_id" required
                    class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
              <option value="">-- Pilih Program Promo --</option>
              <option v-for="promo in programPromos" :key="promo.id" :value="promo.id">{{ promo.nama_promo }} ({{ promo.jenis_diskon === 'persentase' ? promo.besar_diskon + '%' : 'Rp ' + formatCurrency(promo.besar_diskon) }})</option>
            </select>
          </div>

          <div class="mb-4">
            <label for="komisiChannel" class="block text-gray-700 text-sm font-bold mb-2">Komisi Channel (%):</label>
            <input type="number" id="komisiChannel" v-model.number="simulasiInputForm.simulated_komisi_channel_persen" min="0" max="100" step="0.01"
                   class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
          </div>
          <div class="mb-6">
            <label for="pajakPersen" class="block text-gray-700 text-sm font-bold mb-2">Pajak (%):</label>
            <input type="number" id="pajakPersen" v-model.number="simulasiInputForm.simulated_pajak_persen" min="0" max="100" step="0.01"
                   class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
          </div>

          <button type="submit"
                  class="w-full text-lime-900 bg-lime-400 hover:bg-lime-500 focus:ring-4 focus:outline-none focus:ring-lime-200 font-medium rounded-lg px-5 py-2.5 text-center">
            Hitung Simulasi
          </button>
        </form>
      </div>

      <div class="bg-blue-50 rounded-lg p-6 shadow-md">
        <h2 class="text-xl font-semibold text-blue-800 mb-4">3. Hasil Perhitungan</h2>
        <div v-if="simulasiResult" class="text-gray-700 text-sm">
          <h3 class="font-bold text-gray-800 mb-2">Data Awal:</h3>
          <p>Nama Menu: <span class="font-semibold">{{ simulasiResult.nama_menu }}</span></p>
          <p>Channel Menu: <span class="font-semibold">{{ simulasiResult.channel_menu }}</span></p>
          <p>Jumlah Porsi Pembelian: <span class="font-semibold">{{ simulasiResult.jumlah_porsi_pembelian }}</span></p>
          <p>HPP Produk Total: <span class="font-semibold">{{ formatCurrency(simulasiResult.hpp_produk_total) }}</span></p>
          <p>Harga Jual Satuan: <span class="font-semibold">{{ formatCurrency(simulasiResult.harga_jual_kotor_produk) }}</span></p>
          <p class="mb-4">Harga Jual Total: <span class="font-bold text-green-600">{{ formatCurrency(simulasiResult.harga_jual_total_kotor) }}</span></p>

          <h3 class="font-bold text-gray-800 mt-4 mb-2">Detail Promo Terpilih:</h3>
          <p v-if="simulasiResult.nama_promo_terpilih">
            Nama Promo: <span class="font-semibold">{{ simulasiResult.nama_promo_terpilih }}</span> ({{ simulasiResult.jenis_diskon_promo === 'persentase' ? simulasiResult.besar_diskon_promo + '%' : formatCurrency(simulasiResult.besar_diskon_promo) }})
          </p>
          <p v-else class="text-gray-500">Tidak menggunakan promo.</p>
          <p v-if="simulasiResult.promo_applied" class="text-green-600">Promo Diterapkan!</p>
          <p v-if="simulasiResult.nama_promo_terpilih">Min Belanja: {{ formatCurrency(simulasiResult.min_belanja_promo) }}</p>
          <p v-if="simulasiResult.nama_promo_terpilih">Maks Potongan: {{ formatCurrency(simulasiResult.maksimal_potongan_promo) }}</p>
          <p v-if="simulasiResult.nama_promo_terpilih">Ditanggung Merchant: {{ formatPercentage(simulasiResult.ditanggung_merchant_promo_persen) }}</p>
          <p v-if="simulasiResult.nama_promo_terpilih" class="mb-4">Catatan: {{ simulasiResult.catatan_promo }}</p>

          <h3 class="font-bold text-gray-800 mt-4 mb-2">Bagi Konsumen:</h3>
          <p>Harga Jual Konsumen: <span class="font-semibold">{{ formatCurrency(simulasiResult.harga_jual_untuk_konsumen) }}</span></p>
          <p>Diskon Promo Konsumen: <span class="font-semibold">{{ formatCurrency(simulasiResult.diskon_promo_konsumen) }}</span></p>
          <p class="mb-4">Harga Akhir Konsumen: <span class="font-bold text-green-600">{{ formatCurrency(simulasiResult.harga_akhir_konsumen) }}</span></p>

          <h3 class="font-bold text-gray-800 mt-4 mb-2">Biaya Promo Channel:</h3>
          <p>Diskon Ditanggung Channel: <span class="font-semibold">{{ formatCurrency(simulasiResult.potongan_promo_ditanggung_channel) }}</span></p>
          <p>Diskon Ditanggung Merchant: <span class="font-semibold">{{ formatCurrency(simulasiResult.potongan_promo_ditanggung_merchant) }}</span></p>
          <p>Biaya Komisi Channel: <span class="font-semibold">{{ formatCurrency(simulasiResult.biaya_komisi_channel) }}</span></p>
          <p>Biaya Pajak: <span class="font-semibold">{{ formatCurrency(simulasiResult.biaya_pajak) }}</span></p>
          <p class="mb-4">Biaya Subsidi Ongkir: <span class="font-semibold">{{ formatCurrency(simulasiResult.biaya_subsidi_ongkir) }}</span></p>

          <h3 class="font-bold text-gray-800 mt-4 mb-2">Perhitungan Net Sales:</h3>
          <p>Sales Sebelum Komisi/Pajak/Ongkir: <span class="font-semibold">{{ formatCurrency(simulasiResult.sales_sebelum_komisi_pajak_ongkir) }}</span></p>
          <p class="mb-4">Net Sales (Rp): <span class="font-bold text-green-600">{{ formatCurrency(simulasiResult.net_sales) }}</span></p>

          <h3 class="font-bold text-gray-800 mt-4 mb-2">Hasil Akhir:</h3>
          <p>Gross Profit (Rp): <span class="font-bold text-green-600">{{ formatCurrency(simulasiResult.gross_profit) }}</span></p>
          <p>HPP terhadap Net Sales (%): <span class="font-semibold">{{ formatPercentage(simulasiResult.hpp_terhadap_net_sales_persen) }}</span></p>
          <p>Gross Profit terhadap Net Sales (%): <span class="font-semibold">{{ formatPercentage(simulasiResult.gross_profit_terhadap_net_sales_persen) }}</span></p>

        </div>
        <div v-else class="text-center text-gray-500">
          Masukkan data di samping dan klik "Hitung Simulasi" untuk melihat hasilnya.
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Tidak ada gaya scoped karena menggunakan Tailwind CSS */
</style>