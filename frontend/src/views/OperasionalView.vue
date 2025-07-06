<script setup>
import { ref, onMounted, computed } from 'vue'
import axios from 'axios'
import TenagaKerjaList from '../components/List/TenagaKerjaList.vue';
import OperasionalList from '../components/List/OperasionalList.vue';
import InputText from '../components/InputText.vue';

const allCostItems = ref([])
// Menggunakan satu objek formModel untuk mode tambah dan edit
const formModel = ref({
  id: '', // Diisi hanya saat mode edit
  nama: '',
  kategori: 'Operasional',
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

const filteredCostItems = computed(() => {
  return allCostItems.value.filter(item => item.kategori === 'Operasional');
});

// Fungsi untuk mengambil semua bahan baku dari backend
const fetchBahanBakus = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/bahan-bakus`)
    allCostItems.value = response.data
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
    kategori: 'Operasional',
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
  <div class="">
    <h2 class="text-3xl font-bold mb-6">Manajemen Operasional</h2>


    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <div class="lg:col-span-2 bg-white rounded-xl p-6 border border-gray-200 h-full flex flex-col">
        <h3 class="text-xl font-semibold mb-4">Daftar Operasional</h3>
        <div class="flex-grow">
          <OperasionalList
            :operasionals="filteredCostItems"
            @edit="editBahanBaku"
            @delete="deleteBahanBaku"
          />
        </div>
      </div>

      <div class="lg:col-span-1 bg-white rounded-xl p-6 border border-gray-200">
        <h2 class="text-2xl font-semibold text-gray-700 mb-5" v-if="!isEditing">Tambah Operasional</h2>
        <h2 class="text-2xl font-semibold text-gray-700 mb-5" v-else>Edit Operasional</h2>
        <form @submit.prevent="handleSubmit()">
          <div class="mb-4">
            <label for="nama" class="block mb-1 text-sm font-medium text-gray-700">Nama Kegiatan:</label>
            <InputText  id="nama" v-model="formModel.nama" placeholder="Contoh: Gas 12kg" required/>
          </div>
          <!-- <div class="mb-4">
            <label for="kategori" class="block mb-1 text-sm font-medium text-gray-700">Kategori:</label>
            <select id="kategori" v-model="formModel.kategori" required
                    class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg block w-full p-2.5">
              <option value="">-- Pilih Kategori --</option>
              <option v-for="cat in categories" :key="cat" :value="cat">{{ cat }}</option>
            </select>
          </div> -->
          <div class="grid grid-cols-2 gap-4 mb-4">
            <div>
              <label for="harga_beli" class="block mb-1 text-sm font-medium text-gray-700">Bayaran (Rp):</label>
              <InputText type="number" id="harga_beli" v-model.number="formModel.harga_beli" min="0" step="0.01" required
                    />
            </div>
            <div>
              <label for="satuan_beli" class="block mb-1 text-sm font-medium text-gray-700">Satuan Pembayaran:</label>
            <InputText type="text" id="satuan_beli" v-model="formModel.satuan_beli" placeholder="Contoh: Tabung" required
                  /> </div>
          </div>
          <div class="mb-4">

          </div>
          <div class="mb-4 border border-slate-300  border-dashed  p-3 rounded-md bg-subtle">
            <h3 class="text-base font-medium text-gray-700 mb-2">Konversi Satuan</h3>
            <div class="grid grid-cols-2 gap-4">
                <div>
                    <label for="netto_per_beli" class="block mb-1 text-sm font-medium text-gray-700">Pemakaian:</label>
                    <InputText type="number" id="netto_per_beli" v-model.number="formModel.netto_per_beli" min="0" step="0.0001" required
                           />
                </div>
                <div>
                    <label for="satuan_pemakaian" class="block mb-1 text-sm font-medium text-gray-700">Satuan:</label>
                    <InputText type="text" id="satuan_pemakaian" v-model="formModel.satuan_pemakaian" placeholder="Contoh: Menit" required
                           />
                </div>
            </div>
          </div>

          <div class="mb-1">
            <label for="catatan" class="block text-sm font-medium text-lime-900">Catatan:</label>
            <textarea id="catatan" v-model="formModel.catatan"  placeholder="Catatan tambahan tentang tenaga kerja ini"
                      class="bg-lime-100 p-3 rounded-lg w-full text-lime-900 placeholder:text-lime-800/50 placeholder:text-sm"></textarea>
          </div>

          <div class="items-center space-y-2.5">
            <button type="submit"
                    class="w-full text-lime-900 bg-lime-400 hover:bg-lime-500 focus:ring-4 focus:outline-none focus:ring-lime-200 font-medium rounded-lg px-5 py-2.5 text-center">

              {{ isEditing ? 'Simpan Perubahan' : 'Simpan Bahan' }}
            </button>
            <button type="button" @click="cancelEdit()"
                    class="w-full text-slate-500 bg-slate-100 hover:bg-slate-200  focus:ring-4 focus:outline-none focus:ring-slate-300 font-semibold rounded-lg text-sm px-5 py-2.5 text-center"
                    v-if="isEditing">
              Batal
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>