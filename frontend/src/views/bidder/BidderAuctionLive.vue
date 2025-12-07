<template>
  <div class="theme-luxury min-h-screen relative">
    <!-- Ambient Background -->
    <div class="fixed inset-0 bg-lux-noir pointer-events-none">
      <!-- Gradient Orbs -->
      <div class="absolute top-0 right-0 w-[800px] h-[800px] bg-lux-gold/3 rounded-full blur-[150px]"></div>
      <div class="absolute bottom-0 left-0 w-[600px] h-[600px] bg-lux-gold/2 rounded-full blur-[120px]"></div>

      <!-- Subtle Grid Pattern -->
      <div class="absolute inset-0 opacity-[0.015]" style="background-image: linear-gradient(rgba(212,175,55,0.4) 1px, transparent 1px), linear-gradient(90deg, rgba(212,175,55,0.4) 1px, transparent 1px); background-size: 80px 80px;"></div>

      <!-- Noise Texture -->
      <div class="absolute inset-0 lux-noise"></div>
    </div>

    <!-- Confetti Effect -->
    <Confetti :active="showConfetti" />

    <!-- Toast Container -->
    <ToastContainer />

    <!-- Loading State -->
    <div v-if="loading" class="relative z-10 flex items-center justify-center min-h-screen">
      <div class="text-center lux-fade-in">
        <div class="relative w-20 h-20 mx-auto mb-6">
          <div class="absolute inset-0 rounded-full border-2 border-lux-gold/20"></div>
          <div class="absolute inset-0 rounded-full border-2 border-transparent border-t-lux-gold animate-spin"></div>
          <div class="absolute inset-2 rounded-full border border-lux-gold/10"></div>
        </div>
        <p class="font-display text-xl text-lux-cream tracking-wide">Loading Auction</p>
        <p class="text-sm text-lux-silver mt-2">Please wait...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="relative z-10 flex items-center justify-center min-h-screen px-4">
      <div class="max-w-md w-full lux-glass-strong rounded-2xl p-8 text-center lux-fade-in">
        <div class="w-16 h-16 mx-auto mb-6 rounded-full bg-red-950/50 flex items-center justify-center">
          <svg class="w-8 h-8 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
        <h2 class="font-display text-2xl text-lux-cream mb-3">Connection Error</h2>
        <p class="text-lux-silver mb-6">{{ error }}</p>
        <button
          @click="handleRetry"
          class="px-8 py-3 rounded-lg lux-btn-gold text-sm tracking-widest"
        >
          Retry
        </button>
      </div>
    </div>

    <!-- Main Content -->
    <div v-else-if="auction" class="relative z-10 min-h-screen pb-28 md:pb-8">
      <!-- Header -->
      <header class="sticky top-0 z-40 header-glass border-b border-lux-gold/20">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex items-center justify-between h-16 sm:h-20">
            <!-- Left: Auction Info -->
            <div class="flex items-center gap-4 min-w-0">
              <div class="hidden sm:flex w-12 h-12 rounded-xl bg-gradient-to-br from-lux-gold/20 to-lux-gold/5 border border-lux-gold/40 flex-shrink-0 items-center justify-center shadow-lg shadow-lux-gold/10">
                <svg class="w-6 h-6 text-lux-gold" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
              </div>
              <div class="min-w-0">
                <h1 class="font-display text-lg sm:text-xl text-lux-cream truncate font-medium">{{ auction.title }}</h1>
                <p v-if="auction.description" class="text-xs text-lux-silver/60 truncate hidden sm:block mt-0.5">{{ auction.description }}</p>
              </div>
            </div>

            <!-- Right: Status Indicators -->
            <div class="flex items-center gap-3 sm:gap-4">
              <!-- Available Points -->
              <div class="points-badge flex items-center gap-2 px-4 py-2 sm:px-5 sm:py-2.5 rounded-xl">
                <svg class="w-5 h-5 text-lux-gold" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span class="text-base sm:text-lg font-bold lux-text-gold tabular-nums">{{ formatNumber(points.available) }}</span>
                <span class="text-xs text-lux-gold/60 hidden sm:inline font-medium">pts</span>
              </div>

              <!-- Participants (hidden on mobile) -->
              <div class="hidden lg:flex items-center gap-2 px-3 py-2 rounded-lg bg-lux-noir-light/50 border border-lux-noir-soft text-sm text-lux-silver">
                <svg class="w-4 h-4 text-lux-silver/70" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M9 6a3 3 0 11-6 0 3 3 0 016 0zM17 6a3 3 0 11-6 0 3 3 0 016 0zM12.93 17c.046-.327.07-.66.07-1a6.97 6.97 0 00-1.5-4.33A5 5 0 0119 16v1h-6.07zM6 11a5 5 0 015 5v1H1v-1a5 5 0 015-5z" />
                </svg>
                <span class="tabular-nums">{{ participantCount }}</span>
              </div>

              <!-- Connection Status -->
              <div
                :class="[
                  'flex items-center gap-2 px-3 py-2 rounded-lg transition-all duration-300',
                  wsConnected ? 'bg-emerald-500/10 border border-emerald-500/30' :
                  wsReconnecting ? 'bg-amber-500/10 border border-amber-500/30' :
                  'bg-red-500/10 border border-red-500/30'
                ]"
              >
                <div
                  :class="[
                    'w-2 h-2 rounded-full transition-colors',
                    wsConnected ? 'bg-emerald-400 shadow-[0_0_8px_rgba(52,211,153,0.6)]' :
                    wsReconnecting ? 'bg-amber-400 animate-pulse' : 'bg-red-400'
                  ]"
                ></div>
                <span
                  :class="[
                    'text-xs font-medium hidden sm:inline',
                    wsConnected ? 'text-emerald-400' : wsReconnecting ? 'text-amber-400' : 'text-red-400'
                  ]"
                >
                  {{ wsConnected ? 'Live' : wsReconnecting ? '再接続中...' : '切断' }}
                </span>
                <span v-if="wsReconnecting" class="text-xs text-amber-400/60 hidden sm:inline">
                  ({{ reconnectAttempt }}/{{ maxReconnectAttempts }})
                </span>
              </div>
            </div>
          </div>
        </div>
      </header>

      <!-- Main Layout -->
      <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 space-y-8">
        <!-- Bid Panel + History Grid -->
        <div v-if="currentItem" class="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <!-- Bid Panel (2 columns) -->
          <div class="lg:col-span-2 lux-fade-in">
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

          <!-- Bid History (1 column) -->
          <div class="lux-fade-in lux-delay-1">
            <BidderBidHistory
              :bids="bids"
              :current-bidder-id="currentBidderId"
            />
          </div>
        </div>

        <!-- Item Tabs -->
        <div class="lux-fade-in lux-delay-2">
          <ItemTabs
            :items="items"
            :selected-item-id="currentItem?.id"
            @select="handleSelectItem"
          />
        </div>
      </main>
    </div>

    <!-- Winning Modal -->
    <Teleport to="body">
      <Transition
        enter-active-class="transition-all duration-500 ease-out"
        enter-from-class="opacity-0"
        enter-to-class="opacity-100"
        leave-active-class="transition-all duration-300 ease-in"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <div v-if="showWinningModal" class="fixed inset-0 z-[100] flex items-center justify-center p-4">
          <!-- Backdrop -->
          <div class="absolute inset-0 bg-lux-noir/90 backdrop-blur-md" @click="closeWinningModal"></div>

          <!-- Modal Content -->
          <div class="relative lux-glass-strong rounded-3xl p-8 sm:p-12 max-w-lg w-full text-center lux-scale-in shadow-2xl">
            <!-- Gold Glow Background -->
            <div class="absolute inset-0 rounded-3xl overflow-hidden">
              <div class="absolute inset-0 bg-gradient-to-br from-lux-gold/10 via-transparent to-lux-gold/5"></div>
            </div>

            <!-- Corner Decorations -->
            <div class="absolute top-4 left-4 w-8 h-8 border-l-2 border-t-2 border-lux-gold/50 rounded-tl-lg"></div>
            <div class="absolute top-4 right-4 w-8 h-8 border-r-2 border-t-2 border-lux-gold/50 rounded-tr-lg"></div>
            <div class="absolute bottom-4 left-4 w-8 h-8 border-l-2 border-b-2 border-lux-gold/50 rounded-bl-lg"></div>
            <div class="absolute bottom-4 right-4 w-8 h-8 border-r-2 border-b-2 border-lux-gold/50 rounded-br-lg"></div>

            <div class="relative z-10">
              <!-- Trophy Icon -->
              <div class="mx-auto w-24 h-24 rounded-full bg-gradient-to-br from-lux-gold to-lux-gold-dark flex items-center justify-center mb-8 shadow-lg shadow-lux-gold/30">
                <svg class="w-12 h-12 text-lux-noir" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2L15.09 8.26L22 9.27L17 14.14L18.18 21.02L12 17.77L5.82 21.02L7 14.14L2 9.27L8.91 8.26L12 2Z" />
                </svg>
              </div>

              <h3 class="font-display text-3xl sm:text-4xl text-lux-cream mb-4 tracking-wide">
                Congratulations!
              </h3>
              <p class="text-lg text-lux-gold mb-2">落札おめでとうございます</p>
              <p class="text-lux-silver mb-8">
                「{{ winningItemName }}」を落札しました
              </p>

              <button
                @click="closeWinningModal"
                class="px-12 py-4 rounded-xl lux-btn-gold text-sm tracking-widest"
              >
                Continue
              </button>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- WebSocket Disconnected Overlay -->
    <Teleport to="body">
      <Transition
        enter-active-class="transition-all duration-300 ease-out"
        enter-from-class="opacity-0"
        enter-to-class="opacity-100"
        leave-active-class="transition-all duration-200 ease-in"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <div
          v-if="!wsConnected && !wsReconnecting && auction"
          class="fixed inset-0 z-[100] flex items-center justify-center p-4"
        >
          <div class="absolute inset-0 bg-lux-noir/95 backdrop-blur-md"></div>

          <div class="relative lux-glass-strong rounded-2xl p-8 max-w-md w-full text-center lux-scale-in">
            <div class="w-20 h-20 mx-auto mb-6 rounded-full bg-red-950/50 border border-red-500/30 flex items-center justify-center">
              <svg class="w-10 h-10 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M18.364 5.636a9 9 0 010 12.728m0 0l-2.829-2.829m2.829 2.829L21 21M15.536 8.464a5 5 0 010 7.072m0 0l-2.829-2.829m-4.243 2.829a5 5 0 010-7.072m0 0L8.464 8.464m-2.829 2.829a9 9 0 010-12.728m0 12.728L3 21" />
              </svg>
            </div>

            <h3 class="font-display text-2xl text-lux-cream mb-3">Connection Lost</h3>
            <p class="text-lux-silver mb-6">
              リアルタイム更新が停止しています。<br />
              入札を続けるには再接続してください。
            </p>

            <button
              @click="handleReconnect"
              class="w-full py-4 rounded-xl lux-btn-gold text-sm tracking-widest flex items-center justify-center gap-2"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              Reconnect
            </button>
          </div>
        </div>
      </Transition>
    </Teleport>
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

// Current bidder ID
const currentBidderId = computed(() => {
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

// Apply theme class to body
onMounted(async () => {
  document.body.classList.add('theme-luxury')

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
  document.body.classList.remove('theme-luxury')
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
/* Header Glass Effect */
.header-glass {
  background: linear-gradient(180deg, rgba(10, 10, 10, 0.95) 0%, rgba(10, 10, 10, 0.9) 100%);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  box-shadow:
    0 4px 30px rgba(0, 0, 0, 0.3),
    inset 0 -1px 0 rgba(212, 175, 55, 0.1);
}

/* Points Badge */
.points-badge {
  background: linear-gradient(135deg, rgba(20, 20, 20, 0.9) 0%, rgba(15, 15, 15, 0.95) 100%);
  border: 1px solid rgba(212, 175, 55, 0.4);
  box-shadow:
    0 0 20px rgba(212, 175, 55, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.05);
}
</style>
