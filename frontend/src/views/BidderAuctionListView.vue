<template>
  <div class="bidder-auction-list-container">
    <!-- ヘッダー -->
    <BidderHeader />

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
      <!-- ページタイトル -->
      <div class="mb-6 sm:mb-8">
        <h1 class="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">
          オークション一覧
        </h1>
        <p class="text-sm sm:text-base text-gray-600">
          開催中および終了したオークションを閲覧できます
        </p>
      </div>

      <!-- 検索バー -->
      <div class="mb-4 sm:mb-6">
        <AuctionSearchBar
          v-model="searchKeyword"
          :loading="auctionStore.loading"
          @search="handleSearch"
          @clear="handleClearSearch"
        />
      </div>

      <!-- フィルタ -->
      <div class="mb-6">
        <AuctionFilters
          :current-status="auctionStore.filters.status"
          :current-sort="auctionStore.filters.sort"
          :status-options="statusOptions"
          :sort-options="sortOptions"
          @update:status="handleStatusFilter"
          @update:sort="handleSortChange"
        />
      </div>

      <!-- エラー表示 -->
      <div v-if="auctionStore.error" class="mb-6 p-4 bg-red-50 border border-red-200 rounded-lg">
        <div class="flex justify-between items-start">
          <p class="text-red-800 text-sm sm:text-base">{{ auctionStore.error }}</p>
          <button
            @click="auctionStore.clearError"
            class="text-red-600 hover:text-red-800 ml-4"
            aria-label="エラーを閉じる"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
      </div>

      <!-- オークションカードグリッド -->
      <AuctionCardGrid
        :auctions="auctionStore.auctions"
        :loading="auctionStore.loading"
        :empty-message="emptyMessage"
        @view-details="handleViewDetails"
        @join-auction="handleJoinAuction"
      >
        <!-- 空の状態のカスタマイズ -->
        <template #empty>
          <div class="text-center">
            <svg class="mx-auto h-12 w-12 sm:h-16 sm:w-16 mb-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"></path>
            </svg>
            <p class="text-gray-500 text-base sm:text-lg mb-4">
              {{ emptyMessage }}
            </p>
            <button
              v-if="hasActiveFilters"
              @click="handleResetFilters"
              class="px-4 py-2 sm:px-6 sm:py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors text-sm sm:text-base"
            >
              フィルタをリセット
            </button>
          </div>
        </template>
      </AuctionCardGrid>

      <!-- 無限スクロールトリガー -->
      <div
        v-if="auctionStore.auctions.length > 0 && auctionStore.pagination.hasMore"
        ref="loadMoreTrigger"
        class="py-6 sm:py-8 flex justify-center"
      >
        <LoadingSpinner
          v-if="auctionStore.loadingMore"
          size="md"
          text="読み込み中..."
          center
        />
        <div v-else class="text-gray-500 text-sm sm:text-base">
          スクロールして続きを読み込む
        </div>
      </div>

      <!-- 読み込み完了メッセージ -->
      <div
        v-if="auctionStore.auctions.length > 0 && !auctionStore.pagination.hasMore"
        class="py-6 sm:py-8 text-center text-gray-500 text-sm sm:text-base"
      >
        すべてのオークションを表示しました（全{{ auctionStore.pagination.total }}件）
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch } from 'vue'
import { useBidderAuctionStore } from '@/stores/bidderAuction'
import BidderHeader from '@/components/bidder/BidderHeader.vue'
import AuctionSearchBar from '@/components/bidder/AuctionSearchBar.vue'
import AuctionFilters from '@/components/bidder/AuctionFilters.vue'
import AuctionCardGrid from '@/components/bidder/AuctionCardGrid.vue'
import LoadingSpinner from '@/components/ui/LoadingSpinner.vue'

const auctionStore = useBidderAuctionStore()

// 検索キーワード
const searchKeyword = ref('')

// ステータスオプション
const statusOptions = [
  { value: 'active', label: '開催中' },
  { value: 'ended', label: '終了' },
  { value: 'cancelled', label: '中止' }
]

// ソートオプション
const sortOptions = [
  { value: 'started_at_desc', label: '開始日時が新しい順' },
  { value: 'started_at_asc', label: '開始日時が古い順' },
  { value: 'updated_at_desc', label: '更新日時が新しい順' },
  { value: 'updated_at_asc', label: '更新日時が古い順' }
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
const handleSearch = (keyword) => {
  auctionStore.searchByKeyword(keyword)
}

const handleClearSearch = () => {
  searchKeyword.value = ''
  auctionStore.searchByKeyword('')
}

const handleStatusFilter = (status) => {
  auctionStore.filterByStatus(status)
}

const handleSortChange = (sort) => {
  auctionStore.changeSort(sort)
}

const handleResetFilters = () => {
  searchKeyword.value = ''
  auctionStore.resetFilters()
}

const handleViewDetails = (auction) => {
  console.log('View details:', auction)
  // TODO: オークション詳細画面への遷移を実装
}

const handleJoinAuction = (auction) => {
  console.log('Join auction:', auction)
  // TODO: オークション参加処理を実装
}

// Intersection Observerのセットアップ（無限スクロール）
const setupIntersectionObserver = () => {
  if (observer) {
    observer.disconnect()
  }

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
      threshold: 0.1
    }
  )

  // トリガー要素を監視
  if (loadMoreTrigger.value) {
    observer.observe(loadMoreTrigger.value)
  }
}

// has_moreが変更されたときにObserverを再設定
watch(() => auctionStore.pagination.hasMore, (hasMore) => {
  if (hasMore) {
    setTimeout(() => {
      setupIntersectionObserver()
    }, 100)
  }
})

// ライフサイクル
onMounted(async () => {
  // 初回データ取得
  await auctionStore.fetchAuctionList()

  // Intersection Observerのセットアップ（次のティックで実行）
  setTimeout(() => {
    setupIntersectionObserver()
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
.bidder-auction-list-container {
  min-height: 100vh;
  background-color: #f9fafb; /* bg-gray-50 */
}

/* レスポンシブデザイン対応 */
@media (max-width: 640px) {
  .bidder-auction-list-container {
    background-color: #ffffff;
  }
}
</style>
