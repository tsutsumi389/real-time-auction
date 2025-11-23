<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuctionStore } from '@/stores/auction'
import Button from '@/components/ui/Button.vue'
import Alert from '@/components/ui/Alert.vue'
import AuctionBasicInfo from '@/components/admin/AuctionBasicInfo.vue'
import ItemList from '@/components/admin/ItemList.vue'

const router = useRouter()
const auctionStore = useAuctionStore()

// Form data
const basicInfo = ref({
  title: '',
  description: ''
})

const items = ref([
  { name: '', description: '', lot_number: 1 }
])

// Validation errors
const basicInfoErrors = ref({
  title: '',
  description: ''
})

const itemErrors = ref([
  { name: '', description: '' }
])

// Form state
const isSubmitting = ref(false)
const submitError = ref('')

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
  }
}

// Validation functions for items
function validateItemField(index, field) {
  if (!itemErrors.value[index]) {
    itemErrors.value[index] = { name: '', description: '' }
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
  const newLotNumber = items.value.length + 1
  items.value.push({
    name: '',
    description: '',
    lot_number: newLotNumber
  })
  itemErrors.value.push({ name: '', description: '' })
}

// Form submission
async function handleSubmit() {
  submitError.value = ''

  if (!validateForm()) {
    submitError.value = '入力内容に誤りがあります。エラーメッセージを確認してください。'
    return
  }

  isSubmitting.value = true

  try {
    const formData = {
      title: basicInfo.value.title,
      description: basicInfo.value.description,
      items: items.value
    }

    const result = await auctionStore.handleCreateAuction(formData)

    if (result) {
      // Success - redirect to auction list
      router.push({ name: 'auction-list' })
    } else {
      // Failed - show error
      submitError.value = auctionStore.error || 'オークションの作成に失敗しました'
    }
  } catch (error) {
    submitError.value = error.message || 'オークションの作成に失敗しました'
  } finally {
    isSubmitting.value = false
  }
}

function handleCancel() {
  if (confirm('入力内容が失われますが、よろしいですか?')) {
    router.push({ name: 'auction-list' })
  }
}
</script>

<template>
  <div class="max-w-4xl mx-auto p-6">
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900 mb-2">オークション作成</h1>
      <p class="text-gray-600">新しいオークションを作成します</p>
    </div>

    <!-- Error Alert -->
    <Alert v-if="submitError" variant="destructive" class="mb-6">
      <p class="font-semibold">エラー</p>
      <p>{{ submitError }}</p>
    </Alert>

    <form @submit.prevent="handleSubmit" class="space-y-6">
      <!-- Basic Information Section -->
      <AuctionBasicInfo
        v-model="basicInfo"
        :errors="basicInfoErrors"
        @validate="validateBasicInfoField"
      />

      <!-- Items Section -->
      <ItemList
        v-model="items"
        :errors="itemErrors"
        @validate="validateItemField"
        @add-item="handleAddItem"
      />

      <!-- Action Buttons -->
      <div class="flex justify-between items-center pt-4">
        <Button type="button" @click="handleCancel" variant="outline">
          キャンセル
        </Button>
        <div class="flex gap-3">
          <Button
            type="submit"
            :disabled="isSubmitting"
            class="min-w-[120px]"
          >
            {{ isSubmitting ? '作成中...' : '作成する' }}
          </Button>
        </div>
      </div>
    </form>
  </div>
</template>
