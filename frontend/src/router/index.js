import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useBidderAuthStore } from '@/stores/bidderAuthStore'
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
      path: '/login',
      name: 'bidder-login',
      component: () => import('../views/bidder/BidderLoginView.vue'),
      meta: { requiresBidderGuest: true, hideLayout: true }
    },
    {
      path: '/auctions',
      name: 'bidder-auction-list',
      component: () => import('../views/BidderAuctionListView.vue'),
      meta: { requiresBidderAuth: true }
    }
  ]
})

// グローバル認証ガード
router.beforeEach(async (to, _from, next) => {
  const authStore = useAuthStore()
  const bidderAuthStore = useBidderAuthStore()

  // ローディング中は待機（管理者）
  if (authStore.loading) {
    // ローディングが完了するまで待つ（最大5秒）
    let attempts = 0
    while (authStore.loading && attempts < 50) {
      await new Promise(resolve => setTimeout(resolve, 100))
      attempts++
    }
  }

  // ローディング中は待機（入札者）
  if (bidderAuthStore.loading) {
    // ローディングが完了するまで待つ（最大5秒）
    let attempts = 0
    while (bidderAuthStore.loading && attempts < 50) {
      await new Promise(resolve => setTimeout(resolve, 100))
      attempts++
    }
  }

  // 入札者認証が必要なルート
  if (to.meta.requiresBidderAuth) {
    if (!bidderAuthStore.isAuthenticated) {
      // 未認証の場合は入札者ログイン画面へリダイレクト
      next({
        name: 'bidder-login',
        query: { redirect: to.fullPath }
      })
      return
    }
  }

  // 入札者がゲスト専用（ログイン済みならリダイレクト）
  if (to.meta.requiresBidderGuest && bidderAuthStore.isAuthenticated) {
    // すでにログイン済みの場合はオークション一覧へリダイレクト
    next({ name: 'bidder-auction-list' })
    return
  }

  // 管理者認証が必要なルート
  if (to.meta.requiresAuth) {
    if (!authStore.isAuthenticated) {
      // 未認証の場合は管理者ログイン画面へリダイレクト
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

  // すでに管理者としてログイン済みの場合、管理者ログイン画面へのアクセスをダッシュボードにリダイレクト
  if (to.name === 'admin-login' && authStore.isAuthenticated) {
    next({ name: 'admin-dashboard' })
    return
  }

  next()
})

export default router
