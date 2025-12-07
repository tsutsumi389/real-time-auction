<template>
  <header class="header-glass border-b border-lux-gold/20 sticky top-0 z-20">
    <div class="px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16 sm:h-20">
        <!-- 左側:ロゴ・タイトル -->
        <div class="flex items-center">
          <div class="flex items-center gap-3">
            <div class="hidden sm:flex w-12 h-12 rounded-xl bg-gradient-to-br from-lux-gold/20 to-lux-gold/5 border border-lux-gold/40 flex-shrink-0 items-center justify-center shadow-lg shadow-lux-gold/10">
              <svg class="w-6 h-6 text-lux-gold" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
              </svg>
            </div>
            <div>
              <h1 class="font-display text-lg sm:text-xl text-lux-cream font-medium tracking-wide">リアルタイムオークション</h1>
            </div>
          </div>
        </div>

        <!-- 右側:ナビゲーション、ユーザーメニュー -->
        <div class="flex items-center gap-4">
          <!-- ナビゲーションリンク(デスクトップ) -->
          <nav class="hidden md:flex items-center gap-6">
            <router-link
              to="/auctions"
              class="text-sm font-medium text-lux-silver hover:text-lux-gold transition-colors"
              :class="{ 'text-lux-gold': isCurrentRoute('/auctions') }"
            >
              オークション一覧
            </router-link>
          </nav>

          <!-- ログイン済みの場合:ユーザーメニュー -->
          <div v-if="bidderAuthStore.isAuthenticated" class="flex items-center gap-3">
            <!-- ユーザーアイコン(ドロップダウントリガー) -->
            <div class="relative">
              <button
                @click="toggleMenu"
                @keydown.enter="toggleMenu"
                @keydown.space.prevent="toggleMenu"
                @keydown.escape="closeMenu"
                :aria-expanded="isMenuOpen"
                aria-haspopup="true"
                aria-label="ユーザーメニューを開く"
                class="h-10 w-10 rounded-full bg-gradient-to-br from-lux-gold via-amber-400 to-lux-gold-dark flex items-center justify-center text-lux-noir font-bold text-base shadow-lg shadow-lux-gold/50 border-2 border-lux-gold/30 transition-all duration-200 ease-out hover:scale-110 hover:shadow-xl hover:shadow-lux-gold/70 hover:border-lux-gold/50 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-lux-gold focus-visible:ring-offset-2 focus-visible:ring-offset-lux-noir"
                :class="{ 'ring-2 ring-lux-gold scale-110 shadow-xl shadow-lux-gold/70': isMenuOpen }"
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
                  class="absolute right-0 top-full mt-2 w-60 lux-glass-strong rounded-xl border border-lux-gold/20 p-2 z-50 shadow-2xl"
                >
                  <!-- ユーザー情報セクション -->
                  <div class="px-3 py-2 flex items-start gap-3 mb-1">
                    <!-- アバター -->
                    <div class="flex-shrink-0">
                      <div
                        class="h-10 w-10 rounded-full bg-gradient-to-br from-lux-gold via-amber-400 to-lux-gold-dark flex items-center justify-center text-lux-noir font-bold text-base shadow-md shadow-lux-gold/40 border-2 border-lux-gold/20"
                      >
                        {{ userInitial }}
                      </div>
                    </div>
                    <!-- 名前とメール -->
                    <div class="flex-1 min-w-0">
                      <p class="text-sm font-semibold text-lux-cream truncate">
                        {{ displayName }}
                      </p>
                      <p class="text-xs text-lux-silver truncate">
                        {{ bidderAuthStore.user?.email || '' }}
                      </p>
                    </div>
                  </div>

                  <!-- セパレーター -->
                  <div class="border-t border-lux-gold/20 my-1"></div>

                  <!-- ログアウト -->
                  <button
                    @click="handleLogout"
                    :disabled="loading"
                    role="menuitem"
                    aria-label="ログアウト"
                    class="w-full flex items-center gap-2 px-3 py-2 text-sm text-red-400 rounded-lg hover:bg-red-950/30 hover:text-red-300 transition-colors duration-150 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <LogOut :size="18" :stroke-width="2" />
                    <span>{{ loading ? 'ログアウト中...' : 'ログアウト' }}</span>
                  </button>
                </div>
              </Transition>
            </div>
          </div>

          <!-- 未ログインの場合:ログイン/新規登録ボタン -->
          <div v-else class="flex items-center gap-3">
            <router-link
              to="/login"
              class="inline-flex items-center px-4 py-2 border border-lux-gold/30 shadow-sm text-sm font-medium rounded-lg text-lux-silver bg-lux-noir-light/50 hover:bg-lux-noir-light/80 hover:text-lux-gold transition-all duration-200"
            >
              ログイン
            </router-link>
            <router-link
              to="/register"
              class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-lg text-lux-noir bg-lux-gold hover:bg-lux-gold-dark transition-all duration-200"
            >
              新規登録
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <!-- モバイルナビゲーション -->
    <nav v-if="showMobileNav" class="md:hidden border-t border-lux-gold/20 lux-glass">
      <div class="px-4 py-3 space-y-2">
        <router-link
          to="/auctions"
          class="block px-3 py-2 rounded-lg text-base font-medium text-lux-silver hover:text-lux-gold hover:bg-lux-noir-light/50 transition-colors"
          :class="{ 'text-lux-gold bg-lux-gold/10': isCurrentRoute('/auctions') }"
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

// ユーザーイニシャル(アバター用)
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
/* Header Glass Effect */
.header-glass {
  background: linear-gradient(180deg, rgba(10, 10, 10, 0.95) 0%, rgba(10, 10, 10, 0.9) 100%);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  box-shadow:
    0 4px 30px rgba(0, 0, 0, 0.3),
    inset 0 -1px 0 rgba(212, 175, 55, 0.1);
}

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
