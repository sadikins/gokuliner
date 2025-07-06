<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import ProgramPromoList from '../components/List/ProgramPromoList.vue'

const API_BASE_URL = 'http://localhost:8080/api'

const programPromos = ref([])
const formModel = ref({
  id: '',
  nama_promo: '',
  channel: '',
  jenis_diskon: 'persentase', // Default
  besar_diskon: 0,
  min_belanja: 0,
  maksimal_potongan: 0,
  ditanggung_merchant_persen: 0,
  catatan: ''
})

const isEditing = ref(false)

// Daftar pilihan channel (sesuai video)
const channels = ['GoFood', 'GrabFood', 'ShopeeFood', 'Internal', 'Lainnya', 'Semua']

// --- Fetch Data ---
const fetchProgramPromos = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/program-promos`)
    programPromos.value = response.data
  } catch (error) {
    console.error('Error fetching program promos:', error)
    alert('Gagal mengambil data program promo.')
  }
}

// --- Handler Form ---
const handleSubmit = async () => {
  // Validasi dasar
  if (formModel.value.besar_diskon <= 0 && formModel.value.jenis_diskon === 'persentase') {
    alert('Besar diskon (persentase) harus lebih dari 0.')
    return
  }
  if (formModel.value.min_belanja < 0 || formModel.value.maksimal_potongan < 0 || formModel.value.ditanggung_merchant_persen < 0) {
    alert('Nilai nominal tidak boleh negatif.')
    return
  }

  try {
    if (isEditing.value) {
      await axios.put(`${API_BASE_URL}/program-promos/${formModel.value.id}`, formModel.value)
      alert('Program promo berhasil diperbarui!')
    } else {
      await axios.post(`${API_BASE_URL}/program-promos`, formModel.value)
      alert('Program promo berhasil ditambahkan!')
    }

    resetForm()
    fetchProgramPromos()
  } catch (error) {
    console.error('Error saving program promo:', error.response ? error.response.data : error)
    const errorMessage = error.response && error.response.data && error.response.data.error
                         ? error.response.data.error
                         : 'Gagal menyimpan program promo. Pastikan nama unik dan data valid.';
    alert(errorMessage);
  }
}

// --- Fungsi Edit/Delete ---
const editProgramPromo = (promo) => {
  isEditing.value = true
  formModel.value = { ...promo } // Salin objek
}

const deleteProgramPromo = async (id, nama) => {
  if (confirm(`Apakah Anda yakin ingin menghapus program promo "${nama}"?`)) {
    try {
      await axios.delete(`${API_BASE_URL}/program-promos/${id}`)
      alert('Program promo berhasil dihapus!')
      fetchProgramPromos()
      resetForm()
    } catch (error) {
      console.error('Error deleting program promo:', error.response ? error.response.data : error)
      const errorMessage = error.response && error.response.data && error.response.data.error
                           ? error.response.data.error
                           : 'Gagal menghapus program promo. Mungkin terkait dengan data lain.';
      alert(errorMessage);
    }
  }
}

const cancelEdit = () => {
  resetForm()
}

const resetForm = () => {
  isEditing.value = false
  formModel.value = {
    id: '',
    nama_promo: '',
    channel: '',
    jenis_diskon: 'persentase',
    besar_diskon: 0,
    min_belanja: 0,
    maksimal_potongan: 0,
    ditanggung_merchant_persen: 0,
    catatan: ''
  }
}

onMounted(fetchProgramPromos)
</script>

<template>
  <div class="">
    <h1 class="text-3xl font-bold text-gray-800 mb-6">Manajemen Program Promo</h1>

    <div class="bg-gray-50 rounded-lg p-6 mb-8 shadow-md">
      <h2 class="text-2xl font-semibold text-gray-700 mb-5" v-if="!isEditing">Tambah Program Promo Baru</h2>
      <h2 class="text-2xl font-semibold text-gray-700 mb-5" v-else>Edit Program Promo: {{ formModel.nama_promo }}</h2>
      <form @submit.prevent="handleSubmit()">
        <div class="mb-4">
          <label for="namaPromo" class="block text-gray-700 text-sm font-bold mb-2">Nama Promo:</label>
          <input type="text" id="namaPromo" v-model="formModel.nama_promo" required
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
        </div>
        <div class="mb-4">
          <label for="channel" class="block text-gray-700 text-sm font-bold mb-2">Channel:</label>
          <select id="channel" v-model="formModel.channel" required
                  class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
            <option value="">-- Pilih Channel --</option>
            <option v-for="ch in channels" :key="ch" :value="ch">{{ ch }}</option>
          </select>
        </div>
        <div class="mb-4">
          <label for="jenisDiskon" class="block text-gray-700 text-sm font-bold mb-2">Jenis Diskon:</label>
          <select id="jenisDiskon" v-model="formModel.jenis_diskon" required
                  class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
            <option value="persentase">Persentase (%)</option>
            <option value="nominal">Nominal (Rp)</option>
          </select>
        </div>
        <div class="mb-4">
          <label for="besarDiskon" class="block text-gray-700 text-sm font-bold mb-2">
            Besar Diskon ({{ formModel.jenis_diskon === 'persentase' ? '%' : 'Rp' }}):
          </label>
          <input type="number" id="besarDiskon" v-model.number="formModel.besar_diskon" min="0" step="0.01" required
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
        </div>
        <div class="mb-4">
          <label for="minBelanja" class="block text-gray-700 text-sm font-bold mb-2">Min Belanja (Rp):</label>
          <input type="number" id="minBelanja" v-model.number="formModel.min_belanja" min="0" step="0.01"
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
        </div>
        <div class="mb-4">
          <label for="maksimalPotongan" class="block text-gray-700 text-sm font-bold mb-2">Maksimal Potongan (Rp):</label>
          <input type="number" id="maksimalPotongan" v-model.number="formModel.maksimal_potongan" min="0" step="0.01"
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
        </div>
        <div class="mb-4">
          <label for="ditanggungMerchantPersen" class="block text-gray-700 text-sm font-bold mb-2">Ditanggung Merchant (%):</label>
          <input type="number" id="ditanggungMerchantPersen" v-model.number="formModel.ditanggung_merchant_persen" min="0" max="100" step="0.01"
                 class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline">
        </div>
        <div class="mb-6">
          <label for="catatan" class="block text-gray-700 text-sm font-bold mb-2">Catatan:</label>
          <textarea id="catatan" v-model="formModel.catatan"
                    class="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline h-20"></textarea>
        </div>

        <button type="submit" v-if="!isEditing"
                class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline">
          Tambah Program Promo
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
      <h2 class="text-2xl font-semibold text-gray-700 mb-5">Daftar Program Promo</h2>
      <ProgramPromoList
        :programPromos="programPromos"
        @edit="editProgramPromo"
        @delete="deleteProgramPromo"
      />
    </div>
  </div>
</template>

<style scoped>
/* Hapus semua gaya scoped */
</style>