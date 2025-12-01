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
      <h1 class="mt-4 text-3xl font-bold text-gray-900">商品編集</h1>
      <p class="mt-2 text-sm text-gray-600">商品情報を編集します</p>
    </div>

    <!-- ローディング -->
    <div v-if="pageLoading" class="bg-white shadow rounded-lg p-8 text-center">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-4 border-blue-500 border-t-transparent"></div>
      <p class="mt-2 text-gray-500">読み込み中...</p>
    </div>

    <!-- 商品が見つからない -->
    <div v-else-if="notFound" class="bg-white shadow rounded-lg p-8 text-center">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">商品が見つかりません</h3>
      <p class="mt-1 text-sm text-gray-500">指定された商品は存在しないか、削除された可能性があります。</p>
      <div class="mt-6">
        <router-link
          to="/admin/items"
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700"
        >
          商品一覧に戻る
        </router-link>
      </div>
    </div>

    <!-- 編集フォーム -->
    <template v-else>
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

      <!-- オークション紐づけ情報 -->
      <div v-if="item" class="mb-6 bg-white shadow rounded-lg p-4">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-sm font-medium text-gray-500">紐づけオークション</h3>
            <div class="mt-1">
              <router-link
                v-if="item.auction_id"
                :to="`/admin/auctions/${item.auction_id}/edit`"
                class="text-blue-600 hover:text-blue-900 font-medium"
              >
                {{ item.auction_title }}
              </router-link>
              <span v-else class="text-gray-400">未割当</span>
            </div>
          </div>
          <div v-if="item.lot_number > 0" class="text-sm text-gray-500">
            ロット番号: {{ item.lot_number }}
          </div>
        </div>
      </div>

      <!-- フォーム -->
      <form @submit.prevent="handleSubmit" class="bg-white shadow rounded-lg p-6">
        <!-- 編集不可の警告 -->
        <div v-if="item && !item.can_edit" class="mb-6 bg-yellow-50 border border-yellow-200 text-yellow-700 px-4 py-3 rounded">
          <p class="text-sm">この商品は既に開始されているため、編集できません。</p>
        </div>

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
            :disabled="item && !item.can_edit"
            class="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm disabled:bg-gray-100 disabled:cursor-not-allowed"
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
            :disabled="item && !item.can_edit"
            class="block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm disabled:bg-gray-100 disabled:cursor-not-allowed"
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
              :disabled="item && !item.can_edit"
              class="block w-full pl-3 pr-10 py-2 border border-gray-300 rounded-md shadow-sm placeholder-gray-400 focus:outline-none focus:ring-blue-500 focus:border-blue-500 sm:text-sm disabled:bg-gray-100 disabled:cursor-not-allowed"
              placeholder="0"
            />
            <div class="absolute inset-y-0 right-0 pr-3 flex items-center pointer-events-none">
              <span class="text-gray-500 sm:text-sm">pt</span>
            </div>
          </div>
          <p class="mt-1 text-sm text-gray-500">オークション開始時のポイント（任意）</p>
        </div>

        <!-- メディア管理セクション -->
        <div class="mb-6">
          <h3 class="text-sm font-medium text-gray-700 mb-3">メディア</h3>
          
          <!-- メディアアップローダー -->
          <MediaUploader
            :item-id="route.params.id"
            :current-image-count="imageCount"
            :current-video-count="videoCount"
            :disabled="item && !item.can_edit"
            @upload-success="handleUploadSuccess"
            @upload-error="handleUploadError"
          />
          
          <!-- メディアギャラリー -->
          <div class="mt-6">
            <div v-if="mediaLoading" class="text-center py-4">
              <div class="inline-block animate-spin rounded-full h-6 w-6 border-2 border-blue-500 border-t-transparent"></div>
              <p class="mt-2 text-sm text-gray-500">メディアを読み込み中...</p>
            </div>
            <MediaGallery
              v-else
              :item-id="route.params.id"
              :media="media"
              :disabled="item && !item.can_edit"
              @update:media="handleMediaUpdate"
              @delete-success="handleDeleteSuccess"
              @delete-error="handleUploadError"
            />
          </div>
        </div>

        <!-- ボタン -->
        <div class="flex justify-between">
          <!-- 削除ボタン -->
          <div>
            <button
              v-if="item"
              type="button"
              :disabled="!item.can_delete || deleteLoading"
              :title="getDeleteTooltip()"
              @click="openDeleteDialog"
              class="inline-flex items-center px-4 py-2 border border-red-300 text-sm font-medium rounded-md text-red-700 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              削除
            </button>
          </div>

          <!-- 保存・キャンセルボタン -->
          <div class="flex space-x-3">
            <router-link
              to="/admin/items"
              class="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              キャンセル
            </router-link>
            <button
              type="submit"
              :disabled="loading || (item && !item.can_edit)"
              class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="loading" class="mr-2">
                <svg class="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </span>
              {{ loading ? '保存中...' : '保存' }}
            </button>
          </div>
        </div>
      </form>
    </template>

    <!-- 削除確認ダイアログ -->
    <div v-if="showDeleteDialog" class="fixed inset-0 z-50 overflow-y-auto">
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <!-- オーバーレイ -->
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="closeDeleteDialog"></div>

        <!-- ダイアログ -->
        <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div class="sm:flex sm:items-start">
              <div class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-red-100 sm:mx-0 sm:h-10 sm:w-10">
                <svg class="h-6 w-6 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
              </div>
              <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
                <h3 class="text-lg leading-6 font-medium text-gray-900">商品を削除</h3>
                <div class="mt-2">
                  <p class="text-sm text-gray-500">
                    「{{ item?.name }}」を削除してもよろしいですか？この操作は取り消せません。
                  </p>
                </div>
              </div>
            </div>
          </div>
          <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
            <button
              type="button"
              :disabled="deleteLoading"
              @click="handleDelete"
              class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-red-600 text-base font-medium text-white hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50"
            >
              <span v-if="deleteLoading" class="mr-2">
                <svg class="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </span>
              {{ deleteLoading ? '削除中...' : '削除' }}
            </button>
            <button
              type="button"
              :disabled="deleteLoading"
              @click="closeDeleteDialog"
              class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm"
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
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useItemStore } from '@/stores/item'
import { getItemMedia } from '@/services/mediaApi'
import MediaUploader from '@/components/admin/MediaUploader.vue'
import MediaGallery from '@/components/admin/MediaGallery.vue'

const route = useRoute()
const router = useRouter()
const itemStore = useItemStore()

// ローカル状態
const pageLoading = ref(true)
const loading = ref(false)
const deleteLoading = ref(false)
const error = ref('')
const successMessage = ref('')
const notFound = ref(false)
const showDeleteDialog = ref(false)

// メディア関連の状態
const media = ref([])
const mediaLoading = ref(false)
const imageCount = computed(() => media.value.filter((m) => m.media_type === 'image').length)
const videoCount = computed(() => media.value.filter((m) => m.media_type === 'video').length)

const form = reactive({
  name: '',
  description: '',
  starting_price: null,
})

// 現在の商品
const item = computed(() => itemStore.currentItem)

// メディアを取得
async function fetchMedia() {
  const itemId = route.params.id
  mediaLoading.value = true
  try {
    const response = await getItemMedia(itemId)
    media.value = response.items || []
  } catch (err) {
    console.error('Failed to fetch media:', err)
    media.value = []
  } finally {
    mediaLoading.value = false
  }
}

// アップロード成功時
function handleUploadSuccess(uploadedMedia) {
  media.value.push(uploadedMedia)
  successMessage.value = 'メディアをアップロードしました'
  setTimeout(() => {
    successMessage.value = ''
  }, 3000)
}

// アップロードエラー時
function handleUploadError(err) {
  error.value = err.message || 'メディアのアップロードに失敗しました'
}

// 削除成功時
function handleDeleteSuccess(mediaId) {
  media.value = media.value.filter((m) => m.id !== mediaId)
}

// 順序変更時
function handleMediaUpdate(newMedia) {
  media.value = newMedia
}

// 初期データ取得
onMounted(async () => {
  const itemId = route.params.id
  const success = await itemStore.fetchItemDetail(itemId)

  if (!success || !itemStore.currentItem) {
    notFound.value = true
  } else {
    // フォームに値をセット
    form.name = itemStore.currentItem.name || ''
    form.description = itemStore.currentItem.description || ''
    form.starting_price = itemStore.currentItem.starting_price || null

    // メディアを取得
    await fetchMedia()
  }

  pageLoading.value = false
})

// 削除ボタンのツールチップ
function getDeleteTooltip() {
  if (!item.value) return ''
  if (item.value.auction_id) return 'オークションに紐づいているため削除できません'
  if (item.value.bid_count > 0) return '入札履歴があるため削除できません'
  return ''
}

// 削除ダイアログを開く
function openDeleteDialog() {
  if (item.value?.can_delete) {
    showDeleteDialog.value = true
  }
}

// 削除ダイアログを閉じる
function closeDeleteDialog() {
  showDeleteDialog.value = false
}

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

    const result = await itemStore.handleUpdateItem(route.params.id, data)

    if (result) {
      successMessage.value = '商品を更新しました'
      setTimeout(() => {
        successMessage.value = ''
      }, 3000)
    } else {
      error.value = itemStore.error || '商品の更新に失敗しました'
    }
  } catch (err) {
    error.value = err.message || '商品の更新に失敗しました'
  } finally {
    loading.value = false
  }
}

// 削除処理
async function handleDelete() {
  deleteLoading.value = true
  error.value = ''

  try {
    const success = await itemStore.handleDeleteItem(route.params.id)

    if (success) {
      // 成功時は一覧画面に戻る
      router.push('/admin/items')
    } else {
      error.value = itemStore.error || '商品の削除に失敗しました'
      showDeleteDialog.value = false
    }
  } catch (err) {
    error.value = err.message || '商品の削除に失敗しました'
    showDeleteDialog.value = false
  } finally {
    deleteLoading.value = false
  }
}
</script>
