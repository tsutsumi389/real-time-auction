<template>
  <div class="bid-history-container rounded-2xl overflow-hidden h-full flex flex-col">
    <!-- Header -->
    <div class="px-6 py-5 border-b border-lux-gold/20 flex items-center justify-between bg-gradient-to-r from-lux-noir-light/80 to-lux-noir/80">
      <div class="flex items-center gap-3">
        <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-lux-gold/20 to-lux-gold/5 border border-lux-gold/30 flex items-center justify-center">
          <svg class="w-5 h-5 text-lux-gold" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 9l2 2 4-4" />
          </svg>
        </div>
        <div>
          <h2 class="font-display text-lg text-lux-cream">入札履歴</h2>
          <p class="text-[10px] text-lux-silver/50 uppercase tracking-widest">Bid History</p>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <span class="px-3 py-1.5 rounded-lg bg-lux-noir-medium border border-lux-gold/20 text-xs font-semibold text-lux-gold tabular-nums">
          {{ bids.length }}
        </span>
      </div>
    </div>

    <!-- Bid List -->
    <div class="flex-1 overflow-y-auto lux-scrollbar p-4 space-y-3 bg-gradient-to-b from-lux-noir-light/30 to-transparent">
      <TransitionGroup name="bid-list">
        <div
          v-for="(bid, index) in bids"
          :key="bid.id"
          :class="[
            'relative rounded-xl p-4 transition-all duration-300',
            bid.is_winning ? 'bid-card-winning' :
            isOwnBid(bid) ? 'bid-card-own' : 'bid-card-other'
          ]"
        >
          <!-- Winning Indicator -->
          <div v-if="bid.is_winning" class="absolute -left-px top-1/2 -translate-y-1/2 w-1.5 h-10 bg-gradient-to-b from-lux-gold-light via-lux-gold to-lux-gold-dark rounded-r-full shadow-lg shadow-lux-gold/30"></div>

          <!-- Rank Badge for top 3 -->
          <div
            v-if="index < 3 && bid.is_winning"
            class="absolute -top-2 -right-2 w-7 h-7 rounded-full bg-gradient-to-br from-lux-gold to-lux-gold-dark flex items-center justify-center shadow-lg shadow-lux-gold/40 text-xs font-bold text-lux-noir"
          >
            {{ index + 1 }}
          </div>

          <div class="flex items-start justify-between gap-4">
            <!-- Left: Bidder Info -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 mb-1.5">
                <!-- Avatar -->
                <div
                  :class="[
                    'w-8 h-8 rounded-full flex items-center justify-center text-xs font-bold',
                    bid.is_winning ? 'bg-lux-gold/20 text-lux-gold border border-lux-gold/40' :
                    isOwnBid(bid) ? 'bg-blue-500/20 text-blue-300 border border-blue-500/30' :
                    'bg-lux-noir-medium text-lux-silver/60 border border-lux-noir-soft'
                  ]"
                >
                  {{ getBidderInitial(bid) }}
                </div>

                <div class="min-w-0 flex-1">
                  <div class="flex items-center gap-2">
                    <span
                      :class="[
                        'font-semibold text-sm truncate',
                        bid.is_winning ? 'text-lux-gold' : isOwnBid(bid) ? 'text-blue-300' : 'text-lux-cream'
                      ]"
                    >
                      {{ getBidderDisplayName(bid) }}
                    </span>

                    <span
                      v-if="isOwnBid(bid)"
                      class="inline-flex items-center px-2 py-0.5 rounded-md text-[9px] font-bold bg-blue-500/25 text-blue-300 uppercase tracking-wider border border-blue-500/30"
                    >
                      You
                    </span>
                  </div>

                  <div class="flex items-center gap-2 mt-0.5">
                    <span class="text-[11px] text-lux-silver/50 tabular-nums">
                      {{ formatDateTime(bid.bid_at) }}
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Right: Price -->
            <div class="text-right flex-shrink-0">
              <div
                :class="[
                  'font-display text-xl font-light tabular-nums',
                  bid.is_winning ? 'lux-text-gold' : 'text-lux-cream'
                ]"
              >
                {{ formatNumber(bid.price) }}
              </div>
              <div class="text-[9px] text-lux-silver/40 uppercase tracking-widest">points</div>
            </div>
          </div>

          <!-- Winning Badge -->
          <div
            v-if="bid.is_winning"
            class="mt-3 pt-3 border-t border-lux-gold/20 flex items-center justify-between"
          >
            <span class="inline-flex items-center gap-1.5 text-xs font-semibold text-lux-gold">
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
              </svg>
              最高入札額
            </span>
            <span class="text-[10px] text-lux-gold/60 uppercase tracking-wider">Leading</span>
          </div>
        </div>
      </TransitionGroup>

      <!-- Empty State -->
      <div
        v-if="bids.length === 0"
        class="flex flex-col items-center justify-center py-16 text-center"
      >
        <div class="relative mb-6">
          <div class="w-20 h-20 rounded-2xl bg-gradient-to-br from-lux-noir-light to-lux-noir-medium border border-lux-gold/20 flex items-center justify-center">
            <svg class="w-10 h-10 text-lux-gold/30" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
            </svg>
          </div>
          <div class="absolute -bottom-1 -right-1 w-6 h-6 rounded-full bg-lux-noir-medium border border-lux-gold/30 flex items-center justify-center">
            <svg class="w-3 h-3 text-lux-gold/50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
          </div>
        </div>
        <p class="font-display text-lg text-lux-cream mb-2">入札履歴がありません</p>
        <p class="text-xs text-lux-silver/50 max-w-[220px] leading-relaxed">
          オークションが開始されると、<br />入札履歴がここに表示されます
        </p>
      </div>
    </div>

    <!-- Legend -->
    <div
      v-if="bids.length > 0"
      class="px-4 py-3 border-t border-lux-gold/15 bg-lux-noir-light/50 flex items-center justify-center gap-6 text-[10px] text-lux-silver/50 uppercase tracking-wider"
    >
      <div class="flex items-center gap-1.5">
        <span class="w-2.5 h-2.5 rounded-full bg-gradient-to-br from-lux-gold-light to-lux-gold shadow-sm shadow-lux-gold/30"></span>
        <span>最高額</span>
      </div>
      <div class="flex items-center gap-1.5">
        <span class="w-2.5 h-2.5 rounded-full bg-blue-400 shadow-sm shadow-blue-400/30"></span>
        <span>あなた</span>
      </div>
      <div class="flex items-center gap-1.5">
        <span class="w-2.5 h-2.5 rounded-full bg-lux-silver/40"></span>
        <span>その他</span>
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
