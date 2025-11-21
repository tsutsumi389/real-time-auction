<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- ヘッダー -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900">ダッシュボード</h1>
      <p class="mt-2 text-sm text-gray-600">管理者ダッシュボード</p>
    </div>

    <!-- メインコンテンツ -->
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
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAuthStore } from '@/stores/auth'
import Card from '@/components/ui/Card.vue'

const authStore = useAuthStore()

// Role label for display
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
</script>
