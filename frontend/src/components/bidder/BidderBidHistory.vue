<template>
  <div class="lux-glass-strong rounded-2xl overflow-hidden h-full flex flex-col">
    <!-- Header -->
    <div class="px-6 py-4 border-b border-lux-gold/10 flex items-center justify-between">
      <h2 class="font-display text-lg text-lux-cream">Bid History</h2>
      <span class="text-xs text-lux-silver px-2 py-1 rounded-full bg-lux-noir-light">
        {{ bids.length }} bids
      </span>
    </div>

    <!-- Bid List -->
    <div class="flex-1 overflow-y-auto lux-scrollbar p-4 space-y-3">
      <TransitionGroup name="bid-list">
        <div
          v-for="(bid, index) in bids"
          :key="bid.id"
          :class="[
            'relative rounded-xl p-4 transition-all duration-300',
            bid.is_winning
              ? 'bg-gradient-to-r from-lux-gold/10 via-lux-gold/5 to-transparent border border-lux-gold/30'
              : isOwnBid(bid)
              ? 'bg-lux-noir-light/80 border border-blue-500/20'
              : 'bg-lux-noir-light/50 border border-lux-noir-soft'
          ]"
        >
          <!-- Winning Indicator -->
          <div v-if="bid.is_winning" class="absolute -left-px top-1/2 -translate-y-1/2 w-1 h-8 bg-lux-gold rounded-r-full"></div>

          <div class="flex items-start justify-between gap-4">
            <!-- Left: Bidder Info -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 mb-1">
                <span
                  :class="[
                    'font-medium text-sm truncate',
                    bid.is_winning ? 'text-lux-gold' : isOwnBid(bid) ? 'text-blue-300' : 'text-lux-cream'
                  ]"
                >
                  {{ getBidderDisplayName(bid) }}
                </span>

                <span
                  v-if="isOwnBid(bid)"
                  class="inline-flex items-center px-1.5 py-0.5 rounded text-[10px] font-semibold bg-blue-500/20 text-blue-300 uppercase tracking-wider"
                >
                  You
                </span>
              </div>

              <div class="flex items-center gap-3 text-xs">
                <span class="text-lux-silver/60">
                  {{ formatDateTime(bid.bid_at) }}
                </span>

                <span
                  v-if="bid.is_winning"
                  class="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-lux-gold/10 text-lux-gold text-[10px] font-semibold uppercase tracking-wider"
                >
                  <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                  </svg>
                  Highest
                </span>
              </div>
            </div>

            <!-- Right: Price -->
            <div class="text-right flex-shrink-0">
              <div
                :class="[
                  'font-display text-xl',
                  bid.is_winning ? 'lux-text-gold' : 'text-lux-cream'
                ]"
              >
                {{ formatNumber(bid.price) }}
              </div>
              <div class="text-[10px] text-lux-silver/60 uppercase tracking-wider">points</div>
            </div>
          </div>
        </div>
      </TransitionGroup>

      <!-- Empty State -->
      <div
        v-if="bids.length === 0"
        class="flex flex-col items-center justify-center py-12 text-center"
      >
        <div class="w-16 h-16 rounded-full bg-lux-noir-light border border-lux-gold/10 flex items-center justify-center mb-4">
          <svg class="w-8 h-8 text-lux-silver/40" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
        </div>
        <p class="font-display text-lux-cream mb-1">No Bids Yet</p>
        <p class="text-xs text-lux-silver/60 max-w-[200px]">
          Bids will appear here once the auction starts
        </p>
      </div>
    </div>

    <!-- Legend -->
    <div
      v-if="bids.length > 0"
      class="px-4 py-3 border-t border-lux-gold/10 flex items-center justify-center gap-6 text-[10px] text-lux-silver/60 uppercase tracking-wider"
    >
      <div class="flex items-center gap-1.5">
        <span class="w-2 h-2 rounded-full bg-lux-gold"></span>
        <span>Highest</span>
      </div>
      <div class="flex items-center gap-1.5">
        <span class="w-2 h-2 rounded-full bg-blue-400"></span>
        <span>Your Bid</span>
      </div>
      <div class="flex items-center gap-1.5">
        <span class="w-2 h-2 rounded-full bg-lux-silver/40"></span>
        <span>Others</span>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  bids: {
    type: Array,
    required: true,
    default: () => [],
  },
  currentBidderId: {
    type: String,
    default: null,
  },
})

// Check if bid belongs to current user
function isOwnBid(bid) {
  // If we have bidder_id in the bid and it matches current user
  if (props.currentBidderId && bid.bidder_id) {
    return bid.bidder_id === props.currentBidderId
  }
  // Fallback: check if bid has is_own flag
  return bid.is_own === true
}

// Get bidder display name
function getBidderDisplayName(bid) {
  if (isOwnBid(bid)) {
    return bid.bidder_display_name || 'You'
  }
  // For privacy, show anonymized name for others
  return bid.bidder_display_name || 'Bidder'
}

// Format number with comma separator
function formatNumber(value) {
  return new Intl.NumberFormat('ja-JP').format(value)
}

// Format datetime
function formatDateTime(dateString) {
  if (!dateString) return ''
  const date = new Date(dateString)

  // Check if date is today
  const today = new Date()
  const isToday =
    date.getDate() === today.getDate() &&
    date.getMonth() === today.getMonth() &&
    date.getFullYear() === today.getFullYear()

  if (isToday) {
    // Show only time for today
    return new Intl.DateTimeFormat('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
      hour12: false,
    }).format(date)
  } else {
    // Show date and time for other days
    return new Intl.DateTimeFormat('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
      hour12: false,
    }).format(date)
  }
}
</script>

<style scoped>
/* Bid List Transitions */
.bid-list-enter-active {
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.bid-list-leave-active {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.bid-list-enter-from {
  opacity: 0;
  transform: translateX(-20px);
}

.bid-list-leave-to {
  opacity: 0;
  transform: translateX(20px);
}

.bid-list-move {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}
</style>
