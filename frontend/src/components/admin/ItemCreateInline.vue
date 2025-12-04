<script setup>
import { ref } from 'vue'
import { createItem, assignItemsToAuction } from '@/services/itemApi'
import Button from '@/components/ui/Button.vue'
import Input from '@/components/ui/Input.vue'
import Textarea from '@/components/ui/Textarea.vue'
import Alert from '@/components/ui/Alert.vue'

const props = defineProps({
  auctionId: {
    type: String,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['item-created'])

// Form data
const formData = ref({
  name: '',
  description: '',
  starting_price: null
})

// Validation errors
const errors = ref({
  name: '',
  description: '',
  starting_price: ''
})

// Loading state
const isCreating = ref(false)
const submitError = ref('')

// Validation functions
function validateField(field) {
  if (field === 'name') {
    if (!formData.value.name.trim()) {
      errors.value.name = '商品名を入力してください'
      return false
    }
    if (formData.value.name.length > 200) {
      errors.value.name = '商品名は200文字以内で入力してください'
      return false
    }
    errors.value.name = ''
    return true
  } else if (field === 'description') {
    if (formData.value.description.length > 2000) {
      errors.value.description = '説明は2000文字以内で入力してください'
      return false
    }
    errors.value.description = ''
    return true
  } else if (field === 'starting_price') {
    const price = formData.value.starting_price
    if (price !== null && price !== '' && price < 1) {
      errors.value.starting_price = '開始価格は1以上で入力してください'
      return false
    }
    errors.value.starting_price = ''
    return true
  }
}

function validateForm() {
  let isValid = true

  if (!validateField('name')) {
    isValid = false
  }
  if (!validateField('description')) {
    isValid = false
  }
  if (!validateField('starting_price')) {
    isValid = false
  }

  return isValid
}

// Reset form
function resetForm() {
  formData.value = {
    name: '',
    description: '',
    starting_price: null
  }
  errors.value = {
    name: '',
    description: '',
    starting_price: ''
  }
  submitError.value = ''
}

// Handle form submission
async function handleSubmit() {
  submitError.value = ''

  if (!validateForm()) {
    submitError.value = '入力内容に誤りがあります。エラーメッセージを確認してください。'
    return
  }

  isCreating.value = true

  try {
    // Step 1: Create item (without auction_id)
    const itemData = {
      name: formData.value.name,
      description: formData.value.description || '',
      starting_price: formData.value.starting_price
    }

    const createdItem = await createItem(itemData)

    // Step 2: Assign item to auction
    await assignItemsToAuction(props.auctionId, [createdItem.id])

    // Success
    resetForm()
    emit('item-created')
  } catch (error) {
    submitError.value = error.response?.data?.error || error.message || '商品の作成に失敗しました'
  } finally {
    isCreating.value = false
  }
}
</script>

<template>
  <div class="space-y-4">
    <!-- Error Alert -->
    <Alert v-if="submitError" variant="destructive" class="mb-4">
      <p class="font-semibold">エラー</p>
      <p>{{ submitError }}</p>
    </Alert>

    <!-- Form -->
    <form @submit.prevent="handleSubmit" class="space-y-4">
      <!-- Item Name -->
      <div>
        <label for="item-name" class="block text-sm font-medium text-gray-700 mb-1">
          商品名 <span class="text-red-500">*</span>
        </label>
        <Input
          id="item-name"
          v-model="formData.name"
          placeholder="商品名を入力してください"
          :disabled="isCreating || loading"
          @blur="validateField('name')"
          :class="{ 'border-red-500': errors.name }"
        />
        <p v-if="errors.name" class="mt-1 text-sm text-red-600">{{ errors.name }}</p>
      </div>

      <!-- Item Description -->
      <div>
        <label for="item-description" class="block text-sm font-medium text-gray-700 mb-1">
          説明
        </label>
        <Textarea
          id="item-description"
          v-model="formData.description"
          placeholder="商品の説明を入力してください"
          :disabled="isCreating || loading"
          @blur="validateField('description')"
          :class="{ 'border-red-500': errors.description }"
          rows="4"
        />
        <p v-if="errors.description" class="mt-1 text-sm text-red-600">{{ errors.description }}</p>
      </div>

      <!-- Starting Price -->
      <div>
        <label for="item-starting-price" class="block text-sm font-medium text-gray-700 mb-1">
          開始価格
        </label>
        <Input
          id="item-starting-price"
          v-model.number="formData.starting_price"
          type="number"
          min="1"
          placeholder="開始価格を入力してください"
          :disabled="isCreating || loading"
          @blur="validateField('starting_price')"
          :class="{ 'border-red-500': errors.starting_price }"
        />
        <p v-if="errors.starting_price" class="mt-1 text-sm text-red-600">{{ errors.starting_price }}</p>
      </div>

      <!-- Submit Button -->
      <div class="flex justify-end gap-3">
        <Button
          type="button"
          variant="outline"
          @click="resetForm"
          :disabled="isCreating || loading"
        >
          クリア
        </Button>
        <Button
          type="submit"
          :disabled="isCreating || loading"
          class="min-w-[120px]"
        >
          {{ isCreating ? '作成中...' : '作成して追加' }}
        </Button>
      </div>
    </form>
  </div>
</template>
