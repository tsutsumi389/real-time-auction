<template>
  <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- ヘッダー -->
    <div class="mb-8">
      <router-link
        to="/admin/items"
        class="inline-flex items-center text-sm text-gray-500 hover:text-gray-700"
      >
        <svg class="mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        商品一覧に戻る
      </router-link>
      <h1 class="mt-4 text-3xl font-bold text-gray-900">新規商品作成</h1>
      <p class="mt-2 text-sm text-gray-600">オークションに出品する商品を登録します</p>
    </div>

    <!-- エラー表示 -->
    <div v-if="error" class="mb-6">
      <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded relative" role="alert">
        <span class="block sm:inline">{{ error }}</span>
        <button
          @click="error = ''"
          class="absolute top-0 bottom-0 right-0 px-4 py-3"
          aria-label="エラーを閉じる"
        >
          <span class="text-2xl">&times;</span>
        </button>
      </div>
    </div>

    <!-- フォーム -->
    <form @submit.prevent="handleSubmit" class="bg-white shadow rounded-lg p-6">
      <!-- 商品名 -->
      <div class="mb-6">
        <label for="name" class="block text-sm font-medium text-gray-700 mb-1">
          商品名 <span class="text-red-500">*</span>
        </label>
        <input
          id="name"
          v-model="form.name"
          type="text"
          required
          maxlength="200"
          class="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
          placeholder="商品名を入力"
        />
        <p class="mt-1 text-sm text-gray-500">最大200文字</p>
      </div>

      <!-- 説明 -->
      <div class="mb-6">
        <label for="description" class="block text-sm font-medium text-gray-700 mb-1">
          説明
        </label>
        <textarea
          id="description"
          v-model="form.description"
          rows="4"
          maxlength="2000"
          class="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
          placeholder="商品の説明を入力（任意）"
        ></textarea>
        <p class="mt-1 text-sm text-gray-500">最大2000文字</p>
      </div>

      <!-- 開始価格 -->
      <div class="mb-6">
        <label for="starting_price" class="block text-sm font-medium text-gray-700 mb-1">
          開始価格
        </label>
        <div class="relative rounded-md shadow-sm">
          <input
            id="starting_price"
            v-model.number="form.starting_price"
            type="number"
            min="1"
            class="block w-full pl-3 pr-10 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm"
            placeholder="0"
          />
          <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
            <span class="text-gray-500 sm:text-sm">pt</span>
          </div>
        </div>
        <p class="mt-1 text-sm text-gray-500">オークション開始時のポイント（任意）</p>
      </div>

      <!-- メディアについての案内 -->
      <div class="mb-6 p-4 bg-blue-50 border border-blue-200 rounded-lg">
        <div class="flex items-start gap-3">
          <svg class="w-5 h-5 text-blue-500 flex-shrink-0 mt-0.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <div>
            <h4 class="text-sm font-medium text-blue-800">画像・動画について</h4>
            <p class="mt-1 text-sm text-blue-600">
              画像・動画は商品の作成後に追加できます。作成後、編集画面からアップロードしてください。
            </p>
          </div>
        </div>
      </div>

      <!-- ボタン -->
      <div class="flex justify-end space-x-3">
        <router-link
          to="/admin/items"
          class="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
        >
          キャンセル
        </router-link>
        <button
          type="submit"
          :disabled="loading"
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <span v-if="loading" class="mr-2">
            <svg class="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </span>
          {{ loading ? '作成中...' : '作成' }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useItemStore } from '@/stores/item'

const router = useRouter()
const itemStore = useItemStore()

// ローカル状態
const loading = ref(false)
const error = ref('')

const form = reactive({
  name: '',
  description: '',
  starting_price: null,
})

// 送信処理
async function handleSubmit() {
  // バリデーション
  if (!form.name.trim()) {
    error.value = '商品名は必須です'
    return
  }

  loading.value = true
  error.value = ''

  try {
    const data = {
      name: form.name.trim(),
      description: form.description.trim(),
    }

    // 開始価格が設定されている場合のみ追加
    if (form.starting_price && form.starting_price > 0) {
      data.starting_price = form.starting_price
    }

    const result = await itemStore.handleCreateItem(data)

    if (result) {
      // 成功時は編集画面に遷移（メディアをすぐにアップロードできるように）
      router.push(`/admin/items/${result.id}/edit`)
    } else {
      error.value = itemStore.error || '商品の作成に失敗しました'
    }
  } catch (err) {
    error.value = err.message || '商品の作成に失敗しました'
  } finally {
    loading.value = false
  }
}
</script>
