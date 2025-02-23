import { createRouter, createWebHistory, RouteRecordRaw } from 'vue-router'
import { useUserStore } from '@/stores/user'

// 创建路由实例
const routes: Array<RouteRecordRaw> = [
  {
    path: '/',
    name: 'Login',
    component: () => import('@/views/login/index.vue')
  },
  {
    path: '/main',
    name: 'MainPage',
    component: () => import('@/views/mainpage/index.vue')
  },
  {
    path: '/:catchAll(.*)',
    redirect: '/main'
  }
]

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes
})

router.beforeEach((to, from, next) => {
  const userStore = useUserStore()
  if (to.path === '/') {
    next()
  } else {
    if (userStore.Token) {
      next()
    } else {
      next('/')
    }
  }
})

export default router
