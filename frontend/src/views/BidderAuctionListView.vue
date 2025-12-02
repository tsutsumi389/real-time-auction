<template>
  <div class="bidder-auction-list-container">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
      <!-- Page title -->
      <div class="mb-6 sm:mb-8">
        <h1 class="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">
          オークション一覧
        </h1>
        <p class="text-sm sm:text-base text-muted-foreground">
          開催中および終了したオークションを閲覧できます
        </p>
      </div>

      <!-- Search bar -->
      <div class="mb-4 sm:mb-6">
        <AuctionSearchBar
          v-model="searchKeyword"
          :loading="auctionStore.loading"
          @search="handleSearch"
          @clear="handleClearSearch"
        />
      </div>

      <!-- Filters -->
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

      <!-- Error display -->
      <Alert v-if="auctionStore.error" variant="destructive" class="mb-6">
        <AlertCircle class="h-4 w-4" />
        <AlertTitle>エラーが発生しました</AlertTitle>
        <AlertDescription class="flex justify-between items-start">
          <span>{{ auctionStore.error }}</span>
          <button
            @click="auctionStore.clearError"
            class="text-destructive hover:text-destructive/80 ml-4 p-1 rounded-full hover:bg-destructive/10 transition-colors"
            aria-label="エラーを閉じる"
          >
            <X class="h-4 w-4" />
          </button>
        </AlertDescription>
      </Alert>

      <!-- Auction card grid -->
      <AuctionCardGrid
        :auctions="auctionStore.auctions"
        :loading="auctionStore.loading"
        :empty-message="emptyMessage"
        @view-details="handleViewDetails"
        @join-auction="handleJoinAuction"
      >
        <!-- Custom empty state -->
        <template #empty>
          <div class="text-center animate-fade-in">
            <Inbox class="mx-auto h-12 w-12 sm:h-16 sm:w-16 mb-4 text-muted-foreground" :stroke-width="1.5" />
            <p class="text-muted-foreground text-base sm:text-lg mb-4">
              {{ emptyMessage }}
            </p>
            <Button
              v-if="hasActiveFilters"
              @click="handleResetFilters"
              variant="outline"
            >
              <RotateCcw class="mr-2 h-4 w-4" />
              フィルタをリセット
            </Button>
          </div>
        </template>
      </AuctionCardGrid>

      <!-- Infinite scroll trigger -->
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
        <div v-else class="text-muted-foreground text-sm sm:text-base flex items-center gap-2">
          <ChevronsDown class="h-4 w-4 animate-bounce" />
          スクロールして続きを読み込む
        </div>
      </div>

      <!-- Load complete message -->
      <div
        v-if="auctionStore.auctions.length > 0 && !auctionStore.pagination.hasMore"
        class="py-6 sm:py-8 text-center text-muted-foreground text-sm sm:text-base"
      >
        <CheckCircle2 class="inline-block h-5 w-5 mr-2 text-green-500" />
        すべてのオークションを表示しました（全{{ auctionStore.pagination.total }}件）
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { AlertCircle, X, Inbox, RotateCcw, ChevronsDown, CheckCircle2 } from 'lucide-vue-next'
import { useBidderAuctionStore } from '@/stores/bidderAuction'
import AuctionSearchBar from '@/components/bidder/AuctionSearchBar.vue'
import AuctionFilters from '@/components/bidder/AuctionFilters.vue'
import AuctionCardGrid from '@/components/bidder/AuctionCardGrid.vue'
import LoadingSpinner from '@/components/ui/LoadingSpinner.vue'
import Alert from '@/components/ui/Alert.vue'
import AlertTitle from '@/components/ui/AlertTitle.vue'
import AlertDescription from '@/components/ui/AlertDescription.vue'
import Button from '@/components/ui/Button.vue'

const router = useRouter()

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
  router.push({
    name: 'bidder-auction-detail',
    params: { id: auction.id }
  })
}

const handleJoinAuction = (auction) => {
  // Navigate to live auction screen
  router.push({
    name: 'bidder-auction-live',
    params: { id: auction.id }
  })
}

// Intersection Observer setup (infinite scroll)
const setupIntersectionObserver = () => {
  if (observer) {
    observer.disconnect()
  }

  observer = new IntersectionObserver(
    (entries) => {
      const entry = entries[0]
      // Load more when trigger element is visible
      if (entry.isIntersecting && auctionStore.pagination.hasMore && !auctionStore.loadingMore) {
        auctionStore.loadMoreAuctions()
      }
    },
    {
      rootMargin: '200px', // Start loading 200px before visible
      threshold: 0.1
    }
  )

  // Observe trigger element
  if (loadMoreTrigger.value) {
    observer.observe(loadMoreTrigger.value)
  }
}

// Re-setup Observer when hasMore changes
watch(() => auctionStore.pagination.hasMore, async (hasMore) => {
  if (hasMore) {
    await nextTick()
    setupIntersectionObserver()
  }
})

// Lifecycle
onMounted(async () => {
  // Initial data fetch
  await auctionStore.fetchAuctionList()

  // Setup Intersection Observer after next tick
  await nextTick()
  setupIntersectionObserver()
})

onUnmounted(() => {
  // Cleanup Observer
  if (observer) {
    observer.disconnect()
  }

  // Reset store
  auctionStore.reset()
})
</script>

<style scoped>
.bidder-auction-list-container {
  min-height: 100vh;
  background-color: hsl(var(--background));
}

/* Responsive design */
@media (max-width: 640px) {
  .bidder-auction-list-container {
    background-color: hsl(var(--background));
  }
}
</style>
