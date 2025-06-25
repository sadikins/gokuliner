<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import ResepList from '../components/ResepList.vue'

const API_BASE_URL = 'http://localhost:8080/api'

// --- State Reaktif ---
const reseps = ref([])
const bahanBakus = ref([])
const existingReseps = ref([])

const isEditing = ref(false)
const selectedResep = ref(null)

const formModel = ref({
  id: '', // Untuk mode edit
  nama: '',
  is_sub_resep: false,
  jumlah_porsi: 1,
  komponen: [
    { komponen_id: '', kuantitas: 0, tipe_komponen: 'bahan_baku' }
  ]
})

// --- Fetch Data Awal ---
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

const fetchExistingReseps = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/reseps`)
    existingReseps.value = response.data.filter(r => r.id !== (formModel.value.id ? formModel.value.id : null));
    console.log('Resep yang sudah ada berhasil diambil:', response.data);
  } catch (error) {
    console.error('Error fetching existing reseps:', error)
    alert('Gagal mengambil resep yang sudah ada untuk pilihan.')
  }
}

// --- Handler Form (add/update) ---
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

// --- Dynamic Komponen Management ---
const addKomponen = () => {
  formModel.value.komponen.push({ komponen_id: '', kuantitas: 0, tipe_komponen: 'bahan_baku' })
}

const removeKomponen = (index) => {
  formModel.value.komponen.splice(index, 1)
}

// --- Edit/Delete Handlers dari ResepList ---
const editResep = (resep) => {
  isEditing.value = true
  selectedResep.value = { ...resep }
  formModel.value = {
    id: resep.id,
    nama: resep.nama,
    is_sub_resep: resep.is_sub_resep,
    jumlah_porsi: resep.jumlah_porsi || 1,
    komponen: resep.komponen.map(c => ({ ...c }))
  }
  fetchExistingReseps();
}

const deleteResep = async (id, nama) => {
  if (confirm(`Apakah Anda yakin ingin menghapus resep "${nama}"? Semua komponen terkait juga akan dihapus.`)) {
    try {
      await axios.delete(`${API_BASE_URL}/reseps/${id}`)
      alert('Resep berhasil dihapus!')
      fetchReseps()
      fetchExistingReseps()
      resetForm()
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
  selectedResep.value = null
  formModel.value = {
    id: '',
    nama: '',
    is_sub_resep: false,
    jumlah_porsi: 1,
    komponen: [{ komponen_id: '', kuantitas: 0, tipe_komponen: 'bahan_baku' }]
  }
  fetchExistingReseps();
}

// --- On Mounted ---
onMounted(() => {
  fetchReseps()
  fetchBahanBakus()
  fetchExistingReseps()
})

// --- Handler Event dari Child Component ResepList (HPP Calculated) ---
const handleHPPCalculated = (hppResult) => {
  console.log('HPP diterima di parent (ResepView):', hppResult)
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
                <option v-for="bb in bahanBakus" :key="bb.id" :value="bb.id">{{ bb.nama }} ({{ bb.satuan }})</option>
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
      />
    </div>
  </div>
</template>

<style scoped>
/* Hapus semua gaya scoped yang ada di sini */
</style>