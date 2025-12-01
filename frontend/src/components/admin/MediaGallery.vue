<script setup>
/**
 * MediaGallery Component
 * メディア一覧表示、並び替え、削除機能を持つギャラリーコンポーネント
 */
import { ref, computed, watch } from 'vue'
import { deleteItemMedia, reorderItemMedia } from '@/services/mediaApi'

const props = defineProps({
  itemId: {
    type: String,
    required: true,
  },
  media: {
    type: Array,
    default: () => [],
  },
  disabled: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['delete-success', 'delete-error', 'reorder-success', 'reorder-error', 'update:media'])

// ローカル状態
const localMedia = ref([...props.media])
const deletingId = ref(null)
const deleteConfirmId = ref(null)
const reordering = ref(false)
const error = ref('')
const successMessage = ref('')
const draggingIndex = ref(null)
const dragOverIndex = ref(null)

// propsの変更を監視
watch(
  () => props.media,
  (newMedia) => {
    localMedia.value = [...newMedia]
  },
  { deep: true }
)

// 画像と動画を分類
const imageMedia = computed(() => localMedia.value.filter((m) => m.media_type === 'image'))
const videoMedia = computed(() => localMedia.value.filter((m) => m.media_type === 'video'))

// 削除確認ダイアログを開く
function openDeleteConfirm(mediaId) {
  if (!props.disabled) {
    deleteConfirmId.value = mediaId
  }
}

// 削除確認ダイアログを閉じる
function closeDeleteConfirm() {
  deleteConfirmId.value = null
}

// メディアを削除
async function handleDelete(mediaId) {
  if (props.disabled) return

  deletingId.value = mediaId
  error.value = ''

  try {
    await deleteItemMedia(props.itemId, mediaId)

    // ローカル配列から削除
    localMedia.value = localMedia.value.filter((m) => m.id !== mediaId)
    emit('update:media', localMedia.value)
    emit('delete-success', mediaId)

    successMessage.value = 'メディアを削除しました'
    setTimeout(() => {
      successMessage.value = ''
    }, 3000)
  } catch (err) {
    error.value = err.message || '削除に失敗しました'
    emit('delete-error', err)
  } finally {
    deletingId.value = null
    deleteConfirmId.value = null
  }
}

// ドラッグ開始
function handleDragStart(index, e) {
  if (props.disabled) return

  draggingIndex.value = index
  e.dataTransfer.effectAllowed = 'move'
  e.dataTransfer.setData('text/plain', index)

  // ドラッグ中の要素に透明度を設定
  setTimeout(() => {
    e.target.style.opacity = '0.4'
  }, 0)
}

// ドラッグ終了
function handleDragEnd(e) {
  e.target.style.opacity = '1'
  draggingIndex.value = null
  dragOverIndex.value = null
}

// ドラッグオーバー
function handleDragOver(index, e) {
  e.preventDefault()
  if (props.disabled) return

  dragOverIndex.value = index
}

// ドラッグリーブ
function handleDragLeave() {
  dragOverIndex.value = null
}

// ドロップ
async function handleDrop(targetIndex, e) {
  e.preventDefault()
  if (props.disabled) return

  const sourceIndex = draggingIndex.value
  if (sourceIndex === null || sourceIndex === targetIndex) {
    draggingIndex.value = null
    dragOverIndex.value = null
    return
  }

  // ローカルで並び替え
  const newMedia = [...localMedia.value]
  const [removed] = newMedia.splice(sourceIndex, 1)
  newMedia.splice(targetIndex, 0, removed)

  // display_orderを更新
  newMedia.forEach((m, i) => {
    m.display_order = i + 1
  })

  localMedia.value = newMedia
  emit('update:media', newMedia)

  draggingIndex.value = null
  dragOverIndex.value = null

  // サーバーに保存
  await saveOrder()
}

// 順序をサーバーに保存
async function saveOrder() {
  if (props.disabled || reordering.value) return

  reordering.value = true
  error.value = ''

  try {
    const mediaOrder = localMedia.value.map((m, index) => ({
      id: m.id,
      display_order: index + 1,
    }))

    await reorderItemMedia(props.itemId, mediaOrder)
    emit('reorder-success')

    successMessage.value = '順序を保存しました'
    setTimeout(() => {
      successMessage.value = ''
    }, 2000)
  } catch (err) {
    error.value = err.message || '順序の保存に失敗しました'
    emit('reorder-error', err)
  } finally {
    reordering.value = false
  }
}

// 矢印ボタンで移動
async function moveMedia(index, direction) {
  if (props.disabled) return

  const newIndex = index + direction
  if (newIndex < 0 || newIndex >= localMedia.value.length) return

  const newMedia = [...localMedia.value]
  const [removed] = newMedia.splice(index, 1)
  newMedia.splice(newIndex, 0, removed)

  // display_orderを更新
  newMedia.forEach((m, i) => {
    m.display_order = i + 1
  })

  localMedia.value = newMedia
  emit('update:media', newMedia)

  await saveOrder()
}

// プレビューモーダル
const previewMedia = ref(null)

function openPreview(media) {
  previewMedia.value = media
}

function closePreview() {
  previewMedia.value = null
}
</script>

<template>
  <div class="media-gallery">
    <!-- 成功メッセージ -->
    <div
      v-if="successMessage"
      class="mb-4 p-3 bg-green-50 border border-green-200 rounded-lg text-green-700 text-sm"
    >
      {{ successMessage }}
    </div>

    <!-- エラーメッセージ -->
    <div
      v-if="error"
      class="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg text-red-700 text-sm"
    >
      {{ error }}
    </div>

    <!-- 空の状態 -->
    <div v-if="localMedia.length === 0" class="text-center py-8 text-gray-500">
      <svg
        class="mx-auto w-12 h-12 text-gray-400 mb-3"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
        />
      </svg>
      <p>メディアが登録されていません</p>
    </div>

    <!-- メディアグリッド -->
    <div
      v-else
      class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4"
    >
      <div
        v-for="(media, index) in localMedia"
        :key="media.id"
        :class="[
          'relative group rounded-lg overflow-hidden bg-gray-100 aspect-square',
          draggingIndex === index ? 'ring-2 ring-blue-500' : '',
          dragOverIndex === index && draggingIndex !== index ? 'ring-2 ring-green-500' : '',
          disabled ? 'cursor-default' : 'cursor-move',
        ]"
        :draggable="!disabled"
        @dragstart="handleDragStart(index, $event)"
        @dragend="handleDragEnd"
        @dragover="handleDragOver(index, $event)"
        @dragleave="handleDragLeave"
        @drop="handleDrop(index, $event)"
      >
        <!-- 画像 -->
        <template v-if="media.media_type === 'image'">
          <img
            :src="media.thumbnail_url || media.url"
            :alt="`Media ${index + 1}`"
            class="w-full h-full object-cover"
            loading="lazy"
            @click="openPreview(media)"
          />
        </template>

        <!-- 動画 -->
        <template v-else-if="media.media_type === 'video'">
          <div
            class="w-full h-full flex items-center justify-center bg-gray-800"
            @click="openPreview(media)"
          >
            <img
              v-if="media.thumbnail_url"
              :src="media.thumbnail_url"
              :alt="`Video ${index + 1}`"
              class="w-full h-full object-cover"
              loading="lazy"
            />
            <div v-else class="text-white text-center">
              <svg
                class="w-12 h-12 mx-auto"
                fill="currentColor"
                viewBox="0 0 24 24"
              >
                <path d="M8 5v14l11-7z" />
              </svg>
              <span class="text-sm">動画</span>
            </div>
          </div>
        </template>

        <!-- オーバーレイ（ホバー時） -->
        <div
          :class="[
            'absolute inset-0 bg-black bg-opacity-0 group-hover:bg-opacity-50 transition-all duration-200',
            'flex flex-col items-center justify-center gap-2 opacity-0 group-hover:opacity-100',
          ]"
        >
          <!-- 順序番号 -->
          <span class="absolute top-2 left-2 bg-black bg-opacity-60 text-white text-xs px-2 py-1 rounded">
            {{ index + 1 }}
          </span>

          <!-- 移動ボタン -->
          <div v-if="!disabled" class="flex gap-2">
            <button
              v-if="index > 0"
              @click.stop="moveMedia(index, -1)"
              class="p-2 bg-white rounded-full hover:bg-gray-100 shadow"
              title="前に移動"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
              </svg>
            </button>
            <button
              v-if="index < localMedia.length - 1"
              @click.stop="moveMedia(index, 1)"
              class="p-2 bg-white rounded-full hover:bg-gray-100 shadow"
              title="後に移動"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
              </svg>
            </button>
          </div>

          <!-- 削除ボタン -->
          <button
            v-if="!disabled"
            @click.stop="openDeleteConfirm(media.id)"
            :disabled="deletingId === media.id"
            class="p-2 bg-red-500 text-white rounded-full hover:bg-red-600 shadow disabled:opacity-50"
            title="削除"
          >
            <svg v-if="deletingId !== media.id" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
            <svg v-else class="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
            </svg>
          </button>
        </div>

        <!-- メディアタイプバッジ -->
        <span
          :class="[
            'absolute bottom-2 right-2 text-xs px-2 py-0.5 rounded',
            media.media_type === 'image' ? 'bg-blue-500 text-white' : 'bg-purple-500 text-white',
          ]"
        >
          {{ media.media_type === 'image' ? '画像' : '動画' }}
        </span>
      </div>
    </div>

    <!-- 並び替え中のインジケータ -->
    <div v-if="reordering" class="mt-4 text-center text-sm text-gray-500">
      <svg class="inline-block w-4 h-4 animate-spin mr-1" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4" />
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z" />
      </svg>
      順序を保存中...
    </div>

    <!-- ドラッグ&ドロップのヒント -->
    <p v-if="!disabled && localMedia.length > 1" class="mt-4 text-xs text-gray-500 text-center">
      ドラッグ&ドロップで並び替えができます
    </p>

    <!-- 削除確認ダイアログ -->
    <div
      v-if="deleteConfirmId !== null"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black bg-opacity-50"
      @click.self="closeDeleteConfirm"
    >
      <div class="bg-white rounded-lg shadow-xl max-w-sm w-full p-6">
        <h3 class="text-lg font-medium text-gray-900 mb-4">削除の確認</h3>
        <p class="text-sm text-gray-600 mb-6">
          このメディアを削除してもよろしいですか？この操作は取り消せません。
        </p>
        <div class="flex justify-end gap-3">
          <button
            @click="closeDeleteConfirm"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
          >
            キャンセル
          </button>
          <button
            @click="handleDelete(deleteConfirmId)"
            :disabled="deletingId !== null"
            class="px-4 py-2 text-sm font-medium text-white bg-red-600 border border-transparent rounded-md hover:bg-red-700 disabled:opacity-50"
          >
            削除
          </button>
        </div>
      </div>
    </div>

    <!-- プレビューモーダル -->
    <div
      v-if="previewMedia"
      class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black bg-opacity-90"
      @click.self="closePreview"
    >
      <button
        @click="closePreview"
        class="absolute top-4 right-4 p-2 text-white hover:text-gray-300"
      >
        <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>

      <!-- 画像プレビュー -->
      <img
        v-if="previewMedia.media_type === 'image'"
        :src="previewMedia.url"
        alt="Preview"
        class="max-w-full max-h-full object-contain"
      />

      <!-- 動画プレビュー -->
      <video
        v-else-if="previewMedia.media_type === 'video'"
        :src="previewMedia.url"
        controls
        class="max-w-full max-h-full"
      />
    </div>
  </div>
</template>

<style scoped>
.media-gallery {
  width: 100%;
}
</style>
