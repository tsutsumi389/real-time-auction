<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuctionStore } from '@/stores/auction'
import Card from '@/components/ui/Card.vue'
import Label from '@/components/ui/Label.vue'
import Input from '@/components/ui/Input.vue'
import Button from '@/components/ui/Button.vue'
import Alert from '@/components/ui/Alert.vue'

const router = useRouter()
const auctionStore = useAuctionStore()

// Form data
const formData = ref({
  title: '',
  description: '',
  items: [
    { name: '', description: '', lot_number: 0 }
  ]
})

// Validation errors
const errors = ref({
  title: '',
  description: '',
  items: []
})

// Form state
const isSubmitting = ref(false)
const submitError = ref('')

// Computed
const hasItems = computed(() => formData.value.items.length > 0)

// Validation functions
function validateTitle() {
  if (!formData.value.title.trim()) {
    errors.value.title = 'タイトルを入力してください'
    return false
  }
  if (formData.value.title.length > 200) {
    errors.value.title = 'タイトルは200文字以内で入力してください'
    return false
  }
  errors.value.title = ''
  return true
}

function validateDescription() {
  if (formData.value.description.length > 2000) {
    errors.value.description = '説明は2000文字以内で入力してください'
    return false
  }
  errors.value.description = ''
  return true
}

function validateItemName(index) {
  if (!errors.value.items[index]) {
    errors.value.items[index] = { name: '', description: '' }
  }

  if (!formData.value.items[index].name.trim()) {
    errors.value.items[index].name = '商品名を入力してください'
    return false
  }
  if (formData.value.items[index].name.length > 200) {
    errors.value.items[index].name = '商品名は200文字以内で入力してください'
    return false
  }
  errors.value.items[index].name = ''
  return true
}

function validateItemDescription(index) {
  if (!errors.value.items[index]) {
    errors.value.items[index] = { name: '', description: '' }
  }

  if (formData.value.items[index].description.length > 2000) {
    errors.value.items[index].description = '説明は2000文字以内で入力してください'
    return false
  }
  errors.value.items[index].description = ''
  return true
}

function validateForm() {
  let isValid = true

  // Validate title
  if (!validateTitle()) {
    isValid = false
  }

  // Validate description
  if (!validateDescription()) {
    isValid = false
  }

  // Validate items
  formData.value.items.forEach((_, index) => {
    if (!validateItemName(index)) {
      isValid = false
    }
    if (!validateItemDescription(index)) {
      isValid = false
    }
  })

  return isValid
}

// Item management
function addItem() {
  const newLotNumber = formData.value.items.length
  formData.value.items.push({
    name: '',
    description: '',
    lot_number: newLotNumber
  })
  errors.value.items.push({ name: '', description: '' })
}

function removeItem(index) {
  if (formData.value.items.length === 1) {
    alert('最低1つの商品が必要です')
    return
  }

  if (confirm('この商品を削除してもよろしいですか?')) {
    formData.value.items.splice(index, 1)
    errors.value.items.splice(index, 1)

    // Recalculate lot numbers
    formData.value.items.forEach((item, idx) => {
      item.lot_number = idx
    })
  }
}

function moveItemUp(index) {
  if (index === 0) return

  const temp = formData.value.items[index]
  formData.value.items[index] = formData.value.items[index - 1]
  formData.value.items[index - 1] = temp

  // Recalculate lot numbers
  formData.value.items.forEach((item, idx) => {
    item.lot_number = idx
  })

  // Swap errors
  const tempError = errors.value.items[index]
  errors.value.items[index] = errors.value.items[index - 1]
  errors.value.items[index - 1] = tempError
}

function moveItemDown(index) {
  if (index === formData.value.items.length - 1) return

  const temp = formData.value.items[index]
  formData.value.items[index] = formData.value.items[index + 1]
  formData.value.items[index + 1] = temp

  // Recalculate lot numbers
  formData.value.items.forEach((item, idx) => {
    item.lot_number = idx
  })

  // Swap errors
  const tempError = errors.value.items[index]
  errors.value.items[index] = errors.value.items[index + 1]
  errors.value.items[index + 1] = tempError
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
    const result = await auctionStore.handleCreateAuction(formData.value)

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
      <Card class="p-6">
        <h2 class="text-xl font-semibold mb-4">基本情報</h2>

        <div class="space-y-4">
          <!-- Title -->
          <div>
            <Label for="title" class="required">タイトル</Label>
            <Input
              id="title"
              v-model="formData.title"
              type="text"
              placeholder="例: 2025年春季競走馬セリ"
              maxlength="200"
              @blur="validateTitle"
              :class="{ 'border-red-500': errors.title }"
            />
            <p v-if="errors.title" class="text-red-500 text-sm mt-1">{{ errors.title }}</p>
            <p class="text-gray-500 text-sm mt-1">{{ formData.title.length }}/200文字</p>
          </div>

          <!-- Description -->
          <div>
            <Label for="description">説明</Label>
            <textarea
              id="description"
              v-model="formData.description"
              placeholder="オークションの概要を入力してください"
              rows="4"
              maxlength="2000"
              @blur="validateDescription"
              :class="[
                'w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500',
                { 'border-red-500': errors.description }
              ]"
            ></textarea>
            <p v-if="errors.description" class="text-red-500 text-sm mt-1">{{ errors.description }}</p>
            <p class="text-gray-500 text-sm mt-1">{{ formData.description.length }}/2000文字</p>
          </div>
        </div>
      </Card>

      <!-- Items Section -->
      <Card class="p-6">
        <div class="flex justify-between items-center mb-4">
          <h2 class="text-xl font-semibold">商品情報</h2>
          <Button type="button" @click="addItem" variant="outline" size="sm">
            + 商品を追加
          </Button>
        </div>

        <div class="space-y-4">
          <div
            v-for="(item, index) in formData.items"
            :key="index"
            class="border border-gray-200 rounded-lg p-4"
          >
            <div class="flex justify-between items-center mb-3">
              <h3 class="font-medium text-gray-700">商品 #{{ index + 1 }}</h3>
              <div class="flex gap-2">
                <Button
                  type="button"
                  @click="moveItemUp(index)"
                  :disabled="index === 0"
                  variant="outline"
                  size="sm"
                >
                  ▲
                </Button>
                <Button
                  type="button"
                  @click="moveItemDown(index)"
                  :disabled="index === formData.items.length - 1"
                  variant="outline"
                  size="sm"
                >
                  ▼
                </Button>
                <Button
                  type="button"
                  @click="removeItem(index)"
                  variant="destructive"
                  size="sm"
                >
                  削除
                </Button>
              </div>
            </div>

            <div class="space-y-3">
              <!-- Item Name -->
              <div>
                <Label :for="`item-name-${index}`" class="required">商品名</Label>
                <Input
                  :id="`item-name-${index}`"
                  v-model="item.name"
                  type="text"
                  placeholder="例: ディープインパクト産駒"
                  maxlength="200"
                  @blur="validateItemName(index)"
                  :class="{ 'border-red-500': errors.items[index]?.name }"
                />
                <p v-if="errors.items[index]?.name" class="text-red-500 text-sm mt-1">
                  {{ errors.items[index].name }}
                </p>
                <p class="text-gray-500 text-sm mt-1">{{ item.name.length }}/200文字</p>
              </div>

              <!-- Item Description -->
              <div>
                <Label :for="`item-description-${index}`">商品説明</Label>
                <textarea
                  :id="`item-description-${index}`"
                  v-model="item.description"
                  placeholder="商品の詳細情報を入力してください"
                  rows="3"
                  maxlength="2000"
                  @blur="validateItemDescription(index)"
                  :class="[
                    'w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500',
                    { 'border-red-500': errors.items[index]?.description }
                  ]"
                ></textarea>
                <p v-if="errors.items[index]?.description" class="text-red-500 text-sm mt-1">
                  {{ errors.items[index].description }}
                </p>
                <p class="text-gray-500 text-sm mt-1">{{ item.description.length }}/2000文字</p>
              </div>
            </div>
          </div>
        </div>
      </Card>

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

<style scoped>
.required::after {
  content: ' *';
  color: #ef4444;
}
</style>
