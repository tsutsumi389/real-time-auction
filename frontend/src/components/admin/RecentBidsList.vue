<template>
  <Card class="p-6">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">最新の入札</h3>

    <!-- 空の状態 -->
    <div
      v-if="!bids || bids.length === 0"
      class="text-center py-8 text-gray-500"
    >
      <p>まだ入札がありません</p>
    </div>

    <!-- 入札リスト -->
    <div v-else class="space-y-3">
      <div
        v-for="bid in bids"
        :key="bid.id"
        class="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
      >
        <div class="flex-1 min-w-0">
          <!-- オークション情報 -->
          <p class="text-sm font-medium text-gray-900 truncate">
            {{ bid.item_title || 'オークション' }}
          </p>
          <!-- 入札者情報 -->
          <p class="text-xs text-gray-600 mt-1">
            入札者: <span class="font-medium">{{ bid.bidder_display_name }}</span>
          </p>
        </div>

        <div class="ml-4 text-right flex-shrink-0">
          <!-- 価格 -->
          <p class="text-sm font-bold text-gray-900">
            {{ formatPrice(bid.price) }}
          </p>
          <!-- 相対時刻 -->
          <p
            class="text-xs text-gray-500 mt-1"
            :title="formatDateTime(bid.bid_at)"
          >
            {{ formatRelativeTime(bid.bid_at) }}
          </p>
        </div>
      </div>
    </div>

    <!-- 全て見るリンク -->
    <div v-if="bids && bids.length > 0" class="mt-4 text-center">
      <router-link
        to="/admin/bids"
        class="text-sm text-blue-600 hover:text-blue-800 font-medium"
      >
        全ての入札を見る →
      </router-link>
    </div>
  </Card>
</template>

<script setup>
import Card from '@/components/ui/Card.vue'
import { formatRelativeTime, formatDateTime } from '@/utils/timeFormatter'

defineProps({
  /** 最新の入札リスト */
  bids: {
    type: Array,
    default: () => [],
  },
})

// 価格フォーマット（ポイント）
const formatPrice = (price) => {
  return `${price.toLocaleString()} pt`
}
</script>
