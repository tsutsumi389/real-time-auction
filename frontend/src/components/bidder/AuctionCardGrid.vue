<template>
  <div class="auction-card-grid">
    <!-- ARIA Live region for screen readers -->
    <div aria-live="polite" aria-atomic="true" class="sr-only">
      <template v-if="loading">読み込み中</template>
      <template v-else-if="auctions.length > 0">{{ auctions.length }}件のオークションを表示中</template>
      <template v-else>オークションが見つかりませんでした</template>
    </div>

    <!-- Skeleton Loading State -->
    <div
      v-if="loading"
      class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5 sm:gap-6"
    >
      <AuctionCardSkeleton
        v-for="i in skeletonCount"
        :key="`skeleton-${i}`"
      />
    </div>

    <!-- Grid Layout with Transition -->
    <TransitionGroup
      v-else-if="auctions.length > 0"
      name="card-list"
      tag="div"
      class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-5 sm:gap-6"
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

    <!-- Empty State -->
    <div v-else class="empty-state py-20 text-center">
      <slot name="empty">
        <div class="lux-fade-in">
          <div class="w-20 h-20 mx-auto mb-6 rounded-2xl bg-lux-noir-light/50 border border-lux-gold/20 flex items-center justify-center">
            <svg class="w-10 h-10 text-lux-gold/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"></path>
            </svg>
          </div>
          <p class="font-display text-xl text-lux-cream/60">{{ emptyMessage }}</p>
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
  transition: all 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  transition-delay: var(--stagger-delay, 0ms);
}

.card-list-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.card-list-enter-from {
  opacity: 0;
  transform: translateY(30px) scale(0.95);
}

.card-list-leave-to {
  opacity: 0;
  transform: scale(0.9);
}

.card-list-move {
  transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
}

/* Luxury color utilities */
.bg-lux-noir-light\/50 {
  background-color: hsl(0 0% 8% / 0.5);
}

.border-lux-gold\/20 {
  border-color: hsl(43 74% 49% / 0.2);
}

.text-lux-gold\/40 {
  color: hsl(43 74% 49% / 0.4);
}

.text-lux-cream\/60 {
  color: hsl(45 30% 96% / 0.6);
}

/* Font display */
.font-display {
  font-family: 'Cormorant Garamond', Georgia, serif;
}
</style>
