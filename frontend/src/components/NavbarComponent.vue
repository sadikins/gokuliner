<script setup>
  import { PhCardholder, PhChartPieSlice, PhCoins, PhDrone, PhHouse, PhReceipt, PhTag, PhTrayArrowDown, PhUsers } from '@phosphor-icons/vue'
  import { computed, ref } from 'vue';

  const active = ref(false);


  const isNavlinkActive = computed( () => active.value ? 'bg-gray-50' : '' )


  console.log(active.value)

 </script>

<template>
  <aside class="w-64 bg-white flex flex-col shadow-lg">
      <div class="px-6 py-4 border-b border-gray-200 mb-2">

          <h1 class="text-xl font-bold flex ">
            <PhDrone class="h-8 w-8  font-medium text-lime-400"/> vinHPP
          </h1>

          <p class="text-sm text-slate-400">Kalkulator HPP Cerdas</p>
      </div>

      <nav class="flex-1 px-4 py-4 space-y-5">
        <router-link to="/" class="nav-link flex  items-center px-4 hover:bg-gray-100 py-2 text-gray-700 rounded-md font-medium " activeClass="bg-lime-200">
          <PhHouse :size="24" class="mr-3"/>
          Dashboard
        </router-link>
        <router-link to="/bahan-baku" class="nav-link flex items-center px-4 hover:bg-gray-100 py-2 text-gray-700 rounded-md font-medium "  activeClass="bg-lime-200">
          <PhCardholder :size="24" class="ph-fill  mr-3"/>
          Bahan Baku
        </router-link>
        <router-link to="/tenaga-kerja" class="nav-link flex items-center px-4 hover:bg-gray-100 py-2 text-gray-700 rounded-md font-medium "  activeClass="bg-lime-200">
          <PhUsers :size="24" class="ph-fill  mr-3"/>
          Tenaga Kerja
        </router-link>
        <router-link to="/operasional" class="nav-link flex items-center px-4 hover:bg-gray-100 py-2 text-gray-700 rounded-md font-medium "  activeClass="bg-lime-200">
          <PhTrayArrowDown :size="24" class="ph-fill  mr-3"/>
          Operasional
        </router-link>
        <router-link to="/reseps" class="nav-link flex items-center px-4 hover:bg-gray-100 py-2 text-gray-700 rounded-md font-medium "  activeClass="bg-lime-200">

          <PhReceipt :size="24" class="mr-3" />
          Resep & HPP
        </router-link>

        <router-link to="/harga-jual" class="nav-link flex items-center px-4 hover:bg-gray-100 py-2 text-gray-700 rounded-md font-medium "  activeClass="bg-lime-200">
          <PhCoins :size="24" class="ph-fill mr-3" />
          Harga Jual
        </router-link>
        <router-link to="/program-promos" class="nav-link flex items-center px-4 hover:bg-gray-100 py-2 text-gray-700 rounded-md font-medium "  activeClass="bg-lime-200">
          <PhTag :size="24" class="mr-3"/>
          Promo
        </router-link>
        <router-link to="/simulasi" class="nav-link flex items-center px-4 hover:bg-gray-100 py-2 text-gray-700 rounded-md font-medium "  activeClass="bg-lime-200">
          <PhChartPieSlice :size="24" class="mr-3"/>
          Simulasi
        </router-link>
      </nav>
      <div class="p-4 text-xs text-gray-400 border-t border-gray-100">
        © 2024 VinHPP. Dibuat sebagai prototype.
      </div>
    </aside>
</template>

