<template>
  <div class="bidder-auction-detail-container">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
      <!-- ナビゲーションヘッダー -->
      <div class="mb-6 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
        <button
          @click="handleBackToList"
          class="inline-flex items-center px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
        >
          <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
          </svg>
          一覧に戻る
        </button>

        <button
          v-if="auction && auction.status === 'active'"
          @click="handleGoToLive"
          class="inline-flex items-center px-6 py-2 text-sm font-medium text-white bg-green-600 rounded-lg hover:bg-green-700 transition-colors"
        >
          ライブ入札へ
          <svg class="h-5 w-5 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6"></path>
          </svg>
        </button>
      </div>

      <!-- エラー表示 -->
      <div v-if="error" class="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
        <div class="flex justify-between items-start">
          <div>
            <p class="text-red-800 text-sm sm:text-base font-medium mb-2">{{ error }}</p>
            <button
              @click="fetchAuctionDetail"
              class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors text-sm"
            >
              再読み込み
            </button>
          </div>
          <button
            @click="error = null"
            class="text-red-600 hover:text-red-800 ml-4"
            aria-label="エラーを閉じる"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
      </div>

      <!-- ローディング状態 -->
      <div v-if="loading" class="space-y-6">
        <!-- オークション概要スケルトン -->
        <div class="bg-white border border-gray-200 rounded-lg p-6 animate-pulse">
          <div class="h-8 bg-gray-200 rounded w-1/2 mb-4"></div>
          <div class="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
          <div class="h-4 bg-gray-200 rounded w-2/3"></div>
        </div>

        <!-- アイテムグリッドスケルトン -->
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          <div v-for="i in 8" :key="i" class="bg-white border border-gray-200 rounded-lg overflow-hidden animate-pulse">
            <div class="h-48 bg-gray-200"></div>
            <div class="p-4 space-y-3">
              <div class="h-4 bg-gray-200 rounded w-3/4"></div>
              <div class="h-4 bg-gray-200 rounded w-1/2"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- オークション詳細表示 -->
      <div v-else-if="auction && !error">
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
              <svg class="h-5 w-5 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"></path>
              </svg>
              <span class="text-sm sm:text-base">出品アイテム: {{ itemCount }}点</span>
            </div>
          </div>
        </div>

        <!-- 出品アイテム一覧 -->
        <div class="mb-6">
          <h2 class="text-xl font-semibold text-gray-900 mb-4">
            出品アイテム一覧 ({{ itemCount }}点)
          </h2>
        </div>

        <!-- アイテムカードグリッド -->
        <div v-if="items.length > 0" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          <div
            v-for="item in items"
            :key="item.id"
            @click="handleItemClick(item)"
            class="bg-white border border-gray-200 rounded-lg overflow-hidden hover:shadow-lg transition-shadow cursor-pointer"
          >
            <!-- サムネイル画像 -->
            <div class="relative w-full h-48 bg-gray-100">
              <img
                v-if="getItemThumbnail(item)"
                :src="getItemThumbnail(item)"
                :alt="item.name"
                class="w-full h-full object-cover"
                @error="handleImageError"
              />
              <div v-else class="w-full h-full flex items-center justify-center text-gray-400">
                <svg class="h-12 w-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                </svg>
              </div>
            </div>

            <!-- カード情報 -->
            <div class="p-4">
              <div class="mb-2">
                <span class="inline-block px-2 py-1 text-xs font-semibold text-blue-800 bg-blue-100 rounded">
                  LOT {{ item.lot_number }}
                </span>
              </div>
              <h3 class="text-base font-semibold text-gray-900 mb-2 line-clamp-2">
                {{ item.name }}
              </h3>
              <p class="text-lg font-bold text-gray-900">
                ¥{{ formatPrice(item.starting_price) }}
              </p>
            </div>
          </div>
        </div>

        <!-- アイテムがない場合 -->
        <div v-else class="text-center py-12">
          <svg class="mx-auto h-12 w-12 text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"></path>
          </svg>
          <p class="text-gray-500 text-lg">
            このオークションにはまだアイテムが登録されていません
          </p>
        </div>
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
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getAuctionDetail } from '@/services/bidderAuctionApi'
import AuctionStatusBadge from '@/components/bidder/AuctionStatusBadge.vue'
import ItemDetailModal from '@/components/bidder/ItemDetailModal.vue'

const router = useRouter()
const route = useRoute()

// State
const auction = ref(null)
const items = ref([])
const loading = ref(true)
const error = ref(null)
const selectedItem = ref(null)
const isModalOpen = ref(false)

// Computed
const itemCount = computed(() => items.value.length)

// Methods
const fetchAuctionDetail = async () => {
  loading.value = true
  error.value = null

  try {
    const auctionId = route.params.id
    const data = await getAuctionDetail(auctionId)

    auction.value = data
    items.value = data.items || []
  } catch (err) {
    console.error('Failed to fetch auction detail:', err)

    if (err.response?.status === 404) {
      error.value = 'オークションが見つかりませんでした'
    } else if (err.response?.status >= 500) {
      error.value = 'サーバーエラーが発生しました'
    } else {
      error.value = '通信エラーが発生しました'
    }
  } finally {
    loading.value = false
  }
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

const formatPrice = (price) => {
  if (!price) return '0'
  return new Intl.NumberFormat('ja-JP').format(price)
}

const getItemThumbnail = (item) => {
  if (!item.media || item.media.length === 0) return null
  return item.media[0].thumbnail_url || item.media[0].url
}

const handleImageError = (event) => {
  event.target.style.display = 'none'
}

const handleBackToList = () => {
  router.push({ name: 'bidder-auction-list' })
}

const handleGoToLive = () => {
  router.push({ name: 'bidder-auction-live', params: { id: auction.value.id } })
}

const handleItemClick = (item) => {
  selectedItem.value = item
  isModalOpen.value = true
}

const handleCloseModal = () => {
  isModalOpen.value = false
  selectedItem.value = null
}

// Lifecycle
onMounted(() => {
  fetchAuctionDetail()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
