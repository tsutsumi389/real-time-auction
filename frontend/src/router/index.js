import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import HomeView from '../views/HomeView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView,
      meta: { requiresAuth: false }
    },
    {
      path: '/admin/login',
      name: 'admin-login',
      component: () => import('../views/admin/LoginView.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/admin/dashboard',
      name: 'admin-dashboard',
      component: () => import('../views/admin/DashboardView.vue'),
      meta: { requiresAuth: true }
    }
  ]
})

// グローバル認証ガード
router.beforeEach(async (to, from, next) => {
  const authStore = useAuthStore()

  // ユーザー情報が未ロードの場合、トークンから復元を試みる
  if (!authStore.user && !authStore.loading) {
    await authStore.restoreUser()
  }

  // 認証が必要なルート
  if (to.meta.requiresAuth) {
    if (!authStore.isAuthenticated) {
      // 未認証の場合はログイン画面へリダイレクト
      next({
        name: 'admin-login',
        query: { redirect: to.fullPath }
      })
      return
    }
  }

  // すでにログイン済みの場合、ログイン画面へのアクセスをダッシュボードにリダイレクト
  if (to.name === 'admin-login' && authStore.isAuthenticated) {
    next({ name: 'admin-dashboard' })
    return
  }

  next()
})

export default router
