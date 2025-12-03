<template>
  <Card class="p-6">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">新規入札者</h3>

    <!-- 空の状態 -->
    <div
      v-if="!bidders || bidders.length === 0"
      class="text-center py-8 text-gray-500"
    >
      <p>新規入札者はいません</p>
    </div>

    <!-- 入札者リスト -->
    <div v-else class="space-y-3">
      <div
        v-for="bidder in bidders"
        :key="bidder.id"
        class="p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
      >
        <div class="flex items-start justify-between">
          <div class="flex-1 min-w-0">
            <!-- 表示名 -->
            <p class="text-sm font-medium text-gray-900 truncate">
              {{ bidder.display_name }}
            </p>
            <!-- メールアドレス -->
            <p class="text-xs text-gray-600 mt-1 truncate">
              {{ bidder.email }}
            </p>
          </div>

          <div class="ml-4 text-right flex-shrink-0">
            <!-- 登録日 -->
            <p
              class="text-xs text-gray-500"
              :title="formatDateTime(bidder.created_at)"
            >
              {{ formatRelativeTime(bidder.created_at) }}
            </p>
          </div>
        </div>

        <!-- ステータスバッジ（オプション） -->
        <div v-if="bidder.status" class="mt-2">
          <span
            :class="[
              'inline-flex items-center px-2 py-1 rounded-full text-xs font-medium',
              bidder.status === 'active'
                ? 'bg-green-100 text-green-800'
                : bidder.status === 'suspended'
                ? 'bg-yellow-100 text-yellow-800'
                : 'bg-gray-100 text-gray-800',
            ]"
          >
            {{ statusLabel(bidder.status) }}
          </span>
        </div>
      </div>
    </div>

    <!-- 全て見るリンク -->
    <div v-if="bidders && bidders.length > 0" class="mt-4 text-center">
      <router-link
        to="/admin/bidders"
        class="text-sm text-blue-600 hover:text-blue-800 font-medium"
      >
        全ての入札者を見る →
      </router-link>
    </div>
  </Card>
</template>

<script setup>
import Card from '@/components/ui/Card.vue'
import { formatRelativeTime, formatDateTime } from '@/utils/timeFormatter'

defineProps({
  /** 新規入札者リスト */
  bidders: {
    type: Array,
    default: () => [],
  },
})

// ステータスラベル
const statusLabel = (status) => {
  const labels = {
    active: 'アクティブ',
    suspended: '一時停止',
    deleted: '削除済み',
  }
  return labels[status] || status
}
</script>
