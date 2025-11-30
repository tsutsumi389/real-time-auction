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

        <!-- 右側：ナビゲーション、ポイント、ユーザーメニュー -->
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
          </nav>

          <!-- ログイン済みの場合：ユーザーメニュー -->
          <div v-if="bidderAuthStore.isAuthenticated" class="flex items-center gap-3">
            <!-- ユーザーアイコン（ドロップダウントリガー） -->
            <div class="relative">
              <button
                @click="toggleMenu"
                @keydown.enter="toggleMenu"
                @keydown.space.prevent="toggleMenu"
                @keydown.escape="closeMenu"
                :aria-expanded="isMenuOpen"
                aria-haspopup="true"
                aria-label="ユーザーメニューを開く"
                class="h-10 w-10 rounded-full bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center text-white font-semibold text-sm transition-transform duration-200 ease-out hover:scale-105 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-green-600 focus-visible:ring-offset-2"
                :class="{ 'ring-2 ring-green-500 scale-105': isMenuOpen }"
              >
                {{ userInitial }}
              </button>

              <!-- ドロップダウンメニュー -->
              <Transition
                enter-active-class="transition duration-150 ease-out"
                enter-from-class="opacity-0 -translate-y-1 scale-98"
                enter-to-class="opacity-100 translate-y-0 scale-100"
                leave-active-class="transition duration-100 ease-in"
                leave-from-class="opacity-100 translate-y-0 scale-100"
                leave-to-class="opacity-0 -translate-y-1 scale-98"
              >
                <div
                  v-if="isMenuOpen"
                  ref="menuRef"
                  role="menu"
                  aria-orientation="vertical"
                  class="absolute right-0 top-full mt-2 w-60 bg-white rounded-lg shadow-lg border border-gray-200 p-2 z-50"
                >
                  <!-- ユーザー情報セクション -->
                  <div class="px-3 py-2 flex items-start gap-3 mb-1">
                    <!-- アバター -->
                    <div class="flex-shrink-0">
                      <div
                        class="h-10 w-10 rounded-full bg-gradient-to-br from-green-400 to-green-600 flex items-center justify-center text-white font-semibold text-sm"
                      >
                        {{ userInitial }}
                      </div>
                    </div>
                    <!-- 名前とメール -->
                    <div class="flex-1 min-w-0">
                      <p class="text-sm font-semibold text-gray-900 truncate">
                        {{ displayName }}
                      </p>
                      <p class="text-xs text-gray-500 truncate">
                        {{ bidderAuthStore.user?.email || '' }}
                      </p>
                    </div>
                  </div>

                  <!-- セパレーター -->
                  <div class="border-t border-gray-200 my-1"></div>

                  <!-- ログアウト -->
                  <button
                    @click="handleLogout"
                    :disabled="loading"
                    role="menuitem"
                    aria-label="ログアウト"
                    class="w-full flex items-center gap-2 px-3 py-2 text-sm text-red-600 rounded-md hover:bg-red-50 hover:text-red-700 transition-colors duration-150 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <LogOut :size="18" :stroke-width="2" />
                    <span>{{ loading ? 'ログアウト中...' : 'ログアウト' }}</span>
                  </button>
                </div>
              </Transition>
            </div>
          </div>

          <!-- 未ログインの場合：ログイン/新規登録ボタン -->
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
      </div>
    </nav>
  </header>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { onClickOutside } from '@vueuse/core'
import { LogOut } from 'lucide-vue-next'
import { useBidderAuthStore } from '@/stores/bidderAuthStore'

const router = useRouter()
const route = useRoute()
const bidderAuthStore = useBidderAuthStore()
const loading = ref(false)
const showMobileNav = ref(false)
const isMenuOpen = ref(false)
const menuRef = ref(null)

// 現在のルートチェック
const isCurrentRoute = (path) => {
  return route.path === path || route.path.startsWith(path + '/')
}

// ユーザーイニシャル（アバター用）
const userInitial = computed(() => {
  const displayName = bidderAuthStore.user?.displayName || bidderAuthStore.user?.email || ''
  return displayName.charAt(0).toUpperCase()
})

// 表示名
const displayName = computed(() => {
  return bidderAuthStore.user?.displayName || bidderAuthStore.user?.email || 'ゲスト'
})

// メニューの開閉
function toggleMenu() {
  isMenuOpen.value = !isMenuOpen.value
}

function closeMenu() {
  isMenuOpen.value = false
}

// メニュー外クリックで閉じる
onClickOutside(menuRef, () => {
  if (isMenuOpen.value) {
    closeMenu()
  }
})

// ログアウト処理
async function handleLogout() {
  loading.value = true
  try {
    await bidderAuthStore.logout()
    closeMenu()
    router.push({ name: 'bidder-login' })
  } catch (error) {
    console.error('Logout error:', error)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
/* モーション低減対応 */
@media (prefers-reduced-motion: reduce) {
  .transition,
  .transition-transform,
  .transition-colors {
    transition: opacity 100ms ease !important;
    transform: none !important;
  }
}
</style>
