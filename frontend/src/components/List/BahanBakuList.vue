<script setup>
import { PhPencil, PhTrash } from '@phosphor-icons/vue';
import { defineProps, defineEmits, ref } from 'vue'

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
  <div class="">
    <table class="w-full text-sm text-left text-gray-500">
      <thead class="text-xs text-gray-700 uppercase bg-gray-50">
        <tr>
          <th class="px-6 py-3">#</th>
          <th class="px-6 py-3">Nama</th>
           <th class="px-6 py-3">Harga Beli</th>
          <th class="px-6 py-3">Netto</th>
          <th class="px-6 py-3">Kategori</th>
          <th class="px-6 py-3">Tanggal Dibuat</th>
          <th class="px-6 py-3">Aksi</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="bahanBakus.length === 0">
          <td colspan="9" class="px-6 py-24 font-medium text-center text-gray-400">
            Tidak ada bahan baku. Tambahkan Bahan Baku.
          </td>
        </tr>
        <tr v-for="(bb, index) in bahanBakus" :key="bb.id" class="border-b border-slate-200" >
          <td class="px-6 py-4 font-medium">{{ index+1 }}</td>
          <td class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap">{{ bb.nama }}</td>
          <td class="px-6 py-4">Rp {{ formatRupiah(bb.harga_beli) }} / {{ bb.satuan_beli }}</td>
          <td class="px-6 py-4 font-medium">{{ bb.netto_per_beli }} / {{ bb.satuan_pemakaian }}</td>
          <td class="px-6 py-4 font-medium" >
            <span :class="bb.kategori === 'Kemasan' ? 'bg-lime-400 p-2 text-white rounded-2xl' : '' ">
              {{ bb.kategori }}
            </span>
          </td>
          <td class="px-6 py-4">{{ formatDate(bb.created_at) }}</td>
          <td class="px-6 py-4 flex items-center gap-4">
            <button @click="emit('edit', bb)" class="text-lime-600 hover:text-lime-800">
                <PhPencil class="ph-bold text-lg"/>
            </button>
            <button @click="emit('delete', bb.id, bb.nama)" class="text-red-600 hover:text-red-800">
                <PhTrash class="ph-bold text-lg"/>
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