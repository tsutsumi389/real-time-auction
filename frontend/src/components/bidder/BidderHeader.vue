<template>
  <header class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-20">
    <div class="px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <!-- 左側：ロゴ・タイトル -->
        <div class="flex items-center">
          <div class="flex items-center gap-2">
            <div class="text-2xl">🏇</div>
            <h1 class="text-xl font-bold text-gray-900">リアルタイムオークション</h1>
          </div>
        </div>

        <!-- 右側：ナビゲーションとユーザー情報 -->
        <div class="flex items-center gap-4">
          <!-- ナビゲーションリンク（デスクトップ） -->
          <nav class="hidden md:flex items-center gap-6">
            <router-link
              to="/auctions"
              class="text-sm font-medium text-gray-700 hover:text-blue-600 transition-colors"
              :class="{ 'text-blue-600': isCurrentRoute('/auctions') }"
            >
              オークション一覧
            </router-link>
            <!-- TODO: 入札履歴ページ実装後に有効化
            <router-link
              v-if="bidderAuthStore.isAuthenticated"
              to="/my/bids"
              class="text-sm font-medium text-gray-700 hover:text-blue-600 transition-colors"
              :class="{ 'text-blue-600': isCurrentRoute('/my/bids') }"
            >
              入札履歴
            </router-link>
            -->
            <!-- TODO: ポイントページ実装後に有効化
            <router-link
              v-if="bidderAuthStore.isAuthenticated"
              to="/my/points"
              class="text-sm font-medium text-gray-700 hover:text-blue-600 transition-colors"
              :class="{ 'text-blue-600': isCurrentRoute('/my/points') }"
            >
              ポイント
            </router-link>
            -->
          </nav>

          <!-- ユーザー情報（ログイン済みの場合） -->
          <div v-if="bidderAuthStore.isAuthenticated" class="flex items-center gap-4">
            <!-- ユーザー情報 -->
            <div class="hidden sm:block text-right">
              <p class="text-sm font-medium text-gray-900">{{ bidderAuthStore.user?.displayName || bidderAuthStore.user?.email }}</p>
              <p class="text-xs text-gray-500">入札者</p>
            </div>

            <!-- ユーザーアバター -->
            <div class="flex-shrink-0">
              <div
                class="h-10 w-10 rounded-full bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center text-white font-semibold text-sm"
              >
                {{ userInitial }}
              </div>
            </div>

            <!-- ログアウトボタン -->
            <button
              @click="handleLogout"
              :disabled="loading"
              class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <svg
                v-if="loading"
                class="animate-spin -ml-0.5 mr-2 h-4 w-4 text-gray-700"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
              {{ loading ? 'ログアウト中...' : 'ログアウト' }}
            </button>
          </div>

          <!-- ログイン/新規登録ボタン（未ログインの場合） -->
          <div v-else class="flex items-center gap-3">
            <router-link
              to="/login"
              class="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
            >
              ログイン
            </router-link>
            <router-link
              to="/register"
              class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
            >
              新規登録
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <!-- モバイルナビゲーション -->
    <nav v-if="showMobileNav" class="md:hidden border-t border-gray-200 bg-white">
      <div class="px-4 py-3 space-y-2">
        <router-link
          to="/auctions"
          class="block px-3 py-2 rounded-md text-base font-medium text-gray-700 hover:text-blue-600 hover:bg-gray-50 transition-colors"
          :class="{ 'text-blue-600 bg-blue-50': isCurrentRoute('/auctions') }"
          @click="showMobileNav = false"
        >
          オークション一覧
        </router-link>
        <!-- TODO: 入札履歴ページ実装後に有効化
        <router-link
          v-if="bidderAuthStore.isAuthenticated"
          to="/my/bids"
          class="block px-3 py-2 rounded-md text-base font-medium text-gray-700 hover:text-blue-600 hover:bg-gray-50 transition-colors"
          :class="{ 'text-blue-600 bg-blue-50': isCurrentRoute('/my/bids') }"
          @click="showMobileNav = false"
        >
          入札履歴
        </router-link>
        -->
        <!-- TODO: ポイントページ実装後に有効化
        <router-link
          v-if="bidderAuthStore.isAuthenticated"
          to="/my/points"
          class="block px-3 py-2 rounded-md text-base font-medium text-gray-700 hover:text-blue-600 hover:bg-gray-50 transition-colors"
          :class="{ 'text-blue-600 bg-blue-50': isCurrentRoute('/my/points') }"
          @click="showMobileNav = false"
        >
          ポイント
        </router-link>
        -->
      </div>
    </nav>
  </header>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useBidderAuthStore } from '@/stores/bidderAuthStore'

const router = useRouter()
const route = useRoute()
const bidderAuthStore = useBidderAuthStore()
const loading = ref(false)
const showMobileNav = ref(false)

// 現在のルートチェック
const isCurrentRoute = (path) => {
  return route.path === path || route.path.startsWith(path + '/')
}

// ユーザーイニシャル（アバター用）
const userInitial = computed(() => {
  const displayName = bidderAuthStore.user?.displayName || bidderAuthStore.user?.email || ''
  return displayName.charAt(0).toUpperCase()
})

// ログアウト処理
async function handleLogout() {
  loading.value = true
  try {
    await bidderAuthStore.logout()
    router.push({ name: 'bidder-login' })
  } catch (error) {
    console.error('Logout error:', error)
  } finally {
    loading.value = false
  }
}
</script>
