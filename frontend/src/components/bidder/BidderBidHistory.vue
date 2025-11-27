<template>
  <div class="bg-white rounded-lg shadow-sm p-6">
    <div class="flex items-center justify-between mb-4">
      <h2 class="text-lg font-semibold text-gray-900">入札履歴</h2>
      <span class="text-sm text-gray-500">{{ bids.length }}件</span>
    </div>

    <!-- Bid List -->
    <div class="space-y-2 max-h-[600px] overflow-y-auto custom-scrollbar">
      <div
        v-for="(bid, index) in bids"
        :key="bid.id"
        :class="[
          'p-4 rounded-lg transition-all duration-200',
          bid.is_winning
            ? 'bg-gradient-to-r from-blue-50 to-indigo-50 border-2 border-blue-200 shadow-sm'
            : isOwnBid(bid)
            ? 'bg-blue-50 border border-blue-100'
            : 'bg-gray-50 border border-gray-100',
          index === 0 ? 'animate-slide-in' : ''
        ]"
      >
        <div class="flex items-start justify-between">
          <!-- Left side: Bidder info -->
          <div class="flex-1">
            <div class="flex items-center gap-2">
              <span
                :class="[
                  'font-semibold',
                  isOwnBid(bid) ? 'text-blue-900' : 'text-gray-900'
                ]"
              >
                {{ getBidderDisplayName(bid) }}
              </span>
              <span
                v-if="isOwnBid(bid)"
                class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-blue-600 text-white"
              >
                あなた
              </span>
            </div>
            <div class="flex items-center gap-2 mt-2">
              <span class="text-xs text-gray-500">
                {{ formatDateTime(bid.bid_at) }}
              </span>
              <span
                v-if="bid.is_winning"
                class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-semibold bg-green-100 text-green-800"
              >
                <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
                最高入札
              </span>
            </div>
          </div>

          <!-- Right side: Price -->
          <div class="text-right ml-4">
            <div
              :class="[
                'text-xl font-bold',
                bid.is_winning ? 'text-green-600' : 'text-gray-900'
              ]"
            >
              {{ formatNumber(bid.price) }}
            </div>
            <div class="text-xs text-gray-500 mt-1">ポイント</div>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div
        v-if="bids.length === 0"
        class="text-center py-12 text-gray-500"
      >
        <svg class="w-16 h-16 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
        </svg>
        <p class="font-medium">入札履歴はありません</p>
        <p class="text-sm mt-2">価格が開示されると入札できるようになります</p>
      </div>
    </div>

    <!-- Legend -->
    <div
      v-if="bids.length > 0"
      class="mt-4 pt-4 border-t border-gray-200 flex flex-wrap gap-4 text-xs text-gray-600"
    >
      <div class="flex items-center">
        <span class="inline-block w-3 h-3 bg-gradient-to-r from-blue-50 to-indigo-50 border border-blue-200 rounded mr-2"></span>
        最高入札
      </div>
      <div class="flex items-center">
        <span class="inline-block w-3 h-3 bg-blue-50 border border-blue-100 rounded mr-2"></span>
        自分の入札
      </div>
      <div class="flex items-center">
        <span class="inline-block w-3 h-3 bg-gray-50 border border-gray-100 rounded mr-2"></span>
        他者の入札
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

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
    return bid.bidder_display_name || 'あなた'
  }
  // For privacy, show anonymized name for others
  return bid.bidder_display_name || '入札者'
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
    return new Intl.DateTimeFormat('ja-JP', {
      hour: '2-digit',
      minute: '2-digit',
      second: '2-digit',
    }).format(date)
  } else {
    // Show date and time for other days
    return new Intl.DateTimeFormat('ja-JP', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    }).format(date)
  }
}
</script>

<style scoped>
/* Custom scrollbar */
.custom-scrollbar {
  scrollbar-width: thin;
  scrollbar-color: #cbd5e0 #f7fafc;
}

.custom-scrollbar::-webkit-scrollbar {
  width: 8px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: #f7fafc;
  border-radius: 4px;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #cbd5e0;
  border-radius: 4px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #a0aec0;
}

/* Slide-in animation for new bids */
@keyframes slide-in {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.animate-slide-in {
  animation: slide-in 0.3s ease-out;
}
</style>
