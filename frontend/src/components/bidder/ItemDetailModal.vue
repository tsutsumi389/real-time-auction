<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="open"
        ref="modalBackdrop"
        class="fixed inset-0 z-50 flex items-center justify-center p-0 sm:p-4 bg-black/80 backdrop-blur-md"
        @click.self="handleClose"
        role="dialog"
        aria-modal="true"
        aria-labelledby="modal-title"
        aria-describedby="modal-description"
      >
        <div
          ref="modalContent"
          class="relative lux-modal-glass border border-lux-gold/20 shadow-2xl w-full h-full sm:h-auto sm:rounded-2xl sm:max-w-6xl sm:max-h-[90vh] overflow-y-auto focus:outline-none"
          tabindex="-1"
        >
          <!-- 閉じるボタン -->
          <button
            ref="closeButton"
            @click="handleClose"
            class="absolute top-3 right-3 sm:top-4 sm:right-4 z-10 p-2.5 text-lux-gold/60 hover:text-lux-gold bg-lux-noir-light/90 backdrop-blur-sm rounded-full shadow-lg hover:shadow-xl border border-lux-gold/20 hover:border-lux-gold/40 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-lux-gold"
            aria-label="モーダルを閉じる"
          >
            <svg class="h-5 w-5 sm:h-6 sm:w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>

          <!-- モーダルコンテンツ -->
          <div class="flex flex-col lg:grid lg:grid-cols-3 gap-4 sm:gap-6 lg:gap-8 p-4 sm:p-6 lg:p-8">
            <!-- 左側: メディアギャラリー -->
            <div class="lg:col-span-2 order-1">
              <!-- メイン画像/動画 -->
              <div class="relative bg-lux-noir-light rounded-xl sm:rounded-2xl overflow-hidden mb-4 h-[50vh] sm:h-[400px] lg:h-[500px] shadow-inner border border-lux-gold/10">
                <Transition name="image-fade" mode="out-in">
                  <img
                    v-if="currentMedia && currentMedia.media_type === 'image'"
                    :key="currentMediaIndex"
                    :src="currentMedia.url"
                    :alt="`${item.name} - 画像 ${currentMediaIndex + 1}/${mediaList.length || 1}`"
                    class="w-full h-full object-contain"
                    @load="handleImageLoad"
                    @error="handleImageError"
                  />
                </Transition>
                <video
                  v-if="currentMedia && currentMedia.media_type === 'video'"
                  :src="currentMedia.url"
                  controls
                  class="w-full h-full object-contain"
                  :title="`${item.name} - 動画`"
                  :aria-label="`${item.name}の動画`"
                >
                  お使いのブラウザは動画タグをサポートしていません。
                </video>
                <div
                  v-if="!currentMedia"
                  class="w-full h-full flex flex-col items-center justify-center text-lux-gold/40 bg-gradient-to-br from-lux-noir-light to-lux-noir"
                  role="img"
                  aria-label="画像がありません"
                >
                  <svg class="h-16 w-16 sm:h-24 sm:w-24 mb-3 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                  </svg>
                  <span class="text-sm font-medium" aria-hidden="true">画像がありません</span>
                </div>

                <!-- ナビゲーションボタン（常に表示、無効化で対応） -->
                <button
                  v-if="hasMultipleMedia"
                  @click="previousMedia"
                  :disabled="currentMediaIndex === 0"
                  class="absolute left-2 sm:left-4 top-1/2 transform -translate-y-1/2 p-2.5 sm:p-3 bg-lux-noir-light/95 backdrop-blur-sm rounded-full shadow-lg border border-lux-gold/20 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-lux-gold disabled:opacity-30 disabled:cursor-not-allowed hover:bg-lux-noir-light hover:border-lux-gold/40 hover:scale-105 active:scale-95"
                  :aria-label="`前の画像へ（現在 ${currentMediaIndex + 1}枚目 / 全${mediaList.length}枚）`"
                >
                  <svg class="h-5 w-5 sm:h-6 sm:w-6 text-lux-gold" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
                  </svg>
                </button>

                <button
                  v-if="hasMultipleMedia"
                  @click="nextMedia"
                  :disabled="currentMediaIndex >= mediaList.length - 1"
                  class="absolute right-2 sm:right-4 top-1/2 transform -translate-y-1/2 p-2.5 sm:p-3 bg-lux-noir-light/95 backdrop-blur-sm rounded-full shadow-lg border border-lux-gold/20 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-lux-gold disabled:opacity-30 disabled:cursor-not-allowed hover:bg-lux-noir-light hover:border-lux-gold/40 hover:scale-105 active:scale-95"
                  :aria-label="`次の画像へ（現在 ${currentMediaIndex + 1}枚目 / 全${mediaList.length}枚）`"
                >
                  <svg class="h-5 w-5 sm:h-6 sm:w-6 text-lux-gold" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
                  </svg>
                </button>

                <!-- メディアカウンター -->
                <div
                  v-if="hasMultipleMedia"
                  class="absolute bottom-3 right-3 px-3 py-1.5 bg-black/70 backdrop-blur-sm text-white text-xs sm:text-sm font-medium rounded-full shadow-lg"
                  aria-live="polite"
                  aria-atomic="true"
                >
                  {{ currentMediaIndex + 1 }} / {{ mediaList.length }}
                </div>
              </div>

              <!-- サムネイル一覧 -->
              <div
                v-if="hasMultipleMedia"
                class="flex gap-2 sm:gap-3 overflow-x-auto pb-2 scrollbar-thin"
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
                    'flex-shrink-0 w-16 h-16 sm:w-20 sm:h-20 rounded-xl overflow-hidden border-2 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-lux-gold',
                    currentMediaIndex === index
                      ? 'border-lux-gold ring-2 ring-lux-gold/30 shadow-lg scale-105'
                      : 'border-lux-gold/20 hover:border-lux-gold/50 opacity-70 hover:opacity-100'
                  ]"
                >
                  <img
                    v-if="media.media_type === 'image'"
                    :src="media.thumbnail_url || media.url"
                    alt=""
                    class="w-full h-full object-cover"
                    loading="lazy"
                  />
                  <div v-else class="w-full h-full bg-lux-noir-light flex items-center justify-center" aria-hidden="true">
                    <svg class="h-6 w-6 sm:h-8 sm:w-8 text-lux-gold/40" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
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
                <div class="mb-4 flex items-center flex-wrap gap-2">
                  <span class="inline-flex items-center px-3 py-1.5 text-sm font-bold text-lux-gold bg-lux-gold/10 border border-lux-gold/30 rounded-full">
                    LOT {{ item.lot_number }}
                  </span>
                  <!-- ステータスバッジ -->
                  <span
                    v-if="item.status === 'active'"
                    class="inline-flex items-center px-3 py-1.5 text-sm font-semibold text-emerald-400 bg-emerald-500/10 border border-emerald-500/30 rounded-full"
                  >
                    <span class="w-2 h-2 bg-emerald-400 rounded-full mr-2 animate-pulse"></span>
                    入札中
                  </span>
                  <span
                    v-else-if="item.status === 'ended'"
                    class="inline-flex items-center px-3 py-1.5 text-sm font-semibold text-lux-silver/60 bg-lux-silver/10 border border-lux-silver/20 rounded-full"
                  >
                    終了
                  </span>
                </div>

                <!-- 商品名 -->
                <h2 id="modal-title" class="font-display text-xl sm:text-2xl lg:text-3xl font-bold text-lux-cream mb-4 leading-tight">
                  {{ item.name }}
                </h2>

                <!-- 価格情報 -->
                <div class="mb-6 p-4 sm:p-5 bg-gradient-to-br from-lux-gold/5 to-lux-gold/10 rounded-2xl border border-lux-gold/20" role="region" aria-label="価格情報">
                  <dl>
                    <div class="mb-2">
                      <dt class="text-xs sm:text-sm text-lux-silver/60 font-medium uppercase tracking-wide mb-1">開始価格</dt>
                      <dd class="text-2xl sm:text-3xl lg:text-4xl font-bold text-lux-cream tabular-nums">
                        {{ formatPrice(item.starting_price) }}pt
                      </dd>
                    </div>
                    <!-- 現在価格（入札がある場合） -->
                    <div v-if="item.current_price && item.current_price > item.starting_price" class="mt-4 pt-4 border-t border-lux-gold/20">
                      <dt class="text-xs sm:text-sm text-emerald-400 font-medium uppercase tracking-wide mb-1">現在価格</dt>
                      <dd class="text-xl sm:text-2xl lg:text-3xl font-bold text-emerald-400 tabular-nums">
                        {{ formatPrice(item.current_price) }}pt
                      </dd>
                    </div>
                  </dl>
                </div>

                <!-- 説明文 -->
                <div class="mb-6">
                  <h3 class="text-base sm:text-lg font-semibold text-lux-cream mb-3 flex items-center">
                    <svg class="h-5 w-5 mr-2 text-lux-gold/60" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                    </svg>
                    商品説明
                  </h3>
                  <div
                    id="modal-description"
                    class="text-lux-silver/80 text-sm leading-relaxed whitespace-pre-wrap max-h-48 sm:max-h-64 overflow-y-auto p-4 bg-lux-gold/5 border border-lux-gold/10 rounded-xl"
                  >
                    {{ item.description || '商品の説明はありません' }}
                  </div>
                </div>

                <!-- メタデータ -->
                <div v-if="hasMetadata" class="mb-6" role="region" aria-label="詳細情報">
                  <h3 class="text-base sm:text-lg font-semibold text-lux-cream mb-3 flex items-center">
                    <svg class="h-5 w-5 mr-2 text-lux-gold/60" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
                    </svg>
                    詳細情報
                  </h3>
                  <dl class="space-y-0.5 max-h-48 overflow-y-auto bg-lux-gold/5 border border-lux-gold/10 rounded-xl p-3">
                    <div
                      v-for="(value, key) in item.metadata"
                      :key="key"
                      class="flex justify-between items-center py-2.5 px-2 hover:bg-lux-gold/10 rounded-lg transition-colors"
                    >
                      <dt class="text-xs sm:text-sm text-lux-silver/60">{{ formatMetadataKey(key) }}</dt>
                      <dd class="text-xs sm:text-sm font-semibold text-lux-cream text-right ml-4">{{ value }}</dd>
                    </div>
                  </dl>
                </div>

                <!-- 閉じるボタン -->
                <button
                  @click="handleClose"
                  class="w-full px-6 py-3 bg-lux-gold/10 border border-lux-gold/30 text-lux-gold rounded-xl hover:bg-lux-gold/20 hover:border-lux-gold/40 transition-all duration-200 font-medium text-sm sm:text-base focus:outline-none focus:ring-2 focus:ring-lux-gold hover:shadow-md active:scale-[0.98]"
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
    currentMediaIndex.value--
  }
}

const nextMedia = () => {
  if (currentMediaIndex.value < mediaList.value.length - 1) {
    currentMediaIndex.value++
  }
}

const selectMedia = (index) => {
  if (currentMediaIndex.value !== index) {
    currentMediaIndex.value = index
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

    // Add inert attribute to main content for better accessibility
    const mainContent = document.getElementById('main-content')
    if (mainContent) {
      mainContent.setAttribute('inert', '')
      mainContent.setAttribute('aria-hidden', 'true')
    }

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

    // Remove inert attribute from main content
    const mainContent = document.getElementById('main-content')
    if (mainContent) {
      mainContent.removeAttribute('inert')
      mainContent.removeAttribute('aria-hidden')
    }

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

  // Remove inert attribute if still present
  const mainContent = document.getElementById('main-content')
  if (mainContent) {
    mainContent.removeAttribute('inert')
    mainContent.removeAttribute('aria-hidden')
  }
})
</script>

<style scoped>
/* Luxury modal glass effect */
.lux-modal-glass {
  background: linear-gradient(
    135deg,
    rgba(15, 15, 15, 0.95) 0%,
    rgba(10, 10, 10, 0.98) 100%
  );
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
}

.border-lux-gold\/20 {
  border-color: hsl(43 74% 49% / 0.2);
}

.border-lux-gold\/40 {
  border-color: hsl(43 74% 49% / 0.4);
}

.bg-lux-noir-light\/90 {
  background-color: hsl(0 0% 8% / 0.9);
}

.text-lux-gold {
  color: hsl(43 74% 49%);
}

.text-lux-gold\/60 {
  color: hsl(43 74% 49% / 0.6);
}

.ring-lux-gold {
  --tw-ring-color: hsl(43 74% 49%);
}

.focus\:ring-lux-gold:focus {
  --tw-ring-color: hsl(43 74% 49%);
}

.bg-lux-noir {
  background-color: hsl(0 0% 4%);
}

.bg-lux-noir-light {
  background-color: hsl(0 0% 8%);
}

.bg-lux-noir-light\/95 {
  background-color: hsl(0 0% 8% / 0.95);
}

.text-lux-cream {
  color: hsl(45 30% 96%);
}

.text-lux-silver\/60 {
  color: hsl(220 10% 70% / 0.6);
}

.text-lux-silver\/80 {
  color: hsl(220 10% 70% / 0.8);
}

.bg-lux-gold\/5 {
  background-color: hsl(43 74% 49% / 0.05);
}

.bg-lux-gold\/10 {
  background-color: hsl(43 74% 49% / 0.1);
}

.from-lux-gold\/5 {
  --tw-gradient-from: hsl(43 74% 49% / 0.05);
}

.to-lux-gold\/10 {
  --tw-gradient-to: hsl(43 74% 49% / 0.1);
}

.border-lux-gold\/10 {
  border-color: hsl(43 74% 49% / 0.1);
}

.bg-lux-silver\/10 {
  background-color: hsl(220 10% 70% / 0.1);
}

.border-lux-silver\/20 {
  border-color: hsl(220 10% 70% / 0.2);
}

.font-display {
  font-family: 'Cormorant Garamond', Georgia, serif;
}

/* Modal animation */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active > div,
.modal-leave-active > div {
  transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1), opacity 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.modal-enter-from > div {
  transform: scale(0.95) translateY(10px);
  opacity: 0;
}

.modal-leave-to > div {
  transform: scale(0.95) translateY(10px);
  opacity: 0;
}

/* Image fade transition */
.image-fade-enter-active,
.image-fade-leave-active {
  transition: opacity 0.25s ease-out;
}

.image-fade-enter-from,
.image-fade-leave-to {
  opacity: 0;
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

  .image-fade-enter-active,
  .image-fade-leave-active {
    transition: none;
  }

  .transition-opacity,
  .transition-all,
  .transition-colors,
  .transition-transform {
    transition: none !important;
  }

  .animate-pulse {
    animation: none;
  }
}

/* Scrollbar styling for thumbnail container */
.scrollbar-thin {
  scrollbar-width: thin;
  scrollbar-color: hsl(var(--muted)) transparent;
}

.scrollbar-thin::-webkit-scrollbar {
  height: 6px;
}

.scrollbar-thin::-webkit-scrollbar-track {
  background: transparent;
}

.scrollbar-thin::-webkit-scrollbar-thumb {
  background-color: hsl(var(--muted));
  border-radius: 3px;
}

.scrollbar-thin::-webkit-scrollbar-thumb:hover {
  background-color: hsl(var(--muted-foreground));
}

/* General scrollbar styling */
.overflow-x-auto::-webkit-scrollbar,
.overflow-y-auto::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

.overflow-x-auto::-webkit-scrollbar-track,
.overflow-y-auto::-webkit-scrollbar-track {
  background: hsl(var(--muted) / 0.5);
  border-radius: 3px;
}

.overflow-x-auto::-webkit-scrollbar-thumb,
.overflow-y-auto::-webkit-scrollbar-thumb {
  background: hsl(var(--muted));
  border-radius: 3px;
}

.overflow-x-auto::-webkit-scrollbar-thumb:hover,
.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: hsl(var(--muted-foreground));
}

/* Tabular number display for prices */
.tabular-nums {
  font-variant-numeric: tabular-nums;
}

/* Dark mode scrollbar */
.dark .scrollbar-thin::-webkit-scrollbar-thumb {
  background-color: hsl(217.2 32.6% 17.5%);
}

.dark .scrollbar-thin::-webkit-scrollbar-thumb:hover {
  background-color: hsl(215 20.2% 65.1%);
}
</style>
