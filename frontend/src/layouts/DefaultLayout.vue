<template>
  <div :class="layoutClasses">
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

// 入札者向けルートの場合はbidder-themeクラスを適用
const isBidderRoute = computed(() => {
  return route.path.startsWith('/bidder')
})

const layoutClasses = computed(() => {
  const classes = ['min-h-screen']
  
  if (isBidderRoute.value) {
    classes.push('bidder-theme', 'bg-background')
  } else {
    classes.push('bg-gray-50')
  }
  
  return classes.join(' ')
})
</script>

<style scoped>
/* DefaultLayoutはシンプルなコンテナのみ提供 */
/* 各ページコンポーネントが独自のヘッダー・フッターを持つ */
</style>
