<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import ResepList from '../components/ResepList.vue'

const API_BASE_URL = 'http://localhost:8080/api'

// --- State Reaktif ---
const reseps = ref([]) // Menyimpan daftar resep yang diambil dari backend
const bahanBakus = ref([]) // Menyimpan daftar bahan baku untuk pilihan komponen (dropdown)
const existingReseps = ref([]) // Menyimpan daftar resep yang sudah ada untuk pilihan komponen (dropdown)

const isEditing = ref(false) // State untuk menandakan apakah sedang dalam mode edit
const selectedResep = ref(null) // Menyimpan data resep yang sedang diedit (untuk tampilan H2)

// Model untuk formulir resep baru/edit
const formModel = ref({
  id: '', // Diisi hanya saat mode edit
  nama: '',
  is_sub_resep: false,
  jumlah_porsi: 1, // Default jumlah porsi adalah 1
  komponen: [ // Inisialisasi dengan satu komponen kosong untuk input awal
    { komponen_id: '', kuantitas: 0, tipe_komponen: 'bahan_baku' }
  ]
})

// --- Fungsi untuk Mengambil Data dari Backend ---

// Mengambil semua resep
const fetchReseps = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/reseps`)
    reseps.value = response.data
    console.log('Resep berhasil diambil:', response.data);
  } catch (error) {
    console.error('Error fetching reseps:', error)
    alert('Gagal mengambil data resep.')
  }
    }

    // --- State untuk Modal Detail Resep ---
const showDetailModal = ref(false)
const currentDetailedResep = ref(null)

// Mengambil semua bahan baku (untuk dropdown saat menambah komponen)
const fetchBahanBakus = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/bahan-bakus`)
    bahanBakus.value = response.data
    console.log('Bahan baku berhasil diambil:', response.data);
  } catch (error) {
    console.error('Error fetching bahan baku:', error)
    alert('Gagal mengambil data bahan baku untuk pilihan.')
  }
}

// Mengambil resep yang sudah ada yang bisa menjadi komponen (sub-resep atau resep jadi)
const fetchExistingReseps = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/reseps`)
    // Filter resep yang sudah ada agar tidak termasuk resep yang sedang diedit (hindari self-reference)
    existingReseps.value = response.data.filter(r => r.id !== (formModel.value.id ? formModel.value.id : null));
    console.log('Resep yang sudah ada berhasil diambil:', response.data);
  } catch (error) {
    console.error('Error fetching existing reseps:', error)
    alert('Gagal mengambil resep yang sudah ada untuk pilihan.')
  }
}

// --- Handler untuk Form Penambahan/Update Resep ---
const handleSubmitResep = async () => {
  // Validasi sederhana pada komponen sebelum mengirim
  if (formModel.value.komponen.some(c => !c.komponen_id || c.kuantitas <= 0)) {
    alert('Pastikan semua komponen dipilih dan kuantitas lebih dari 0.')
    return
  }
  // Validasi jumlah porsi
  if (formModel.value.jumlah_porsi <= 0) {
    alert('Jumlah porsi harus lebih dari 0.')
    return
  }

  try {
    if (isEditing.value) {
      // Jika mode edit, kirim PUT request
      await axios.put(`${API_BASE_URL}/reseps/${formModel.value.id}`, formModel.value)
      alert('Resep berhasil diperbarui!')
    } else {
      // Jika mode tambah, kirim POST request
      await axios.post(`${API_BASE_URL}/reseps`, formModel.value)
      alert('Resep berhasil ditambahkan!')
    }

    resetForm() // Reset formulir setelah berhasil
    fetchReseps() // Muat ulang daftar resep
    fetchExistingReseps() // Muat ulang pilihan resep yang ada
  } catch (error) {
    console.error('Error saving resep:', error.response ? error.response.data : error)
    const errorMessage = error.response && error.response.data && error.response.data.error
                         ? error.response.data.error
                         : 'Gagal menyimpan resep. Pastikan nama unik dan data valid.';
    alert(errorMessage);
  }
}

// --- Fungsi untuk Mengelola Komponen Resep Secara Dinamis di Form ---
const addKomponen = () => {
  // Menambahkan komponen kosong baru ke array komponen
  formModel.value.komponen.push({ komponen_id: '', kuantitas: 0, tipe_komponen: 'bahan_baku' })
}

const removeKomponen = (index) => {
  // Menghapus komponen dari array berdasarkan indeks
  formModel.value.komponen.splice(index, 1)
}

// --- Edit/Delete Handlers yang Dipicu dari ResepList ---
const editResep = (resep) => {
  isEditing.value = true
  selectedResep.value = { ...resep } // Salin objek untuk referensi di H2
  // Isi formModel dengan data dari resep yang akan diedit (salinan objek)
  formModel.value = {
    id: resep.id,
    nama: resep.nama,
    is_sub_resep: resep.is_sub_resep,
    jumlah_porsi: resep.jumlah_porsi || 1, // Pastikan ada default jika null/0
    komponen: resep.komponen.map(c => ({ ...c })) // Salin setiap objek komponen
  }
  // Refresh pilihan resep yang ada, pastikan resep yang diedit tidak muncul sebagai komponen dirinya sendiri
  fetchExistingReseps();
}

const deleteResep = async (id, nama) => {
  if (confirm(`Apakah Anda yakin ingin menghapus resep "${nama}"? Semua komponen terkait juga akan dihapus.`)) {
    try {
      await axios.delete(`${API_BASE_URL}/reseps/${id}`)
      alert('Resep berhasil dihapus!')
      fetchReseps() // Muat ulang daftar resep
      fetchExistingReseps() // Muat ulang pilihan resep yang ada
      resetForm() // Reset form jika resep yang sedang diedit dihapus
    } catch (error) {
      console.error('Error deleting resep:', error.response ? error.response.data : error)
      const errorMessage = error.response && error.response.data && error.response.data.error
                           ? error.response.data.error
                           : 'Gagal menghapus resep. Mungkin sedang digunakan sebagai komponen di resep lain.';
      alert(errorMessage);
    }
  }
}

// Fungsi untuk membatalkan mode edit dan mereset form
const cancelEdit = () => {
  resetForm()
    }

    // --- Fungsi untuk Menampilkan Detail Resep ---
const showResepDetail = async (resepId) => {
  try {
    const response = await axios.get(`${API_BASE_URL}/reseps/${resepId}`)
    currentDetailedResep.value = response.data // Data detail dari backend DTO
    showDetailModal.value = true // Tampilkan modal
  } catch (error) {
    console.error('Error fetching resep detail:', error.response ? error.response.data : error)
    alert('Gagal mengambil detail resep.')
  }
}

// --- Fungsi untuk Menutup Modal Detail Resep ---
const closeDetailModal = () => {
  showDetailModal.value = false
  currentDetailedResep.value = null
}

// --- Fungsi untuk Duplikat Resep ---
const duplicateResep = async (resepId, resepNama) => {
  if (confirm(`Apakah Anda yakin ingin menduplikasi resep "${resepNama}"?`)) {
    try {
      const response = await axios.post(`${API_BASE_URL}/reseps/${resepId}/duplicate`) // API baru untuk duplikasi
      alert(`Resep "${resepNama}" berhasil diduplikasi menjadi "${response.data.nama_resep_baru}"!`)
      fetchReseps() // Muat ulang daftar resep
      fetchExistingReseps() // Muat ulang pilihan resep yang ada
    } catch (error) {
      console.error('Error duplicating resep:', error.response ? error.response.data : error)
      const errorMessage = error.response && error.response.data && error.response.data.error
                           ? error.response.data.error
                           : 'Gagal menduplikasi resep.';
      alert(errorMessage);
    }
  }
}

// Fungsi untuk mereset form ke kondisi awal (kosong atau mode tambah)
const resetForm = () => {
  isEditing.value = false
  selectedResep.value = null
  formModel.value = {
    id: '',
    nama: '',
    is_sub_resep: false,
    jumlah_porsi: 1,
    komponen: [{ komponen_id: '', kuantitas: 0, tipe_komponen: 'bahan_baku' }]
  }
  fetchExistingReseps(); // Panggil ulang untuk menghapus filter dari selectedResep
}

// --- Lifecycle Hook: onMounted ---
// Fungsi ini akan dipanggil saat komponen selesai dimuat ke DOM
onMounted(() => {
  fetchReseps() // Ambil daftar resep yang ada
  fetchBahanBakus() // Ambil daftar bahan baku untuk dropdown
  fetchExistingReseps() // Ambil daftar resep yang sudah ada untuk dropdown
})

// --- Handler Event dari Child Component ResepList (HPP Calculated) ---
// (Opsional, untuk menerima notifikasi HPP dari ResepList jika diperlukan)
const handleHPPCalculated = (hppResult) => {
  console.log('HPP diterima di parent (ResepView):', hppResult)
  // Anda bisa melakukan sesuatu di sini, misalnya menyimpan hasil HPP ke state lain
  // atau menampilkannya di bagian terpisah dari halaman ResepView.
}
</script>

<template>
  <div class="p-5 max-w-5xl mx-auto font-sans">
    <h1 class="text-3xl font-bold text-gray-800 mb-6">Manajemen Resep</h1>

    <div class="bg-gray-50 rounded-lg p-6 mb-8 shadow-md">
      <h2 class="text-2xl font-semibold text-gray-700 mb-5" v-if="!isEditing">Tambah Resep Baru</h2>
      <h2 class="text-2xl font-semibold text-gray-700 mb-5" v-else>Edit Resep: {{ selectedResep?.nama }}</h2>
      <form @submit.prevent="handleSubmitResep()">
        <div class="mb-4">
          <label for="resepNama" class="block text-gray-700 text-sm font-bold mb-2">Nama Resep:</label>
          <input type="text" id="resepNama" v-model="formModel.nama" required
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
        </div>
        <div class="mb-4 flex items-center">
          <input type="checkbox" id="isSubResep" v-model="formModel.is_sub_resep"
                 class="mr-2 h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded">
          <label for="isSubResep" class="text-gray-700 text-sm font-bold">Ini adalah Sub-Resep (dapat digunakan sebagai komponen di resep lain)</label>
        </div>
        <div class="mb-6">
          <label for="jumlahPorsi" class="block text-gray-700 text-sm font-bold mb-2">Jumlah Porsi Dihasilkan:</label>
          <input type="number" id="jumlahPorsi" v-model.number="formModel.jumlah_porsi" min="1" step="0.0001" required
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
        </div>

        <h3 class="text-xl font-semibold text-gray-700 mb-4">Komponen Resep</h3>
        <div v-for="(komp, index) in formModel.komponen" :key="index"
             class="border border-dashed border-gray-400 rounded-lg p-4 mb-4 bg-gray-100 relative">
          <div class="mb-4">
            <label :for="'tipeKomponen' + index" class="block text-gray-700 text-sm font-bold mb-2">Tipe Komponen:</label>
            <select :id="'tipeKomponen' + index" v-model="komp.tipe_komponen" @change="komp.komponen_id = ''"
                    class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
              <option value="bahan_baku">Bahan Baku</option>
              <option value="resep">Resep Lain</option>
            </select>
          </div>

          <div class="mb-4">
            <label :for="'komponenId' + index" class="block text-gray-700 text-sm font-bold mb-2">Pilih Komponen:</label>
            <select :id="'komponenId' + index" v-model="komp.komponen_id" required
                    class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
              <option value="">-- Pilih --</option>
              <template v-if="komp.tipe_komponen === 'bahan_baku'">
                <option v-for="bb in bahanBakus" :key="bb.id" :value="bb.id">{{ bb.nama }} ({{ bb.satuan_pemakaian }})</option>
              </template>
              <template v-else-if="komp.tipe_komponen === 'resep'">
                <option v-for="r in existingReseps" :key="r.id" :value="r.id">{{ r.nama }}</option>
              </template>
            </select>
          </div>

          <div class="mb-4">
            <label :for="'kuantitas' + index" class="block text-gray-700 text-sm font-bold mb-2">Kuantitas Penggunaan:</label>
            <input type="number" :id="'kuantitas' + index" v-model.number="komp.kuantitas" step="0.0001" min="0" required
                   class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
          </div>

          <button type="button" @click="removeKomponen(index)"
                  class="absolute top-2 right-2 bg-red-500 hover:bg-red-700 text-white font-bold py-1 px-2 text-xs rounded-full"
                  v-if="formModel.komponen.length > 1">
            &times;
          </button>
        </div>

        <button type="button" @click="addKomponen"
                class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline mr-4">
          Tambah Komponen Lain
        </button>

        <div class="flex items-center space-x-4 mt-6">
          <button type="submit"
                  class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
            {{ isEditing ? 'Simpan Perubahan Resep' : 'Simpan Resep Baru' }}
          </button>
          <button type="button" @click="cancelEdit()"
                  class="bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline"
                  v-if="isEditing">
            Batal Edit
          </button>
        </div>
      </form>
    </div>

    <div class="bg-gray-50 rounded-lg p-6 shadow-md">
      <h2 class="text-2xl font-semibold text-gray-700 mb-5">Daftar Resep</h2>
      <ResepList
        :reseps="reseps"
        @hppCalculated="handleHPPCalculated"
        @edit="editResep"
        @delete="deleteResep"
        @showDetail="showResepDetail"
        @duplicate="duplicateResep"      />
    </div>

    <div v-if="showDetailModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 flex justify-center items-center p-4 z-50">
      <div class="bg-white rounded-lg shadow-xl max-w-2xl w-full p-6">
        <div class="flex justify-between items-center border-b pb-3 mb-4">
          <h3 class="text-2xl font-bold text-gray-800">Detail Resep: {{ currentDetailedResep?.nama }}</h3>
          <button @click="closeDetailModal" class="text-gray-500 hover:text-gray-700 text-3xl font-semibold">&times;</button>
        </div>

        <div v-if="currentDetailedResep" class="text-gray-700 text-sm">
          <p class="mb-2"><strong>ID:</strong> {{ currentDetailedResep.id }}</p>
          <p class="mb-2"><strong>Nama Resep:</strong> {{ currentDetailedResep.nama }}</p>
          <p class="mb-2"><strong>Sub-Resep:</strong> {{ currentDetailedResep.is_sub_resep ? 'Ya' : 'Tidak' }}</p>
          <p class="mb-2"><strong>Jumlah Porsi:</strong> {{ currentDetailedResep.jumlah_porsi }}</p>
          <p class="mb-2"><strong>Dibuat:</strong> {{ currentDetailedResep.created_at }}</p>
          <p class="mb-4"><strong>Diperbarui:</strong> {{ currentDetailedResep.updated_at }}</p>

          <h4 class="text-lg font-semibold text-gray-700 mb-3">Daftar Komponen:</h4>
          <ul class="list-disc pl-5">
            <li v-for="(comp, idx) in currentDetailedResep.komponen" :key="idx" class="mb-1">
              {{ comp.kuantitas }} {{ comp.tipe === 'bahan_baku' ? comp.satuan : '' }} dari <strong>{{ comp.nama }}</strong> (Tipe: {{ comp.tipe }})
              <span v-if="comp.tipe === 'bahan_baku' && comp.harga_unit"> (Harga Unit Pakai: Rp {{ comp.harga_unit.toLocaleString('id-ID', { minimumFractionDigits: 2, maximumFractionDigits: 4 }) }})</span>
            </li>
          </ul>
        </div>
        <div v-else class="text-center text-gray-500">Memuat detail resep...</div>

        <div class="flex justify-end mt-6">
          <button @click="closeDetailModal" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Tidak ada gaya scoped karena menggunakan Tailwind CSS */
</style>