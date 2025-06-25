<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import BahanBakuList from '../components/BahanBakuList.vue'

const bahanBakus = ref([])
// Gunakan SATU objek untuk binding form, baik itu untuk CREATE atau UPDATE
const formModel = ref({
  id: '', // Tambahkan ID untuk kasus update
  nama: '',
  satuan: '',
  harga_beli: 0
})

const isEditing = ref(false)

const API_BASE_URL = 'http://localhost:8080/api'

const fetchBahanBakus = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/bahan-bakus`)
    bahanBakus.value = response.data
  } catch (error) {
    console.error('Error fetching bahan baku:', error)
    alert('Gagal mengambil data bahan baku.')
  }
}

const handleSubmit = async () => { // Fungsi ini akan menangani CREATE dan UPDATE
  try {
    if (isEditing.value) {
      // Logic untuk UPDATE
      await axios.put(`${API_BASE_URL}/bahan-bakus/${formModel.value.id}`, formModel.value)
      alert('Bahan baku berhasil diperbarui!')
    } else {
      // Logic untuk CREATE
      await axios.post(`${API_BASE_URL}/bahan-bakus`, formModel.value)
      alert('Bahan baku berhasil ditambahkan!')
    }

    resetForm() // Reset form setelah berhasil
    fetchBahanBakus() // Refresh list
  } catch (error) {
    console.error('Error saving bahan baku:', error.response ? error.response.data : error)
    const errorMessage = error.response && error.response.data && error.response.data.error
                         ? error.response.data.error
                         : 'Gagal menyimpan bahan baku. Pastikan nama unik dan data valid.';
    alert(errorMessage);
  }
}

const editBahanBaku = (bahanBaku) => {
  isEditing.value = true
  // Isi formModel dengan data dari bahanBaku yang dipilih
  formModel.value = { ...bahanBaku } // Salin objek
}

const cancelEdit = () => {
  resetForm()
}

const resetForm = () => {
  isEditing.value = false
  formModel.value = { id: '', nama: '', satuan: '', harga_beli: 0 } // Reset ke nilai awal
}

const deleteBahanBaku = async (id, nama) => {
  if (confirm(`Apakah Anda yakin ingin menghapus bahan baku "${nama}"?`)) {
    try {
      await axios.delete(`${API_BASE_URL}/bahan-bakus/${id}`)
      alert('Bahan baku berhasil dihapus!')
      fetchBahanBakus() // Refresh list
      resetForm() // Reset form jika bahan baku yang diedit dihapus
    } catch (error) {
      console.error('Error deleting bahan baku:', error.response ? error.response.data : error)
      const errorMessage = error.response && error.response.data && error.response.data.error
                           ? error.response.data.error
                           : 'Gagal menghapus bahan baku. Mungkin sedang digunakan dalam resep lain.';
      alert(errorMessage);
    }
  }
}

onMounted(fetchBahanBakus)
</script>

<template>
  <div class="p-5 max-w-4xl mx-auto font-sans">
    <h1 class="text-3xl font-bold text-gray-800 mb-6">Manajemen Bahan Baku</h1>

    <div class="bg-gray-50 rounded-lg p-6 mb-8 shadow-md">
      <h2 class="text-2xl font-semibold text-gray-700 mb-5" v-if="!isEditing">Tambah Bahan Baku Baru</h2>
      <h2 class="text-2xl font-semibold text-gray-700 mb-5" v-else>Edit Bahan Baku: {{ formModel.nama }}</h2>
      <form @submit.prevent="handleSubmit()">
        <div class="mb-4">
          <label for="nama" class="block text-gray-700 text-sm font-bold mb-2">Nama:</label>
          <input type="text" id="nama" v-model="formModel.nama" required
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
        </div>
        <div class="mb-4">
          <label for="satuan" class="block text-gray-700 text-sm font-bold mb-2">Satuan:</label>
          <input type="text" id="satuan" v-model="formModel.satuan" required
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
        </div>
        <div class="mb-6">
          <label for="harga_beli" class="block text-gray-700 text-sm font-bold mb-2">Harga Beli:</label>
          <input type="number" id="harga_beli" v-model.number="formModel.harga_beli" step="0.01" required
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
        </div>

        <button type="submit" v-if="!isEditing"
                class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
          Tambah
        </button>
        <div v-else class="flex items-center space-x-4">
          <button type="submit"
                  class="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
            Simpan Perubahan
          </button>
          <button type="button" @click="cancelEdit()"
                  class="bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
            Batal
          </button>
        </div>
      </form>
    </div>

    <div class="bg-gray-50 rounded-lg p-6 shadow-md">
      <h2 class="text-2xl font-semibold text-gray-700 mb-5">Daftar Bahan Baku</h2>
      <BahanBakuList
        :bahanBakus="bahanBakus"
        @edit="editBahanBaku"
        @delete="deleteBahanBaku"
      />
    </div>
  </div>
</template>

<style scoped>
/* Gaya CSS Anda yang sudah ada */

.bahan-baku-container {
  padding: 20px;
  max-width: 800px;
  margin: 0 auto;
  font-family: sans-serif;
}

h1, h2 {
  color: #333;
}

.form-section, .list-section {
  background-color: #f9f9f9;
  border-radius: 8px;
  padding: 20px;
  margin-bottom: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.form-group {
  margin-bottom: 15px;
}

label {
  display: block;
  margin-bottom: 5px;
  font-weight: bold;
}

input[type="text"],
input[type="number"] {
  width: calc(100% - 20px);
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
}

button {
  background-color: #007bff;
  color: white;
  padding: 10px 15px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 16px;
}

button:hover {
  background-color: #0056b3;
}

/* Gaya baru untuk tombol edit/hapus */
.edit-buttons button {
  margin-right: 10px;
  width: auto; /* Override default width */
}
.save-btn {
  background-color: #28a745;
}
.save-btn:hover {
  background-color: #218838;
}
.cancel-btn {
  background-color: #6c757d;
}
.cancel-btn:hover {
  background-color: #5a6268;
}
</style>