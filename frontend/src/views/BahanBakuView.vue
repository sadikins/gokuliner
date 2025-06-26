<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import BahanBakuList from '../components/BahanBakuList.vue'

const bahanBakus = ref([])
// Menggunakan satu objek formModel untuk mode tambah dan edit
const formModel = ref({
  id: '', // Diisi hanya saat mode edit
  nama: '',
  kategori: '',
  harga_beli: 0,        // Akan dikirim sebagai number, backend akan convert ke decimal
  satuan_beli: '',
  netto_per_beli: 0,    // Akan dikirim sebagai number, backend akan convert ke decimal
  satuan_pemakaian: '',
  catatan: ''
})

const isEditing = ref(false) // State untuk menandakan apakah sedang dalam mode edit

const API_BASE_URL = 'http://localhost:8080/api'

// Daftar kategori yang bisa dipilih (sesuai contoh gambar Kalkuliner)
const categories = [
  'Bahan Kering', 'Bahan Cair', 'Bumbu', 'Kemasan',
  'Topping', 'Pelengkap', 'Energi', 'Bahan Protein', 'Uncategorized'
]

// Fungsi untuk mengambil semua bahan baku dari backend
const fetchBahanBakus = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/bahan-bakus`)
    bahanBakus.value = response.data
    console.log('Bahan baku berhasil diambil:', response.data);
  } catch (error) {
    console.error('Error fetching bahan baku:', error)
    alert('Gagal mengambil data bahan baku.')
  }
}

// Fungsi untuk menangani submit form (tambah atau update)
const handleSubmit = async () => {
  // Validasi sederhana untuk harga dan netto
  if (formModel.value.harga_beli <= 0 || formModel.value.netto_per_beli <= 0) {
    alert('Harga Beli dan Netto harus lebih dari 0.')
    return;
  }

  try {
    if (isEditing.value) {
      // Jika mode edit, kirim PUT request
      await axios.put(`${API_BASE_URL}/bahan-bakus/${formModel.value.id}`, formModel.value)
      alert('Bahan baku berhasil diperbarui!')
    } else {
      // Jika mode tambah, kirim POST request
      await axios.post(`${API_BASE_URL}/bahan-bakus`, formModel.value)
      alert('Bahan baku berhasil ditambahkan!')
    }

    resetForm() // Reset form setelah sukses
    fetchBahanBakus() // Muat ulang daftar bahan baku
  } catch (error) {
    console.error('Error saving bahan baku:', error.response ? error.response.data : error)
    const errorMessage = error.response && error.response.data && error.response.data.error
                         ? error.response.data.error
                         : 'Gagal menyimpan bahan baku. Pastikan nama unik dan data valid.';
    alert(errorMessage);
  }
}

// Fungsi untuk mengaktifkan mode edit
const editBahanBaku = (bahanBaku) => {
  isEditing.value = true
  // Isi formModel dengan data dari bahan baku yang akan diedit (salinan objek)
  formModel.value = { ...bahanBaku }
}

// Fungsi untuk membatalkan mode edit dan mereset form
const cancelEdit = () => {
  resetForm()
}

// Fungsi untuk mereset form ke kondisi awal (kosong atau mode tambah)
const resetForm = () => {
  isEditing.value = false
  formModel.value = {
    id: '',
    nama: '',
    kategori: '',
    harga_beli: 0,
    satuan_beli: '',
    netto_per_beli: 0,
    satuan_pemakaian: '',
    catatan: ''
  }
}

// Fungsi untuk menghapus bahan baku
const deleteBahanBaku = async (id, nama) => {
  if (confirm(`Apakah Anda yakin ingin menghapus bahan baku "${nama}"?`)) {
    try {
      await axios.delete(`${API_BASE_URL}/bahan-bakus/${id}`)
      alert('Bahan baku berhasil dihapus!')
      fetchBahanBakus() // Muat ulang daftar bahan baku
      resetForm() // Reset form jika bahan baku yang sedang diedit dihapus
    } catch (error) {
      console.error('Error deleting bahan baku:', error.response ? error.response.data : error)
      const errorMessage = error.response && error.response.data && error.response.data.error
                           ? error.response.data.error
                           : 'Gagal menghapus bahan baku. Mungkin sedang digunakan dalam resep lain.';
      alert(errorMessage);
    }
  }
}

// Panggil fetchBahanBakus saat komponen pertama kali dimuat
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
          <label for="kategori" class="block text-gray-700 text-sm font-bold mb-2">Kategori:</label>
          <select id="kategori" v-model="formModel.kategori" required
                  class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
            <option value="">-- Pilih Kategori --</option>
            <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
          </select>
        </div>
        <div class="mb-4">
          <label for="harga_beli" class="block text-gray-700 text-sm font-bold mb-2">Harga Beli (per satuan beli):</label>
          <input type="number" id="harga_beli" v-model.number="formModel.harga_beli" step="0.0001" min="0" required
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
        </div>
        <div class="mb-4">
          <label for="satuan_beli" class="block text-gray-700 text-sm font-bold mb-2">Satuan Beli:</label>
          <input type="text" id="satuan_beli" v-model="formModel.satuan_beli" required
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
        </div>
        <div class="mb-4">
          <label for="netto_per_beli" class="block text-gray-700 text-sm font-bold mb-2">Netto (Satuan Pemakaian per Satuan Beli):</label>
          <input type="number" id="netto_per_beli" v-model.number="formModel.netto_per_beli" step="0.0001" min="0" required
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
        </div>
        <div class="mb-4">
          <label for="satuan_pemakaian" class="block text-gray-700 text-sm font-bold mb-2">Satuan Pemakaian (di resep):</label>
          <input type="text" id="satuan_pemakaian" v-model="formModel.satuan_pemakaian" required
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
        </div>
        <div class="mb-6">
          <label for="catatan" class="block text-gray-700 text-sm font-bold mb-2">Catatan:</label>
          <textarea id="catatan" v-model="formModel.catatan"
                    class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline h-20"></textarea>
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
/* Tidak ada gaya scoped karena menggunakan Tailwind CSS */
</style>