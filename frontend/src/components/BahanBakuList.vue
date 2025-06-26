<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  bahanBakus: Array
})

const emit = defineEmits(['edit', 'delete'])

const formatRupiah = (value) => {
  return value.toLocaleString('id-ID', { minimumFractionDigits: 0, maximumFractionDigits: 2 });
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' });
}
</script>

<template>
  <div class="overflow-x-auto">
    <table class="min-w-full bg-white border border-gray-200">
      <thead class="bg-gray-100">
        <tr>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Nama</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Kategori</th> <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Harga Beli</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Satuan Beli</th> <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Netto</th>      <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Satuan Pakai</th> <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Catatan</th>    <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Tanggal Dibuat</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Aksi</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="bahanBakus.length === 0">
          <td colspan="9" class="py-4 px-4 text-center text-gray-500">Tidak ada bahan baku. Tambahkan yang baru di atas.</td> </tr>
        <tr v-for="(bb, index) in bahanBakus" :key="bb.id" :class="index % 2 === 0 ? 'bg-white' : 'bg-gray-50'">
          <td class="py-3 px-4 border-b text-sm text-gray-800">{{ bb.nama }}</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800">{{ bb.kategori }}</td>           <td class="py-3 px-4 border-b text-sm text-gray-800">Rp {{ formatRupiah(bb.harga_beli) }}</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800">{{ bb.satuan_beli }}</td>        <td class="py-3 px-4 border-b text-sm text-gray-800">{{ bb.netto_per_beli }}</td>     <td class="py-3 px-4 border-b text-sm text-gray-800">{{ bb.satuan_pemakaian }}</td>   <td class="py-3 px-4 border-b text-sm text-gray-800">{{ bb.catatan }}</td>            <td class="py-3 px-4 border-b text-sm text-gray-800">{{ formatDate(bb.created_at) }}</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800">
            <button @click="emit('edit', bb)"
                    class="bg-yellow-500 hover:bg-yellow-600 text-white font-bold py-1 px-3 rounded-md text-xs mr-2">
              Edit
            </button>
            <button @click="emit('delete', bb.id, bb.nama)"
                    class="bg-red-500 hover:bg-red-600 text-white font-bold py-1 px-3 rounded-md text-xs">
              Hapus
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
/* Hapus semua gaya scoped yang ada di sini */
</style>