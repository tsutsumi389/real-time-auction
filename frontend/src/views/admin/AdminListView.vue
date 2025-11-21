<template>
  <div class="min-h-screen bg-gray-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- ヘッダー -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">管理者一覧</h1>
        <p class="mt-2 text-sm text-gray-600">
          システム管理者とオークショニアのアカウント管理
        </p>
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
      <div class="bg-white shadow rounded-lg p-6 mb-6">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <!-- 検索 -->
          <div>
            <label for="search" class="block text-sm font-medium text-gray-700 mb-1">
              メールアドレス検索
            </label>
            <input
              id="search"
              v-model="searchInput"
              type="text"
              placeholder="例: admin@example.com"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              @keyup.enter="handleSearch"
            />
          </div>

          <!-- ロールフィルタ -->
          <div>
            <label for="role" class="block text-sm font-medium text-gray-700 mb-1">
              ロール
            </label>
            <select
              id="role"
              v-model="adminStore.filters.role"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              @change="handleFilterChange"
            >
              <option value="">すべて</option>
              <option value="system_admin">システム管理者</option>
              <option value="auctioneer">オークショニア</option>
            </select>
          </div>

          <!-- 状態フィルタ -->
          <div>
            <label for="status" class="block text-sm font-medium text-gray-700 mb-1">
              状態
            </label>
            <select
              id="status"
              v-model="adminStore.filters.status"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              @change="handleFilterChange"
            >
              <option value="">すべて</option>
              <option value="active">有効</option>
              <option value="inactive">停止中</option>
            </select>
          </div>

          <!-- アクションボタン -->
          <div class="flex items-end gap-2">
            <button
              @click="handleSearch"
              class="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md transition-colors"
              :disabled="adminStore.loading"
            >
              検索
            </button>
            <button
              @click="handleReset"
              class="flex-1 bg-gray-200 hover:bg-gray-300 text-gray-700 font-medium py-2 px-4 rounded-md transition-colors"
              :disabled="adminStore.loading"
            >
              リセット
            </button>
          </div>
        </div>
      </div>

      <!-- 一覧テーブル -->
      <div class="bg-white shadow rounded-lg overflow-hidden">
        <!-- ローディング表示 -->
        <div v-if="adminStore.loading" class="p-8 text-center">
          <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <p class="mt-2 text-gray-600">読み込み中...</p>
        </div>

        <!-- データ表示 -->
        <div v-else-if="adminStore.admins.length > 0">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th
                  scope="col"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                  @click="adminStore.changeSort('id')"
                >
                  <div class="flex items-center gap-1">
                    ID
                    <span v-if="adminStore.filters.sort === 'id'">
                      {{ adminStore.filters.order === 'asc' ? '↑' : '↓' }}
                    </span>
                  </div>
                </th>
                <th
                  scope="col"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                  @click="adminStore.changeSort('email')"
                >
                  <div class="flex items-center gap-1">
                    メールアドレス
                    <span v-if="adminStore.filters.sort === 'email'">
                      {{ adminStore.filters.order === 'asc' ? '↑' : '↓' }}
                    </span>
                  </div>
                </th>
                <th
                  scope="col"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                  @click="adminStore.changeSort('role')"
                >
                  <div class="flex items-center gap-1">
                    ロール
                    <span v-if="adminStore.filters.sort === 'role'">
                      {{ adminStore.filters.order === 'asc' ? '↑' : '↓' }}
                    </span>
                  </div>
                </th>
                <th
                  scope="col"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                  @click="adminStore.changeSort('status')"
                >
                  <div class="flex items-center gap-1">
                    状態
                    <span v-if="adminStore.filters.sort === 'status'">
                      {{ adminStore.filters.order === 'asc' ? '↑' : '↓' }}
                    </span>
                  </div>
                </th>
                <th
                  scope="col"
                  class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider cursor-pointer hover:bg-gray-100"
                  @click="adminStore.changeSort('created_at')"
                >
                  <div class="flex items-center gap-1">
                    作成日
                    <span v-if="adminStore.filters.sort === 'created_at'">
                      {{ adminStore.filters.order === 'asc' ? '↑' : '↓' }}
                    </span>
                  </div>
                </th>
                <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  アクション
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="admin in adminStore.admins" :key="admin.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ admin.id }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ admin.email }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm">
                  <span
                    :class="[
                      'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
                      admin.role === 'system_admin'
                        ? 'bg-purple-100 text-purple-800'
                        : 'bg-blue-100 text-blue-800',
                    ]"
                  >
                    {{ getRoleLabel(admin.role) }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm">
                  <span
                    :class="[
                      'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
                      admin.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800',
                    ]"
                  >
                    {{ getStatusLabel(admin.status) }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {{ formatDate(admin.created_at) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <div class="flex gap-2">
                    <button
                      @click="handleEdit(admin.id)"
                      class="text-blue-600 hover:text-blue-900"
                      :aria-label="`${admin.email}のアカウントを編集`"
                    >
                      編集
                    </button>
                    <button
                      @click="openStatusChangeDialog(admin)"
                      class="text-red-600 hover:text-red-900"
                      :aria-label="`${admin.email}のアカウントを${admin.status === 'active' ? '停止' : '復活'}`"
                    >
                      {{ admin.status === 'active' ? '停止' : '復活' }}
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>

          <!-- ページネーション -->
          <div class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6">
            <div class="flex-1 flex justify-between sm:hidden">
              <button
                @click="adminStore.changePage(adminStore.pagination.currentPage - 1)"
                :disabled="adminStore.pagination.currentPage === 1"
                class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                前へ
              </button>
              <button
                @click="adminStore.changePage(adminStore.pagination.currentPage + 1)"
                :disabled="adminStore.pagination.currentPage === adminStore.pagination.totalPages"
                class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                次へ
              </button>
            </div>
            <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
              <div>
                <p class="text-sm text-gray-700">
                  全 <span class="font-medium">{{ adminStore.pagination.totalItems }}</span> 件中
                  <span class="font-medium">
                    {{ (adminStore.pagination.currentPage - 1) * adminStore.pagination.itemsPerPage + 1 }}
                  </span>
                  -
                  <span class="font-medium">
                    {{
                      Math.min(
                        adminStore.pagination.currentPage * adminStore.pagination.itemsPerPage,
                        adminStore.pagination.totalItems
                      )
                    }}
                  </span>
                  件を表示
                </p>
              </div>
              <div>
                <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
                  <button
                    @click="adminStore.changePage(adminStore.pagination.currentPage - 1)"
                    :disabled="adminStore.pagination.currentPage === 1"
                    class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <span class="sr-only">前へ</span>
                    ‹
                  </button>
                  <button
                    v-for="page in getPageNumbers()"
                    :key="page"
                    @click="page !== '...' && adminStore.changePage(page)"
                    :class="[
                      page === adminStore.pagination.currentPage
                        ? 'z-10 bg-blue-50 border-blue-500 text-blue-600'
                        : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50',
                      'relative inline-flex items-center px-4 py-2 border text-sm font-medium',
                      page === '...' ? 'cursor-default' : 'cursor-pointer',
                    ]"
                    :disabled="page === '...'"
                  >
                    {{ page }}
                  </button>
                  <button
                    @click="adminStore.changePage(adminStore.pagination.currentPage + 1)"
                    :disabled="adminStore.pagination.currentPage === adminStore.pagination.totalPages"
                    class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <span class="sr-only">次へ</span>
                    ›
                  </button>
                </nav>
              </div>
            </div>
          </div>
        </div>

        <!-- データなし表示 -->
        <div v-else class="p-8 text-center text-gray-500">
          <p>管理者が見つかりませんでした</p>
        </div>
      </div>
    </div>

    <!-- 状態変更モーダル -->
    <div
      v-if="showStatusDialog"
      class="fixed z-10 inset-0 overflow-y-auto"
      aria-labelledby="modal-title"
      role="dialog"
      aria-modal="true"
    >
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="closeStatusChangeDialog"></div>

        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

        <div
          class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full"
        >
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="sm:flex sm:items-start">
              <div
                :class="[
                  'mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full sm:mx-0 sm:h-10 sm:w-10',
                  selectedAdmin?.status === 'active' ? 'bg-red-100' : 'bg-green-100',
                ]"
              >
                <span :class="[selectedAdmin?.status === 'active' ? 'text-red-600' : 'text-green-600', 'text-xl']">
                  {{ selectedAdmin?.status === 'active' ? '⚠' : '✓' }}
                </span>
              </div>
              <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
                <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                  アカウント{{ selectedAdmin?.status === 'active' ? '停止' : '復活' }}の確認
                </h3>
                <div class="mt-2">
                  <p class="text-sm text-gray-500">
                    {{ selectedAdmin?.email }} のアカウントを{{
                      selectedAdmin?.status === 'active' ? '停止' : '復活'
                    }}してもよろしいですか？
                  </p>
                  <p v-if="selectedAdmin?.status === 'active'" class="mt-2 text-sm text-red-600">
                    停止後、このアカウントではログインできなくなります。
                  </p>
                </div>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button
              @click="handleStatusChange"
              :disabled="statusChanging"
              :class="[
                selectedAdmin?.status === 'active'
                  ? 'bg-red-600 hover:bg-red-700'
                  : 'bg-green-600 hover:bg-green-700',
                'w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 text-base font-medium text-white focus:outline-none focus:ring-2 focus:ring-offset-2 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed',
              ]"
            >
              <span v-if="statusChanging" class="inline-block animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></span>
              {{ selectedAdmin?.status === 'active' ? '停止する' : '復活する' }}
            </button>
            <button
              @click="closeStatusChangeDialog"
              :disabled="statusChanging"
              class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
            >
              キャンセル
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin'

const router = useRouter()
const adminStore = useAdminStore()

// ローカル状態
const searchInput = ref('')
const showStatusDialog = ref(false)
const selectedAdmin = ref(null)
const statusChanging = ref(false)

// 初期データ取得
onMounted(async () => {
  await adminStore.fetchAdminList()
})

// 検索実行
function handleSearch() {
  adminStore.setFiltersAndFetch({ search: searchInput.value })
}

// フィルタ変更
function handleFilterChange() {
  adminStore.setFiltersAndFetch()
}

// リセット
function handleReset() {
  searchInput.value = ''
  adminStore.resetFilters()
}

// 編集画面へ遷移
function handleEdit(adminId) {
  router.push({ name: 'admin-edit', params: { id: adminId } })
}

// 状態変更モーダルを開く
function openStatusChangeDialog(admin) {
  selectedAdmin.value = admin
  showStatusDialog.value = true
}

// 状態変更モーダルを閉じる
function closeStatusChangeDialog() {
  if (!statusChanging.value) {
    showStatusDialog.value = false
    selectedAdmin.value = null
  }
}

// 状態変更実行
async function handleStatusChange() {
  if (!selectedAdmin.value) return

  statusChanging.value = true
  const newStatus = selectedAdmin.value.status === 'active' ? 'inactive' : 'active'

  const success = await adminStore.changeAdminStatus(selectedAdmin.value.id, newStatus)

  if (success) {
    statusChanging.value = false
    closeStatusChangeDialog()
    // 成功メッセージは省略（必要に応じてToast通知を追加）
  } else {
    statusChanging.value = false
    // エラーメッセージはstoreのerrorに設定されている
  }
}

// ロールラベル取得
function getRoleLabel(role) {
  return role === 'system_admin' ? 'システム管理者' : 'オークショニア'
}

// 状態ラベル取得
function getStatusLabel(status) {
  return status === 'active' ? '有効' : '停止中'
}

// 日付フォーマット
function formatDate(dateString) {
  const date = new Date(dateString)
  return date.toLocaleDateString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
  })
}

// ページ番号配列取得（省略記号付き）
function getPageNumbers() {
  const current = adminStore.pagination.currentPage
  const total = adminStore.pagination.totalPages
  const delta = 2 // 現在ページの前後に表示するページ数

  if (total <= 7) {
    // 総ページ数が7以下の場合は全て表示
    return Array.from({ length: total }, (_, i) => i + 1)
  }

  const pages = []
  const left = current - delta
  const right = current + delta

  for (let i = 1; i <= total; i++) {
    if (i === 1 || i === total || (i >= left && i <= right)) {
      pages.push(i)
    } else if (pages[pages.length - 1] !== '...') {
      pages.push('...')
    }
  }

  return pages
}
</script>
