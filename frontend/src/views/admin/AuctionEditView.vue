<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuctionStore } from '@/stores/auction'
import Button from '@/components/ui/Button.vue'
import Alert from '@/components/ui/Alert.vue'
import Card from '@/components/ui/Card.vue'
import AuctionBasicInfo from '@/components/admin/AuctionBasicInfo.vue'
import EditItemForm from '@/components/admin/EditItemForm.vue'
import AuctionStatusBadge from '@/components/admin/AuctionStatusBadge.vue'
import DeleteConfirmModal from '@/components/admin/DeleteConfirmModal.vue'
import ItemCreateInline from '@/components/admin/ItemCreateInline.vue'
import ItemSelectList from '@/components/admin/ItemSelectList.vue'
import Tabs from '@/components/ui/Tabs.vue'
import TabsList from '@/components/ui/TabsList.vue'
import TabsTrigger from '@/components/ui/TabsTrigger.vue'
import TabsContent from '@/components/ui/TabsContent.vue'

const route = useRoute()
const router = useRouter()
const auctionStore = useAuctionStore()

// Get auction ID from route
const auctionId = computed(() => route.params.id)

// Form data
const basicInfo = ref({
  title: '',
  description: '',
  started_at: '',
})

const items = ref([])

// Original data for change detection
const originalData = ref(null)

// Validation errors
const basicInfoErrors = ref({
  title: '',
  description: '',
  started_at: '',
})

const itemErrors = ref([])

// Form state
const isLoading = ref(true)
const isSubmitting = ref(false)
const submitError = ref('')
const canEdit = ref(false)
const canEditReason = ref('')
const auctionStatus = ref('pending')

// Delete modal state
const showDeleteModal = ref(false)
const deleteItemIndex = ref(-1)
const isDeleting = ref(false)

// Success message state
const showSuccessMessage = ref(false)

// Tab state
const activeTab = ref('assigned')
const selectedItemIds = ref([])

// Drag and drop state
const draggedIndex = ref(null)

// Computed for change detection
const hasChanges = computed(() => {
  if (!originalData.value) return false

  // Check basic info changes
  if (basicInfo.value.title !== originalData.value.title) return true
  if (basicInfo.value.description !== originalData.value.description) return true

  // Check items changes
  if (items.value.length !== originalData.value.items.length) return true

  for (let i = 0; i < items.value.length; i++) {
    const current = items.value[i]
    const original = originalData.value.items.find((item) => item.id === current.id)
    if (!original) return true
    if (current.name !== original.name) return true
    if (current.description !== original.description) return true
    if (current.lot_number !== original.lot_number) return true
  }

  return false
})

// Load auction data on mount
onMounted(async () => {
  await loadAuctionData()
})

// Watch for route changes
watch(
  () => route.params.id,
  async () => {
    await loadAuctionData()
  }
)

async function loadAuctionData() {
  isLoading.value = true
  submitError.value = ''

  // Check if redirected from auction creation
  if (route.query.created === 'true') {
    showSuccessMessage.value = true
    activeTab.value = 'create' // Default to create tab
    // Auto-hide after 5 seconds
    setTimeout(() => {
      showSuccessMessage.value = false
    }, 5000)
  }

  try {
    const auction = await auctionStore.fetchAuctionForEdit(auctionId.value)

    if (!auction) {
      submitError.value = 'オークションが見つかりません'
      isLoading.value = false
      return
    }

    // Set form data
    basicInfo.value = {
      title: auction.title || '',
      description: auction.description || '',
      started_at: auction.started_at ? formatDateTimeLocal(auction.started_at) : '',
    }

    items.value = auction.items.map((item) => ({
      id: item.id,
      name: item.name || '',
      description: item.description || '',
      lot_number: item.lot_number,
      starting_price: item.starting_price,
      current_price: item.current_price,
      started_at: item.started_at,
      ended_at: item.ended_at,
      can_edit: item.can_edit,
      can_delete: item.can_delete,
      bid_count: item.bid_count,
    }))

    // Initialize item errors
    itemErrors.value = items.value.map(() => ({ name: '', description: '', starting_price: '' }))

    // Store original data for change detection
    originalData.value = {
      title: auction.title || '',
      description: auction.description || '',
      items: auction.items.map((item) => ({
        id: item.id,
        name: item.name || '',
        description: item.description || '',
        lot_number: item.lot_number,
      })),
    }

    canEdit.value = auction.can_edit
    canEditReason.value = auction.can_edit_reason || ''
    auctionStatus.value = auction.status
  } catch (err) {
    submitError.value = err.message || 'オークション情報の取得に失敗しました'
  } finally {
    isLoading.value = false
  }
}

function formatDateTimeLocal(dateString) {
  if (!dateString) return ''
  const date = new Date(dateString)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  return `${year}-${month}-${day}T${hours}:${minutes}`
}

// Validation functions for basic info
function validateBasicInfoField(field) {
  if (field === 'title') {
    if (!basicInfo.value.title.trim()) {
      basicInfoErrors.value.title = 'タイトルを入力してください'
      return false
    }
    if (basicInfo.value.title.length > 200) {
      basicInfoErrors.value.title = 'タイトルは200文字以内で入力してください'
      return false
    }
    basicInfoErrors.value.title = ''
    return true
  } else if (field === 'description') {
    if (basicInfo.value.description.length > 2000) {
      basicInfoErrors.value.description = '説明は2000文字以内で入力してください'
      return false
    }
    basicInfoErrors.value.description = ''
    return true
  } else if (field === 'started_at') {
    basicInfoErrors.value.started_at = ''
    return true
  }
}

// Validation functions for items
function validateItemField(index, field) {
  if (!itemErrors.value[index]) {
    itemErrors.value[index] = { name: '', description: '', starting_price: '' }
  }

  if (field === 'name') {
    if (!items.value[index].name.trim()) {
      itemErrors.value[index].name = '商品名を入力してください'
      return false
    }
    if (items.value[index].name.length > 200) {
      itemErrors.value[index].name = '商品名は200文字以内で入力してください'
      return false
    }
    itemErrors.value[index].name = ''
    return true
  } else if (field === 'description') {
    if (items.value[index].description.length > 2000) {
      itemErrors.value[index].description = '説明は2000文字以内で入力してください'
      return false
    }
    itemErrors.value[index].description = ''
    return true
  } else if (field === 'starting_price') {
    const price = items.value[index].starting_price
    if (price !== null && price !== '' && price < 1) {
      itemErrors.value[index].starting_price = '開始価格は1以上で入力してください'
      return false
    }
    itemErrors.value[index].starting_price = ''
    return true
  }
}

function validateForm() {
  let isValid = true

  // Validate title
  if (!validateBasicInfoField('title')) {
    isValid = false
  }

  // Validate description
  if (!validateBasicInfoField('description')) {
    isValid = false
  }

  // Validate items
  items.value.forEach((_, index) => {
    if (!validateItemField(index, 'name')) {
      isValid = false
    }
    if (!validateItemField(index, 'description')) {
      isValid = false
    }
  })

  return isValid
}

// Item management
function handleAddItem() {
  // Add new item locally (will be saved when user clicks save)
  const newLotNumber = items.value.length + 1
  items.value.push({
    id: null, // null indicates new item not yet saved
    name: '',
    description: '',
    lot_number: newLotNumber,
    starting_price: null,
    current_price: null,
    started_at: null,
    ended_at: null,
    can_edit: true,
    can_delete: true,
    bid_count: 0,
    isNew: true, // flag to indicate this is a new unsaved item
  })
  itemErrors.value.push({ name: '', description: '', starting_price: '' })
}

function requestDeleteItem(index) {
  deleteItemIndex.value = index
  showDeleteModal.value = true
}

async function confirmDeleteItem() {
  if (deleteItemIndex.value < 0) return

  const item = items.value[deleteItemIndex.value]

  // If item is new (not saved to server yet), just remove locally
  if (item.isNew || !item.id) {
    items.value.splice(deleteItemIndex.value, 1)
    itemErrors.value.splice(deleteItemIndex.value, 1)

    // Recalculate lot numbers
    items.value.forEach((item, idx) => {
      item.lot_number = idx + 1
    })

    showDeleteModal.value = false
    deleteItemIndex.value = -1
    return
  }

  isDeleting.value = true

  const success = await auctionStore.handleDeleteItem(auctionId.value, item.id)
  if (success) {
    items.value.splice(deleteItemIndex.value, 1)
    itemErrors.value.splice(deleteItemIndex.value, 1)

    // Recalculate lot numbers
    items.value.forEach((item, idx) => {
      item.lot_number = idx + 1
    })
  } else {
    submitError.value = auctionStore.error || '商品の削除に失敗しました'
  }

  isDeleting.value = false
  showDeleteModal.value = false
  deleteItemIndex.value = -1
}

function cancelDeleteItem() {
  showDeleteModal.value = false
  deleteItemIndex.value = -1
}

async function handleItemUpdate(index, newValue) {
  // Only update if the item is editable
  if (!items.value[index].can_edit) return

  const item = items.value[index]
  if (item.id && originalData.value) {
    const original = originalData.value.items.find((i) => i.id === item.id)
    if (original && (newValue.name !== original.name || newValue.description !== original.description)) {
      // Update item on server
      const result = await auctionStore.handleUpdateItem(auctionId.value, item.id, {
        name: newValue.name,
        description: newValue.description,
      })
      if (!result) {
        submitError.value = auctionStore.error || '商品の更新に失敗しました'
      }
    }
  }
}

// Form submission
async function handleSubmit() {
  submitError.value = ''

  if (!validateForm()) {
    submitError.value = '入力内容に誤りがあります。エラーメッセージを確認してください。'
    return
  }

  if (!hasChanges.value) {
    router.push({ name: 'auction-list' })
    return
  }

  isSubmitting.value = true

  try {
    // Update auction basic info
    const auctionData = {
      title: basicInfo.value.title,
      description: basicInfo.value.description,
    }

    const result = await auctionStore.handleUpdateAuction(auctionId.value, auctionData)

    if (!result) {
      submitError.value = auctionStore.error || 'オークションの更新に失敗しました'
      isSubmitting.value = false
      return
    }

    // Save new items
    const newItems = items.value.filter((item) => item.isNew || !item.id)
    for (const item of newItems) {
      if (!item.name.trim()) {
        submitError.value = '商品名を入力してください'
        isSubmitting.value = false
        return
      }

      const addResult = await auctionStore.handleAddItem(auctionId.value, {
        name: item.name,
        description: item.description || '',
        starting_price: item.starting_price,
      })

      if (!addResult) {
        submitError.value = auctionStore.error || '商品の追加に失敗しました'
        isSubmitting.value = false
        return
      }

      // Update local item with server-assigned id
      item.id = addResult.id
      item.isNew = false
    }

    // Check if item order changed (only for existing items)
    const existingItems = items.value.filter((item) => item.id && !item.isNew)
    const originalItemOrder = originalData.value.items.map((item) => item.id)
    const currentItemOrder = existingItems.map((item) => item.id)
    const orderChanged = JSON.stringify(originalItemOrder) !== JSON.stringify(currentItemOrder)

    if (orderChanged || newItems.length > 0) {
      // If there are new items or order changed, reorder all items
      const allItemIds = items.value.map((item) => item.id).filter((id) => id)
      const reorderSuccess = await auctionStore.handleReorderItems(auctionId.value, allItemIds)
      if (!reorderSuccess) {
        submitError.value = auctionStore.error || '商品順序の更新に失敗しました'
        isSubmitting.value = false
        return
      }
    }

    // Success - redirect to auction list
    router.push({ name: 'auction-list' })
  } catch (error) {
    submitError.value = error.message || 'オークションの更新に失敗しました'
  } finally {
    isSubmitting.value = false
  }
}

function handleCancel() {
  if (hasChanges.value) {
    if (!confirm('変更内容が失われますが、よろしいですか?')) {
      return
    }
  }
  router.push({ name: 'auction-list' })
}

function handleBack() {
  if (hasChanges.value) {
    if (!confirm('変更内容が失われますが、よろしいですか?')) {
      return
    }
  }
  router.push({ name: 'auction-list' })
}

// Handle item created from inline form
async function handleItemCreated() {
  // Switch to assigned tab
  activeTab.value = 'assigned'
  // Reload auction data to show new item
  await loadAuctionData()
}

// Handle items assigned from select list
async function handleItemsAssigned() {
  // Switch to assigned tab
  activeTab.value = 'assigned'
  // Reload auction data to show new items
  await loadAuctionData()
}

// Handle move up/down for items
function handleMoveUp(index) {
  if (index === 0) return

  const newItems = [...items.value]
  const temp = newItems[index]
  newItems[index] = newItems[index - 1]
  newItems[index - 1] = temp

  // Recalculate lot numbers
  newItems.forEach((item, idx) => {
    item.lot_number = idx + 1
  })

  items.value = newItems

  // Swap errors
  const newErrors = [...itemErrors.value]
  const tempError = newErrors[index]
  newErrors[index] = newErrors[index - 1]
  newErrors[index - 1] = tempError
  itemErrors.value = newErrors
}

function handleMoveDown(index) {
  if (index === items.value.length - 1) return

  const newItems = [...items.value]
  const temp = newItems[index]
  newItems[index] = newItems[index + 1]
  newItems[index + 1] = temp

  // Recalculate lot numbers
  newItems.forEach((item, idx) => {
    item.lot_number = idx + 1
  })

  items.value = newItems

  // Swap errors
  const newErrors = [...itemErrors.value]
  const tempError = newErrors[index]
  newErrors[index] = newErrors[index + 1]
  newErrors[index + 1] = tempError
  itemErrors.value = newErrors
}

// Drag and drop handlers
function handleDragStart(index) {
  draggedIndex.value = index
}

function handleDrop({ fromIndex, toIndex }) {
  if (fromIndex === toIndex) return

  const newItems = [...items.value]
  const newErrors = [...itemErrors.value]

  // Remove dragged item
  const [draggedItem] = newItems.splice(fromIndex, 1)
  const [draggedError] = newErrors.splice(fromIndex, 1)

  // Insert at new position
  newItems.splice(toIndex, 0, draggedItem)
  newErrors.splice(toIndex, 0, draggedError)

  // Recalculate lot numbers
  newItems.forEach((item, idx) => {
    item.lot_number = idx + 1
  })

  items.value = newItems
  itemErrors.value = newErrors
  draggedIndex.value = null
}
</script>

<template>
  <div class="max-w-4xl mx-auto p-6">
    <!-- Header -->
    <div class="mb-8">
      <button
        @click="handleBack"
        class="flex items-center text-gray-600 hover:text-gray-900 mb-4"
      >
        <svg class="w-5 h-5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        戻る
      </button>
      <h1 class="text-3xl font-bold text-gray-900 mb-2">オークション編集</h1>
      <p class="text-gray-600">オークション情報を編集します</p>
    </div>

    <!-- Success Message Banner -->
    <div
      v-if="showSuccessMessage"
      class="mb-6 p-4 bg-green-50 border border-green-200 rounded-lg shadow-sm"
    >
      <div class="flex items-start justify-between">
        <div class="flex items-start gap-3">
          <svg
            class="w-6 h-6 text-green-600 mt-0.5 flex-shrink-0"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <div>
            <p class="font-semibold text-green-900">オークションを作成しました</p>
            <p class="text-sm text-green-800 mt-1">
              次に商品を追加してオークションを完成させましょう。下の「商品紐づけ管理」ボタンから商品を追加できます。
            </p>
          </div>
        </div>
        <button
          @click="showSuccessMessage = false"
          class="text-green-600 hover:text-green-800 flex-shrink-0 ml-3"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
    </div>

    <!-- Content -->
    <template v-else>
      <!-- Error Alert -->
      <Alert v-if="submitError" variant="destructive" class="mb-6">
        <p class="font-semibold">エラー</p>
        <p>{{ submitError }}</p>
      </Alert>

      <!-- Status Card -->
      <Card class="p-4 mb-6">
        <div class="flex items-center justify-between">
          <div class="flex items-center gap-3">
            <span class="text-gray-700">ステータス:</span>
            <AuctionStatusBadge :status="auctionStatus" />
          </div>
          <div v-if="!canEdit" class="text-sm text-yellow-700 bg-yellow-50 px-3 py-1 rounded">
            {{ canEditReason || '編集できません' }}
          </div>
          <div v-else class="text-sm text-green-700 bg-green-50 px-3 py-1 rounded">
            編集可能
          </div>
        </div>
      </Card>

      <!-- Quick Actions Card -->
      <Card class="p-4 mb-6">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-sm font-medium text-gray-700">商品管理</h3>
            <p class="text-xs text-gray-500 mt-1">未割当商品の追加・解除を行います</p>
          </div>
          <router-link
            :to="{ name: 'auction-items-assign', params: { id: auctionId } }"
            class="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
            </svg>
            商品紐づけ管理
          </router-link>
        </div>
      </Card>

      <form @submit.prevent="handleSubmit" class="space-y-6">
        <!-- Basic Information Section -->
        <AuctionBasicInfo
          v-model="basicInfo"
          :errors="basicInfoErrors"
          @validate="validateBasicInfoField"
          :disabled="!canEdit"
        />

        <!-- Items Section with Tabs -->
        <Card class="p-6">
          <h2 class="text-xl font-semibold mb-4">商品管理</h2>

          <Tabs v-model="activeTab" default-value="assigned">
            <TabsList class="mb-6">
              <TabsTrigger value="assigned">
                割当済み商品 ({{ items.length }})
              </TabsTrigger>
              <TabsTrigger value="create">
                新規作成
              </TabsTrigger>
              <TabsTrigger value="select">
                既存から選択
              </TabsTrigger>
            </TabsList>

            <!-- Assigned Items Tab -->
            <TabsContent value="assigned">
              <div>
                <div class="flex justify-between items-center mb-4">
                  <p class="text-sm text-gray-600">オークションに割り当てられた商品の一覧です</p>
                  <Button
                    v-if="canEdit"
                    type="button"
                    @click="handleAddItem"
                    variant="outline"
                    size="sm"
                  >
                    + 商品を追加
                  </Button>
                </div>

                <div v-if="items.length === 0" class="text-center py-8 text-gray-500 border border-dashed border-gray-300 rounded-lg">
                  商品がありません。「新規作成」または「既存から選択」タブから商品を追加してください。
                </div>

                <div v-else class="space-y-4">
                  <EditItemForm
                    v-for="(item, index) in items"
                    :key="item.id || index"
                    :model-value="item"
                    @update:model-value="(newValue) => { items[index] = newValue; handleItemUpdate(index, newValue) }"
                    :index="index"
                    :errors="itemErrors[index] || { name: '', description: '' }"
                    @validate="(field) => validateItemField(index, field)"
                    :can-move-up="index > 0 && canEdit"
                    :can-move-down="index < items.length - 1 && canEdit"
                    :can-delete="item.can_delete && items.length > 1 && canEdit"
                    :can-edit="item.can_edit && canEdit"
                    @move-up="handleMoveUp(index)"
                    @move-down="handleMoveDown(index)"
                    @delete="requestDeleteItem(index)"
                    @drag-start="handleDragStart"
                    @drop="handleDrop"
                  />
                </div>
              </div>
            </TabsContent>

            <!-- Create New Item Tab -->
            <TabsContent value="create">
              <ItemCreateInline
                :auction-id="auctionId"
                :loading="isLoading"
                @item-created="handleItemCreated"
              />
            </TabsContent>

            <!-- Select Existing Items Tab -->
            <TabsContent value="select">
              <ItemSelectList
                :auction-id="auctionId"
                v-model:selected-ids="selectedItemIds"
                @items-assigned="handleItemsAssigned"
              />
            </TabsContent>
          </Tabs>
        </Card>

        <!-- Action Buttons -->
        <div class="flex justify-between items-center pt-4">
          <Button type="button" @click="handleCancel" variant="outline">
            キャンセル
          </Button>
          <div class="flex gap-3">
            <Button
              type="submit"
              :disabled="isSubmitting || !canEdit"
              class="min-w-[120px]"
            >
              {{ isSubmitting ? '更新中...' : '更新する' }}
            </Button>
          </div>
        </div>
      </form>
    </template>

    <!-- Delete Confirmation Modal -->
    <DeleteConfirmModal
      :is-open="showDeleteModal"
      title="商品削除の確認"
      :message="`「${items[deleteItemIndex]?.name || ''}」を削除してもよろしいですか？この操作は取り消せません。`"
      :is-deleting="isDeleting"
      @confirm="confirmDeleteItem"
      @cancel="cancelDeleteItem"
    />
  </div>
</template>
