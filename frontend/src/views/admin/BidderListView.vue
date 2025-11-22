<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- ヘッダー -->
    <div class="mb-8 flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">入札者一覧</h1>
        <p class="mt-2 text-sm text-gray-600">入札者アカウントの管理とポイント付与</p>
      </div>
      <Button @click="handleRegister" class="ml-4">
        <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        新規登録
      </Button>
    </div>

    <!-- エラー表示 -->
    <div v-if="bidderStore.error" class="mb-6">
      <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded relative" role="alert">
        <span class="block sm:inline">{{ bidderStore.error }}</span>
        <button
          @click="bidderStore.clearError"
          class="absolute top-0 bottom-0 right-0 px-4 py-3"
          aria-label="エラーを閉じる"
        >
          <span class="text-2xl">&times;</span>
        </button>
      </div>
    </div>

    <!-- フィルタエリア -->
    <BidderFilters
      v-model="filters"
      :loading="bidderStore.loading"
      @search="handleSearch"
      @reset="handleReset"
      @filter-change="handleFilterChange"
      class="mb-6"
    />

    <!-- 一覧テーブル -->
    <BidderTable
      :bidders="bidderStore.bidders"
      :loading="bidderStore.loading"
      :sort-field="bidderStore.filters.sort"
      :sort-order="bidderStore.filters.order"
      @sort="handleSort"
      @edit="handleEdit"
      @grant-points="openGrantPointsDialog"
      @show-history="openPointHistoryDialog"
      @status-change="openStatusChangeDialog"
    />

    <!-- ページネーション -->
    <Pagination
      v-if="bidderStore.bidders.length > 0"
      :current-page="bidderStore.pagination.currentPage"
      :total-pages="bidderStore.pagination.totalPages"
      :total-items="bidderStore.pagination.totalItems"
      :items-per-page="bidderStore.pagination.itemsPerPage"
      @change-page="bidderStore.changePage"
    />

    <!-- ポイント付与モーダル -->
    <GrantPointsDialog
      v-model="showGrantPointsDialog"
      :bidder="selectedBidder"
      :loading="pointsGranting"
      @confirm="handleGrantPoints"
    />

    <!-- ポイント履歴モーダル -->
    <PointHistoryDialog
      v-model="showPointHistoryDialog"
      :bidder="selectedBidder"
    />

    <!-- 状態変更モーダル -->
    <BidderStatusChangeDialog
      v-model="showStatusDialog"
      :bidder="selectedBidder"
      :loading="statusChanging"
      @confirm="handleStatusChange"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useBidderStore } from '@/stores/bidder'
import BidderFilters from '@/components/bidder/BidderFilters.vue'
import BidderTable from '@/components/bidder/BidderTable.vue'
import GrantPointsDialog from '@/components/bidder/GrantPointsDialog.vue'
import PointHistoryDialog from '@/components/bidder/PointHistoryDialog.vue'
import BidderStatusChangeDialog from '@/components/bidder/BidderStatusChangeDialog.vue'
import Pagination from '@/components/Pagination.vue'
import Button from '@/components/ui/Button.vue'

const router = useRouter()
const bidderStore = useBidderStore()

// ローカル状態
const filters = ref({
  keyword: '',
  status: ['active', 'suspended'],
})
const showGrantPointsDialog = ref(false)
const showPointHistoryDialog = ref(false)
const showStatusDialog = ref(false)
const selectedBidder = ref(null)
const pointsGranting = ref(false)
const statusChanging = ref(false)

// 初期データ取得
onMounted(async () => {
  await bidderStore.fetchBidderList()
})

// 検索実行
function handleSearch() {
  bidderStore.setFiltersAndFetch(filters.value)
}

// フィルタ変更（状態の即時適用）
function handleFilterChange() {
  bidderStore.setFiltersAndFetch(filters.value)
}

// リセット
function handleReset() {
  filters.value = {
    keyword: '',
    status: ['active', 'suspended'],
  }
  bidderStore.resetFilters()
}

// ソート変更
function handleSort(field) {
  bidderStore.changeSort(field)
}

// 編集画面へ遷移
function handleEdit(bidderId) {
  router.push({ name: 'bidder-edit', params: { id: bidderId } })
}

// 新規登録画面へ遷移
function handleRegister() {
  router.push({ name: 'bidder-register' })
}

// ポイント付与モーダルを開く
function openGrantPointsDialog(bidder) {
  selectedBidder.value = bidder
  showGrantPointsDialog.value = true
}

// ポイント付与実行
async function handleGrantPoints(points) {
  if (!selectedBidder.value) return

  pointsGranting.value = true

  const success = await bidderStore.addPoints(selectedBidder.value.id, points)

  if (success) {
    pointsGranting.value = false
    showGrantPointsDialog.value = false
    selectedBidder.value = null
  } else {
    pointsGranting.value = false
  }
}

// ポイント履歴モーダルを開く
function openPointHistoryDialog(bidder) {
  selectedBidder.value = bidder
  showPointHistoryDialog.value = true
}

// 状態変更モーダルを開く
function openStatusChangeDialog(bidder) {
  selectedBidder.value = bidder
  showStatusDialog.value = true
}

// 状態変更実行
async function handleStatusChange() {
  if (!selectedBidder.value) return

  statusChanging.value = true
  const newStatus = selectedBidder.value.status === 'active' ? 'suspended' : 'active'

  const success = await bidderStore.changeBidderStatus(selectedBidder.value.id, newStatus)

  if (success) {
    statusChanging.value = false
    showStatusDialog.value = false
    selectedBidder.value = null
  } else {
    statusChanging.value = false
  }
}
</script>
