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

    <!-- Main Content -->
    <div class="relative z-10">
      <!-- Header -->
      <header class="sticky top-0 z-40 header-glass border-b border-lux-gold/20">
        <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div class="flex items-center justify-between h-16 sm:h-20">
            <!-- Left: Page Title -->
            <div class="flex items-center gap-4 min-w-0">
              <div class="hidden sm:flex w-12 h-12 rounded-xl bg-gradient-to-br from-lux-gold/20 to-lux-gold/5 border border-lux-gold/40 flex-shrink-0 items-center justify-center shadow-lg shadow-lux-gold/10">
                <svg class="w-6 h-6 text-lux-gold" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                </svg>
              </div>
              <div class="min-w-0">
                <h1 class="font-display text-lg sm:text-2xl text-lux-cream truncate font-medium tracking-wide">オークション一覧</h1>
                <p class="text-xs sm:text-sm text-lux-silver/60 truncate hidden sm:block mt-0.5">開催中および終了したオークションを閲覧できます</p>
              </div>
            </div>
          </div>
        </div>
      </header>

      <!-- Content Area -->
      <main class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 sm:py-8">
        <!-- Search Bar -->
        <div class="mb-6 lux-fade-in">
          <AuctionSearchBar
            v-model="searchKeyword"
            :loading="auctionStore.loading"
            @search="handleSearch"
            @clear="handleClearSearch"
          />
        </div>

        <!-- Filters -->
        <div class="mb-8 lux-fade-in lux-delay-1">
          <AuctionFilters
            :current-status="auctionStore.filters.status"
            :current-sort="auctionStore.filters.sort"
            :status-options="statusOptions"
            :sort-options="sortOptions"
            @update:status="handleStatusFilter"
            @update:sort="handleSortChange"
          />
        </div>

        <!-- Error Display -->
        <Transition
          enter-active-class="transition-all duration-300 ease-out"
          enter-from-class="opacity-0 -translate-y-2"
          enter-to-class="opacity-100 translate-y-0"
          leave-active-class="transition-all duration-200 ease-in"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <div v-if="auctionStore.error" class="mb-6 lux-glass-strong rounded-xl p-4 border border-red-500/30 bg-red-950/20">
            <div class="flex items-start gap-3">
              <div class="flex-shrink-0 w-10 h-10 rounded-lg bg-red-500/10 flex items-center justify-center">
                <AlertCircle class="h-5 w-5 text-red-400" />
              </div>
              <div class="flex-1 min-w-0">
                <h3 class="font-display text-base text-red-400 mb-1">エラーが発生しました</h3>
                <p class="text-sm text-red-300/80">{{ auctionStore.error }}</p>
              </div>
              <button
                @click="auctionStore.clearError"
                class="flex-shrink-0 p-2 rounded-lg hover:bg-red-500/10 transition-colors"
                aria-label="エラーを閉じる"
              >
                <X class="h-4 w-4 text-red-400" />
              </button>
            </div>
          </div>
        </Transition>

        <!-- Auction Card Grid -->
        <div class="lux-fade-in lux-delay-2">
          <AuctionCardGrid
            :auctions="auctionStore.auctions"
            :loading="auctionStore.loading"
            :empty-message="emptyMessage"
            @view-details="handleViewDetails"
            @join-auction="handleJoinAuction"
          >
            <!-- Custom Empty State -->
            <template #empty>
              <div class="text-center py-16 lux-fade-in">
                <div class="w-20 h-20 mx-auto mb-6 rounded-2xl bg-lux-noir-light/50 border border-lux-gold/20 flex items-center justify-center">
                  <Inbox class="h-10 w-10 text-lux-gold/40" :stroke-width="1.5" />
                </div>
                <p class="font-display text-xl text-lux-cream/60 mb-6">
                  {{ emptyMessage }}
                </p>
                <button
                  v-if="hasActiveFilters"
                  @click="handleResetFilters"
                  class="inline-flex items-center gap-2 px-6 py-3 rounded-xl lux-glass border border-lux-gold/30 text-lux-gold text-sm font-medium tracking-wide hover:bg-lux-gold/10 transition-all duration-300"
                >
                  <RotateCcw class="h-4 w-4" />
                  フィルタをリセット
                </button>
              </div>
            </template>
          </AuctionCardGrid>
        </div>

        <!-- Infinite Scroll Trigger -->
        <div
          v-if="auctionStore.auctions.length > 0 && auctionStore.pagination.hasMore"
          ref="loadMoreTrigger"
          class="py-8 flex justify-center"
        >
          <div v-if="auctionStore.loadingMore" class="flex items-center gap-3 text-lux-silver">
            <div class="relative w-8 h-8">
              <div class="absolute inset-0 rounded-full border-2 border-lux-gold/20"></div>
              <div class="absolute inset-0 rounded-full border-2 border-transparent border-t-lux-gold animate-spin"></div>
            </div>
            <span class="text-sm font-medium">読み込み中...</span>
          </div>
          <div v-else class="flex items-center gap-2 text-lux-silver/60 text-sm">
            <ChevronsDown class="h-4 w-4 animate-bounce" />
            スクロールして続きを読み込む
          </div>
        </div>

        <!-- Load Complete Message -->
        <div
          v-if="auctionStore.auctions.length > 0 && !auctionStore.pagination.hasMore"
          class="py-8 text-center"
        >
          <div class="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-emerald-500/10 border border-emerald-500/20 text-emerald-400 text-sm">
            <CheckCircle2 class="h-4 w-4" />
            すべてのオークションを表示しました（全{{ auctionStore.pagination.total }}件）
          </div>
        </div>
      </main>
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

const router = useRouter()
const auctionStore = useBidderAuctionStore()

// Search keyword
const searchKeyword = ref('')

// Status options
const statusOptions = [
  { value: 'active', label: '開催中' },
  { value: 'ended', label: '終了' },
  { value: 'cancelled', label: '中止' }
]

// Sort options
const sortOptions = [
  { value: 'started_at_desc', label: '開始日時が新しい順' },
  { value: 'started_at_asc', label: '開始日時が古い順' },
  { value: 'updated_at_desc', label: '更新日時が新しい順' },
  { value: 'updated_at_asc', label: '更新日時が古い順' }
]

// Infinite scroll trigger element
const loadMoreTrigger = ref(null)
let observer = null

// Check if filters are active
const hasActiveFilters = computed(() => {
  return auctionStore.filters.keyword !== '' ||
         auctionStore.filters.status !== 'active' ||
         auctionStore.filters.sort !== 'started_at_desc'
})

// Empty state message
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

// Event handlers
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
      if (entry.isIntersecting && auctionStore.pagination.hasMore && !auctionStore.loadingMore) {
        auctionStore.loadMoreAuctions()
      }
    },
    {
      rootMargin: '200px',
      threshold: 0.1
    }
  )

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
  document.body.classList.add('theme-luxury')

  await auctionStore.fetchAuctionList()
  await nextTick()
  setupIntersectionObserver()
})

onUnmounted(() => {
  document.body.classList.remove('theme-luxury')

  if (observer) {
    observer.disconnect()
  }

  auctionStore.reset()
})
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

.border-lux-gold\/20 {
  border-color: hsl(43 74% 49% / 0.2);
}

.border-lux-gold\/30 {
  border-color: hsl(43 74% 49% / 0.3);
}

.border-lux-gold\/40 {
  border-color: hsl(43 74% 49% / 0.4);
}

.bg-lux-noir-light\/50 {
  background-color: hsl(0 0% 8% / 0.5);
}

.shadow-lux-gold\/10 {
  --tw-shadow-color: hsl(43 74% 49% / 0.1);
}

.from-lux-gold\/20 {
  --tw-gradient-from: hsl(43 74% 49% / 0.2);
}

.to-lux-gold\/5 {
  --tw-gradient-to: hsl(43 74% 49% / 0.05);
}

/* Font display */
.font-display {
  font-family: 'Cormorant Garamond', Georgia, serif;
}
</style>
