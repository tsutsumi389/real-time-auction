<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuctionStore } from '@/stores/auction'
import Button from '@/components/ui/Button.vue'
import Alert from '@/components/ui/Alert.vue'
import AuctionBasicInfo from '@/components/admin/AuctionBasicInfo.vue'

const router = useRouter()
const auctionStore = useAuctionStore()

// Form data
const basicInfo = ref({
  title: '',
  description: '',
  started_at: ''
})

// Validation errors
const basicInfoErrors = ref({
  title: '',
  description: '',
  started_at: ''
})

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
  } else if (field === 'started_at') {
    // started_at is optional, no validation needed
    basicInfoErrors.value.started_at = ''
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

  // Validate started_at
  if (!validateBasicInfoField('started_at')) {
    isValid = false
  }

  return isValid
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
    // Format started_at to ISO 8601 if provided
    let startedAt = null
    if (basicInfo.value.started_at) {
      startedAt = new Date(basicInfo.value.started_at).toISOString()
    }

    const formData = {
      title: basicInfo.value.title,
      description: basicInfo.value.description,
      started_at: startedAt,
      items: []
    }

    const result = await auctionStore.handleCreateAuction(formData)

    if (result) {
      // Success - redirect to auction edit page with created flag
      router.push({
        name: 'auction-edit',
        params: { id: result.id },
        query: { created: 'true' }
      })
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

      <!-- Info Message -->
      <div class="rounded-lg border border-blue-200 bg-blue-50 p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-blue-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-blue-800">
              商品はオークション作成後に追加できます
            </h3>
            <div class="mt-2 text-sm text-blue-700">
              <p>
                オークションを作成すると、編集画面に移動します。そこで商品を追加してオークションを完成させましょう。
              </p>
            </div>
          </div>
        </div>
      </div>

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
