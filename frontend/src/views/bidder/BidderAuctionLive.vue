<template>
  <div class="min-h-screen bg-gray-50 p-4 pb-24 md:pb-4 relative">
    <!-- Confetti Effect -->
    <Confetti :active="showConfetti" />

    <!-- Toast Container -->
    <ToastContainer />

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
    <div v-else-if="auction" class="max-w-7xl mx-auto space-y-4 relative z-10">
      <!-- Header -->
      <div class="bg-white rounded-lg shadow-sm p-6">
        <div class="flex items-start justify-between">
          <div>
            <h1 class="text-3xl font-bold text-gray-900">{{ auction.title }}</h1>
            <p v-if="auction.description" class="mt-2 text-gray-600">{{ auction.description }}</p>
          </div>
          <!-- Status & Points -->
          <div class="flex items-center gap-4">
            <!-- Available Points -->
            <div class="flex items-center gap-1.5 bg-green-50 border border-green-200 rounded-lg px-3 py-1.5">
              <svg class="w-4 h-4 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span class="text-sm font-semibold text-green-700">{{ formatNumber(points.available) }}</span>
              <span class="text-xs text-green-600">pts</span>
            </div>
            <!-- Participant Count -->
            <div class="flex items-center gap-1.5 text-sm text-gray-600">
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM17 6a3 3 0 11-6 0 3 3 0 016 0zM12.93 17c.046-.327.07-.66.07-1a6.97 6.97 0 00-1.5-4.33A5 5 0 0119 16v1h-6.07zM6 11a5 5 0 015 5v1H1v-1a5 5 0 015-5z" />
              </svg>
              <span>{{ participantCount }}人</span>
            </div>
            <!-- Connection Status -->
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
      </div>

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

    <!-- Winning Modal -->
    <div v-if="showWinningModal" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm">
      <div class="bg-white rounded-2xl shadow-2xl max-w-md w-full p-8 text-center transform transition-all animate-bounce-in">
        <div class="mx-auto flex items-center justify-center h-20 w-20 rounded-full bg-green-100 mb-6">
          <svg class="h-10 w-10 text-green-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
        <h3 class="text-2xl font-bold text-gray-900 mb-2">落札おめでとうございます！</h3>
        <p class="text-gray-600 mb-6">
          商品「{{ winningItemName }}」を落札しました。
        </p>
        <button
          @click="closeWinningModal"
          class="w-full inline-flex justify-center rounded-xl border border-transparent shadow-sm px-4 py-3 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:text-sm"
        >
          閉じる
        </button>
      </div>
    </div>

    <!-- WebSocket Disconnected Overlay -->
    <div 
      v-if="!wsConnected && !wsReconnecting && auction" 
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/70 backdrop-blur-sm"
    >
      <div class="bg-white rounded-2xl shadow-2xl max-w-md w-full p-8 text-center">
        <div class="mx-auto flex items-center justify-center h-20 w-20 rounded-full bg-red-100 mb-6">
          <svg class="h-10 w-10 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
        <h3 class="text-2xl font-bold text-gray-900 mb-2">接続が切断されました</h3>
        <p class="text-gray-600 mb-6">
          リアルタイム更新が停止しています。<br />
          入札を続けるには再接続してください。
        </p>
        <button
          @click="handleReconnect"
          class="w-full inline-flex justify-center rounded-xl border border-transparent shadow-sm px-4 py-3 bg-blue-600 text-base font-medium text-white hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:text-sm"
        >
          <svg class="w-5 h-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          再接続する
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted, computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useBidderAuctionLiveStore } from '@/stores/bidderAuctionLive'
import { useToast } from '@/composables/useToast'
import ItemTabs from '@/components/bidder/ItemTabs.vue'
import BidPanel from '@/components/bidder/BidPanel.vue'
import BidderBidHistory from '@/components/bidder/BidderBidHistory.vue'
import ToastContainer from '@/components/ui/ToastContainer.vue'
import Confetti from '@/components/ui/Confetti.vue'

const route = useRoute()
const store = useBidderAuctionLiveStore()
const toast = useToast()

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
  participantCount,
} = storeToRefs(store)

// Local state for winning modal
const showWinningModal = ref(false)
const showConfetti = ref(false)
const winningItemName = ref('')

function triggerWinningEffect(itemName) {
  winningItemName.value = itemName
  showConfetti.value = true
  setTimeout(() => {
    showWinningModal.value = true
  }, 500)
  
  // Stop confetti after 5 seconds
  setTimeout(() => {
    showConfetti.value = false
  }, 5000)
}

function closeWinningModal() {
  showWinningModal.value = false
}

// Current bidder ID (could be retrieved from bidder auth store)
const currentBidderId = computed(() => {
  // TODO: Get from bidder auth store if available
  // For now, we rely on the store's isOwnBidWinning logic which checks token internally
  // But for the watcher above, we need the ID. 
  // Let's try to get it from the store if exposed, or parse token again.
  // Since store doesn't expose it directly as a ref, we might need to rely on isOwnBidWinning check in a different way
  // or import the token service here.
  // For simplicity, let's use a workaround: check if the winning bid in bids array is 'is_own'
  return null 
})

// Improved winning check for watcher
watch(
  () => currentItem.value?.status,
  (newStatus, oldStatus) => {
    if (newStatus === 'ended' && oldStatus === 'active') {
      // Check if the winning bid is ours
      const winningBid = bids.value.find(b => b.is_winning)
      if (winningBid && (winningBid.is_own || winningBid.bidder_id === currentBidderId.value)) {
        triggerWinningEffect(currentItem.value.name)
      }
    }
  }
)


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

  // Initialize audio on first user interaction
  if (store.audioNotification && !store.audioNotification.isEnabled) {
    store.audioNotification.initAudio()
  }

  const result = await store.placeBid(currentItem.value.id, currentPrice.value)

  if (result.success) {
    // Play success sound
    if (store.audioNotification) {
      store.audioNotification.playBidSuccessSound()
    }

    toast.success(
      '入札成功',
      `${formatNumber(currentPrice.value)}ポイントで入札しました`
    )
  } else {
    toast.error(
      '入札失敗',
      result.error || '入札処理中にエラーが発生しました'
    )
  }
}

function handleReconnect() {
  const auctionId = route.params.id
  store.connectWebSocket(auctionId)
}

function formatNumber(value) {
  return new Intl.NumberFormat('ja-JP').format(value)
}

function handleRetry() {
  const auctionId = route.params.id
  store.clearError()
  store.initialize(auctionId)
}
</script>

<style scoped>
@keyframes bounce-in {
  0% { transform: scale(0.3); opacity: 0; }
  50% { transform: scale(1.05); opacity: 1; }
  70% { transform: scale(0.9); }
  100% { transform: scale(1); }
}
.animate-bounce-in {
  animation: bounce-in 0.5s cubic-bezier(0.215, 0.610, 0.355, 1.000) both;
}
</style>
