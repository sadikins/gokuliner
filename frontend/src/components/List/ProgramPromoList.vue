<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  programPromos: Array
})

const emit = defineEmits(['edit', 'delete'])

const formatRupiah = (value) => {
  return parseFloat(value).toLocaleString('id-ID', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
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
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Nama Promo</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Channel</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Jenis Diskon</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Besar</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Min Belanja</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Maks Potongan</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Ditanggung Merchant (%)</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Catatan</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Aksi</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="programPromos.length === 0">
          <td colspan="9" class="py-4 px-4 text-center text-gray-500">Tidak ada program promo. Tambahkan yang baru.</td>
        </tr>
        <tr v-for="(promo, index) in programPromos" :key="promo.id" :class="index % 2 === 0 ? 'bg-white' : 'bg-gray-50'">
          <td class="py-3 px-4 border-b text-sm text-gray-800">{{ promo.nama_promo }}</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800">{{ promo.channel }}</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800">{{ promo.jenis_diskon }}</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800">
            {{ promo.jenis_diskon === 'persentase' ? promo.besar_diskon + '%' : 'Rp ' + formatRupiah(promo.besar_diskon) }}
          </td>
          <td class="py-3 px-4 border-b text-sm text-gray-800">Rp {{ formatRupiah(promo.min_belanja) }}</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800">Rp {{ formatRupiah(promo.maksimal_potongan) }}</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800">{{ promo.ditanggung_merchant_persen }}%</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800">{{ promo.catatan }}</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800 whitespace-nowrap">
            <button @click="emit('edit', promo)"
                    class="bg-yellow-500 hover:bg-yellow-600 text-white font-bold py-1 px-3 rounded-md text-xs mr-2">
              Edit
            </button>
            <button @click="emit('delete', promo.id, promo.nama_promo)"
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
/* Hapus semua gaya scoped */
</style>