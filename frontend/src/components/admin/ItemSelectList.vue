<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { getItemList, assignItemsToAuction } from '@/services/itemApi'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Alert from '@/components/ui/Alert.vue'

const props = defineProps({
  auctionId: {
    type: String,
    required: true
  },
  selectedIds: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:selectedIds', 'items-assigned'])

// State
const items = ref([])
const isLoading = ref(false)
const searchQuery = ref('')
const pagination = ref({
  currentPage: 1,
  totalPages: 1,
  totalItems: 0,
  limit: 20
})

// Local selected IDs
const localSelectedIds = computed({
  get: () => props.selectedIds,
  set: (value) => emit('update:selectedIds', value)
})

// Assignment state
const isAssigning = ref(false)
const submitError = ref('')
const successMessage = ref('')

// Load items on mount
onMounted(() => {
  loadItems()
})

// Watch search query with debounce
let searchTimeout = null
watch(searchQuery, () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    loadItems(1)
  }, 300)
})

// Load unassigned items
async function loadItems(page = 1) {
  isLoading.value = true
  submitError.value = ''

  try {
    const response = await getItemList({
      status: 'unassigned',
      search: searchQuery.value,
      page: page,
      limit: pagination.value.limit,
    })

    items.value = response.items || []
    pagination.value = {
      currentPage: response.page || 1,
      totalPages: Math.ceil((response.total || 0) / pagination.value.limit),
      totalItems: response.total || 0,
      limit: pagination.value.limit
    }
  } catch (error) {
    submitError.value = error.response?.data?.error || error.message || '商品の取得に失敗しました'
  } finally {
    isLoading.value = false
  }
}

// Toggle item selection
function toggleItemSelection(itemId) {
  const newSelectedIds = [...localSelectedIds.value]
  const index = newSelectedIds.indexOf(itemId)

  if (index === -1) {
    newSelectedIds.push(itemId)
  } else {
    newSelectedIds.splice(index, 1)
  }

  localSelectedIds.value = newSelectedIds
}

// Toggle all items on current page
function toggleAllItems() {
  const allSelected = items.value.every(item => localSelectedIds.value.includes(item.id))

  if (allSelected) {
    // Deselect all on current page
    const itemIdsOnPage = items.value.map(item => item.id)
    localSelectedIds.value = localSelectedIds.value.filter(id => !itemIdsOnPage.includes(id))
  } else {
    // Select all on current page
    const itemIdsOnPage = items.value.map(item => item.id)
    const newSelectedIds = [...localSelectedIds.value]
    itemIdsOnPage.forEach(id => {
      if (!newSelectedIds.includes(id)) {
        newSelectedIds.push(id)
      }
    })
    localSelectedIds.value = newSelectedIds
  }
}

// Check if all items on current page are selected
const allItemsSelected = computed(() => {
  if (items.value.length === 0) return false
  return items.value.every(item => localSelectedIds.value.includes(item.id))
})

// Assign selected items
async function handleAssignItems() {
  if (localSelectedIds.value.length === 0) return

  isAssigning.value = true
  submitError.value = ''
  successMessage.value = ''

  try {
    await assignItemsToAuction(props.auctionId, localSelectedIds.value)

    successMessage.value = `${localSelectedIds.value.length}件の商品を追加しました`

    // Clear selection
    localSelectedIds.value = []

    // Emit success event
    emit('items-assigned')

    // Reload items
    await loadItems(pagination.value.currentPage)

    // Auto-hide success message after 3 seconds
    setTimeout(() => {
      successMessage.value = ''
    }, 3000)
  } catch (error) {
    submitError.value = error.response?.data?.error || error.message || '商品の割り当てに失敗しました'
  } finally {
    isAssigning.value = false
  }
}

// Change page
function changePage(page) {
  if (page < 1 || page > pagination.value.totalPages) return
  loadItems(page)
}

// Format price
function formatPrice(price) {
  if (price === null || price === undefined) return '-'
  return `¥${price.toLocaleString()}`
}
</script>

<template>
  <div class="space-y-4">
    <!-- Success Message -->
    <Alert v-if="successMessage" variant="success" class="mb-4">
      <p class="font-semibold">成功</p>
      <p>{{ successMessage }}</p>
    </Alert>

    <!-- Error Alert -->
    <Alert v-if="submitError" variant="destructive" class="mb-4">
      <p class="font-semibold">エラー</p>
      <p>{{ submitError }}</p>
    </Alert>

    <!-- Search Bar -->
    <div class="flex gap-3">
      <div class="flex-1">
        <Input
          v-model="searchQuery"
          placeholder="商品名で検索..."
          :disabled="isLoading"
        />
      </div>
      <Button
        @click="handleAssignItems"
        :disabled="localSelectedIds.length === 0 || isAssigning"
        class="min-w-[140px]"
      >
        {{ isAssigning ? '追加中...' : `選択した商品を追加 (${localSelectedIds.length})` }}
      </Button>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
    </div>

    <!-- Empty State -->
    <div v-else-if="items.length === 0" class="text-center py-12 text-gray-500">
      <svg class="w-12 h-12 mx-auto mb-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
      </svg>
      <p class="text-lg font-medium">未割当の商品がありません</p>
      <p class="text-sm mt-2">新しい商品を作成するか、他のオークションから商品を解除してください</p>
    </div>

    <!-- Items List -->
    <template v-else>
      <!-- Select All Checkbox -->
      <div class="flex items-center gap-2 px-4 py-2 bg-gray-50 rounded-md">
        <input
          type="checkbox"
          :checked="allItemsSelected"
          @change="toggleAllItems"
          class="w-4 h-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
        />
        <label class="text-sm text-gray-700 cursor-pointer" @click="toggleAllItems">
          すべて選択
        </label>
        <span class="text-sm text-gray-500 ml-2">
          ({{ items.length }}件中{{ items.filter(item => localSelectedIds.includes(item.id)).length }}件選択)
        </span>
      </div>

      <!-- Items -->
      <div class="border border-gray-200 rounded-lg overflow-hidden">
        <ul class="divide-y divide-gray-200">
          <li
            v-for="item in items"
            :key="item.id"
            class="hover:bg-gray-50 transition-colors"
          >
            <label
              :for="`item-${item.id}`"
              class="flex items-center gap-4 p-4 cursor-pointer"
            >
              <!-- Checkbox -->
              <input
                :id="`item-${item.id}`"
                type="checkbox"
                :checked="localSelectedIds.includes(item.id)"
                @change="toggleItemSelection(item.id)"
                class="w-5 h-5 text-blue-600 border-gray-300 rounded focus:ring-blue-500 flex-shrink-0"
              />

              <!-- Item Info -->
              <div class="flex-1 min-w-0">
                <div class="text-sm font-medium text-gray-900 truncate">
                  {{ item.name }}
                </div>
                <div v-if="item.description" class="text-sm text-gray-500 truncate mt-1">
                  {{ item.description }}
                </div>
              </div>

              <!-- Starting Price -->
              <div class="flex-shrink-0 text-sm text-gray-700">
                {{ formatPrice(item.starting_price) }}
              </div>
            </label>
          </li>
        </ul>
      </div>

      <!-- Pagination -->
      <div v-if="pagination.totalPages > 1" class="flex items-center justify-between">
        <div class="text-sm text-gray-500">
          全{{ pagination.totalItems }}件中
          {{ (pagination.currentPage - 1) * pagination.limit + 1 }}-{{ Math.min(pagination.currentPage * pagination.limit, pagination.totalItems) }}件を表示
        </div>
        <div class="flex gap-2">
          <Button
            variant="outline"
            size="sm"
            @click="changePage(pagination.currentPage - 1)"
            :disabled="pagination.currentPage === 1"
          >
            前へ
          </Button>
          <div class="flex items-center gap-1">
            <span class="text-sm text-gray-700">
              {{ pagination.currentPage }} / {{ pagination.totalPages }}
            </span>
          </div>
          <Button
            variant="outline"
            size="sm"
            @click="changePage(pagination.currentPage + 1)"
            :disabled="pagination.currentPage === pagination.totalPages"
          >
            次へ
          </Button>
        </div>
      </div>
    </template>
  </div>
</template>
