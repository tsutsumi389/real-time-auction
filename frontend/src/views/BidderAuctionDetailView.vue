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

    <!-- Skip link for keyboard users -->
    <a
      href="#main-content"
      class="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4 focus:z-50 focus:px-4 focus:py-2 focus:bg-lux-gold focus:text-lux-noir focus:rounded-lg focus:outline-none"
    >
      メインコンテンツへスキップ
    </a>

    <!-- Main Content -->
    <div class="relative z-10">
      <main
        id="main-content"
        class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8"
        aria-label="オークション詳細"
      >
        <!-- ナビゲーションヘッダー -->
        <nav
          class="mb-6 sm:mb-8 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 lux-fade-in"
          aria-label="ページナビゲーション"
        >
          <button
            @click="handleBackToList"
            class="group inline-flex items-center px-4 py-2.5 text-sm font-medium text-lux-cream lux-glass border border-lux-gold/30 rounded-xl shadow-sm hover:bg-lux-gold/10 hover:border-lux-gold/50 hover:shadow transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-lux-gold focus:ring-offset-lux-noir"
            aria-label="オークション一覧に戻る"
          >
            <svg class="h-5 w-5 mr-2 text-lux-gold/60 group-hover:text-lux-gold transition-colors" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
            </svg>
            一覧に戻る
          </button>

          <button
            v-if="auction && auction.status === 'active'"
            @click="handleGoToLive"
            class="group inline-flex items-center px-6 py-3 text-sm font-semibold text-lux-noir bg-gradient-to-r from-lux-gold via-yellow-400 to-lux-gold rounded-xl shadow-lg shadow-lux-gold/30 hover:shadow-xl hover:shadow-lux-gold/50 transform hover:scale-105 transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-lux-gold focus:ring-offset-lux-noir border-2 border-lux-gold/50 hover:border-lux-gold animate-shimmer"
            aria-label="ライブ入札画面へ移動"
          >
            <span class="relative flex h-2.5 w-2.5 mr-2.5">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-500 opacity-75"></span>
              <span class="relative inline-flex rounded-full h-2.5 w-2.5 bg-emerald-600"></span>
            </span>
            ライブ入札へ
            <svg class="h-5 w-5 ml-2.5 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true" stroke-width="2.5">
              <path stroke-linecap="round" stroke-linejoin="round" d="M13 7l5 5m0 0l-5 5m5-5H6"></path>
            </svg>
          </button>
        </nav>

        <!-- エラー表示 -->
        <Transition
          enter-active-class="transition-all duration-300 ease-out"
          enter-from-class="opacity-0 -translate-y-2"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition-all duration-200 ease-in"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <div v-if="errorState.hasError" class="mb-6 lux-glass-strong rounded-xl p-4 border border-red-500/30 bg-red-950/20" role="alert" aria-live="polite">
            <div class="flex items-start gap-3">
              <!-- エラーアイコン -->
              <div class="flex-shrink-0 w-10 h-10 rounded-lg bg-red-500/10 flex items-center justify-center">
                <svg
                  v-if="errorState.type === 'notFound'"
                  class="h-5 w-5 text-red-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                <svg
                  v-else-if="errorState.type === 'server'"
                  class="h-5 w-5 text-red-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
                </svg>
                <svg
                  v-else
                  class="h-5 w-5 text-red-400"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                  aria-hidden="true"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
              </div>

              <!-- エラーメッセージ -->
              <div class="flex-1 min-w-0">
                <h3 class="font-display text-base text-red-400 mb-1">
                  {{ errorState.title }}
                </h3>
                <p class="text-sm text-red-300/80">
                  {{ errorState.message }}
                </p>

                <!-- アクションボタン -->
                <div class="mt-4 flex flex-wrap gap-3">
                  <button
                    v-if="errorState.type !== 'notFound'"
                    @click="handleRetry"
                    :disabled="retrying"
                    class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-red-500/80 rounded-lg hover:bg-red-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 focus:ring-offset-lux-noir"
                  >
                    <svg
                      v-if="retrying"
                      class="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                      fill="none"
                      viewBox="0 0 24 24"
                      aria-hidden="true"
                    >
                      <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                      <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    <svg
                      v-else
                      class="-ml-1 mr-2 h-4 w-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                      aria-hidden="true"
                    >
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                    </svg>
                    {{ retrying ? '再読み込み中...' : '再読み込み' }}
                  </button>
                  <button
                    @click="handleBackToList"
                    class="inline-flex items-center px-4 py-2 text-sm font-medium text-red-300 bg-red-500/10 rounded-lg hover:bg-red-500/20 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 focus:ring-offset-lux-noir"
                  >
                    <svg class="-ml-1 mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
                    </svg>
                    一覧に戻る
                  </button>
                </div>
              </div>

              <!-- 閉じるボタン（404エラー以外の場合のみ） -->
              <button
                v-if="errorState.type !== 'notFound'"
                @click="clearError"
                class="flex-shrink-0 p-2 rounded-lg hover:bg-red-500/10 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 focus:ring-offset-lux-noir"
                aria-label="エラーメッセージを閉じる"
              >
                <svg class="h-4 w-4 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                </svg>
              </button>
            </div>
          </div>
        </Transition>

        <!-- ローディング状態 -->
        <div v-if="loading" class="space-y-6 sm:space-y-8 lux-fade-in" aria-busy="true" aria-label="オークション情報を読み込み中">
          <!-- ローディングスピナー（中央表示） -->
          <div class="flex flex-col items-center justify-center py-12">
            <div class="relative w-16 h-16 mb-4">
              <div class="absolute inset-0 rounded-full border-2 border-lux-gold/20"></div>
              <div class="absolute inset-0 rounded-full border-2 border-transparent border-t-lux-gold animate-spin"></div>
            </div>
            <span class="text-lux-silver text-sm font-medium">オークション情報を読み込み中...</span>
          </div>

          <!-- オークション概要スケルトン -->
          <div class="lux-glass border border-lux-gold/20 rounded-2xl p-6 sm:p-8 animate-pulse" aria-hidden="true">
            <div class="flex flex-col sm:flex-row justify-between gap-4 sm:gap-6">
              <div class="flex-1">
                <div class="h-9 bg-lux-gold/10 rounded-lg w-3/4 mb-4"></div>
                <div class="h-4 bg-lux-gold/10 rounded w-full mb-2"></div>
                <div class="h-4 bg-lux-gold/10 rounded w-2/3"></div>
              </div>
              <div class="h-7 bg-lux-gold/10 rounded-full w-24"></div>
            </div>
            <div class="mt-6 pt-6 border-t border-lux-gold/20 grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div class="h-16 bg-lux-gold/5 rounded-xl"></div>
              <div class="h-16 bg-lux-gold/5 rounded-xl"></div>
            </div>
          </div>

          <!-- アイテムグリッドスケルトン -->
          <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6" aria-hidden="true">
            <div v-for="i in 8" :key="i" class="lux-glass border border-lux-gold/20 rounded-2xl overflow-hidden animate-pulse">
              <div class="aspect-[4/3] bg-lux-gold/10"></div>
              <div class="p-4 space-y-3">
                <div class="h-5 bg-lux-gold/10 rounded-full w-16"></div>
                <div class="h-5 bg-lux-gold/10 rounded w-3/4"></div>
                <div class="h-4 bg-lux-gold/10 rounded w-1/2"></div>
              </div>
            </div>
          </div>
        </div>

        <!-- オークション詳細表示 -->
        <article v-else-if="auction && !errorState.hasError" class="auction-detail-content lux-fade-in" aria-label="オークション詳細情報">
          <!-- オークション概要セクション -->
          <header class="lux-glass border border-lux-gold/20 rounded-2xl shadow-sm hover:shadow-lg hover:shadow-lux-gold/5 transition-all duration-300 p-6 sm:p-8 mb-8">
            <div class="flex flex-col sm:flex-row justify-between items-start gap-4 sm:gap-6">
              <div class="flex-1 min-w-0">
                <h1 class="font-display text-2xl sm:text-3xl lg:text-4xl font-bold text-lux-cream mb-3 tracking-tight">
                  {{ auction.title }}
                </h1>
                <p class="text-lux-silver/80 text-sm sm:text-base leading-relaxed">
                  {{ auction.description || 'オークションの説明はありません' }}
                </p>
              </div>
              <div class="flex-shrink-0">
                <AuctionStatusBadge :status="auction.status" />
              </div>
            </div>

            <div class="mt-6 pt-6 border-t border-lux-gold/20">
              <dl class="grid grid-cols-1 sm:grid-cols-2 gap-4 sm:gap-6">
                <div v-if="auction.started_at" class="flex items-center p-3 bg-lux-gold/5 rounded-xl border border-lux-gold/10">
                  <div class="flex-shrink-0 p-2 bg-lux-gold/10 rounded-lg mr-3">
                    <svg class="h-5 w-5 text-lux-gold" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                    </svg>
                  </div>
                  <div>
                    <dt class="text-xs text-lux-silver/60 font-medium uppercase tracking-wide">開始予定日時</dt>
                    <dd class="text-sm sm:text-base font-semibold text-lux-cream mt-0.5">
                      <time :datetime="auction.started_at">{{ formatDate(auction.started_at) }}</time>
                    </dd>
                  </div>
                </div>
                <div class="flex items-center p-3 bg-lux-gold/5 rounded-xl border border-lux-gold/10">
                  <div class="flex-shrink-0 p-2 bg-lux-gold/10 rounded-lg mr-3">
                    <svg class="h-5 w-5 text-lux-gold" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"></path>
                    </svg>
                  </div>
                  <div>
                    <dt class="text-xs text-lux-silver/60 font-medium uppercase tracking-wide">出品アイテム数</dt>
                    <dd class="text-sm sm:text-base font-semibold text-lux-cream mt-0.5">{{ itemCount }}点</dd>
                  </div>
                </div>
              </dl>
            </div>
          </header>

          <!-- 出品アイテム一覧 -->
          <section aria-labelledby="items-heading" class="lux-fade-in lux-delay-1">
            <div class="mb-6">
              <h2 id="items-heading" class="font-display text-xl sm:text-2xl font-bold text-lux-cream flex items-center">
                <span>出品アイテム一覧</span>
                <span class="ml-3 px-3 py-1 text-sm font-medium bg-lux-gold/10 border border-lux-gold/30 text-lux-gold rounded-full">{{ itemCount }}点</span>
              </h2>
            </div>

            <!-- アイテムがない場合の空状態表示 -->
            <div v-if="itemCount === 0" class="lux-glass border border-lux-gold/20 border-dashed rounded-2xl p-12 text-center" role="status">
              <div class="flex flex-col items-center">
                <div class="p-4 bg-lux-gold/10 rounded-full mb-4">
                  <svg
                    class="h-12 w-12 text-lux-gold/40"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                    aria-hidden="true"
                  >
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"></path>
                  </svg>
                </div>
                <h3 class="font-display text-lg font-semibold text-lux-cream mb-2">
                  出品アイテムがありません
                </h3>
                <p class="text-sm text-lux-silver/60 max-w-sm leading-relaxed">
                  このオークションにはまだアイテムが登録されていません。アイテムが追加されるまでお待ちください。
                </p>
              </div>
            </div>

            <!-- アイテムカードグリッド -->
            <ItemCardGrid
              v-else
              :items="items"
              @item-click="handleItemClick"
            />
          </section>
        </article>
      </main>
    </div>

    <!-- アイテム詳細モーダル -->
    <ItemDetailModal
      v-if="selectedItem"
      :item="selectedItem"
      :open="isModalOpen"
      @close="handleCloseModal"
    />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getAuctionDetail } from '@/services/bidderAuctionApi'
import AuctionStatusBadge from '@/components/bidder/AuctionStatusBadge.vue'
import ItemCardGrid from '@/components/bidder/ItemCardGrid.vue'
import ItemDetailModal from '@/components/bidder/ItemDetailModal.vue'

const router = useRouter()
const route = useRoute()

// State
const auction = ref(null)
const items = ref([])
const loading = ref(true)
const retrying = ref(false)
const selectedItem = ref(null)
const isModalOpen = ref(false)

// Error state management
const errorState = reactive({
  hasError: false,
  type: null, // 'notFound', 'server', 'network'
  title: '',
  message: ''
})

// Computed
const itemCount = computed(() => items.value.length)

// Error handling helpers
const setError = (type, title, message) => {
  errorState.hasError = true
  errorState.type = type
  errorState.title = title
  errorState.message = message
}

const clearError = () => {
  errorState.hasError = false
  errorState.type = null
  errorState.title = ''
  errorState.message = ''
}

// Error messages configuration
const errorMessages = {
  notFound: {
    title: 'オークションが見つかりません',
    message: '指定されたオークションは存在しないか、削除された可能性があります。URLをご確認いただくか、オークション一覧から再度お探しください。'
  },
  server: {
    title: 'サーバーエラーが発生しました',
    message: 'サーバーで問題が発生しました。しばらく時間をおいてから再度お試しください。問題が解決しない場合は、管理者にお問い合わせください。'
  },
  network: {
    title: '通信エラーが発生しました',
    message: 'ネットワーク接続に問題があります。インターネット接続をご確認の上、再度お試しください。'
  },
  unknown: {
    title: '予期しないエラーが発生しました',
    message: '処理中にエラーが発生しました。再度お試しください。問題が解決しない場合は、管理者にお問い合わせください。'
  }
}

// Methods
const fetchAuctionDetail = async () => {
  loading.value = true
  clearError()

  try {
    const auctionId = route.params.id
    
    // Validate UUID format
    if (!isValidUUID(auctionId)) {
      setError('notFound', errorMessages.notFound.title, errorMessages.notFound.message)
      loading.value = false
      return
    }

    const data = await getAuctionDetail(auctionId)

    auction.value = data
    items.value = data.items || []
  } catch (err) {
    console.error('Failed to fetch auction detail:', err)
    handleApiError(err)
  } finally {
    loading.value = false
  }
}

const handleApiError = (err) => {
  const status = err.response?.status

  if (status === 404) {
    setError('notFound', errorMessages.notFound.title, errorMessages.notFound.message)
  } else if (status >= 500) {
    setError('server', errorMessages.server.title, errorMessages.server.message)
  } else if (err.code === 'ERR_NETWORK' || err.message === 'Network Error' || !navigator.onLine) {
    setError('network', errorMessages.network.title, errorMessages.network.message)
  } else {
    setError('unknown', errorMessages.unknown.title, errorMessages.unknown.message)
  }
}

const handleRetry = async () => {
  retrying.value = true
  clearError()

  try {
    await fetchAuctionDetail()
  } finally {
    retrying.value = false
  }
}

const isValidUUID = (str) => {
  const uuidRegex = /^[0-9a-f]{8}-[0-9a-f]{4}-[1-5][0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/i
  return uuidRegex.test(str)
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

const handleBackToList = () => {
  router.push({ name: 'bidder-auction-list' })
}

const handleGoToLive = () => {
  router.push({ name: 'bidder-auction-live', params: { id: auction.value.id } })
}

// Flag to track if modal was closed via back button
let closedViaBackButton = false

const handleItemClick = (item) => {
  selectedItem.value = item
  isModalOpen.value = true
  // Add modal state to browser history for back button support
  window.history.pushState({ modal: true, itemId: item.id }, '')
}

const handleCloseModal = () => {
  // If not closed via back button, we need to go back in history
  if (!closedViaBackButton && window.history.state?.modal) {
    window.history.back()
  }
  closedViaBackButton = false
  isModalOpen.value = false
  selectedItem.value = null
}

// Browser back button handler
const handlePopState = () => {
  // When back button is pressed and modal is open, close modal
  if (isModalOpen.value) {
    closedViaBackButton = true
    isModalOpen.value = false
    selectedItem.value = null
  }
}

// Lifecycle
onMounted(() => {
  document.body.classList.add('theme-luxury')
  fetchAuctionDetail()
  // Add browser back button listener
  window.addEventListener('popstate', handlePopState)
})

onUnmounted(() => {
  document.body.classList.remove('theme-luxury')
  // Clean up browser back button listener
  window.removeEventListener('popstate', handlePopState)
})
</script>

<style scoped>
/* Luxury color utilities */
.bg-lux-noir {
  background-color: hsl(0 0% 4%);
}

.bg-lux-gold\/3 {
  background-color: hsl(43 74% 49% / 0.03);
}

.bg-lux-gold\/2 {
  background-color: hsl(43 74% 49% / 0.02);
}

.text-lux-gold {
  color: hsl(43 74% 49%);
}

.text-lux-cream {
  color: hsl(45 30% 96%);
}

.text-lux-silver {
  color: hsl(220 10% 70%);
}

.border-lux-gold\/10 {
  border-color: hsl(43 74% 49% / 0.1);
}

.border-lux-gold\/20 {
  border-color: hsl(43 74% 49% / 0.2);
}

.border-lux-gold\/30 {
  border-color: hsl(43 74% 49% / 0.3);
}

.border-lux-gold\/50 {
  border-color: hsl(43 74% 49% / 0.5);
}

.bg-lux-gold\/5 {
  background-color: hsl(43 74% 49% / 0.05);
}

.bg-lux-gold\/10 {
  background-color: hsl(43 74% 49% / 0.1);
}

.text-lux-noir {
  color: hsl(0 0% 4%);
}

.from-lux-gold {
  --tw-gradient-from: hsl(43 74% 49%);
}

.focus\:ring-lux-gold:focus {
  --tw-ring-color: hsl(43 74% 49%);
}

.focus\:ring-offset-lux-noir:focus {
  --tw-ring-offset-color: hsl(0 0% 4%);
}

.ring-lux-gold {
  --tw-ring-color: hsl(43 74% 49%);
}

.ring-offset-lux-noir {
  --tw-ring-offset-color: hsl(0 0% 4%);
}

/* Glass morphism effect */
.lux-glass {
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.03) 0%,
    rgba(255, 255, 255, 0.01) 100%
  );
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
}

.lux-glass-strong {
  background: linear-gradient(
    135deg,
    rgba(255, 255, 255, 0.05) 0%,
    rgba(255, 255, 255, 0.02) 100%
  );
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
}

/* Noise texture */
.lux-noise {
  background-image: url("data:image/svg+xml,%3Csvg viewBox='0 0 400 400' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noiseFilter'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noiseFilter)' opacity='0.03'/%3E%3C/svg%3E");
}

/* Fade in animations */
.lux-fade-in {
  animation: lux-fade-in 0.6s ease-out;
}

.lux-delay-1 {
  animation-delay: 0.1s;
  opacity: 0;
  animation-fill-mode: forwards;
}

.lux-delay-2 {
  animation-delay: 0.2s;
  opacity: 0;
  animation-fill-mode: forwards;
}

@keyframes lux-fade-in {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Page entrance animation */
.auction-detail-content {
  animation: page-fade-in 0.4s ease-out;
}

@keyframes page-fade-in {
  from {
    opacity: 0;
    transform: translateY(8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Skeleton pulse animation enhancement */
.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

/* Font display */
.font-display {
  font-family: 'Cormorant Garamond', Georgia, serif;
}

/* Shimmer animation for live button */
@keyframes shimmer {
  0% {
    background-position: -200% center;
  }
  100% {
    background-position: 200% center;
  }
}

.animate-shimmer {
  background-size: 200% auto;
  animation: shimmer 3s linear infinite;
}

/* Respect reduced motion preference */
@media (prefers-reduced-motion: reduce) {
  .auction-detail-content,
  .lux-fade-in {
    animation: none;
    opacity: 1;
  }

  .animate-pulse,
  .animate-shimmer {
    animation: none;
  }
}
</style>
