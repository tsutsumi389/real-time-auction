<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- ヘッダー -->
    <div class="mb-8">
      <button
        @click="handleBack"
        class="flex items-center text-gray-600 hover:text-gray-900 mb-4"
      >
        <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        オークション編集に戻る
      </button>
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-3xl font-bold text-gray-900">商品紐づけ管理</h1>
          <p v-if="auctionTitle" class="mt-2 text-sm text-gray-600">
            オークション: <span class="font-medium">{{ auctionTitle }}</span>
          </p>
        </div>
        <div v-if="!canModify" class="text-sm text-yellow-700 bg-yellow-50 px-4 py-2 rounded-md">
          <svg class="inline-block w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
          オークションが開始済みのため、商品の解除はできません
        </div>
      </div>
    </div>

    <!-- エラー表示 -->
    <div v-if="errorMessage" class="mb-6">
      <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded relative" role="alert">
        <span class="block sm:inline">{{ errorMessage }}</span>
        <button
          @click="errorMessage = ''"
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

    <!-- ローディング -->
    <div v-if="isLoading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
    </div>

    <!-- 2カラムレイアウト -->
    <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- 左カラム: 未割当商品一覧 -->
      <div class="bg-white shadow rounded-lg overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-lg font-medium text-gray-900">未割当商品</h2>
          <p class="mt-1 text-sm text-gray-500">オークションに追加する商品を選択してください</p>
        </div>

        <!-- 検索ボックス -->
        <div class="px-6 py-4 border-b border-gray-200">
          <div class="relative">
            <input
              v-model="unassignedSearch"
              type="text"
              placeholder="商品名で検索..."
              class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-md leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
              @keyup.enter="handleSearchUnassigned"
            />
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <svg class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
          </div>
        </div>

        <!-- 未割当商品リスト -->
        <div class="max-h-96 overflow-y-auto">
          <div v-if="unassignedLoading" class="p-8 text-center">
            <div class="inline-block animate-spin rounded-full h-6 w-6 border-4 border-blue-500 border-t-transparent"></div>
          </div>

          <div v-else-if="unassignedItems.length === 0" class="p-8 text-center text-gray-500">
            未割当の商品がありません
          </div>

          <ul v-else class="divide-y divide-gray-200">
            <li
              v-for="item in unassignedItems"
              :key="item.id"
              class="px-6 py-4 hover:bg-gray-50 cursor-pointer"
              @click="toggleItemSelection(item.id)"
            >
              <div class="flex items-center">
                <input
                  type="checkbox"
                  :checked="selectedItemIds.includes(item.id)"
                  class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                  @click.stop
                  @change="toggleItemSelection(item.id)"
                />
                <div class="ml-3 flex-1">
                  <div class="text-sm font-medium text-gray-900">{{ item.name }}</div>
                  <div v-if="item.description" class="text-sm text-gray-500 truncate">
                    {{ item.description }}
                  </div>
                </div>
                <div v-if="item.starting_price" class="text-sm text-gray-500">
                  {{ formatPrice(item.starting_price) }}
                </div>
              </div>
            </li>
          </ul>
        </div>

        <!-- ページネーション（未割当） -->
        <div v-if="unassignedPagination.totalPages > 1" class="px-6 py-4 border-t border-gray-200 flex justify-between items-center">
          <span class="text-sm text-gray-500">
            {{ unassignedPagination.totalItems }}件中 {{ (unassignedPagination.currentPage - 1) * 20 + 1 }}-{{ Math.min(unassignedPagination.currentPage * 20, unassignedPagination.totalItems) }}件
          </span>
          <div class="flex gap-2">
            <button
              :disabled="unassignedPagination.currentPage <= 1"
              class="px-3 py-1 text-sm border rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
              @click="changeUnassignedPage(unassignedPagination.currentPage - 1)"
            >
              前へ
            </button>
            <button
              :disabled="unassignedPagination.currentPage >= unassignedPagination.totalPages"
              class="px-3 py-1 text-sm border rounded disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-50"
              @click="changeUnassignedPage(unassignedPagination.currentPage + 1)"
            >
              次へ
            </button>
          </div>
        </div>

        <!-- 追加ボタン -->
        <div class="px-6 py-4 border-t border-gray-200 bg-gray-50">
          <button
            :disabled="selectedItemIds.length === 0 || isAssigning"
            class="w-full inline-flex justify-center items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            @click="handleAssignItems"
          >
            <svg v-if="isAssigning" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ isAssigning ? '追加中...' : `選択した商品を追加 (${selectedItemIds.length}件)` }}
          </button>
        </div>
      </div>

      <!-- 右カラム: オークション内商品一覧 -->
      <div class="bg-white shadow rounded-lg overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-lg font-medium text-gray-900">オークション内商品</h2>
          <p class="mt-1 text-sm text-gray-500">
            {{ assignedItems.length }}件の商品が登録されています
            <span v-if="assignedItems.length > 0" class="text-xs text-gray-400">（ドラッグ&ドロップで順序変更可能）</span>
          </p>
        </div>

        <!-- オークション内商品リスト（ドラッグ&ドロップ対応） -->
        <div class="max-h-[500px] overflow-y-auto relative">
          <!-- ローディングオーバーレイ -->
          <div v-if="isReordering" class="absolute inset-0 bg-white bg-opacity-75 flex items-center justify-center z-50">
            <div class="flex flex-col items-center gap-2">
              <div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-blue-500 border-t-transparent"></div>
              <p class="text-sm text-gray-600">順序を更新中...</p>
            </div>
          </div>

          <div v-if="assignedLoading" class="p-8 text-center">
            <div class="inline-block animate-spin rounded-full h-6 w-6 border-4 border-blue-500 border-t-transparent"></div>
          </div>

          <div v-else-if="assignedItems.length === 0" class="p-8 text-center text-gray-500">
            商品がまだ登録されていません
          </div>

          <ul v-else class="divide-y divide-gray-200" ref="sortableListRef">
            <li
              v-for="(item, index) in assignedItems"
              :key="item.id"
              :data-id="item.id"
              :class="[
                'px-6 py-4 flex items-center gap-4 transition-all duration-200',
                {
                  'opacity-50 scale-95 shadow-lg z-50': draggedIndex === index,
                  'border-blue-400 bg-blue-50 border-2 border-dashed': dragOverIndex === index && draggedIndex !== index,
                  'hover:bg-blue-50': draggedIndex === null && !isReordering,
                  'cursor-move': !isReordering,
                  'cursor-not-allowed': isReordering
                }
              ]"
              draggable="true"
              @dragstart="handleDragStart($event, index)"
              @dragover.prevent="handleDragOver($event, index)"
              @dragenter.prevent
              @drop="handleDrop($event, index)"
              @dragend="handleDragEnd"
              @touchstart="handleTouchStart($event, index)"
              @touchmove="handleTouchMove($event, index)"
              @touchend="handleTouchEnd($event, index)"
            >
              <!-- ドラッグハンドル -->
              <div class="flex-shrink-0 text-gray-400 cursor-move hover:text-gray-600 transition-colors duration-200 md:w-5 md:h-5 w-8 h-8 flex items-center justify-center">
                <svg class="w-full h-full" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                </svg>
              </div>

              <!-- ロット番号 -->
              <div class="flex-shrink-0 w-12 h-8 bg-blue-100 text-blue-800 rounded flex items-center justify-center text-sm font-medium">
                #{{ item.lot_number }}
              </div>

              <!-- 商品情報 -->
              <div class="flex-1 min-w-0">
                <div class="text-sm font-medium text-gray-900 truncate">{{ item.name }}</div>
                <div v-if="item.starting_price" class="text-sm text-gray-500">
                  {{ formatPrice(item.starting_price) }}
                </div>
              </div>

              <!-- ステータス表示 -->
              <div class="flex-shrink-0">
                <span
                  v-if="item.started_at"
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800"
                >
                  開始済み
                </span>
                <span
                  v-else
                  class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800"
                >
                  未開始
                </span>
              </div>

              <!-- 解除ボタン -->
              <div class="flex-shrink-0">
                <button
                  :disabled="!canModify || item.started_at || isUnassigning === item.id"
                  :title="getUnassignButtonTitle(item)"
                  class="inline-flex items-center px-3 py-1 border border-red-300 text-sm font-medium rounded-md text-red-700 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed"
                  @click="handleUnassignItem(item)"
                >
                  <svg v-if="isUnassigning === item.id" class="animate-spin -ml-0.5 mr-1.5 h-4 w-4" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  解除
                </button>
              </div>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getItemList, assignItemsToAuction, unassignItemFromAuction } from '@/services/itemApi'
import { getAuctionForEdit, reorderItems } from '@/services/auctionApi'

const route = useRoute()
const router = useRouter()

// Computed auction ID from route
const auctionId = computed(() => route.params.id)

// State
const isLoading = ref(true)
const errorMessage = ref('')
const successMessage = ref('')
const auctionTitle = ref('')
const auctionStatus = ref('pending')

// Unassigned items state
const unassignedItems = ref([])
const unassignedSearch = ref('')
const unassignedLoading = ref(false)
const unassignedPagination = ref({
  currentPage: 1,
  totalPages: 1,
  totalItems: 0,
})
const selectedItemIds = ref([])

// Assigned items state
const assignedItems = ref([])
const assignedLoading = ref(false)

// Operation state
const isAssigning = ref(false)
const isUnassigning = ref(null) // Holds item ID being unassigned

// Drag and drop state
const draggedIndex = ref(null)
const dragOverIndex = ref(null)
const sortableListRef = ref(null)
const isReordering = ref(false)

// Touch event state
const touchStartY = ref(0)
const touchCurrentY = ref(0)
const isTouching = ref(false)

// Computed: Can modify (add/remove items) - only before auction starts
const canModify = computed(() => {
  // Check if any item has been started
  return !assignedItems.value.some(item => item.started_at)
})

// Load auction data and items on mount
onMounted(async () => {
  await loadData()
})

// Watch for route changes
watch(() => route.params.id, async () => {
  await loadData()
})

// Load all data
async function loadData() {
  isLoading.value = true
  errorMessage.value = ''
  
  try {
    await Promise.all([
      loadAuctionData(),
      loadUnassignedItems(),
    ])
  } catch (err) {
    errorMessage.value = err.message || 'データの読み込みに失敗しました'
  } finally {
    isLoading.value = false
  }
}

// Load auction data including assigned items
async function loadAuctionData() {
  assignedLoading.value = true
  
  try {
    const response = await getAuctionForEdit(auctionId.value)
    auctionTitle.value = response.title
    auctionStatus.value = response.status
    assignedItems.value = (response.items || []).map(item => ({
      id: item.id,
      name: item.name,
      description: item.description,
      starting_price: item.starting_price,
      lot_number: item.lot_number,
      started_at: item.started_at,
      ended_at: item.ended_at,
    }))
  } catch (err) {
    throw new Error('オークション情報の取得に失敗しました')
  } finally {
    assignedLoading.value = false
  }
}

// Load unassigned items
async function loadUnassignedItems(page = 1) {
  unassignedLoading.value = true
  
  try {
    const response = await getItemList({
      status: 'unassigned',
      search: unassignedSearch.value,
      page: page,
      limit: 20,
    })
    
    unassignedItems.value = response.items || []
    unassignedPagination.value = {
      currentPage: response.page || 1,
      totalPages: Math.ceil((response.total || 0) / 20),
      totalItems: response.total || 0,
    }
  } catch (err) {
    throw new Error('未割当商品の取得に失敗しました')
  } finally {
    unassignedLoading.value = false
  }
}

// Search unassigned items
function handleSearchUnassigned() {
  loadUnassignedItems(1)
}

// Change unassigned items page
function changeUnassignedPage(page) {
  loadUnassignedItems(page)
}

// Toggle item selection
function toggleItemSelection(itemId) {
  const index = selectedItemIds.value.indexOf(itemId)
  if (index === -1) {
    selectedItemIds.value.push(itemId)
  } else {
    selectedItemIds.value.splice(index, 1)
  }
}

// Assign selected items to auction
async function handleAssignItems() {
  if (selectedItemIds.value.length === 0) return
  
  isAssigning.value = true
  errorMessage.value = ''
  
  try {
    await assignItemsToAuction(auctionId.value, selectedItemIds.value)
    
    successMessage.value = `${selectedItemIds.value.length}件の商品を追加しました`
    selectedItemIds.value = []
    
    // Reload both lists
    await Promise.all([
      loadAuctionData(),
      loadUnassignedItems(unassignedPagination.value.currentPage),
    ])
  } catch (err) {
    errorMessage.value = err.response?.data?.error || err.message || '商品の追加に失敗しました'
  } finally {
    isAssigning.value = false
  }
}

// Unassign item from auction
async function handleUnassignItem(item) {
  if (!canModify.value || item.started_at) return
  
  if (!confirm(`「${item.name}」をオークションから解除してもよろしいですか？`)) {
    return
  }
  
  isUnassigning.value = item.id
  errorMessage.value = ''
  
  try {
    await unassignItemFromAuction(auctionId.value, item.id)
    
    successMessage.value = `「${item.name}」をオークションから解除しました`
    
    // Reload both lists
    await Promise.all([
      loadAuctionData(),
      loadUnassignedItems(unassignedPagination.value.currentPage),
    ])
  } catch (err) {
    errorMessage.value = err.response?.data?.error || err.message || '商品の解除に失敗しました'
  } finally {
    isUnassigning.value = null
  }
}

// Get unassign button title (tooltip)
function getUnassignButtonTitle(item) {
  if (item.started_at) {
    return '開始済みの商品は解除できません'
  }
  if (!canModify.value) {
    return 'オークションが開始済みのため解除できません'
  }
  return 'オークションから解除'
}

// Drag and drop handlers
function handleDragStart(event, index) {
  if (isReordering.value) return
  draggedIndex.value = index
  event.dataTransfer.effectAllowed = 'move'
}

function handleDragOver(event, index) {
  if (draggedIndex.value === null) return
  if (index === draggedIndex.value) return

  event.dataTransfer.dropEffect = 'move'
  dragOverIndex.value = index
}

async function handleDrop(event, targetIndex) {
  if (draggedIndex.value === null) return
  if (targetIndex === draggedIndex.value) return

  dragOverIndex.value = null
  isReordering.value = true

  // Reorder items locally
  const itemsCopy = [...assignedItems.value]
  const [removed] = itemsCopy.splice(draggedIndex.value, 1)
  itemsCopy.splice(targetIndex, 0, removed)
  assignedItems.value = itemsCopy

  // Update lot numbers
  assignedItems.value.forEach((item, idx) => {
    item.lot_number = idx + 1
  })

  // Send reorder request to backend
  try {
    const itemIds = assignedItems.value.map(item => item.id)
    await reorderItems(auctionId.value, itemIds)
    successMessage.value = '商品の順序を更新しました'
    setTimeout(() => {
      successMessage.value = ''
    }, 3000)
  } catch (err) {
    errorMessage.value = err.response?.data?.error || err.message || '順序の更新に失敗しました'
    // Reload to restore original order
    await loadAuctionData()
  } finally {
    isReordering.value = false
  }

  draggedIndex.value = null
}

function handleDragEnd(event) {
  draggedIndex.value = null
  dragOverIndex.value = null
}

// Touch event handlers for mobile support
function handleTouchStart(event, index) {
  if (isReordering.value) return
  isTouching.value = true
  draggedIndex.value = index
  const touch = event.touches[0]
  touchStartY.value = touch.clientY
  touchCurrentY.value = touch.clientY
}

function handleTouchMove(event, index) {
  if (!isTouching.value || draggedIndex.value === null) return
  event.preventDefault()

  const touch = event.touches[0]
  touchCurrentY.value = touch.clientY

  // Calculate which item we're hovering over based on touch position
  const listItems = sortableListRef.value?.children
  if (!listItems) return

  for (let i = 0; i < listItems.length; i++) {
    const rect = listItems[i].getBoundingClientRect()
    if (touch.clientY >= rect.top && touch.clientY <= rect.bottom) {
      if (i !== draggedIndex.value) {
        dragOverIndex.value = i
      }
      break
    }
  }
}

async function handleTouchEnd(event, index) {
  if (!isTouching.value || draggedIndex.value === null) return

  const targetIndex = dragOverIndex.value

  // Reset touch state
  isTouching.value = false
  touchStartY.value = 0
  touchCurrentY.value = 0

  // If we have a valid drop target, perform the reorder
  if (targetIndex !== null && targetIndex !== draggedIndex.value) {
    await handleDrop(event, targetIndex)
  } else {
    draggedIndex.value = null
    dragOverIndex.value = null
  }
}

// Format price
function formatPrice(price) {
  if (price === null || price === undefined) return '-'
  return new Intl.NumberFormat('ja-JP').format(price) + ' pt'
}

// Navigate back to auction edit
function handleBack() {
  router.push({ name: 'auction-edit', params: { id: auctionId.value } })
}
</script>
