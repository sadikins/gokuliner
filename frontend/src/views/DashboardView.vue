<script setup>
import { ref, onMounted, computed } from 'vue';
import axios from 'axios';
import { formatCurrency } from '@/utils/formatters';

const API_BASE_URL = 'http://localhost:8080/api';

const dashboardSummary = ref({
  total_bahan_baku: 0,
  total_resep: 0,
  total_biaya_operasional: 0,
  top_reseps_hpp: [],
});

const fetchDashboardSummary = async () => {
  try {
    const response = await axios.get(`${API_BASE_URL}/dashboard`);
    dashboardSummary.value = {
      ...dashboardSummary.value,
      ...response.data
    };
    console.log('Dashboard data loaded:', dashboardSummary.value);
  } catch (error) {
    console.error('Error fetching dashboard summary:', error);
    alert('Gagal mengambil data ringkasan dashboard.');
  }
};

onMounted(fetchDashboardSummary);



</script>

<template>
  <div class="p-6">
    <h1 class="text-3xl font-bold text-gray-800 mb-6">Selamat Datang!</h1>
    <p class="text-gray-700 mb-8">Ini adalah prototipe fungsional aplikasi vinHPP. Mulai kelola HPP Anda dengan memilih menu di samping.</p>

    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <div class="bg-white p-6 rounded-lg shadow flex items-center space-x-4">
        <div class="bg-blue-100 p-3 rounded-full text-blue-600">
          <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 20 20"><path d="M4 3a2 2 0 100 4h12a2 2 0 100-4H4zm0 8a2 2 0 100 4h12a2 2 0 100-4H4z"></path></svg>
        </div>
        <div>
          <p class="text-gray-500 text-sm">Total Bahan Baku</p>
          <p class="text-2xl font-bold text-gray-800">{{ dashboardSummary.total_bahan_baku }}</p>
        </div>
      </div>

      <div class="bg-white p-6 rounded-lg shadow flex items-center space-x-4">
        <div class="bg-green-100 p-3 rounded-full text-green-600">
          <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M4 5a2 2 0 012-2h8a2 2 0 012 2v10a2 2 0 01-2 2H6a2 2 0 01-2-2V5zm3 1a1 1 0 000 2h.01a1 1 0 000-2H7zm3 0a1 1 0 000 2h.01a1 1 0 000-2H10zm3 0a1 1 0 000 2h.01a1 1 0 000-2H13zm-3 4a1 1 0 000 2h.01a1 1 0 000-2H10zm3 0a1 1 0 000 2h.01a1 1 0 000-2H13z" clip-rule="evenodd"></path></svg>
        </div>
        <div>
          <p class="text-gray-500 text-sm">Total Resep</p>
          <p class="text-2xl font-bold text-gray-800">{{ dashboardSummary.total_resep }}</p>
        </div>
      </div>

      <div class="bg-white p-6 rounded-lg shadow flex items-center space-x-4">
        <div class="bg-yellow-100 p-3 rounded-full text-yellow-600">
          <svg class="h-6 w-6" fill="currentColor" viewBox="0 0 20 20"><path d="M10 2a1 1 0 00-1 1v1a1 1 0 002 0V3a1 1 0 00-1-1zm4 8a4 4 0 11-8 0 4 4 0 018 0zm-.459 4.293a1 1 0 00-.707.293l-1.414 1.414a1 1 0 01-1.414 0l-1.414-1.414a1 1 0 00-.707-.293H7.5a1 1 0 00-.707.293l-1.414 1.414a1 1 0 01-1.414 0L2.293 14.541A1 1 0 001 15v1a1 1 0 001 1h.459a1 1 0 00.707.293l1.414 1.414a1 1 0 011.414 0l1.414-1.414a1 1 0 00.707-.293H12.5a1 1 0 00.707-.293l1.414-1.414a1 1 0 011.414 0l1.414 1.414A1 1 0 0019 16v-1a1 1 0 00-1-1h-.459a1 1 0 00-.707-.293zM10 8a2 2 0 100 4 2 2 0 000-4z"></path></svg>
        </div>
        <div>
          <p class="text-gray-500 text-sm">Total Biaya Operasional</p>
          <!-- <p class="text-2xl font-bold text-gray-800">{{ formatCurrency(dashboardSummary.value.total_biaya_operasional) }}</p> -->
        </div>
      </div>
    </div>

    <div class="bg-white p-6 rounded-lg shadow">
      <h2 class="text-xl font-semibold text-gray-800 mb-4">Top 5 Resep dengan HPP per Porsi Tertinggi</h2>
      <div v-if="dashboardSummary.top_reseps_hpp.length > 0" class="overflow-x-auto">
        <table class="min-w-full bg-white border border-gray-200">
          <thead class="bg-gray-100">
            <tr>
              <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">Nama Resep</th>
              <th class="py-3 px-4 border-b text-left text-sm font-semibold text-gray-600">HPP per Porsi</th>
              </tr>
          </thead>
          <tbody>
            <tr v-for="(resep, index) in dashboardSummary.top_reseps_hpp" :key="resep.resep_id" :class="index % 2 === 0 ? 'bg-white' : 'bg-gray-50'">
              <td class="py-3 px-4 border-b text-sm text-gray-800">{{ resep.resep_nama }}</td>
              <td class="py-3 px-4 border-b text-sm text-gray-800">{{ formatCurrency(resep.hpp_per_porsi) }}</td>
              </tr>
          </tbody>
        </table>
      </div>
      <div v-else class="text-gray-500 text-center py-4">
        Tidak ada data resep HPP tertinggi. Hitung HPP di modul Resep & HPP.
      </div>
    </div>
  </div>
</template>

<style scoped>
/* Hapus semua gaya scoped */
</style>