<template>
  <div class="auction-card-grid">
    <!-- ARIA Live region for screen readers -->
    <div aria-live="polite" aria-atomic="true" class="sr-only">
      <template v-if="loading">読み込み中</template>
      <template v-else-if="auctions.length > 0">{{ auctions.length }}件のオークションを表示中</template>
      <template v-else>オークションが見つかりませんでした</template>
    </div>

    <!-- Skeleton loading state -->
    <div
      v-if="loading"
      class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6"
    >
      <AuctionCardSkeleton
        v-for="i in skeletonCount"
        :key="`skeleton-${i}`"
      />
    </div>

    <!-- Grid layout with transition -->
    <TransitionGroup
      v-else-if="auctions.length > 0"
      name="card-list"
      tag="div"
      class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6"
    >
      <AuctionCard
        v-for="(auction, index) in auctions"
        :key="auction.id"
        :auction="auction"
        :style="{ '--stagger-delay': `${index * 50}ms` }"
        @view-details="handleViewDetails"
        @join-auction="handleJoinAuction"
      />
    </TransitionGroup>

    <!-- Empty state -->
    <div v-else class="empty-state py-20 text-center">
      <slot name="empty">
        <div class="text-gray-500">
          <svg class="mx-auto h-12 w-12 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"></path>
          </svg>
          <p class="text-lg">{{ emptyMessage }}</p>
        </div>
      </slot>
    </div>
  </div>
</template>

<script setup>
import AuctionCard from './AuctionCard.vue'
import AuctionCardSkeleton from './AuctionCardSkeleton.vue'

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
  },
  skeletonCount: {
    type: Number,
    default: 6
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

/* Card list transition animations */
.card-list-enter-active {
  transition: all 0.4s ease-out;
  transition-delay: var(--stagger-delay, 0ms);
}

.card-list-leave-active {
  transition: all 0.3s ease-in;
}

.card-list-enter-from {
  opacity: 0;
  transform: translateY(20px);
}

.card-list-leave-to {
  opacity: 0;
  transform: scale(0.95);
}

.card-list-move {
  transition: transform 0.4s ease;
}
</style>
