<template>
  <div class="min-h-screen bg-gray-50 p-4">
    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center min-h-[400px]">
      <div class="text-center">
        <div class="inline-block animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        <p class="mt-4 text-gray-600">読み込み中...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="max-w-2xl mx-auto mt-8">
      <div class="bg-red-50 border border-red-200 rounded-lg p-6">
        <h2 class="text-lg font-semibold text-red-900 mb-2">エラーが発生しました</h2>
        <p class="text-red-700">{{ error }}</p>
        <button
          @click="handleRetry"
          class="mt-4 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
        >
          再試行
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <div v-else-if="auction" class="max-w-7xl mx-auto space-y-4">
      <!-- Header -->
      <div class="bg-white rounded-lg shadow-sm p-6">
        <div class="flex items-start justify-between">
          <div>
            <h1 class="text-3xl font-bold text-gray-900">{{ auction.title }}</h1>
            <p v-if="auction.description" class="mt-2 text-gray-600">{{ auction.description }}</p>
          </div>
          <!-- WebSocket Status -->
          <div class="flex items-center gap-2">
            <div
              :class="[
                'w-3 h-3 rounded-full',
                wsConnected ? 'bg-green-500' : wsReconnecting ? 'bg-yellow-500 animate-pulse' : 'bg-red-500'
              ]"
            ></div>
            <span class="text-sm text-gray-600">
              {{ wsConnected ? '接続中' : wsReconnecting ? '再接続中...' : '切断' }}
            </span>
            <span v-if="wsReconnecting" class="text-sm text-gray-500">
              ({{ reconnectAttempt }}/{{ maxReconnectAttempts }})
            </span>
          </div>
        </div>
      </div>

      <!-- Points Display Component -->
      <PointsDisplay :points="points" />

      <!-- Item Tabs Component -->
      <ItemTabs
        :items="items"
        :selected-item-id="currentItem?.id"
        @select="handleSelectItem"
      />

      <!-- Main Layout: Bid Panel + History (3 columns on desktop, stacked on mobile) -->
      <div v-if="currentItem" class="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <!-- Bid Panel Component (takes 2 columns on desktop) -->
        <div class="lg:col-span-2">
          <BidPanel
            :item="currentItem"
            :current-price="currentPrice"
            :available-points="points.available"
            :is-loading="loadingBid"
            :can-bid="canBid"
            :disabled-reason="cannotBidReason"
            :is-winning="isOwnBidWinning"
            @bid="handlePlaceBid"
          />
        </div>

        <!-- Bid History Component (takes 1 column on desktop) -->
        <div>
          <BidderBidHistory
            :bids="bids"
            :current-bidder-id="currentBidderId"
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useBidderAuctionLiveStore } from '@/stores/bidderAuctionLive'
import PointsDisplay from '@/components/bidder/PointsDisplay.vue'
import ItemTabs from '@/components/bidder/ItemTabs.vue'
import BidPanel from '@/components/bidder/BidPanel.vue'
import BidderBidHistory from '@/components/bidder/BidderBidHistory.vue'

const route = useRoute()
const store = useBidderAuctionLiveStore()

// State from store
const {
  auction,
  items,
  currentItem,
  bids,
  points,
  loading,
  loadingBid,
  error,
  wsConnected,
  wsReconnecting,
  reconnectAttempt,
  maxReconnectAttempts,
  currentPrice,
  canBid,
  cannotBidReason,
  isOwnBidWinning,
} = storeToRefs(store)

// Current bidder ID (could be retrieved from bidder auth store)
const currentBidderId = computed(() => {
  // TODO: Get from bidder auth store if available
  return null
})

// Initialization
onMounted(async () => {
  const auctionId = route.params.id

  // Initialize auction data
  const success = await store.initialize(auctionId)

  if (success) {
    // Connect WebSocket
    store.connectWebSocket(auctionId)
  }
})

// Cleanup
onUnmounted(() => {
  store.disconnectWebSocket()
  store.reset()
})

// Handlers
function handleSelectItem(itemId) {
  store.selectItem(itemId)
}

async function handlePlaceBid() {
  if (!currentItem.value || !canBid.value) {
    return
  }

  const result = await store.placeBid(currentItem.value.id, currentPrice.value)

  if (result.success) {
    console.log('入札成功')
  } else {
    console.error('入札失敗:', result.error)
  }
}

function handleRetry() {
  const auctionId = route.params.id
  store.clearError()
  store.initialize(auctionId)
}
</script>
