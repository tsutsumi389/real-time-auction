<template>
  <div class="min-h-screen" :class="layoutClass">
    <!-- 入札者用ヘッダー（入札者認証が必要なルートのみ） -->
    <BidderHeader v-if="shouldShowBidderHeader" />

    <slot />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import BidderHeader from '@/components/bidder/BidderHeader.vue'

const route = useRoute()

// 入札者用ヘッダーを表示するかどうか
const shouldShowBidderHeader = computed(() => {
  return route.meta.requiresBidderAuth === true
})

// 入札者用ルートの場合はラグジュアリーテーマ、それ以外は通常の背景
const layoutClass = computed(() => {
  return route.meta.requiresBidderAuth === true ? 'bg-lux-noir' : 'bg-gray-50'
})
</script>

<style scoped>
/* DefaultLayoutはシンプルなコンテナのみ提供 */
/* 各ページコンポーネントが独自のヘッダー・フッターを持つ */
</style>
