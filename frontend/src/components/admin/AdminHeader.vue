<template>
  <header class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-30">
    <div class="px-4 sm:px-6 lg:px-8">
      <div class="flex items-center justify-between h-16">
        <!-- 左側：ハンバーガーメニューとロゴ -->
        <div class="flex items-center">
          <!-- ハンバーガーメニュー（モバイル） -->
          <button
            @click="$emit('toggle-sidebar')"
            class="lg:hidden p-2 rounded-md text-gray-400 hover:text-gray-500 hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-inset focus:ring-blue-500"
            aria-label="メニューを開く"
          >
            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>

          <!-- ロゴ・タイトル -->
          <div class="flex items-center ml-4 lg:ml-0">
            <h1 class="text-xl font-bold text-gray-900">オークション管理</h1>
          </div>
        </div>

        <!-- 右側：ユーザー情報とログアウト -->
        <div class="flex items-center gap-4">
          <!-- ユーザー情報 -->
          <div class="hidden sm:block text-right">
            <p class="text-sm font-medium text-gray-900">{{ authStore.user?.email }}</p>
            <p class="text-xs text-gray-500">{{ roleLabel }}</p>
          </div>

          <!-- ユーザーアバター -->
          <div class="flex-shrink-0">
            <div
              class="h-10 w-10 rounded-full bg-gradient-to-br from-blue-400 to-blue-600 flex items-center justify-center text-white font-semibold text-sm"
            >
              {{ userInitial }}
            </div>
          </div>

          <!-- ログアウトボタン -->
          <button
            @click="handleLogout"
            :disabled="loading"
            class="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm leading-4 font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
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
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)

defineEmits(['toggle-sidebar'])

// ロールラベル
const roleLabel = computed(() => {
  switch (authStore.user?.role) {
    case 'system_admin':
      return 'システム管理者'
    case 'auctioneer':
      return 'オークショニア'
    default:
      return '不明'
  }
})

// ユーザーイニシャル（アバター用）
const userInitial = computed(() => {
  const email = authStore.user?.email || ''
  return email.charAt(0).toUpperCase()
})

// ログアウト処理
async function handleLogout() {
  loading.value = true
  try {
    await authStore.logout()
    router.push({ name: 'admin-login' })
  } catch (error) {
    console.error('Logout error:', error)
  } finally {
    loading.value = false
  }
}
</script>
