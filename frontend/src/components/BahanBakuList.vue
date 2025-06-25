<script setup>
import { defineProps, defineEmits } from 'vue'

const props = defineProps({
  bahanBakus: Array
})

// Mendefinisikan event yang akan di-emit ke parent component
const emit = defineEmits(['edit', 'delete']) // Event 'edit' dan 'delete' baru

const formatRupiah = (value) => {
  return value.toLocaleString('id-ID', { minimumFractionDigits: 2, maximumFractionDigits: 2 });
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
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Satuan</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Harga Beli</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Tanggal Dibuat</th>
          <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Aksi</th>
        </tr>
      </thead>
      <tbody>
        <tr v-if="bahanBakus.length === 0">
          <td colspan="5" class="py-4 px-4 text-center text-gray-500">Tidak ada bahan baku. Tambahkan yang baru di atas.</td>
        </tr>
        <tr v-for="(bb, index) in bahanBakus" :key="bb.id" :class="index % 2 === 0 ? 'bg-white' : 'bg-gray-50'">
          <td class="py-3 px-4 border-b text-sm text-gray-800">{{ bb.nama }}</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800">{{ bb.satuan }}</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800">Rp {{ formatRupiah(bb.harga_beli) }}</td>
          <td class="py-3 px-4 border-b text-sm text-gray-800">{{ formatDate(bb.created_at) }}</td>
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
/* Gaya yang sudah ada */
.bahan-baku-table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
}

.bahan-baku-table th,
.bahan-baku-table td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}

.bahan-baku-table th {
  background-color: #e2e2e2;
  color: #333;
}

.bahan-baku-table tbody tr:nth-child(even) {
  background-color: #f2f2f2;
}

.bahan-baku-table tbody tr:hover {
  background-color: #ddd;
}

/* Gaya baru untuk tombol Aksi */
.action-btn {
  padding: 5px 10px;
  border: none;
  border-radius: 4px;
  cursor: pointer;
  font-size: 0.85em;
  margin-right: 5px;
}

.edit-btn {
  background-color: #ffc107; /* Kuning */
  color: #333;
}

.edit-btn:hover {
  background-color: #e0a800;
}

.delete-btn {
  background-color: #dc3545; /* Merah */
  color: white;
}

.delete-btn:hover {
  background-color: #c82333;
}
</style>