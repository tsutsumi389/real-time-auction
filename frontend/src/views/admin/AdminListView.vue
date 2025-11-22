<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- ヘッダー -->
      <div class="mb-8 flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">管理者一覧</h1>
          <p class="mt-2 text-sm text-gray-600">システム管理者とオークショニアのアカウント管理</p>
        </div>
        <Button @click="handleRegister" class="ml-4">
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          新規登録
        </Button>
      </div>

      <!-- エラー表示 -->
      <div v-if="adminStore.error" class="mb-6">
        <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded relative" role="alert">
          <span class="block sm:inline">{{ adminStore.error }}</span>
          <button
            @click="adminStore.clearError"
            class="absolute top-0 bottom-0 right-0 px-4 py-3"
            aria-label="エラーを閉じる"
          >
            <span class="text-2xl">&times;</span>
          </button>
        </div>
      </div>

      <!-- フィルタエリア -->
      <AdminFilters
        v-model="filters"
        :loading="adminStore.loading"
        @search="handleSearch"
        @reset="handleReset"
        @filter-change="handleFilterChange"
        class="mb-6"
      />

      <!-- 一覧テーブル -->
      <AdminTable
        :admins="adminStore.admins"
        :loading="adminStore.loading"
        :sort-field="adminStore.filters.sort"
        :sort-order="adminStore.filters.order"
        @sort="handleSort"
        @edit="handleEdit"
        @status-change="openStatusChangeDialog"
      />

      <!-- ページネーション -->
      <Pagination
        v-if="adminStore.admins.length > 0"
        :current-page="adminStore.pagination.currentPage"
        :total-pages="adminStore.pagination.totalPages"
        :total-items="adminStore.pagination.totalItems"
        :items-per-page="adminStore.pagination.itemsPerPage"
        @change-page="adminStore.changePage"
      />

    <!-- 状態変更モーダル -->
    <AdminStatusChangeDialog
      v-model="showStatusDialog"
      :admin="selectedAdmin"
      :loading="statusChanging"
      @confirm="handleStatusChange"
    />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import AdminFilters from '@/components/admin/AdminFilters.vue'
import AdminTable from '@/components/admin/AdminTable.vue'
import AdminStatusChangeDialog from '@/components/admin/AdminStatusChangeDialog.vue'
import Pagination from '@/components/Pagination.vue'
import Button from '@/components/ui/Button.vue'

const router = useRouter()
const adminStore = useAdminStore()

// ローカル状態
const filters = ref({
  search: '',
  role: '',
  status: '',
})
const showStatusDialog = ref(false)
const selectedAdmin = ref(null)
const statusChanging = ref(false)

// 初期データ取得
onMounted(async () => {
  await adminStore.fetchAdminList()
})

// 検索実行
function handleSearch() {
  adminStore.setFiltersAndFetch(filters.value)
}

// フィルタ変更（ロール・状態の即時適用）
function handleFilterChange() {
  adminStore.setFiltersAndFetch(filters.value)
}

// リセット
function handleReset() {
  filters.value = {
    search: '',
    role: '',
    status: '',
  }
  adminStore.resetFilters()
}

// ソート変更
function handleSort(field) {
  adminStore.changeSort(field)
}

// 編集画面へ遷移
function handleEdit(adminId) {
  router.push({ name: 'admin-edit', params: { id: adminId } })
}

// 新規登録画面へ遷移
function handleRegister() {
  router.push({ name: 'admin-register' })
}

// 状態変更モーダルを開く
function openStatusChangeDialog(admin) {
  selectedAdmin.value = admin
  showStatusDialog.value = true
}

// 状態変更実行
async function handleStatusChange() {
  if (!selectedAdmin.value) return

  statusChanging.value = true
  const newStatus = selectedAdmin.value.status === 'active' ? 'suspended' : 'active'

  const success = await adminStore.changeAdminStatus(selectedAdmin.value.id, newStatus)

  if (success) {
    statusChanging.value = false
    showStatusDialog.value = false
    selectedAdmin.value = null
    // 成功メッセージは省略（必要に応じてToast通知を追加）
  } else {
    statusChanging.value = false
    // エラーメッセージはstoreのerrorに設定されている
  }
}
</script>
