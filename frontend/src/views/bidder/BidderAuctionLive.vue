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
    <div v-else-if="auction" class="max-w-7xl mx-auto">
      <!-- Header -->
      <div class="bg-white rounded-lg shadow-sm p-6 mb-4">
        <h1 class="text-2xl font-bold text-gray-900">{{ auction.title }}</h1>
        <p v-if="auction.description" class="mt-2 text-gray-600">{{ auction.description }}</p>

        <!-- WebSocket Status -->
        <div class="mt-4 flex items-center gap-2">
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

      <!-- Points Display -->
      <div class="bg-white rounded-lg shadow-sm p-6 mb-4">
        <h2 class="text-lg font-semibold text-gray-900 mb-4">ポイント残高</h2>
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="text-center p-4 bg-blue-50 rounded-lg">
            <div class="text-sm text-gray-600">合計</div>
            <div class="text-2xl font-bold text-blue-600">
              {{ formatNumber(points.total) }}
            </div>
          </div>
          <div class="text-center p-4 bg-green-50 rounded-lg">
            <div class="text-sm text-gray-600">利用可能</div>
            <div class="text-2xl font-bold text-green-600">
              {{ formatNumber(points.available) }}
            </div>
          </div>
          <div class="text-center p-4 bg-yellow-50 rounded-lg">
            <div class="text-sm text-gray-600">予約済み</div>
            <div class="text-2xl font-bold text-yellow-600">
              {{ formatNumber(points.reserved) }}
            </div>
          </div>
        </div>
      </div>

      <!-- Items Tabs -->
      <div class="bg-white rounded-lg shadow-sm p-6 mb-4">
        <h2 class="text-lg font-semibold text-gray-900 mb-4">商品一覧</h2>
        <div class="flex gap-2 overflow-x-auto">
          <button
            v-for="item in items"
            :key="item.id"
            @click="handleSelectItem(item.id)"
            :class="[
              'px-4 py-2 rounded-lg whitespace-nowrap transition-colors',
              currentItem?.id === item.id
                ? 'bg-blue-600 text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            ]"
          >
            {{ item.name }}
            <span
              :class="[
                'ml-2 text-xs px-2 py-1 rounded',
                item.status === 'started' ? 'bg-green-500 text-white' :
                item.status === 'ended' ? 'bg-gray-500 text-white' :
                'bg-yellow-500 text-white'
              ]"
            >
              {{ getStatusLabel(item.status) }}
            </span>
          </button>
        </div>
      </div>

      <!-- Current Item & Bid Panel -->
      <div v-if="currentItem" class="grid grid-cols-1 lg:grid-cols-3 gap-4">
        <!-- Bid Panel -->
        <div class="lg:col-span-2 bg-white rounded-lg shadow-sm p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">{{ currentItem.name }}</h2>

          <!-- Current Price -->
          <div class="text-center p-8 bg-gray-50 rounded-lg mb-6">
            <div class="text-sm text-gray-600 mb-2">現在価格</div>
            <div class="text-5xl font-bold text-gray-900">
              {{ formatNumber(currentPrice) }}
            </div>
            <div class="text-sm text-gray-500 mt-1">ポイント</div>
          </div>

          <!-- Bid Button -->
          <button
            @click="handlePlaceBid"
            :disabled="!canBid"
            :class="[
              'w-full py-4 px-6 rounded-lg font-semibold text-lg transition-colors',
              canBid
                ? 'bg-blue-600 text-white hover:bg-blue-700 active:bg-blue-800'
                : 'bg-gray-200 text-gray-500 cursor-not-allowed'
            ]"
          >
            {{ loadingBid ? '処理中...' : canBid ? '入札する' : cannotBidReason }}
          </button>
        </div>

        <!-- Bid History -->
        <div class="bg-white rounded-lg shadow-sm p-6">
          <h2 class="text-lg font-semibold text-gray-900 mb-4">入札履歴</h2>
          <div class="space-y-2 max-h-96 overflow-y-auto">
            <div
              v-for="bid in bids"
              :key="bid.id"
              :class="[
                'p-3 rounded-lg',
                bid.is_winning ? 'bg-blue-50 border border-blue-200' : 'bg-gray-50'
              ]"
            >
              <div class="flex justify-between items-start">
                <div>
                  <div class="font-semibold text-gray-900">
                    {{ bid.bidder_display_name || '入札者' }}
                  </div>
                  <div class="text-xs text-gray-500 mt-1">
                    {{ formatDateTime(bid.bid_at) }}
                  </div>
                </div>
                <div class="text-right">
                  <div class="font-bold text-gray-900">
                    {{ formatNumber(bid.price) }}
                  </div>
                  <div v-if="bid.is_winning" class="text-xs text-blue-600 font-semibold mt-1">
                    最高入札
                  </div>
                </div>
              </div>
            </div>
            <div v-if="bids.length === 0" class="text-center text-gray-500 py-8">
              入札履歴はありません
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useBidderAuctionLiveStore } from '@/stores/bidderAuctionLive'

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
} = storeToRefs(store)

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

// Formatters
function formatNumber(value) {
  return new Intl.NumberFormat('ja-JP').format(value)
}

function formatDateTime(dateString) {
  if (!dateString) return ''
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('ja-JP', {
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
  }).format(date)
}

function getStatusLabel(status) {
  const labels = {
    pending: '待機中',
    started: '開始',
    ended: '終了',
  }
  return labels[status] || status
}
</script>
