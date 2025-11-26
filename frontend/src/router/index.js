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
      meta: { requiresAuth: false, hideLayout: true }
    },
    {
      path: '/admin/dashboard',
      name: 'admin-dashboard',
      component: () => import('../views/admin/DashboardView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/admin/admins',
      name: 'admin-list',
      component: () => import('../views/admin/AdminListView.vue'),
      meta: { requiresAuth: true, requireSystemAdmin: true }
    },
    {
      path: '/admin/admins/:id/edit',
      name: 'admin-edit',
      component: () => import('../views/admin/AdminEditView.vue'),
      meta: { requiresAuth: true, requireSystemAdmin: true }
    },
    {
      path: '/admin/admins/new',
      name: 'admin-register',
      component: () => import('../views/admin/AdminRegisterView.vue'),
      meta: { requiresAuth: true, requireSystemAdmin: true }
    },
    {
      path: '/admin/bidders',
      name: 'bidder-list',
      component: () => import('../views/admin/BidderListView.vue'),
      meta: { requiresAuth: true, requireSystemAdmin: true }
    },
    {
      path: '/admin/bidders/new',
      name: 'bidder-register',
      component: () => import('../views/admin/BidderRegisterView.vue'),
      meta: { requiresAuth: true, requireSystemAdmin: true }
    },
    {
      path: '/admin/auctions',
      name: 'auction-list',
      component: () => import('../views/admin/AuctionListView.vue'),
      meta: { requiresAuth: true, requireAdminOrAuctioneer: true }
    },
    {
      path: '/admin/auctions/new',
      name: 'auction-new',
      component: () => import('../views/admin/AuctionNewView.vue'),
      meta: { requiresAuth: true, requireAdminOrAuctioneer: true }
    },
    {
      path: '/admin/auctions/:id/live',
      name: 'auction-live',
      component: () => import('../views/admin/AuctionLive.vue'),
      meta: { requiresAuth: true, requireAdminOrAuctioneer: true }
    },
    {
      path: '/auctions',
      name: 'bidder-auction-list',
      component: () => import('../views/BidderAuctionListView.vue'),
      meta: { requiresAuth: false }
    }
  ]
})

// グローバル認証ガード
router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()

  // ローディング中は待機
  if (authStore.loading) {
    // ローディングが完了するまで待つ（最大5秒）
    let attempts = 0
    while (authStore.loading && attempts < 50) {
      await new Promise(resolve => setTimeout(resolve, 100))
      attempts++
    }
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

    // システム管理者権限が必要なルート
    if (to.meta.requireSystemAdmin && !authStore.isSystemAdmin) {
      // 権限不足の場合はダッシュボードへリダイレクト
      next({ name: 'admin-dashboard' })
      return
    }

    // システム管理者またはオークショニア権限が必要なルート
    if (to.meta.requireAdminOrAuctioneer && !authStore.isSystemAdmin && !authStore.isAuctioneer) {
      // 権限不足の場合はダッシュボードへリダイレクト
      next({ name: 'admin-dashboard' })
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
