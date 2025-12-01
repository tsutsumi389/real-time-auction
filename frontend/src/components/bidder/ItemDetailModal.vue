<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="open"
        ref="modalBackdrop"
        class="fixed inset-0 z-50 flex items-center justify-center p-0 sm:p-4 bg-black bg-opacity-50"
        @click.self="handleClose"
        role="dialog"
        aria-modal="true"
        aria-labelledby="modal-title"
        aria-describedby="modal-description"
      >
        <div
          ref="modalContent"
          class="relative bg-white shadow-xl w-full h-full sm:h-auto sm:rounded-lg sm:max-w-6xl sm:max-h-[90vh] overflow-y-auto focus:outline-none"
          tabindex="-1"
        >
          <!-- 閉じるボタン -->
          <button
            ref="closeButton"
            @click="handleClose"
            class="absolute top-3 right-3 sm:top-4 sm:right-4 z-10 p-2 text-gray-400 hover:text-gray-600 bg-white rounded-full shadow-md transition-colors focus:outline-none focus:ring-2 focus:ring-blue-500"
            aria-label="モーダルを閉じる"
          >
            <svg class="h-5 w-5 sm:h-6 sm:w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>

          <!-- モーダルコンテンツ -->
          <div class="flex flex-col lg:grid lg:grid-cols-3 gap-4 sm:gap-6 p-4 sm:p-6">
            <!-- 左側: メディアギャラリー -->
            <div class="lg:col-span-2 order-1">
              <!-- メイン画像/動画 -->
              <div class="relative bg-gray-100 rounded-lg overflow-hidden mb-4 h-[50vh] sm:h-[400px] lg:h-[500px]">
                <img
                  v-if="currentMedia && currentMedia.media_type === 'image'"
                  :src="currentMedia.url"
                  :alt="item.name"
                  class="w-full h-full object-contain transition-opacity duration-200"
                  :class="{ 'opacity-0': isImageChanging }"
                  @load="handleImageLoad"
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
                <div v-else class="w-full h-full flex flex-col items-center justify-center text-gray-400">
                  <svg class="h-16 w-16 sm:h-24 sm:w-24 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                  </svg>
                  <span class="text-sm">画像がありません</span>
                </div>

                <!-- ナビゲーションボタン（常に表示、無効化で対応） -->
                <button
                  v-if="hasMultipleMedia"
                  @click="previousMedia"
                  :disabled="currentMediaIndex === 0"
                  class="absolute left-2 sm:left-4 top-1/2 transform -translate-y-1/2 p-2 sm:p-3 bg-white/90 rounded-full shadow-md transition-all focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-30 disabled:cursor-not-allowed hover:bg-white"
                  aria-label="前の画像"
                >
                  <svg class="h-5 w-5 sm:h-6 sm:w-6 text-gray-800" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
                  </svg>
                </button>

                <button
                  v-if="hasMultipleMedia"
                  @click="nextMedia"
                  :disabled="currentMediaIndex >= mediaList.length - 1"
                  class="absolute right-2 sm:right-4 top-1/2 transform -translate-y-1/2 p-2 sm:p-3 bg-white/90 rounded-full shadow-md transition-all focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-30 disabled:cursor-not-allowed hover:bg-white"
                  aria-label="次の画像"
                >
                  <svg class="h-5 w-5 sm:h-6 sm:w-6 text-gray-800" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
                  </svg>
                </button>

                <!-- メディアカウンター -->
                <div
                  v-if="hasMultipleMedia"
                  class="absolute bottom-3 right-3 px-3 py-1.5 bg-black/60 text-white text-xs sm:text-sm font-medium rounded-full backdrop-blur-sm"
                  aria-live="polite"
                  aria-atomic="true"
                >
                  {{ currentMediaIndex + 1 }} / {{ mediaList.length }}
                </div>
              </div>

              <!-- サムネイル一覧 -->
              <div
                v-if="hasMultipleMedia"
                class="flex gap-2 overflow-x-auto pb-2 scrollbar-thin"
                role="tablist"
                aria-label="メディアサムネイル"
              >
                <button
                  v-for="(media, index) in mediaList"
                  :key="media.id"
                  @click="selectMedia(index)"
                  role="tab"
                  :aria-selected="currentMediaIndex === index"
                  :aria-label="`${media.media_type === 'video' ? '動画' : '画像'} ${index + 1}枚目`"
                  :class="[
                    'flex-shrink-0 w-16 h-16 sm:w-20 sm:h-20 rounded-lg overflow-hidden border-2 transition-all focus:outline-none focus:ring-2 focus:ring-blue-500',
                    currentMediaIndex === index
                      ? 'border-blue-600 ring-2 ring-blue-300'
                      : 'border-gray-300 hover:border-gray-400 opacity-70 hover:opacity-100'
                  ]"
                >
                  <img
                    v-if="media.media_type === 'image'"
                    :src="media.thumbnail_url || media.url"
                    :alt="`サムネイル ${index + 1}`"
                    class="w-full h-full object-cover"
                    loading="lazy"
                  />
                  <div v-else class="w-full h-full bg-gray-200 flex items-center justify-center">
                    <svg class="h-6 w-6 sm:h-8 sm:w-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"></path>
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                  </div>
                </button>
              </div>
            </div>

            <!-- 右側: アイテム詳細情報 -->
            <div class="lg:col-span-1 order-2">
              <div class="lg:sticky lg:top-0">
                <!-- LOT番号 -->
                <div class="mb-3 sm:mb-4">
                  <span class="inline-block px-2 sm:px-3 py-1 text-xs sm:text-sm font-semibold text-blue-800 bg-blue-100 rounded">
                    LOT {{ item.lot_number }}
                  </span>
                  <!-- ステータスバッジ -->
                  <span
                    v-if="item.status === 'active'"
                    class="inline-block ml-2 px-2 py-1 text-xs font-semibold text-green-800 bg-green-100 rounded"
                  >
                    入札中
                  </span>
                  <span
                    v-else-if="item.status === 'ended'"
                    class="inline-block ml-2 px-2 py-1 text-xs font-semibold text-gray-600 bg-gray-100 rounded"
                  >
                    終了
                  </span>
                </div>

                <!-- 商品名 -->
                <h2 id="modal-title" class="text-xl sm:text-2xl font-bold text-gray-900 mb-3 sm:mb-4">
                  {{ item.name }}
                </h2>

                <!-- 価格情報 -->
                <div class="mb-4 sm:mb-6 p-3 sm:p-4 bg-gray-50 rounded-lg">
                  <div class="mb-2">
                    <p class="text-xs sm:text-sm text-gray-600 mb-0.5">開始価格</p>
                    <p class="text-2xl sm:text-3xl font-bold text-gray-900">
                      ¥{{ formatPrice(item.starting_price) }}
                    </p>
                  </div>
                  <!-- 現在価格（入札がある場合） -->
                  <div v-if="item.current_price && item.current_price > item.starting_price" class="mt-3 pt-3 border-t border-gray-200">
                    <p class="text-xs sm:text-sm text-green-600 mb-0.5">現在価格</p>
                    <p class="text-xl sm:text-2xl font-bold text-green-600">
                      ¥{{ formatPrice(item.current_price) }}
                    </p>
                  </div>
                </div>

                <!-- 説明文 -->
                <div class="mb-4 sm:mb-6">
                  <h3 class="text-base sm:text-lg font-semibold text-gray-900 mb-2">商品説明</h3>
                  <div
                    id="modal-description"
                    class="text-gray-600 text-sm whitespace-pre-wrap max-h-48 sm:max-h-64 overflow-y-auto"
                  >
                    {{ item.description || '商品の説明はありません' }}
                  </div>
                </div>

                <!-- メタデータ -->
                <div v-if="hasMetadata" class="mb-4 sm:mb-6">
                  <h3 class="text-base sm:text-lg font-semibold text-gray-900 mb-2">詳細情報</h3>
                  <div class="space-y-2 max-h-48 overflow-y-auto">
                    <div
                      v-for="(value, key) in item.metadata"
                      :key="key"
                      class="flex justify-between py-2 border-b border-gray-200 last:border-b-0"
                    >
                      <span class="text-xs sm:text-sm text-gray-600">{{ formatMetadataKey(key) }}</span>
                      <span class="text-xs sm:text-sm font-medium text-gray-900 text-right ml-2">{{ value }}</span>
                    </div>
                  </div>
                </div>

                <!-- 閉じるボタン -->
                <button
                  @click="handleClose"
                  class="w-full px-4 sm:px-6 py-2.5 sm:py-3 bg-gray-200 text-gray-800 rounded-lg hover:bg-gray-300 transition-colors font-medium text-sm sm:text-base focus:outline-none focus:ring-2 focus:ring-gray-400"
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
import { ref, computed, watch, onMounted, onUnmounted, nextTick } from 'vue'

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

// Refs
const modalBackdrop = ref(null)
const modalContent = ref(null)
const closeButton = ref(null)
const previousActiveElement = ref(null)

// State
const currentMediaIndex = ref(0)
const isImageChanging = ref(false)

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
  isImageChanging.value = false
}

const handleImageLoad = () => {
  isImageChanging.value = false
}

const previousMedia = () => {
  if (currentMediaIndex.value > 0) {
    isImageChanging.value = true
    setTimeout(() => {
      currentMediaIndex.value--
    }, 100)
  }
}

const nextMedia = () => {
  if (currentMediaIndex.value < mediaList.value.length - 1) {
    isImageChanging.value = true
    setTimeout(() => {
      currentMediaIndex.value++
    }, 100)
  }
}

const selectMedia = (index) => {
  if (currentMediaIndex.value !== index) {
    isImageChanging.value = true
    setTimeout(() => {
      currentMediaIndex.value = index
    }, 100)
  }
}

const handleClose = () => {
  emit('close')
}

const handleKeydown = (event) => {
  if (!props.open) return

  switch (event.key) {
    case 'Escape':
      event.preventDefault()
      handleClose()
      break
    case 'ArrowLeft':
      event.preventDefault()
      previousMedia()
      break
    case 'ArrowRight':
      event.preventDefault()
      nextMedia()
      break
  }
}

// Focus trap for accessibility
const trapFocus = (event) => {
  if (!props.open || !modalContent.value) return

  const focusableElements = modalContent.value.querySelectorAll(
    'button, [href], input, select, textarea, [tabindex]:not([tabindex="-1"])'
  )
  const firstElement = focusableElements[0]
  const lastElement = focusableElements[focusableElements.length - 1]

  if (event.shiftKey && document.activeElement === firstElement) {
    event.preventDefault()
    lastElement.focus()
  } else if (!event.shiftKey && document.activeElement === lastElement) {
    event.preventDefault()
    firstElement.focus()
  }
}

const handleTabKey = (event) => {
  if (event.key === 'Tab') {
    trapFocus(event)
  }
}

// Watchers
watch(() => props.open, async (newValue) => {
  if (newValue) {
    // Save the previously focused element
    previousActiveElement.value = document.activeElement

    currentMediaIndex.value = 0
    isImageChanging.value = false

    // Disable body scroll
    document.body.style.overflow = 'hidden'

    // Focus the close button when modal opens
    await nextTick()
    if (closeButton.value) {
      closeButton.value.focus()
    }

    // Add tab key listener for focus trap
    window.addEventListener('keydown', handleTabKey)
  } else {
    // Enable body scroll
    document.body.style.overflow = ''

    // Remove tab key listener
    window.removeEventListener('keydown', handleTabKey)

    // Restore focus to the previously focused element
    await nextTick()
    if (previousActiveElement.value && typeof previousActiveElement.value.focus === 'function') {
      previousActiveElement.value.focus()
    }
  }
})

// Lifecycle
onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
  window.removeEventListener('keydown', handleTabKey)
  // Cleanup
  document.body.style.overflow = ''
  previousActiveElement.value = null
})
</script>

<style scoped>
/* Modal animation */
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
  transform: scale(0.95);
}

/* Respect reduced motion preference */
@media (prefers-reduced-motion: reduce) {
  .modal-enter-active,
  .modal-leave-active {
    transition: none;
  }

  .modal-enter-active > div,
  .modal-leave-active > div {
    transition: none;
  }

  .transition-opacity,
  .transition-all,
  .transition-colors {
    transition: none !important;
  }
}

/* Scrollbar styling for thumbnail container */
.scrollbar-thin {
  scrollbar-width: thin;
  scrollbar-color: #cbd5e1 transparent;
}

.scrollbar-thin::-webkit-scrollbar {
  height: 6px;
}

.scrollbar-thin::-webkit-scrollbar-track {
  background: transparent;
}

.scrollbar-thin::-webkit-scrollbar-thumb {
  background-color: #cbd5e1;
  border-radius: 3px;
}

.scrollbar-thin::-webkit-scrollbar-thumb:hover {
  background-color: #94a3b8;
}

/* General scrollbar styling */
.overflow-x-auto::-webkit-scrollbar,
.overflow-y-auto::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

.overflow-x-auto::-webkit-scrollbar-track,
.overflow-y-auto::-webkit-scrollbar-track {
  background: #f1f5f9;
  border-radius: 3px;
}

.overflow-x-auto::-webkit-scrollbar-thumb,
.overflow-y-auto::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 3px;
}

.overflow-x-auto::-webkit-scrollbar-thumb:hover,
.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}
</style>
