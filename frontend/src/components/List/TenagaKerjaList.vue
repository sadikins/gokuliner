<script setup>
import { PhPencil, PhTrash } from '@phosphor-icons/vue';
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  tenagaKerjas: Array
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
          <th class="px-6 py-3">Nama Pekerjaan</th>
           <th class="px-6 py-3">Pembayaran</th>
          <th class="px-6 py-3">Lama Kegiatan</th>
          <th class="px-6 py-3">Tanggal Dibuat</th>
          <th class="px-6 py-3">Aksi</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="tenagaKerjas.length === 0">
            <td colspan="9" class="px-6 py-24 font-medium text-center text-gray-400">
            Tidak ada data tenaga kerja. Tambahkan data tenaga kerja.
          </td>
        </tr>
        <tr v-for="(bb, index) in tenagaKerjas" :key="bb.id" class="border-b border-slate-200" >
          <td class="px-6 py-4 font-medium text-gray-900 whitespace-nowrap">{{ bb.nama }}</td>
          <td class="px-6 py-4">Rp {{ formatRupiah(bb.harga_beli) }} / {{ bb.satuan_beli }}</td>
          <td class="px-6 py-4 font-medium">{{ bb.netto_per_beli }} / {{ bb.satuan_pemakaian }}</td>
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

