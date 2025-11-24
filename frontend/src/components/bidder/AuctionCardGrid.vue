<template>
  <div class="auction-card-grid">
    <!-- グリッドレイアウト -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6">
      <AuctionCard
        v-for="auction in auctions"
        :key="auction.id"
        :auction="auction"
        @view-details="handleViewDetails"
        @join-auction="handleJoinAuction"
      />
    </div>

    <!-- 空の状態 -->
    <div v-if="auctions.length === 0 && !loading" class="empty-state py-20 text-center">
      <slot name="empty">
        <div class="text-gray-500">
          <svg class="mx-auto h-12 w-12 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"></path>
          </svg>
          <p class="text-lg">{{ emptyMessage }}</p>
        </div>
      </slot>
    </div>

    <!-- ローディング -->
    <div v-if="loading" class="py-20">
      <LoadingSpinner size="lg" text="読み込み中..." center />
    </div>
  </div>
</template>

<script setup>
import AuctionCard from './AuctionCard.vue'
import LoadingSpinner from '../ui/LoadingSpinner.vue'

const props = defineProps({
  auctions: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  emptyMessage: {
    type: String,
    default: 'オークションが見つかりませんでした'
  }
})

const emit = defineEmits(['view-details', 'join-auction'])

const handleViewDetails = (auction) => {
  emit('view-details', auction)
}

const handleJoinAuction = (auction) => {
  emit('join-auction', auction)
}
</script>

<style scoped>
.auction-card-grid {
  width: 100%;
}
</style>
