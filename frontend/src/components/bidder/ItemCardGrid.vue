<template>
  <div class="item-card-grid">
    <!-- アイテムがある場合 -->
    <div
      v-if="items.length > 0"
      class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6"
      role="list"
      :aria-label="`出品アイテム一覧 ${items.length}点`"
    >
      <article
        v-for="item in items"
        :key="item.id"
        role="listitem"
        tabindex="0"
        @click="handleItemClick(item)"
        @keydown.enter="handleItemClick(item)"
        @keydown.space.prevent="handleItemClick(item)"
        class="item-card group bg-white border border-gray-200 rounded-lg overflow-hidden cursor-pointer
               hover:shadow-lg hover:border-gray-300 hover:-translate-y-1
               focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2
               transition-all duration-200"
        :aria-label="`LOT ${item.lot_number}: ${item.name}, 開始価格 ${formatPrice(item.starting_price)}円`"
      >
        <!-- サムネイル画像 -->
        <div class="relative w-full aspect-[4/3] bg-gray-100 overflow-hidden">
          <img
            v-if="getItemThumbnail(item) && !imageErrors[item.id]"
            :src="getItemThumbnail(item)"
            :alt="item.name"
            class="w-full h-full object-cover transition-transform duration-200 group-hover:scale-105"
            loading="lazy"
            @error="handleImageError(item.id, $event)"
          />
          <!-- プレースホルダー画像 -->
          <div
            v-else
            class="w-full h-full flex flex-col items-center justify-center text-gray-400"
            role="img"
            aria-label="画像なし"
          >
            <svg class="h-10 w-10 sm:h-12 sm:w-12 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
            </svg>
            <span class="text-xs sm:text-sm" aria-hidden="true">画像なし</span>
          </div>

          <!-- 画像枚数バッジ（複数メディアがある場合） -->
          <div
            v-if="getMediaCount(item) > 1"
            class="absolute bottom-2 right-2 px-2 py-1 bg-black/60 text-white text-xs rounded-md flex items-center"
            :aria-label="`${getMediaCount(item)}枚の画像があります`"
          >
            <svg class="h-3 w-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
            </svg>
            <span aria-hidden="true">{{ getMediaCount(item) }}</span>
          </div>
        </div>

        <!-- カード情報 -->
        <div class="p-3 sm:p-4">
          <!-- LOT番号バッジ -->
          <div class="mb-2">
            <span class="inline-block px-2 py-1 text-xs font-semibold text-blue-800 bg-blue-100 rounded">
              LOT {{ item.lot_number }}
            </span>
            <!-- ステータスバッジ（オプション） -->
            <span
              v-if="item.status === 'active'"
              class="inline-block ml-2 px-2 py-1 text-xs font-semibold text-green-800 bg-green-100 rounded"
            >
              入札中
            </span>
            <span
              v-else-if="item.status === 'ended'"
              class="inline-block ml-2 px-2 py-1 text-xs font-semibold text-gray-800 bg-gray-100 rounded"
            >
              終了
            </span>
          </div>

          <!-- 商品名 -->
          <h3 class="text-sm sm:text-base font-semibold text-gray-900 mb-2 line-clamp-2 min-h-[2.5rem] sm:min-h-[3rem]">
            {{ item.name }}
          </h3>

          <!-- 開始価格 -->
          <div class="flex items-baseline gap-1">
            <span class="text-xs text-gray-500">開始価格</span>
            <span class="text-base sm:text-lg font-bold text-gray-900">
              ¥{{ formatPrice(item.starting_price) }}
            </span>
          </div>

          <!-- 現在価格（入札がある場合） -->
          <div
            v-if="item.current_price && item.current_price > item.starting_price"
            class="flex items-baseline gap-1 mt-1"
          >
            <span class="text-xs text-green-600">現在価格</span>
            <span class="text-base sm:text-lg font-bold text-green-600">
              ¥{{ formatPrice(item.current_price) }}
            </span>
          </div>
        </div>

        <!-- ホバー時のオーバーレイ（クリック促進） -->
        <div class="absolute inset-0 bg-black/0 group-hover:bg-black/5 transition-colors pointer-events-none"></div>
      </article>
    </div>

    <!-- アイテムがない場合 -->
    <div
      v-else
      class="text-center py-12 px-4"
      role="status"
      aria-label="出品アイテムがありません"
    >
      <svg class="mx-auto h-12 w-12 text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"></path>
      </svg>
      <p class="text-gray-500 text-base sm:text-lg">
        このオークションにはまだアイテムが登録されていません
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'

// Props
const props = defineProps({
  items: {
    type: Array,
    required: true,
    default: () => []
  }
})

// Emits
const emit = defineEmits(['item-click'])

// State for tracking image load errors
const imageErrors = ref({})

// Methods

/**
 * Format price with Japanese locale
 * @param {number} price - The price to format
 * @returns {string} Formatted price string
 */
const formatPrice = (price) => {
  if (!price && price !== 0) return '0'
  return new Intl.NumberFormat('ja-JP').format(price)
}

/**
 * Get the thumbnail URL for an item
 * @param {Object} item - The item object
 * @returns {string|null} Thumbnail URL or null
 */
const getItemThumbnail = (item) => {
  if (!item.media || item.media.length === 0) return null
  // Prefer thumbnail_url, fallback to main url
  const firstMedia = item.media[0]
  return firstMedia.thumbnail_url || firstMedia.url
}

/**
 * Get the media count for an item
 * @param {Object} item - The item object
 * @returns {number} Number of media items
 */
const getMediaCount = (item) => {
  return item.media?.length || 0
}

/**
 * Handle image load error
 * @param {string} itemId - The item ID
 * @param {Event} event - The error event
 */
const handleImageError = (itemId, event) => {
  imageErrors.value[itemId] = true
  // Hide the broken image
  event.target.style.display = 'none'
}

/**
 * Handle item card click
 * @param {Object} item - The clicked item
 */
const handleItemClick = (item) => {
  emit('item-click', item)
}
</script>

<style scoped>
/* Line clamp for item name */
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Ensure smooth transitions respect user preferences */
@media (prefers-reduced-motion: reduce) {
  .item-card {
    transition: none;
  }
  .item-card img {
    transition: none;
  }
}

/* Focus visible styles for better accessibility */
.item-card:focus-visible {
  outline: 2px solid #3b82f6;
  outline-offset: 2px;
}
</style>
