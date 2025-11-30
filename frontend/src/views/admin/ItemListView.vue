<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- ヘッダー -->
    <div class="mb-8 flex justify-between items-center">
      <div>
        <h1 class="text-3xl font-bold text-gray-900">商品一覧</h1>
        <p class="mt-2 text-sm text-gray-600">商品の管理とオークションへの紐づけ</p>
      </div>
      <router-link
        to="/admin/items/new"
        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
      >
        + 新規商品作成
      </router-link>
    </div>

    <!-- エラー表示 -->
    <div v-if="itemStore.error" class="mb-6">
      <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded relative" role="alert">
        <span class="block sm:inline">{{ itemStore.error }}</span>
        <button
          @click="itemStore.clearError"
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

    <!-- 検索・フィルタエリア -->
    <div class="mb-6 bg-white shadow rounded-lg p-4">
      <div class="flex flex-col sm:flex-row gap-4">
        <!-- 検索ボックス -->
        <div class="flex-1">
          <label for="search" class="sr-only">商品名で検索</label>
          <div class="relative">
            <input
              id="search"
              v-model="searchKeyword"
              type="text"
              placeholder="商品名で検索..."
              class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              @keyup.enter="handleSearch"
            />
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
          </div>
        </div>

        <!-- ステータスフィルタ -->
        <div class="sm:w-48">
          <label for="status" class="sr-only">ステータス</label>
          <select
            id="status"
            v-model="statusFilter"
            class="block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm rounded-md"
            @change="handleFilterChange"
          >
            <option value="all">すべて</option>
            <option value="unassigned">未割当</option>
            <option value="assigned">割当済み</option>
          </select>
        </div>

        <!-- 検索ボタン -->
        <button
          @click="handleSearch"
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          検索
        </button>

        <!-- リセットボタン -->
        <button
          @click="handleReset"
          class="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          リセット
        </button>
      </div>
    </div>

    <!-- 商品テーブル -->
    <div class="bg-white shadow rounded-lg overflow-hidden">
      <div v-if="itemStore.loading" class="p-8 text-center">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-blue-500 border-t-transparent"></div>
        <p class="mt-2 text-gray-500">読み込み中...</p>
      </div>

      <div v-else-if="itemStore.items.length === 0" class="p-8 text-center text-gray-500">
        商品が見つかりません
      </div>

      <table v-else class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              ID
            </th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              商品名
            </th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              オークション
            </th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              開始価格
            </th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              作成日
            </th>
            <th scope="col" class="relative px-6 py-3">
              <span class="sr-only">操作</span>
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="item in itemStore.items" :key="item.id" class="hover:bg-gray-50">
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ shortId(item.id) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-medium text-gray-900">{{ item.name }}</div>
              <div v-if="item.description" class="text-sm text-gray-500 truncate max-w-xs">
                {{ item.description }}
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <router-link
                v-if="item.auction_id"
                :to="`/admin/auctions/${item.auction_id}/edit`"
                class="text-sm text-blue-600 hover:text-blue-900"
              >
                {{ item.auction_title }}
              </router-link>
              <span v-else class="text-sm text-gray-400">未割当</span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
              {{ item.starting_price ? formatPrice(item.starting_price) : '-' }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ formatDate(item.created_at) }}
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <router-link
                :to="`/admin/items/${item.id}/edit`"
                class="text-blue-600 hover:text-blue-900"
              >
                編集
              </router-link>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- ページネーション -->
    <Pagination
      v-if="itemStore.items.length > 0"
      :current-page="itemStore.pagination.currentPage"
      :total-pages="itemStore.pagination.totalPages"
      :total-items="itemStore.pagination.totalItems"
      :items-per-page="itemStore.pagination.itemsPerPage"
      @change-page="itemStore.changePage"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useItemStore } from '@/stores/item'
import Pagination from '@/components/Pagination.vue'

const itemStore = useItemStore()

// ローカル状態
const searchKeyword = ref('')
const statusFilter = ref('all')
const successMessage = ref('')

// 初期データ取得
onMounted(async () => {
  await itemStore.fetchItemList()
})

// 検索実行
function handleSearch() {
  itemStore.setFiltersAndFetch({
    search: searchKeyword.value,
    status: statusFilter.value,
  })
}

// フィルタ変更
function handleFilterChange() {
  itemStore.setFiltersAndFetch({
    search: searchKeyword.value,
    status: statusFilter.value,
  })
}

// リセット
function handleReset() {
  searchKeyword.value = ''
  statusFilter.value = 'all'
  itemStore.resetFilters()
}

// IDを短縮表示
function shortId(id) {
  if (!id) return '-'
  return id.substring(0, 8) + '...'
}

// 価格フォーマット（ポイント表示）
function formatPrice(price) {
  if (price === null || price === undefined) return '-'
  return new Intl.NumberFormat('ja-JP').format(price) + ' pt'
}

// 日付フォーマット
function formatDate(dateString) {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}
</script>
