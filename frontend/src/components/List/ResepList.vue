<script setup>
import { defineProps, defineEmits } from 'vue'
    import axios from 'axios'
import {formatCurrency} from '@/utils/formatters'

const props = defineProps({
  reseps: Array
})

const emit = defineEmits(['hppCalculated', 'edit', 'delete', 'showDetail', 'duplicate'])

const API_BASE_URL = 'http://localhost:8080/api'

const getKomponenSummary = (komponen) => {
  if (!komponen || komponen.length === 0) return 'Tidak ada komponen';

  const limited = komponen.slice(0, 3);
  const names = limited.map(c => {
    return `${c.kuantitas}x ${c.tipe_komponen === 'bahan_baku' ? 'B.Baku' : 'Resep'} ID:${c.komponen_id.substring(0, 4)}...`;
  });

  if (komponen.length > 3) {
    names.push(`...dan ${komponen.length - 3} lainnya`);
  }
  return names.join(', ');
}

const calculateHPP = async (resepId, resepNama) => {
  try {
    const response = await axios.get(`${API_BASE_URL}/hpp/${resepId}`)
    const hppResult = response.data

    let alertMessage = `HPP untuk "${hppResult.resep_nama}":\n`;
    alertMessage += `Per Unit Resep: Rp ${formatCurrency(hppResult.hpp_per_unit) })}\n`;
    alertMessage += `Per Porsi (${props.reseps.find(r => r.id === resepId)?.jumlah_porsi || '1'} porsi): Rp ${formatCurrency(hppResult.hpp_per_porsi) })}`;

    alert(alertMessage)

    emit('hppCalculated', hppResult)

  } catch (error) {
    console.error('Error calculating HPP:', error.response ? error.response.data : error)
    const errorMessage = error.response && error.response.data && error.response.data.error
                         ? error.response.data.error
                         : 'Gagal menghitung HPP. Periksa konfigurasi resep.';
    alert(errorMessage);
  }
}

const editResep = (resep) => {
  emit('edit', resep)
}

const deleteResep = async (id, nama) => {
  if (confirm(`Apakah Anda yakin ingin menghapus resep "${nama}"? Semua komponen terkait juga akan dihapus.`)) {
    emit('delete', id, nama)
  }
    }

const showResepDetail = (resepId) => {
  emit('showDetail', resepId) // Emit event ke parent
}
const duplicateResep = (resepId, resepNama) => {
  emit('duplicate', resepId, resepNama) // Emit event ke parent
}
</script>

<template>
  <div class="overflow-x-auto">
    <table class="min-w-full bg-white border border-gray-200">
      <thead class="bg-gray-100">
        <tr>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Nama Resep</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Sub-Resep?</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Jumlah Porsi</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Komponen</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Tanggal Dibuat</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Aksi</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="reseps.length === 0">
          <td colspan="6" class="py-4 px-4 text-center text-gray-500">Tidak ada resep. Tambahkan yang baru.</td>
        </tr>
        <tr v-for="(resep, index) in reseps" :key="resep.id" :class="index % 2 === 0 ? 'bg-white' : 'bg-gray-50'">
          <td class="py-3 px-4 border-b text-sm text-gray-800">{{ resep.nama }}</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800">{{ resep.is_sub_resep ? 'Ya' : 'Tidak' }}</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800">{{ resep.jumlah_porsi }}</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800">{{ getKomponenSummary(resep.komponen) }}</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800">{{ new Date(resep.created_at).toLocaleDateString('id-ID', { year: 'numeric', month: 'long', day: 'numeric' }) }}</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800 space-x-2 whitespace-nowrap">
            <button @click="calculateHPP(resep.id, resep.nama)"
                    class="bg-gray-600 hover:bg-gray-700 text-white font-bold py-1 px-3 rounded-md text-xs">
              Hitung HPP
            </button>
            <button @click="showResepDetail(resep.id)"
                    class="bg-blue-600 hover:bg-blue-700 text-white font-bold py-1 px-3 rounded-md text-xs">
              Detail
            </button>
            <button @click="duplicateResep(resep.id, resep.nama)"
                    class="bg-purple-600 hover:bg-purple-700 text-white font-bold py-1 px-3 rounded-md text-xs">
              Duplikat
            </button>
            <button @click="editResep(resep)"
                    class="bg-yellow-500 hover:bg-yellow-600 text-white font-bold py-1 px-3 rounded-md text-xs">
              Edit
            </button>
            <button @click="deleteResep(resep.id, resep.nama)"
                    class="bg-red-500 hover:bg-red-600 text-white font-bold py-1 px-3 rounded-md text-xs">
              Hapus
            </button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
