<script setup>
import { ref, onMounted, watch } from 'vue'
  import axios from 'axios'
import { formatCurrency, formatPercentage } from '../utils/formatters'

const API_BASE_URL = 'http://localhost:8080/api'

// --- State Data ---
const reseps = ref([]) // Daftar Resep untuk dropdown
const selectedResepId = ref('') // ID Resep yang dipilih di form utama
const hppResep = ref(0) // HPP total dari resep yang dipilih (untuk form utama)
const savedHargaJuals = ref([]) // Daftar harga jual yang tersimpan

// --- State Form Utama (menggabungkan semua input) ---
const formModel = ref({
  id: '', // Diisi hanya saat mode edit
  resep_id: '',
  nama_produk: '',
  channel: '',
  jumlah_porsi_produk: 1,

  // --- Kriteria Perhitungan Optimal (digabungkan ke formModel) ---
  selectedCriteria: '', // 'min_profit_net_sales_persen', 'min_profit_rp_hpp', etc.
  min_profit_net_sales_persen: null, // Nullable number
  min_profit_rp_hpp: null,
  min_profit_persen_hpp: null,
  min_profit_x_lipat_hpp: null,
  max_hpp_net_sales_persen: null,
  target_net_sales_x_lipat_hpp: null,
  target_net_sales_rp: null,
  target_harga_jual_rp: null,
  consumer_pays_including_tax_rp: null,
  target_harga_jual_excl_tax_rp: null,

  // --- Biaya Operasional (juga akan digunakan untuk kalkulasi optimal) ---
  pajak_persen: 0,
  komisi_channel_persen: 0,

  // --- Field untuk Disimpan (Akan diisi otomatis setelah kalkulasi optimal di backend) ---
  metode_perhitungan: '', // Akan diisi dengan metodeTerkalkulasi dari backend
  nilai_kriteria: 0,      // Akan diisi dengan calculated_harga_jual_kotor dari backend
  harga_bulat: false      // Ini adalah bagian dari kalkulator optimal
})

// --- State UI ---
const isEditing = ref(false) // Mode edit/tambah
const originalResepNameForEdit = ref('') // Menyimpan nama resep asli saat edit
const showDetailModal = ref(false) // Kontrol visibilitas modal detail
const currentDetailedHargaJual = ref(null) // Data harga jual yang ditampilkan di modal detail

// --- Hasil Perhitungan Optimal (untuk display di UI setelah submit form utama) ---
const optimalCalculationDisplayResult = ref(null) // Akan diisi dari response API utama

// Daftar pilihan channel
const channels = ['GoFood', 'GrabFood', 'ShopeeFood', 'Internal', 'Lainnya', 'Semua']


// --- Fetch Data Awal ---
const fetchReseps = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/reseps`)
    reseps.value = response.data
  } catch (error) {
    console.error('Error fetching reseps:', error)
    alert('Gagal mengambil daftar resep.')
  }
}

const fetchHPPForSelectedResep = async (resepId) => {
  if (!resepId) {
    hppResep.value = 0
    return
  }
  try {
    const response = await axios.get(`${API_BASE_URL}/hpp/${resepId}`)
    hppResep.value = parseFloat(response.data.hpp_per_porsi)
  } catch (error) {
    console.error('Error fetching HPP:', error.response ? error.response.data : error)
    hppResep.value = 0
    alert('Gagal mengambil HPP resep. Pastikan HPP sudah dihitung untuk resep ini.')
  }
}

const fetchSavedHargaJuals = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/harga-juals`)
    savedHargaJuals.value = response.data
  }
  catch (error) {
    console.error('Error fetching saved harga juals:', error)
    alert('Gagal mengambil daftar harga jual yang tersimpan.')
  }
}


// --- Handler Submit Form Utama (Sekarang Menggabungkan Kalkulasi Optimal) ---
const handleSubmitCalculation = async () => {
  // Validasi dasar form
  if (!selectedResepId.value) {
    alert('Pilih resep terlebih dahulu.')
    return
  }
  if (hppResep.value === 0 && !isEditing.value) {
    alert('HPP resep belum tersedia atau 0. Harap hitung HPP resep terlebih dahulu.')
    return
  }
  if (!formModel.value.selectedCriteria) {
    alert('Pilih setidaknya satu kriteria perhitungan harga jual optimal.')
    return
  }

  // --- Persiapan Payload Komprehensif ---
  const payload = {
    resep_id: selectedResepId.value,
    nama_produk: formModel.value.nama_produk,
    channel: formModel.value.channel,
    jumlah_porsi_produk: parseFloat(formModel.value.jumlah_porsi_produk),
    pajak_persen: parseFloat(formModel.value.pajak_persen),
    komisi_channel_persen: parseFloat(formModel.value.komisi_channel_persen),

    // Kriteria Optimal yang akan dikirim (hanya satu yang non-null)
    selectedCriteria: formModel.value.selectedCriteria,
    min_profit_net_sales_persen: formModel.value.min_profit_net_sales_persen !== null ? parseFloat(formModel.value.min_profit_net_sales_persen) : null,
    min_profit_rp_hpp: formModel.value.min_profit_rp_hpp !== null ? parseFloat(formModel.value.min_profit_rp_hpp) : null,
    min_profit_persen_hpp: formModel.value.min_profit_persen_hpp !== null ? parseFloat(formModel.value.min_profit_persen_hpp) : null,
    min_profit_x_lipat_hpp: formModel.value.min_profit_x_lipat_hpp !== null ? parseFloat(formModel.value.min_profit_x_lipat_hpp) : null,
    max_hpp_net_sales_persen: formModel.value.max_hpp_net_sales_persen !== null ? parseFloat(formModel.value.max_hpp_net_sales_persen) : null,
    target_net_sales_x_lipat_hpp: formModel.value.target_net_sales_x_lipat_hpp !== null ? parseFloat(formModel.value.target_net_sales_x_lipat_hpp) : null,
    target_net_sales_rp: formModel.value.target_net_sales_rp !== null ? parseFloat(formModel.value.target_net_sales_rp) : null,
    target_harga_jual_rp: formModel.value.target_harga_jual_rp !== null ? parseFloat(formModel.value.target_harga_jual_rp) : null,
    consumer_pays_including_tax_rp: formModel.value.consumer_pays_including_tax_rp !== null ? parseFloat(formModel.value.consumer_pays_including_tax_rp) : null,
    target_harga_jual_excl_tax_rp: formModel.value.target_harga_jual_excl_tax_rp !== null ? parseFloat(formModel.value.target_harga_jual_excl_tax_rp) : null,

    // Field ini akan diisi oleh backend
    metode_perhitungan: '', // Akan diisi oleh backend dengan metode optimal yang digunakan
    nilai_kriteria: 0,      // Akan diisi oleh backend dengan calculated_harga_jual_kotor
    harga_bulat: formModel.value.harga_bulat
  }

  // --- Tahap 1: Panggil API Penyimpanan (Backend akan melakukan kalkulasi optimal) ---
  try {
    let response;
    if (isEditing.value) {
      response = await axios.put(`${API_BASE_URL}/harga-juals/${formModel.value.id}`, payload)
      alert('Harga jual berhasil diperbarui dengan hasil optimal!')
    } else {
      response = await axios.post(`${API_BASE_URL}/harga-juals/calculate`, payload)
      alert('Harga jual optimal berhasil dihitung dan disimpan!')
    }

    optimalCalculationDisplayResult.value = response.data; // Tampilkan hasil kalkulasi dari response utama
    resetForm(); // Reset form setelah berhasil
    fetchSavedHargaJuals(); // Muat ulang daftar yang tersimpan
  } catch (error) {
    console.error('Error during optimal calculation or saving:', error.response ? error.response.data : error)
    const errorMessage = error.response && error.response.data && error.response.data.error
                         ? error.response.data.error
                         : 'Gagal menghitung/menyimpan harga jual optimal. Periksa input.';
    alert(errorMessage);
  }
}

// --- Fungsi Reset Form Utama ---
const resetForm = () => {
  isEditing.value = false
  selectedResepId.value = ''
  hppResep.value = 0
  originalResepNameForEdit.value = ''
  optimalCalculationDisplayResult.value = null; // Reset hasil optimal

  formModel.value = {
    id: '',
    resep_id: '',
    nama_produk: '',
    channel: '',
    jumlah_porsi_produk: 1,
    // Reset semua kriteria optimal
    selectedCriteria: '',
    min_profit_net_sales_persen: null, min_profit_rp_hpp: null,
    min_profit_persen_hpp: null, min_profit_x_lipat_hpp: null,
    max_hpp_net_sales_persen: null,
    target_net_sales_x_lipat_hpp: null, target_net_sales_rp: null,
    target_harga_jual_rp: null, consumer_pays_including_tax_rp: null,
    target_harga_jual_excl_tax_rp: null,
    // Reset field yang diisi backend
    metode_perhitungan: '',
    nilai_kriteria: 0,
    harga_bulat: false,     // Ini tetap dari user input
    pajak_persen: 0,
    komisi_channel_persen: 0,
  }
}

// --- Fungsi Edit Harga Jual (Diperbarui untuk mengisi kriteria optimal) ---
const editHargaJual = (hargaJual) => {
  isEditing.value = true
  // Isi formModel dengan data hargaJual yang akan diedit
  formModel.value = {
    id: hargaJual.id,
    resep_id: hargaJual.resep_id,
    nama_produk: hargaJual.nama_produk,
    channel: hargaJual.channel,
    jumlah_porsi_produk: parseFloat(hargaJual.jumlah_porsi_produk),

    // Ketika edit, kita asumsikan ini adalah target_harga_jual_rp
    selectedCriteria: 'target_harga_jual_rp',
    target_harga_jual_rp: parseFloat(hargaJual.harga_jual_kotor), // HJKotor sebagai nilai kriteria

    // Pastikan semua kriteria lain null
    min_profit_net_sales_persen: null, min_profit_rp_hpp: null,
    min_profit_persen_hpp: null, min_profit_x_lipat_hpp: null,
    max_hpp_net_sales_persen: null,
    target_net_sales_x_lipat_hpp: null, target_net_sales_rp: null,
    consumer_pays_including_tax_rp: null,
    target_harga_jual_excl_tax_rp: null,

    pajak_persen: parseFloat(hargaJual.pajak_persen),
    komisi_channel_persen: parseFloat(hargaJual.komisi_channel_persen),
    metode_perhitungan: hargaJual.metode_perhitungan, // Dari data yang disimpan
    nilai_kriteria: parseFloat(hargaJual.nilai_kriteria), // Dari data yang disimpan
    harga_bulat: hargaJual.harga_bulat || false,
  }
  selectedResepId.value = hargaJual.resep_id
  hppResep.value = parseFloat(hargaJual.hpp); // HPP Resep dari data yang tersimpan
  originalResepNameForEdit.value = hargaJual.resep_nama

  window.scrollTo({ top: 0, behavior: 'smooth' }) // Scroll ke atas form
}

// --- Fungsi Delete, Show Detail, Close Detail (tidak berubah) ---
const deleteHargaJual = async (id, namaProduk) => { /* ... */ }
const showHargaJualDetail = (hargaJual) => { /* ... */ }
const closeDetailModal = () => { /* ... */ }


// --- Watcher untuk selectedResepId (Form Utama) ---
watch(selectedResepId, (newVal) => {
  if (newVal) {
    fetchHPPForSelectedResep(newVal);
  } else {
    hppResep.value = 0;
  }
});


  console.log( fetchReseps() );


onMounted(() => {
  fetchReseps()
  fetchSavedHargaJuals()
})


</script>

<template>
  <div class="">
    <h1 class="text-3xl font-bold text-gray-800 mb-6">Kalkulator Harga Jual</h1>

    <div class="bg-gray-50 rounded-lg p-6 mb-8 shadow-md">
      <h2 class="text-2xl font-semibold text-gray-700 mb-5" v-if="!isEditing">Hitung & Simpan Harga Jual Baru</h2>
      <h2 class="text-2xl font-semibold text-gray-700 mb-5" v-else>Edit Harga Jual: {{ formModel.nama_produk }} <span class="text-sm text-gray-500"> (Resep Asli: {{ originalResepNameForEdit || 'Memuat...' }})</span></h2>
      <form @submit.prevent="handleSubmitCalculation()">
        <div class="mb-4">
            <label for="selectResep" class="block text-gray-700 text-sm font-bold mb-2">Pilih Resep:</label>
            <select id="selectResep" v-model="selectedResepId" :disabled="isEditing" required
                    class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline disabled:bg-gray-200">
                <option value="">-- Pilih Resep --</option>
                <option v-for="resep in reseps" :key="resep.id" :value="resep.id">{{ resep.nama }}</option>
            </select>
            <p v-if="selectedResepId && hppResep > 0" class="text-sm text-gray-600 mt-2">
                HPP Resep (Total) Terpilih: Rp {{ hppResep.toLocaleString('id-ID', { minimumFractionDigits: 2, maximumFractionDigits: 4 }) }}
            </p>
            <p v-else-if="selectedResepId && hppResep === 0 && !isEditing" class="text-sm text-red-500 mt-2">
                HPP resep belum tersedia. Pastikan Anda sudah menghitung HPP-nya di halaman Resep.
            </p>
        </div>

        <div class="mb-4">
            <label for="namaProduk" class="block text-gray-700 text-sm font-bold mb-2">Nama Produk (untuk harga jual):</label>
            <input type="text" id="namaProduk" v-model="formModel.nama_produk" required
                   class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
        </div>

        <div class="mb-4">
            <label for="channel" class="block text-gray-700 text-sm font-bold mb-2">Channel Penjualan:</label>
            <select id="channel" v-model="formModel.channel" required
                    class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
                <option value="">-- Pilih Channel --</option>
                <option v-for="ch in channels" :key="ch" :value="ch">{{ ch }}</option>
            </select>
        </div>

        <div class="mb-4">
            <label for="jumlahPorsiProduk" class="block text-gray-700 text-sm font-bold mb-2">Jumlah Porsi Produk yang Dijual:</label>
            <input type="number" id="jumlahPorsiProduk" v-model.number="formModel.jumlah_porsi_produk" min="0.0001" step="0.0001" required
                   class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
        </div>

        <div class="grid grid-cols-2 gap-4 mb-6">
            <div>
                <label for="pajakPersen" class="block text-gray-700 text-sm font-bold mb-2">Pajak (%):</label>
                <input type="number" id="pajakPersen" v-model.number="formModel.pajak_persen" min="0" max="100" step="0.01" required
                       class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
            </div>
            <div>
                <label for="komisiChannelPersen" class="block text-gray-700 text-sm font-bold mb-2">Komisi Channel (%):</label>
                <input type="number" id="komisiChannelPersen" v-model.number="formModel.komisi_channel_persen" min="0" max="100" step="0.01" required
                       class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
            </div>
        </div>

        <div class="mb-6 space-y-4 border p-4 rounded-lg bg-blue-50">
            <h3 class="text-xl font-semibold text-blue-800 mb-3">Pilih Kriteria Perhitungan Harga Jual Optimal:</h3>

            <div class="flex items-start">
                <input type="radio" id="min_profit_net_sales_persen_radio" value="min_profit_net_sales_persen" v-model="formModel.selectedCriteria" class="mt-1 mr-2">
                <label for="min_profit_net_sales_persen_radio" class="block text-gray-700">Saya ingin keuntungan saya minimal <input type="number" v-model.number="formModel.min_profit_net_sales_persen" min="0" step="0.01" :disabled="formModel.selectedCriteria !== 'min_profit_net_sales_persen'" class="ml-2 w-24 border rounded py-1 px-2 text-sm"> % dari net sales</label>
            </div>
            <div class="flex items-start">
                <input type="radio" id="min_profit_rp_hpp_radio" value="min_profit_rp_hpp" v-model="formModel.selectedCriteria" class="mt-1 mr-2">
                <label for="min_profit_rp_hpp_radio" class="block text-gray-700">Saya ingin keuntungan saya minimal Rp <input type="number" v-model.number="formModel.min_profit_rp_hpp" min="0" step="0.01" :disabled="formModel.selectedCriteria !== 'min_profit_rp_hpp'" class="ml-2 w-24 border rounded py-1 px-2 text-sm"> dari HPP</label>
            </div>
            <div class="flex items-start">
                <input type="radio" id="min_profit_persen_hpp_radio" value="min_profit_persen_hpp" v-model="formModel.selectedCriteria" class="mt-1 mr-2">
                <label for="min_profit_persen_hpp_radio" class="block text-gray-700">Saya ingin keuntungan saya minimal <input type="number" v-model.number="formModel.min_profit_persen_hpp" min="0" step="0.01" :disabled="formModel.selectedCriteria !== 'min_profit_persen_hpp'" class="ml-2 w-24 border rounded py-1 px-2 text-sm"> % dari HPP</label>
            </div>
            <div class="flex items-start">
                <input type="radio" id="min_profit_x_lipat_hpp_radio" value="min_profit_x_lipat_hpp" v-model="formModel.selectedCriteria" class="mt-1 mr-2">
                <label for="min_profit_x_lipat_hpp_radio" class="block text-gray-700">Saya ingin keuntungan saya minimal <input type="number" v-model.number="formModel.min_profit_x_lipat_hpp" min="0" step="0.01" :disabled="formModel.selectedCriteria !== 'min_profit_x_lipat_hpp'" class="ml-2 w-24 border rounded py-1 px-2 text-sm"> x lipat dari HPP</label>
            </div>

            <div class="flex items-start">
                <input type="radio" id="max_hpp_net_sales_persen_radio" value="max_hpp_net_sales_persen" v-model="formModel.selectedCriteria" class="mt-1 mr-2">
                <label for="max_hpp_net_sales_persen_radio" class="block text-gray-700">Saya ingin HPP saya maksimal <input type="number" v-model.number="formModel.max_hpp_net_sales_persen" min="0" max="100" step="0.01" :disabled="formModel.selectedCriteria !== 'max_hpp_net_sales_persen'" class="ml-2 w-24 border rounded py-1 px-2 text-sm"> % dari net sales</label>
            </div>

            <div class="flex items-start">
                <input type="radio" id="target_net_sales_x_lipat_hpp_radio" value="target_net_sales_x_lipat_hpp" v-model="formModel.selectedCriteria" class="mt-1 mr-2">
                <label for="target_net_sales_x_lipat_hpp_radio" class="block text-gray-700">Saya ingin net sales saya <input type="number" v-model.number="formModel.target_net_sales_x_lipat_hpp" min="0" step="0.01" :disabled="formModel.selectedCriteria !== 'target_net_sales_x_lipat_hpp'" class="ml-2 w-24 border rounded py-1 px-2 text-sm"> x lipat dari HPP</label>
            </div>
            <div class="flex items-start">
                <input type="radio" id="target_net_sales_rp_radio" value="target_net_sales_rp" v-model="formModel.selectedCriteria" class="mt-1 mr-2">
                <label for="target_net_sales_rp_radio" class="block text-gray-700">Saya ingin net sales saya Rp <input type="number" v-model.number="formModel.target_net_sales_rp" min="0" step="0.01" :disabled="formModel.selectedCriteria !== 'target_net_sales_rp'" class="ml-2 w-24 border rounded py-1 px-2 text-sm"></label>
            </div>

            <div class="flex items-start">
                <input type="radio" id="target_harga_jual_rp_radio" value="target_harga_jual_rp" v-model="formModel.selectedCriteria" class="mt-1 mr-2">
                <label for="target_harga_jual_rp_radio" class="block text-gray-700">Saya ingin harga jual menu ini Rp <input type="number" v-model.number="formModel.target_harga_jual_rp" min="0" step="0.01" :disabled="formModel.selectedCriteria !== 'target_harga_jual_rp'" class="ml-2 w-24 border rounded py-1 px-2 text-sm"></label>
            </div>

            <div class="flex items-start">
                <input type="radio" id="consumer_pays_including_tax_rp_radio" value="consumer_pays_including_tax_rp" v-model="formModel.selectedCriteria" class="mt-1 mr-2">
                <label for="consumer_pays_including_tax_rp_radio" class="block text-gray-700">Saya ingin konsumen membayar menu ini Rp <input type="number" v-model.number="formModel.consumer_pays_including_tax_rp" min="0" step="0.01" :disabled="formModel.selectedCriteria !== 'consumer_pays_including_tax_rp'" class="ml-2 w-24 border rounded py-1 px-2 text-sm"> (sudah termasuk pajak)</label>
            </div>

            <div class="flex items-start">
                <input type="radio" id="target_harga_jual_excl_tax_rp_radio" value="target_harga_jual_excl_tax_rp" v-model="formModel.selectedCriteria" class="mt-1 mr-2">
                <label for="target_harga_jual_excl_tax_rp_radio" class="block text-gray-700">Saya ingin harga jual menu ini Rp <input type="number" v-model.number="formModel.target_harga_jual_excl_tax_rp" min="0" step="0.01" :disabled="formModel.selectedCriteria !== 'target_harga_jual_excl_tax_rp'" class="ml-2 w-24 border rounded py-1 px-2 text-sm"> (belum termasuk pajak)</label>
            </div>
        </div>

        <div class="mb-6 flex items-center">
          <input type="checkbox" id="hargaBulat" v-model="formModel.harga_bulat"
                 class="mr-2 h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded">
          <label for="hargaBulat" class="text-gray-700 text-sm font-bold">Bulatkan Harga Jual Kotor (ke ratusan terdekat ke atas)</label>
        </div>

        <div class="flex items-center space-x-4">
          <button type="submit"
                  class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
            {{ isEditing ? 'Simpan Perubahan Harga Jual' : 'Hitung & Simpan Harga Jual' }}
          </button>
          <button type="button" @click="resetForm()"
                  class="bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
            {{ isEditing ? 'Batal Edit' : 'Reset Form' }}
          </button>
        </div>
      </form>
    </div>

    <div v-if="optimalCalculationDisplayResult" class="mt-6 p-4 bg-blue-100 border border-blue-300 rounded-lg shadow-md">
        <h3 class="text-xl font-semibold text-blue-800 mb-3">Hasil Perhitungan Terbaru:</h3>
        <p class="mb-2"><strong>Metode Terkalkulasi:</strong> {{ optimalCalculationDisplayResult.metode_terkalkulasi }}</p>
        <p class="mb-2"><strong>HPP Produk Total:</strong> {{ formatCurrency(optimalCalculationDisplayResult.hpp_produk_total) }}</p>
        <p class="mb-2"><strong>Komisi Channel:</strong> {{ formatCurrency(optimalCalculationDisplayResult.total_komisi) }}</p> <p class="mb-2"><strong>Pajak:</strong> {{ formatCurrency(optimalCalculationDisplayResult.total_pajak) }}</p>       <p class="mb-2"><strong>Harga Jual Bersih:</strong> {{ formatCurrency(optimalCalculationDisplayResult.harga_jual_bersih) }}</p>
        <p class="mb-2"><strong>Profit:</strong> {{ formatCurrency(optimalCalculationDisplayResult.profit) }}</p>
        <p class="mb-2"><strong>Profit (%):</strong> {{ formatPercentage(optimalCalculationDisplayResult.profit_persen) }}</p>
        <p class="mb-2"><strong>Net Sales:</strong> {{ formatCurrency(optimalCalculationDisplayResult.net_sales) }}</p>
        <p class="mb-2"><strong>HPP terhadap Net Sales (%):</strong> {{ formatPercentage(optimalCalculationDisplayResult.hpp_terhadap_net_sales_persen) }}</p>
        <p class="mb-4"><strong>Gross Profit terhadap Net Sales (%):</strong> {{ formatPercentage(optimalCalculationDisplayResult.gross_profit_terhadap_net_sales_persen) }}</p>

        <h4 class="text-xl font-bold text-green-700">Harga Jual Kotor Optimal: {{ formatCurrency(optimalCalculationDisplayResult.harga_jual_kotor) }}</h4> </div>


    <div class="bg-gray-50 rounded-lg p-6 shadow-md">
        <h2 class="text-2xl font-semibold text-gray-700 mb-5">Daftar Harga Jual Tersimpan</h2>
        <div class="overflow-x-auto">
            <table class="min-w-full bg-white border border-gray-200">
                <thead>
                    <tr>
                        <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Nama Produk</th>
                        <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Channel</th>
                        <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Resep</th>
                        <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">HPP Produk</th>
                        <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Harga Jual Kotor</th>
                        <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Profit (%)</th>
                        <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Metode</th>
                        <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Aksi</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-if="!savedHargaJuals || savedHargaJuals.length === 0">
                        <td colspan="8" class="py-4 px-4 text-center text-gray-500">Belum ada harga jual yang tersimpan.</td>
                    </tr>
                    <tr v-for="(hj, index) in savedHargaJuals" :key="hj.id" :class="index % 2 === 0 ? 'bg-white' : 'bg-gray-50'">
                        <td class="py-3 px-4 border-b text-sm text-gray-800">{{ hj.nama_produk }}</td>
                        <td class="py-3 px-4 border-b text-sm text-gray-800">{{ hj.channel || 'N/A' }}</td>
                        <td class="py-3 px-4 border-b text-sm text-gray-800">{{ hj.resep_nama || 'N/A' }}</td>
                        <td class="py-3 px-4 border-b text-sm text-gray-800">Rp {{ parseFloat(hj.hpp).toLocaleString('id-ID', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}</td>
                        <td class="py-3 px-4 border-b text-sm text-gray-800">Rp {{ parseFloat(hj.harga_jual_kotor).toLocaleString('id-ID', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}</td>
                        <td class="py-3 px-4 border-b text-sm text-gray-800">{{ parseFloat(hj.profit_persen).toLocaleString('id-ID', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }}%</td>
                        <td class="py-3 px-4 border-b text-sm text-gray-800">{{ hj.metode_perhitungan }}</td>
                        <td class="py-3 px-4 border-b text-sm text-gray-800 whitespace-nowrap">
                            <button @click="showHargaJualDetail(hj)"
                                    class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-1 px-3 rounded-md text-xs mr-2">
                                Detail
                            </button>
                            <button @click="editHargaJual(hj)"
                                    class="bg-yellow-500 hover:bg-yellow-600 text-white font-bold py-1 px-3 rounded-md text-xs mr-2">
                                Edit
                            </button>
                            <button @click="deleteHargaJual(hj.id, hj.nama_produk)"
                                    class="bg-red-500 hover:bg-red-600 text-white font-bold py-1 px-3 rounded-md text-xs">
                                Hapus
                            </button>
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>

    <div v-if="showDetailModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 flex justify-center items-center p-4 z-50">
        </div>
  </div>
</template>

<style scoped>
/* Hapus semua gaya scoped karena menggunakan Tailwind CSS */
</style>