<template>
  <div class="bidder-auction-detail-container">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
      <!-- ナビゲーションヘッダー -->
      <div class="mb-6 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <button
          @click="handleBackToList"
          class="inline-flex items-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          aria-label="オークション一覧に戻る"
        >
          <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
          </svg>
          一覧に戻る
        </button>

        <button
          v-if="auction && auction.status === 'active'"
          @click="handleGoToLive"
          class="inline-flex items-center px-6 py-2 text-sm font-medium text-white bg-green-600 rounded-lg hover:bg-green-700 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
          aria-label="ライブ入札画面へ移動"
        >
          ライブ入札へ
          <svg class="h-5 w-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6"></path>
          </svg>
        </button>
      </div>

      <!-- エラー表示 -->
      <Alert
        v-if="errorState.hasError"
        variant="destructive"
        class="mb-6"
        role="alert"
        aria-live="polite"
      >
        <div class="flex items-start gap-3">
          <!-- エラーアイコン -->
          <div class="flex-shrink-0">
            <svg
              v-if="errorState.type === 'notFound'"
              class="h-5 w-5 text-red-600"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            <svg
              v-else-if="errorState.type === 'server'"
              class="h-5 w-5 text-red-600"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
              aria-hidden="true"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
            </svg>
            <svg
              v-else
              class="h-5 w-5 text-red-600"
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
            <h3 class="text-sm font-semibold text-red-800">
              {{ errorState.title }}
            </h3>
            <p class="mt-1 text-sm text-red-700">
              {{ errorState.message }}
            </p>

            <!-- アクションボタン -->
            <div class="mt-4 flex flex-wrap gap-3">
              <button
                v-if="errorState.type !== 'notFound'"
                @click="handleRetry"
                :disabled="retrying"
                class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-red-600 rounded-lg hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
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
                class="inline-flex items-center px-4 py-2 text-sm font-medium text-red-700 bg-red-100 rounded-lg hover:bg-red-200 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
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
            class="flex-shrink-0 text-red-500 hover:text-red-700 transition-colors focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 rounded"
            aria-label="エラーメッセージを閉じる"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
      </Alert>

      <!-- ローディング状態 -->
      <div v-if="loading" class="space-y-6" aria-busy="true" aria-label="オークション情報を読み込み中">
        <!-- ローディングスピナー（中央表示） -->
        <div class="flex flex-col items-center justify-center py-12">
          <LoadingSpinner size="lg" text="オークション情報を読み込み中..." center />
        </div>

        <!-- オークション概要スケルトン -->
        <div class="bg-white border border-gray-200 rounded-lg p-6 animate-pulse" aria-hidden="true">
          <div class="flex flex-col sm:flex-row justify-between gap-4">
            <div class="flex-1">
              <div class="h-8 bg-gray-200 rounded w-3/4 mb-4"></div>
              <div class="h-4 bg-gray-200 rounded w-full mb-2"></div>
              <div class="h-4 bg-gray-200 rounded w-2/3"></div>
            </div>
            <div class="h-6 bg-gray-200 rounded w-24"></div>
          </div>
          <div class="mt-6 grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div class="h-5 bg-gray-200 rounded w-48"></div>
            <div class="h-5 bg-gray-200 rounded w-32"></div>
          </div>
        </div>

        <!-- アイテムグリッドスケルトン -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6" aria-hidden="true">
          <div v-for="i in 8" :key="i" class="bg-white border border-gray-200 rounded-lg overflow-hidden animate-pulse">
            <div class="h-48 bg-gray-200"></div>
            <div class="p-4 space-y-3">
              <div class="h-4 bg-gray-200 rounded w-1/4"></div>
              <div class="h-5 bg-gray-200 rounded w-3/4"></div>
              <div class="h-4 bg-gray-200 rounded w-1/2"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- オークション詳細表示 -->
      <div v-else-if="auction && !errorState.hasError">
        <!-- オークション概要セクション -->
        <div class="bg-white border border-gray-200 rounded-lg p-6 mb-8">
          <div class="flex flex-col sm:flex-row justify-between items-start gap-4">
            <div class="flex-1">
              <h1 class="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">
                {{ auction.title }}
              </h1>
              <p class="text-gray-600 text-sm sm:text-base mb-4">
                {{ auction.description || 'オークションの説明はありません' }}
              </p>
            </div>
            <div class="flex-shrink-0">
              <AuctionStatusBadge :status="auction.status" />
            </div>
          </div>

          <div class="mt-4 grid grid-cols-1 sm:grid-cols-2 gap-4">
            <div v-if="auction.started_at" class="flex items-center text-gray-600">
              <svg class="h-5 w-5 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
              </svg>
              <span class="text-sm sm:text-base">開始予定: {{ formatDate(auction.started_at) }}</span>
            </div>
            <div class="flex items-center text-gray-600">
              <svg class="h-5 w-5 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"></path>
              </svg>
              <span class="text-sm sm:text-base">出品アイテム: {{ itemCount }}点</span>
            </div>
          </div>
        </div>

        <!-- 出品アイテム一覧 -->
        <section aria-labelledby="items-heading">
          <div class="mb-6">
            <h2 id="items-heading" class="text-xl font-semibold text-gray-900 mb-4">
              出品アイテム一覧 ({{ itemCount }}点)
            </h2>
          </div>

          <!-- アイテムがない場合の空状態表示 -->
          <div v-if="itemCount === 0" class="bg-white border border-gray-200 rounded-lg p-12 text-center">
            <div class="flex flex-col items-center">
              <svg
                class="h-16 w-16 text-gray-300 mb-4"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
                aria-hidden="true"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"></path>
              </svg>
              <h3 class="text-lg font-medium text-gray-900 mb-2">
                出品アイテムがありません
              </h3>
              <p class="text-sm text-gray-500 max-w-sm">
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
      </div>
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
import Alert from '@/components/ui/Alert.vue'
import LoadingSpinner from '@/components/ui/LoadingSpinner.vue'

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
  fetchAuctionDetail()
  // Add browser back button listener
  window.addEventListener('popstate', handlePopState)
})

onUnmounted(() => {
  // Clean up browser back button listener
  window.removeEventListener('popstate', handlePopState)
})
</script>

<style scoped>
/* Styles are now handled by ItemCardGrid component */
</style>
