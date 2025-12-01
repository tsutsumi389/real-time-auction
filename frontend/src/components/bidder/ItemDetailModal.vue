<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="open"
        class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black bg-opacity-50"
        @click.self="handleClose"
      >
        <div
          class="relative bg-white rounded-lg shadow-xl max-w-6xl w-full max-h-[90vh] overflow-y-auto"
          role="dialog"
          aria-labelledby="modal-title"
          aria-describedby="modal-description"
        >
          <!-- 閉じるボタン -->
          <button
            @click="handleClose"
            class="absolute top-4 right-4 z-10 p-2 text-gray-400 hover:text-gray-600 bg-white rounded-full shadow-md transition-colors"
            aria-label="モーダルを閉じる"
          >
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>

          <!-- モーダルコンテンツ -->
          <div class="grid grid-cols-1 lg:grid-cols-3 gap-6 p-6">
            <!-- 左側: メディアギャラリー -->
            <div class="lg:col-span-2">
              <!-- メイン画像/動画 -->
              <div class="relative bg-gray-100 rounded-lg overflow-hidden mb-4" style="height: 500px;">
                <img
                  v-if="currentMedia && currentMedia.media_type === 'image'"
                  :src="currentMedia.url"
                  :alt="item.name"
                  class="w-full h-full object-contain"
                  @error="handleImageError"
                />
                <video
                  v-else-if="currentMedia && currentMedia.media_type === 'video'"
                  :src="currentMedia.url"
                  controls
                  class="w-full h-full object-contain"
                >
                  お使いのブラウザは動画タグをサポートしていません。
                </video>
                <div v-else class="w-full h-full flex items-center justify-center text-gray-400">
                  <svg class="h-24 w-24" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                  </svg>
                </div>

                <!-- ナビゲーションボタン -->
                <button
                  v-if="hasMultipleMedia && currentMediaIndex > 0"
                  @click="previousMedia"
                  class="absolute left-4 top-1/2 transform -translate-y-1/2 p-2 bg-white bg-opacity-80 rounded-full shadow-md hover:bg-opacity-100 transition-all"
                  aria-label="前の画像"
                >
                  <svg class="h-6 w-6 text-gray-800" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
                  </svg>
                </button>

                <button
                  v-if="hasMultipleMedia && currentMediaIndex < mediaList.length - 1"
                  @click="nextMedia"
                  class="absolute right-4 top-1/2 transform -translate-y-1/2 p-2 bg-white bg-opacity-80 rounded-full shadow-md hover:bg-opacity-100 transition-all"
                  aria-label="次の画像"
                >
                  <svg class="h-6 w-6 text-gray-800" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
                  </svg>
                </button>
              </div>

              <!-- サムネイル一覧 -->
              <div v-if="hasMultipleMedia" class="flex gap-2 overflow-x-auto pb-2">
                <button
                  v-for="(media, index) in mediaList"
                  :key="media.id"
                  @click="selectMedia(index)"
                  :class="[
                    'flex-shrink-0 w-20 h-20 rounded-lg overflow-hidden border-2 transition-all',
                    currentMediaIndex === index
                      ? 'border-blue-600 ring-2 ring-blue-300'
                      : 'border-gray-300 hover:border-gray-400'
                  ]"
                >
                  <img
                    v-if="media.media_type === 'image'"
                    :src="media.thumbnail_url || media.url"
                    :alt="`サムネイル ${index + 1}`"
                    class="w-full h-full object-cover"
                  />
                  <div v-else class="w-full h-full bg-gray-200 flex items-center justify-center">
                    <svg class="h-8 w-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"></path>
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                  </div>
                </button>
              </div>
            </div>

            <!-- 右側: アイテム詳細情報 -->
            <div class="lg:col-span-1">
              <div class="sticky top-0">
                <!-- LOT番号 -->
                <div class="mb-4">
                  <span class="inline-block px-3 py-1 text-sm font-semibold text-blue-800 bg-blue-100 rounded">
                    LOT {{ item.lot_number }}
                  </span>
                </div>

                <!-- 商品名 -->
                <h2 id="modal-title" class="text-2xl font-bold text-gray-900 mb-4">
                  {{ item.name }}
                </h2>

                <!-- 開始価格 -->
                <div class="mb-6">
                  <p class="text-sm text-gray-600 mb-1">開始価格</p>
                  <p class="text-3xl font-bold text-gray-900">
                    ¥{{ formatPrice(item.starting_price) }}
                  </p>
                </div>

                <!-- 説明文 -->
                <div class="mb-6">
                  <h3 class="text-lg font-semibold text-gray-900 mb-2">商品説明</h3>
                  <div id="modal-description" class="text-gray-600 text-sm whitespace-pre-wrap">
                    {{ item.description || '商品の説明はありません' }}
                  </div>
                </div>

                <!-- メタデータ -->
                <div v-if="hasMetadata" class="mb-6">
                  <h3 class="text-lg font-semibold text-gray-900 mb-2">詳細情報</h3>
                  <div class="space-y-2">
                    <div
                      v-for="(value, key) in item.metadata"
                      :key="key"
                      class="flex justify-between py-2 border-b border-gray-200"
                    >
                      <span class="text-sm text-gray-600">{{ formatMetadataKey(key) }}</span>
                      <span class="text-sm font-medium text-gray-900">{{ value }}</span>
                    </div>
                  </div>
                </div>

                <!-- 閉じるボタン -->
                <button
                  @click="handleClose"
                  class="w-full px-6 py-3 bg-gray-200 text-gray-800 rounded-lg hover:bg-gray-300 transition-colors font-medium"
                >
                  閉じる
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  item: {
    type: Object,
    required: true
  },
  open: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['close'])

// State
const currentMediaIndex = ref(0)

// Computed
const mediaList = computed(() => {
  if (!props.item.media || props.item.media.length === 0) return []
  return [...props.item.media].sort((a, b) => a.display_order - b.display_order)
})

const currentMedia = computed(() => {
  if (mediaList.value.length === 0) return null
  return mediaList.value[currentMediaIndex.value]
})

const hasMultipleMedia = computed(() => mediaList.value.length > 1)

const hasMetadata = computed(() => {
  return props.item.metadata && Object.keys(props.item.metadata).length > 0
})

// Methods
const formatPrice = (price) => {
  if (!price) return '0'
  return new Intl.NumberFormat('ja-JP').format(price)
}

const formatMetadataKey = (key) => {
  // キャメルケースやスネークケースを人間が読みやすい形式に変換
  return key
    .replace(/_/g, ' ')
    .replace(/([A-Z])/g, ' $1')
    .replace(/^./, (str) => str.toUpperCase())
}

const handleImageError = (event) => {
  event.target.style.display = 'none'
}

const previousMedia = () => {
  if (currentMediaIndex.value > 0) {
    currentMediaIndex.value--
  }
}

const nextMedia = () => {
  if (currentMediaIndex.value < mediaList.value.length - 1) {
    currentMediaIndex.value++
  }
}

const selectMedia = (index) => {
  currentMediaIndex.value = index
}

const handleClose = () => {
  emit('close')
}

const handleKeydown = (event) => {
  if (!props.open) return

  switch (event.key) {
    case 'Escape':
      handleClose()
      break
    case 'ArrowLeft':
      previousMedia()
      break
    case 'ArrowRight':
      nextMedia()
      break
  }
}

// Watchers
watch(() => props.open, (newValue) => {
  if (newValue) {
    currentMediaIndex.value = 0
    // モーダルが開いた時にbodyのスクロールを無効化
    document.body.style.overflow = 'hidden'
  } else {
    // モーダルが閉じた時にbodyのスクロールを有効化
    document.body.style.overflow = ''
  }
})

// Lifecycle
onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
  // クリーンアップ
  document.body.style.overflow = ''
})
</script>

<style scoped>
/* モーダルアニメーション */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active > div,
.modal-leave-active > div {
  transition: transform 0.3s ease;
}

.modal-enter-from > div,
.modal-leave-to > div {
  transform: scale(0.9);
}

/* スクロールバーのスタイリング */
.overflow-x-auto::-webkit-scrollbar {
  height: 8px;
}

.overflow-x-auto::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 4px;
}

.overflow-x-auto::-webkit-scrollbar-thumb {
  background: #888;
  border-radius: 4px;
}

.overflow-x-auto::-webkit-scrollbar-thumb:hover {
  background: #555;
}
</style>
