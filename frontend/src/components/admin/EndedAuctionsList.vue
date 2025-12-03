<template>
  <Card class="p-6">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">終了したオークション</h3>

    <!-- 空の状態 -->
    <div
      v-if="!auctions || auctions.length === 0"
      class="text-center py-8 text-gray-500"
    >
      <p>終了したオークションはありません</p>
    </div>

    <!-- オークションリスト -->
    <div v-else class="space-y-3">
      <div
        v-for="auction in auctions"
        :key="auction.id"
        class="p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors"
      >
        <div class="flex items-start justify-between">
          <div class="flex-1 min-w-0">
            <!-- オークションタイトル -->
            <p class="text-sm font-medium text-gray-900 truncate">
              {{ auction.item_title || auction.title }}
            </p>
            <!-- 落札者情報 -->
            <p class="text-xs text-gray-600 mt-1">
              落札者:
              <span v-if="auction.winner_display_name" class="font-medium">
                {{ auction.winner_display_name }}
              </span>
              <span v-else class="text-gray-400">なし</span>
            </p>
          </div>

          <div class="ml-4 text-right flex-shrink-0">
            <!-- 最終価格 -->
            <p class="text-sm font-bold text-gray-900">
              {{ formatPrice(auction.current_price) }}
            </p>
            <!-- 終了時刻 -->
            <p
              class="text-xs text-gray-500 mt-1"
              :title="formatDateTime(auction.ended_at || auction.updated_at)"
            >
              {{ formatRelativeTime(auction.ended_at || auction.updated_at) }}
            </p>
          </div>
        </div>

        <!-- ステータスバッジ -->
        <div class="mt-2">
          <span
            class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800"
          >
            終了
          </span>
        </div>
      </div>
    </div>

    <!-- 全て見るリンク -->
    <div v-if="auctions && auctions.length > 0" class="mt-4 text-center">
      <router-link
        to="/admin/auctions"
        class="text-sm text-blue-600 hover:text-blue-800 font-medium"
      >
        全てのオークションを見る →
      </router-link>
    </div>
  </Card>
</template>

<script setup>
import Card from '@/components/ui/Card.vue'
import { formatRelativeTime, formatDateTime } from '@/utils/timeFormatter'

defineProps({
  /** 終了したオークションリスト */
  auctions: {
    type: Array,
    default: () => [],
  },
})

// 価格フォーマット
const formatPrice = (price) => {
  if (price === null || price === undefined) {
    return '¥0'
  }
  return `¥${price.toLocaleString()}`
}
</script>
