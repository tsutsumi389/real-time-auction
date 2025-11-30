<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- ヘッダー -->
    <div class="mb-8 flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">オークション一覧</h1>
        <p class="mt-2 text-sm text-gray-600">オークションの管理と状態変更</p>
      </div>
      <router-link
        to="/admin/auctions/new"
        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
      >
        + 新規オークション作成
      </router-link>
    </div>

    <!-- エラー表示 -->
    <div v-if="auctionStore.error" class="mb-6">
      <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded relative" role="alert">
        <span class="block sm:inline">{{ auctionStore.error }}</span>
        <button
          @click="auctionStore.clearError"
          class="absolute top-0 bottom-0 right-0 px-4 py-3"
          aria-label="エラーを閉じる"
        >
          <span class="text-2xl">&times;</span>
        </button>
      </div>
    </div>

    <!-- 成功メッセージ -->
    <div v-if="successMessage" class="mb-6">
      <div class="bg-green-50 border border-green-200 text-green-700 px-4 py-3 rounded relative" role="alert">
        <span class="block sm:inline">{{ successMessage }}</span>
        <button
          @click="successMessage = ''"
          class="absolute top-0 bottom-0 right-0 px-4 py-3"
          aria-label="メッセージを閉じる"
        >
          <span class="text-2xl">&times;</span>
        </button>
      </div>
    </div>

    <!-- 検索バー -->
    <div class="mb-6">
      <AuctionSearchBar
        v-model="auctionStore.filters.keyword"
        @search="handleSearch"
      />
    </div>

    <!-- フィルタエリア -->
    <AuctionFilters
      v-model="auctionStore.filters"
      :loading="auctionStore.loading"
      @filter-change="handleFilterChange"
      @reset="handleReset"
      class="mb-6"
    />

    <!-- 一覧テーブル -->
    <AuctionTable
      :auctions="auctionStore.auctions"
      :loading="auctionStore.loading"
      :sort-field="auctionStore.filters.sort"
      :sort-order="auctionStore.filters.order"
      :is-system-admin="authStore.isSystemAdmin"
      @sort="handleSort"
      @view-details="handleViewDetails"
      @view-live="handleViewLive"
      @start="openStartDialog"
      @end="openEndDialog"
      @cancel="openCancelDialog"
      @edit="handleEdit"
    />

    <!-- ページネーション -->
    <Pagination
      v-if="auctionStore.auctions.length > 0"
      :current-page="auctionStore.pagination.currentPage"
      :total-pages="auctionStore.pagination.totalPages"
      :total-items="auctionStore.pagination.totalItems"
      :items-per-page="auctionStore.pagination.itemsPerPage"
      @change-page="auctionStore.changePage"
    />

    <!-- オークション開始モーダル -->
    <AuctionStartDialog
      v-model="showStartDialog"
      :auction="selectedAuction"
      :loading="actionLoading"
      @confirm="handleStartAuction"
    />

    <!-- オークション終了モーダル -->
    <AuctionEndDialog
      v-model="showEndDialog"
      :auction="selectedAuction"
      :loading="actionLoading"
      @confirm="handleEndAuction"
    />

    <!-- オークション中止モーダル -->
    <AuctionCancelDialog
      v-model="showCancelDialog"
      :auction="selectedAuction"
      :loading="actionLoading"
      @confirm="handleCancelAuction"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuctionStore } from '@/stores/auction'
import { useAuthStore } from '@/stores/auth'
import AuctionSearchBar from '@/components/admin/AuctionSearchBar.vue'
import AuctionFilters from '@/components/admin/AuctionFilters.vue'
import AuctionTable from '@/components/admin/AuctionTable.vue'
import AuctionStartDialog from '@/components/admin/AuctionStartDialog.vue'
import AuctionEndDialog from '@/components/admin/AuctionEndDialog.vue'
import AuctionCancelDialog from '@/components/admin/AuctionCancelDialog.vue'
import Pagination from '@/components/Pagination.vue'

const router = useRouter()
const auctionStore = useAuctionStore()
const authStore = useAuthStore()

// ローカル状態
const showStartDialog = ref(false)
const showEndDialog = ref(false)
const showCancelDialog = ref(false)
const selectedAuction = ref(null)
const actionLoading = ref(false)
const successMessage = ref('')

// 初期データ取得
onMounted(async () => {
  await auctionStore.fetchAuctionList()
})

// 検索実行
function handleSearch() {
  auctionStore.setFiltersAndFetch({ keyword: auctionStore.filters.keyword })
}

// フィルタ変更
function handleFilterChange() {
  auctionStore.fetchAuctionList()
}

// フィルタリセット
function handleReset() {
  auctionStore.resetFilters()
}

// ソート変更
function handleSort(field) {
  auctionStore.changeSort(field)
}

// 詳細表示
function handleViewDetails(auctionId) {
  router.push(`/admin/auctions/${auctionId}/edit`)
}

// 編集画面へ
function handleEdit(auctionId) {
  router.push(`/admin/auctions/${auctionId}/edit`)
}

// 開催中画面へ
function handleViewLive(auctionId) {
  router.push(`/admin/auctions/${auctionId}/live`)
}

// オークション開始ダイアログを開く
function openStartDialog(auction) {
  selectedAuction.value = auction
  showStartDialog.value = true
}

// オークション開始処理
async function handleStartAuction() {
  actionLoading.value = true
  const success = await auctionStore.handleStartAuction(selectedAuction.value.id)
  actionLoading.value = false

  if (success) {
    showStartDialog.value = false
    successMessage.value = 'オークションを公開しました'
    setTimeout(() => {
      successMessage.value = ''
    }, 3000)
  }
}

// オークション終了ダイアログを開く
function openEndDialog(auction) {
  selectedAuction.value = auction
  showEndDialog.value = true
}

// オークション終了処理
async function handleEndAuction() {
  actionLoading.value = true
  const success = await auctionStore.handleEndAuction(selectedAuction.value.id)
  actionLoading.value = false

  if (success) {
    showEndDialog.value = false
    successMessage.value = 'オークションを終了しました'
    setTimeout(() => {
      successMessage.value = ''
    }, 3000)
  }
}

// オークション中止ダイアログを開く
function openCancelDialog(auction) {
  selectedAuction.value = auction
  showCancelDialog.value = true
}

// オークション中止処理
async function handleCancelAuction() {
  actionLoading.value = true
  const success = await auctionStore.handleCancelAuction(selectedAuction.value.id)
  actionLoading.value = false

  if (success) {
    showCancelDialog.value = false
    successMessage.value = 'オークションを中止しました'
    setTimeout(() => {
      successMessage.value = ''
    }, 3000)
  }
}
</script>
