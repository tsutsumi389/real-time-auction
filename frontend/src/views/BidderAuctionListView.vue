<template>
  <div class="bidder-auction-list-container">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- ヘッダー -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">
          オークション一覧
        </h1>
        <p class="text-gray-600">
          開催中および終了したオークションを閲覧できます
        </p>
      </div>

      <!-- 検索・フィルタエリア（Phase 4で実装） -->
      <div class="mb-6 space-y-4">
        <!-- 検索バー -->
        <div class="flex gap-4">
          <input
            v-model="searchKeyword"
            type="text"
            placeholder="オークションタイトルで検索..."
            class="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            @keyup.enter="handleSearch"
          />
          <button
            @click="handleSearch"
            class="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
          >
            検索
          </button>
        </div>

        <!-- フィルタ -->
        <div class="flex gap-4 items-center">
          <span class="text-sm font-medium text-gray-700">表示:</span>
          <button
            v-for="status in statusOptions"
            :key="status.value"
            @click="handleStatusFilter(status.value)"
            :class="[
              'px-4 py-2 rounded-lg text-sm font-medium transition-colors',
              auctionStore.filters.status === status.value
                ? 'bg-blue-600 text-white'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            ]"
          >
            {{ status.label }}
          </button>
        </div>

        <!-- ソート -->
        <div class="flex gap-4 items-center">
          <span class="text-sm font-medium text-gray-700">並び替え:</span>
          <select
            v-model="selectedSort"
            @change="handleSortChange"
            class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          >
            <option v-for="sort in sortOptions" :key="sort.value" :value="sort.value">
              {{ sort.label }}
            </option>
          </select>
        </div>
      </div>

      <!-- エラー表示 -->
      <div v-if="auctionStore.error" class="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
        <p class="text-red-800">{{ auctionStore.error }}</p>
        <button
          @click="auctionStore.clearError"
          class="mt-2 text-sm text-red-600 hover:text-red-800 underline"
        >
          閉じる
        </button>
      </div>

      <!-- ローディング（初回） -->
      <div v-if="auctionStore.loading" class="flex justify-center items-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>

      <!-- オークションリスト -->
      <div v-else>
        <!-- オークションカード（Phase 4で実装） -->
        <div v-if="auctionStore.auctions.length > 0" class="space-y-6">
          <!-- 簡易カード表示 -->
          <div
            v-for="auction in auctionStore.auctions"
            :key="auction.id"
            class="bg-white border border-gray-200 rounded-lg p-6 hover:shadow-lg transition-shadow"
          >
            <div class="flex justify-between items-start mb-4">
              <div class="flex-1">
                <h3 class="text-xl font-semibold text-gray-900 mb-2">
                  {{ auction.title }}
                </h3>
                <p class="text-gray-600 text-sm line-clamp-2">
                  {{ auction.description }}
                </p>
              </div>
              <span
                :class="[
                  'px-3 py-1 rounded-full text-sm font-medium',
                  getStatusClass(auction.status)
                ]"
              >
                {{ getStatusLabel(auction.status) }}
              </span>
            </div>

            <div class="flex items-center justify-between text-sm text-gray-500">
              <div class="flex gap-4">
                <span>出品物: {{ auction.item_count }}点</span>
                <span v-if="auction.started_at">
                  開始: {{ formatDate(auction.started_at) }}
                </span>
              </div>
              <button
                class="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm"
              >
                詳細を見る
              </button>
            </div>
          </div>

          <!-- 無限スクロールトリガー（Phase 4で実装） -->
          <div
            v-if="auctionStore.pagination.hasMore"
            ref="loadMoreTrigger"
            class="py-8 flex justify-center"
          >
            <div v-if="auctionStore.loadingMore" class="flex items-center gap-2 text-gray-600">
              <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
              <span>読み込み中...</span>
            </div>
            <div v-else class="text-gray-500">
              スクロールして続きを読み込む
            </div>
          </div>

          <!-- 読み込み完了 -->
          <div v-else class="py-8 text-center text-gray-500">
            すべてのオークションを表示しました（全{{ auctionStore.pagination.total }}件）
          </div>
        </div>

        <!-- 空の状態 -->
        <div v-else class="py-20 text-center">
          <p class="text-gray-500 text-lg">
            {{ emptyMessage }}
          </p>
          <button
            v-if="hasActiveFilters"
            @click="handleResetFilters"
            class="mt-4 px-6 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors"
          >
            フィルタをリセット
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useBidderAuctionStore } from '@/stores/bidderAuction'

const auctionStore = useBidderAuctionStore()

// フィルタ・検索の状態
const searchKeyword = ref('')
const selectedSort = ref('started_at_desc')

// ステータスオプション
const statusOptions = [
  { value: 'active', label: '開催中' },
  { value: 'ended', label: '終了' },
  { value: 'cancelled', label: '中止' },
]

// ソートオプション
const sortOptions = [
  { value: 'started_at_desc', label: '開始日時が新しい順' },
  { value: 'started_at_asc', label: '開始日時が古い順' },
  { value: 'updated_at_desc', label: '更新日時が新しい順' },
  { value: 'updated_at_asc', label: '更新日時が古い順' },
]

// 無限スクロールトリガー要素
const loadMoreTrigger = ref(null)
let observer = null

// フィルタが適用されているか
const hasActiveFilters = computed(() => {
  return auctionStore.filters.keyword !== '' ||
         auctionStore.filters.status !== 'active' ||
         auctionStore.filters.sort !== 'started_at_desc'
})

// 空の状態のメッセージ
const emptyMessage = computed(() => {
  if (auctionStore.filters.keyword) {
    return `「${auctionStore.filters.keyword}」に一致するオークションが見つかりませんでした`
  }
  if (auctionStore.filters.status === 'ended') {
    return '終了したオークションはまだありません'
  }
  if (auctionStore.filters.status === 'cancelled') {
    return '中止されたオークションはありません'
  }
  return '開催中のオークションはまだありません'
})

// イベントハンドラ
const handleSearch = () => {
  auctionStore.searchByKeyword(searchKeyword.value)
}

const handleStatusFilter = (status) => {
  auctionStore.filterByStatus(status)
}

const handleSortChange = () => {
  auctionStore.changeSort(selectedSort.value)
}

const handleResetFilters = () => {
  searchKeyword.value = ''
  selectedSort.value = 'started_at_desc'
  auctionStore.resetFilters()
}

// ステータスのスタイル
const getStatusClass = (status) => {
  switch (status) {
    case 'active':
      return 'bg-green-100 text-green-800'
    case 'ended':
      return 'bg-gray-100 text-gray-800'
    case 'cancelled':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const getStatusLabel = (status) => {
  switch (status) {
    case 'active':
      return '開催中'
    case 'ended':
      return '終了'
    case 'cancelled':
      return '中止'
    case 'pending':
      return '準備中'
    default:
      return status
  }
}

// 日付フォーマット
const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  }).format(date)
}

// Intersection Observerのセットアップ（無限スクロール）
const setupIntersectionObserver = () => {
  observer = new IntersectionObserver(
    (entries) => {
      const entry = entries[0]
      // トリガー要素が表示されたら追加読み込み
      if (entry.isIntersecting && auctionStore.pagination.hasMore && !auctionStore.loadingMore) {
        auctionStore.loadMoreAuctions()
      }
    },
    {
      rootMargin: '100px', // 100px手前で読み込み開始
      threshold: 0.1,
    }
  )

  // トリガー要素を監視
  if (loadMoreTrigger.value) {
    observer.observe(loadMoreTrigger.value)
  }
}

// ライフサイクル
onMounted(async () => {
  // 初回データ取得
  await auctionStore.fetchAuctionList()

  // Intersection Observerのセットアップ（次のティックで実行）
  setTimeout(() => {
    if (loadMoreTrigger.value) {
      setupIntersectionObserver()
    }
  }, 100)
})

onUnmounted(() => {
  // Observerのクリーンアップ
  if (observer) {
    observer.disconnect()
  }

  // ストアのリセット
  auctionStore.reset()
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
