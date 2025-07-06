<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import ResepList from '../components/List/ResepList.vue' // Komponen untuk menampilkan daftar

const API_BASE_URL = 'http://localhost:8080/api'

// --- State Data ---
const reseps = ref([])
const bahanBakus = ref([])
const existingReseps = ref([])

// --- State Form Utama (Kiri Atas) ---
const isEditing = ref(false)
const selectedResepForForm = ref(null) // Menyimpan data resep yang sedang diedit (untuk tampilan H2)

const formModel = ref({
  id: '',
  nama: '',
  is_sub_resep: false,
  jumlah_porsi: 1,
  komponen: [
    { komponen_id: '', kuantitas: 0, tipe_komponen: 'bahan_baku' }
  ]
})

// --- State untuk Panel Rincian Perhitungan (Kanan Atas) ---
const selectedResepForDetails = ref(null) // Resep yang detailnya ditampilkan
const hppResultForDetails = ref(null) // Hasil HPP untuk resep yang detailnya ditampilkan

// --- State untuk Modal Detail Resep (yang lama, bisa dipertimbangkan untuk diganti oleh panel rincian) ---
const showDetailModal = ref(false) // Ini modal yang muncul di tengah layar
const currentDetailedResep = ref(null) // Data untuk modal detail


// --- Fetch Data Awal (tidak berubah) ---
const fetchReseps = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/reseps`)
    reseps.value = response.data
  } catch (error) {
    console.error('Error fetching reseps:', error)
    alert('Gagal mengambil data resep.')
  }
}

const fetchBahanBakus = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/bahan-bakus`)
    bahanBakus.value = response.data
  } catch (error) {
    console.error('Error fetching bahan baku:', error)
    alert('Gagal mengambil data bahan baku untuk pilihan.')
  }
}

const fetchExistingReseps = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/reseps`)
    existingReseps.value = response.data.filter(r => r.id !== (formModel.value.id ? formModel.value.id : null));
  } catch (error) {
    console.error('Error fetching existing reseps:', error)
    alert('Gagal mengambil resep yang sudah ada untuk pilihan.')
  }
}

// --- Handler untuk Form Penambahan/Update Resep (tidak banyak berubah) ---
const handleSubmitResep = async () => {
  if (formModel.value.komponen.some(c => !c.komponen_id || c.kuantitas <= 0)) {
    alert('Pastikan semua komponen dipilih dan kuantitas lebih dari 0.')
    return
  }
  if (formModel.value.jumlah_porsi <= 0) {
    alert('Jumlah porsi harus lebih dari 0.')
    return
  }

  try {
    if (isEditing.value) {
      await axios.put(`${API_BASE_URL}/reseps/${formModel.value.id}`, formModel.value)
      alert('Resep berhasil diperbarui!')
    } else {
      await axios.post(`${API_BASE_URL}/reseps`, formModel.value)
      alert('Resep berhasil ditambahkan!')
    }

    resetForm()
    fetchReseps()
    fetchExistingReseps()
  } catch (error) {
    console.error('Error saving resep:', error.response ? error.response.data : error)
    const errorMessage = error.response && error.response.data && error.response.data.error
                         ? error.response.data.error
                         : 'Gagal menyimpan resep. Pastikan nama unik dan data valid.';
    alert(errorMessage);
  }
}

// --- Fungsi untuk Mengelola Komponen Resep Secara Dinamis di Form (tidak berubah) ---
const addKomponen = () => {
  formModel.value.komponen.push({ komponen_id: '', kuantitas: 0, tipe_komponen: 'bahan_baku' })
}

const removeKomponen = (index) => {
  formModel.value.komponen.splice(index, 1)
}

// --- Fungsi Edit/Delete Resep (Diperbarui untuk mengisi form saja, tidak mengisi detail panel) ---
const editResep = (resep) => {
  isEditing.value = true
  selectedResepForForm.value = { ...resep } // Untuk H2 form
  formModel.value = { // Isi formModel
    id: resep.id,
    nama: resep.nama,
    is_sub_resep: resep.is_sub_resep,
    jumlah_porsi: resep.jumlah_porsi || 1,
    komponen: resep.komponen.map(c => ({ ...c }))
  }
  fetchExistingReseps();
  // Tidak mengisi selectedResepForDetails / hppResultForDetails di sini
  window.scrollTo({ top: 0, behavior: 'smooth' }) // Scroll ke atas form
}

const deleteResep = async (id, nama) => {
  if (confirm(`Apakah Anda yakin ingin menghapus resep "${nama}"? Semua komponen terkait juga akan dihapus.`)) {
    try {
      await axios.delete(`${API_BASE_URL}/reseps/${id}`)
      alert('Resep berhasil dihapus!')
      fetchReseps()
      fetchExistingReseps()
      resetForm() // Reset form jika resep yang diedit/dihapus
      selectedResepForDetails.value = null; // Clear detail panel
      hppResultForDetails.value = null; // Clear HPP detail panel
    } catch (error) {
      console.error('Error deleting resep:', error.response ? error.response.data : error)
      const errorMessage = error.response && error.response.data && error.response.data.error
                           ? error.response.data.error
                           : 'Gagal menghapus resep. Mungkin sedang digunakan sebagai komponen di resep lain.';
      alert(errorMessage);
    }
  }
}

const cancelEdit = () => {
  resetForm()
}

const resetForm = () => {
  isEditing.value = false
  selectedResepForForm.value = null
  formModel.value = {
    id: '',
    nama: '',
    is_sub_resep: false,
    jumlah_porsi: 1,
    komponen: [{ komponen_id: '', kuantitas: 0, tipe_komponen: 'bahan_baku' }]
  }
  fetchExistingReseps();
  // Tidak mempengaruhi selectedResepForDetails / hppResultForDetails di sini
}


// --- Fungsi untuk Menampilkan Detail di Panel Kanan Atas ---
const showResepCalculationDetails = async (resepId) => {
  try {
    // Ambil detail resep
    const resepResponse = await axios.get(`${API_BASE_URL}/reseps/${resepId}`)
    selectedResepForDetails.value = resepResponse.data

    // Ambil HPP terbaru untuk resep ini
    const hppResponse = await axios.get(`${API_BASE_URL}/hpp/${resepId}`)
    hppResultForDetails.value = hppResponse.data

  } catch (error) {
    console.error('Error fetching resep details or HPP for panel:', error.response ? error.response.data : error)
    alert('Gagal mengambil detail resep atau HPP. Pastikan HPP sudah dihitung untuk resep ini.')
    selectedResepForDetails.value = null
    hppResultForDetails.value = null
  }
}

// --- Fungsi untuk Menampilkan/Menutup Modal Detail Resep (yang lama) ---
// Ini adalah modal yang akan muncul di tengah layar.
// Jika ingin semua detail hanya di panel kanan atas, fungsi ini dan modalnya bisa dihapus.
// Saya biarkan untuk demonstrasi bahwa Anda memiliki kedua opsi.
const showDetailModalLegacy = (resep) => {
  currentDetailedResep.value = resep;
  showDetailModal.value = true;
};
const closeDetailModalLegacy = () => {
  showDetailModal.value = false;
  currentDetailedResep.value = null;
};


// --- Fungsi untuk Duplikat Resep (tidak berubah) ---
const duplicateResep = async (resepId, resepNama) => {
  if (confirm(`Apakah Anda yakin ingin menduplikasi resep "${resepNama}"?`)) {
    try {
      const response = await axios.post(`${API_BASE_URL}/reseps/${resepId}/duplicate`)
      alert(`Resep "${resepNama}" berhasil diduplikasi menjadi "${response.data.nama_resep_baru}"!`)
      fetchReseps()
      fetchExistingReseps()
    } catch (error) {
      console.error('Error duplicating resep:', error.response ? error.response.data : error)
      const errorMessage = error.response && error.response.data && error.response.data.error
                           ? error.response.data.error
                           : 'Gagal menduplikasi resep.';
      alert(errorMessage);
    }
  }
}

// --- On Mounted (tidak berubah) ---
onMounted(() => {
  fetchReseps()
  fetchBahanBakus()
  fetchExistingReseps()
} )


const handleHPPCalculated = (hppResult) => {
  console.log('HPP diterima di parent (ResepView):', hppResult)
  // Anda bisa melakukan sesuatu di sini, misalnya menyimpan hasil HPP ke state lain
  // atau menampilkannya di bagian terpisah dari halaman ResepView.
}

// --- Helper untuk Formatting (bisa dipindahkan ke file utilitas frontend) ---
const formatRupiah = (value) => {
  const num = parseFloat(value);
  if (isNaN(num)) return 'Rp 0.00';
  return 'Rp ' + num.toLocaleString('id-ID', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
}

const formatPersen = (value) => {
  const num = parseFloat(value);
  if (isNaN(num)) return '0.00%';
  return num.toLocaleString('id-ID', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) + '%';
}
</script>

<template>
  <div class="">
    <h1 class="text-3xl font-bold text-gray-800 mb-6">Manajemen Resep</h1>

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-6">
      <div class="bg-gray-50 rounded-lg p-6 shadow-md">
        <h2 class="text-2xl font-semibold text-gray-700 mb-5" v-if="!isEditing">Tambah Resep Baru</h2>
        <h2 class="text-2xl font-semibold text-gray-700 mb-5" v-else>Edit Resep: {{ selectedResepForForm?.nama }}</h2>
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
            <input type="number" id="jumlahPorsi" v-model.number="formModel.jumlah_porsi" min="1" step="0.01" required
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
              <input type="number" :id="'kuantitas' + index" v-model.number="komp.kuantitas" min="0" step="0.0001" required
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

      <div class="bg-blue-50 rounded-lg p-6 shadow-md">
        <h2 class="text-2xl font-semibold text-blue-800 mb-5">Rincian Perhitungan Resep</h2>
        <div v-if="selectedResepForDetails">
          <p class="mb-2"><strong class="text-blue-700">Nama Resep:</strong> {{ selectedResepForDetails.nama }}</p>
          <p class="mb-2"><strong class="text-blue-700">Sub-Resep:</strong> {{ selectedResepForDetails.is_sub_resep ? 'Ya' : 'Tidak' }}</p>
          <p class="mb-4"><strong class="text-blue-700">Jumlah Porsi:</strong> {{ selectedResepForDetails.jumlah_porsi }}</p>

          <h3 class="font-bold text-gray-700 mb-3">Komponen:</h3>
          <ul class="list-disc pl-5 mb-4">
            <li v-for="(comp, idx) in selectedResepForDetails.komponen" :key="idx" class="text-sm">
              {{ comp.kuantitas }} {{ comp.tipe === 'bahan_baku' ? comp.satuan : '' }} dari <strong>{{ comp.nama }}</strong> (Tipe: {{ comp.tipe }})
              <span v-if="comp.tipe === 'bahan_baku' && comp.harga_unit > 0"> (Harga Unit Pakai: {{ formatRupiah(comp.harga_unit) }})</span>
            </li>
          </ul>

          <h3 class="font-bold text-gray-700 mb-3">HPP Resep:</h3>
          <div v-if="hppResultForDetails">
            <p><strong class="text-green-700">HPP per Unit Resep:</strong> {{ formatRupiah(hppResultForDetails.hpp_per_unit) }}</p>
            <p><strong class="text-green-700">HPP per Porsi:</strong> {{ formatRupiah(hppResultForDetails.hpp_per_porsi) }}</p>
            <p class="text-xs text-gray-500 mt-2">Dihitung pada: {{ new Date(hppResultForDetails.created_at).toLocaleString() }}</p>
          </div>
          <p v-else class="text-gray-500">HPP belum dihitung untuk resep ini. Klik "Hitung HPP" di bawah.</p>

        </div>
        <div v-else class="text-center text-gray-500 py-4">
          Pilih resep dari daftar di bawah untuk melihat rincian perhitungan.
        </div>
      </div>
    </div>

    <div class="bg-gray-50 rounded-lg p-6 shadow-md">
      <h2 class="text-2xl font-semibold text-gray-700 mb-5">Daftar Resep</h2>
      <div class="flex-grow">
        <ResepList
          :reseps="reseps"
          @hppCalculated="handleHPPCalculated"
          @edit="editResep"
          @delete="deleteResep"
          @showDetail="showResepCalculationDetails"
          @duplicate="duplicateResep"
        />
      </div>
    </div>

    <div v-if="showDetailModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 flex justify-center items-center p-4 z-50">
      <div class="bg-white rounded-lg shadow-xl max-w-2xl w-full p-6">
        <div class="flex justify-between items-center border-b pb-3 mb-4">
          <h3 class="text-2xl font-bold text-gray-800">Detail Resep: {{ currentDetailedResep?.nama }}</h3>
          <button @click="closeDetailModalLegacy" class="text-gray-500 hover:text-gray-700 text-3xl font-semibold">&times;</button>
        </div>

        <div v-if="currentDetailedResep" class="text-gray-700 text-sm grid grid-cols-2 gap-y-2 gap-x-4">
          <p class="col-span-2"><strong class="text-gray-700">ID:</strong> {{ currentDetailedResep.id }}</p>
          <p><strong class="text-gray-700">Nama Resep:</strong> {{ currentDetailedResep.nama }}</p>
          <p><strong class="text-gray-700">Sub-Resep:</strong> {{ currentDetailedResep.is_sub_resep ? 'Ya' : 'Tidak' }}</p>
          <p><strong class="text-gray-700">Jumlah Porsi:</strong> {{ currentDetailedResep.jumlah_porsi }}</p>
          <p><strong class="text-gray-700">Dibuat:</strong> {{ currentDetailedResep.created_at }}</p>
          <p class="mb-4"><strong class="text-gray-700">Diperbarui:</strong> {{ currentDetailedResep.updated_at }}</p>

          <h4 class="text-lg font-semibold text-gray-700 col-span-2 mb-3">Daftar Komponen:</h4>
          <ul class="list-disc pl-5 col-span-2">
            <li v-for="(comp, idx) in currentDetailedResep.komponen" :key="idx" class="mb-1">
              {{ comp.kuantitas }} {{ comp.tipe === 'bahan_baku' ? comp.satuan : '' }} dari <strong>{{ comp.nama }}</strong> (Tipe: {{ comp.tipe }})
              <span v-if="comp.tipe === 'bahan_baku' && comp.harga_unit > 0"> (Harga Unit Pakai: {{ formatRupiah(comp.harga_unit) }})</span>
            </li>
          </ul>
        </div>
        <div v-else class="text-center text-gray-500">Memuat detail resep...</div>

        <div class="flex justify-end mt-6">
          <button @click="closeDetailModalLegacy" class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">
            Tutup
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Hapus semua gaya scoped */
</style>