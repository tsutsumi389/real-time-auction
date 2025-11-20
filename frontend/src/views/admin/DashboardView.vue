<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-white shadow">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex items-center justify-between">
        <h1 class="text-2xl font-bold text-gray-900">管理者ダッシュボード</h1>
        <div class="flex items-center space-x-4">
          <span class="text-sm text-gray-600">
            {{ authStore.user?.email }}
            <span class="text-xs text-gray-400 ml-1">({{ roleLabel }})</span>
          </span>
          <Button variant="outline" size="sm" @click="handleLogout" :disabled="loading">
            {{ loading ? 'ログアウト中...' : 'ログアウト' }}
          </Button>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <Card class="p-6">
        <h2 class="text-xl font-semibold mb-4">ようこそ</h2>
        <p class="text-gray-600">
          管理者ダッシュボードへようこそ。このページは今後実装予定です。
        </p>
        <div class="mt-6 space-y-2 text-sm text-gray-500">
          <p><strong>ユーザーID:</strong> {{ authStore.user?.adminId }}</p>
          <p><strong>メールアドレス:</strong> {{ authStore.user?.email }}</p>
          <p><strong>権限:</strong> {{ roleLabel }}</p>
          <p><strong>認証状態:</strong> {{ authStore.isAuthenticated ? '認証済み' : '未認証' }}</p>
        </div>
      </Card>
    </main>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Card from '@/components/ui/Card.vue'
import Button from '@/components/ui/Button.vue'

const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)

// Role label for display
const roleLabel = computed(() => {
  switch (authStore.user?.role) {
    case 'system_admin':
      return 'システム管理者'
    case 'auctioneer':
      return 'オークション主催者'
    default:
      return '不明'
  }
})

/**
 * Handle logout
 */
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
