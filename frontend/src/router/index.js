import { createRouter, createWebHistory } from 'vue-router'
import BahanBakuView from '../views/BahanBakuView.vue'
import ResepView from '../views/ResepView.vue' // <-- Import ResepView

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      redirect: '/bahan-baku' // Atau bisa juga redirect ke '/reseps'
    },
    {
      path: '/bahan-baku',
      name: 'bahan-baku',
      component: BahanBakuView
    },
    {
      path: '/reseps', // <-- Tambahkan route ini
      name: 'reseps',
      component: ResepView
    }
  ]
})

export default router