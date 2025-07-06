import { createRouter, createWebHistory } from 'vue-router'
import { Home, Download, Setting, MangaDetail } from '../views'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: Home,
    },
    {
      path: '/download',
      name: 'download',
      component: Download,
    },
    {
      path: '/setting',
      name: 'setting',
      component: Setting,
    },
    {
      path: '/manga/:path',
      name: 'manga',
      component: MangaDetail,
    },
  ],
})

export default router
