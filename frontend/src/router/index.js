import { createRouter, createWebHistory } from 'vue-router'
import BahanBakuView from '../views/BahanBakuView.vue'
import ResepView from '../views/ResepsView.vue' // <-- Import ResepView
import HargaJualView from '../views/HargaJualView.vue'; // <-- Import ResepView
import ProgramPromoView from '../views/ProgramPromoView.vue'; // <<< Import ProgramPromoView
import SimulasiView from '../views/SimulasiView.vue'; // <<< Import SimulasiView
import DashboardView from '../views/DashboardView.vue'; // <<< Import SimulasiView
import TenagaKerjaView from '../views/TenagaKerjaView.vue'; // <<< Import SimulasiView
import OperasionalView from '../views/OperasionalView.vue'; // <<< Import SimulasiView


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'dashboard',
      component: DashboardView// Atau bisa juga redirect ke '/reseps'
    },
    {
      path: '/bahan-baku',
      name: 'bahan-baku',
      component: BahanBakuView
    },
    {
      path: '/tenaga-kerja',
      name: 'tenaga-kerja',
      component: TenagaKerjaView
    },
    {
      path: '/operasional',
      name: 'Operasional',
      component: OperasionalView
    },
    {
      path: '/reseps', // <-- Tambahkan route ini
      name: 'reseps',
      component: ResepView
    },
    {
      path: '/harga-jual', // <-- Tambahkan route ini
      name: 'harga-jual',
      component: HargaJualView
    },
    {
      path: '/program-promos', // <<< Route Baru
      name: 'program-promos',
      component: ProgramPromoView
    },
    {
      path: '/simulasi', // <<< Route Baru
      name: 'simulasi',
      component: SimulasiView
    }
  ]
})

export default router