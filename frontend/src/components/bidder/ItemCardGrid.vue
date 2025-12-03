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
        class="item-card group relative bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-2xl overflow-hidden cursor-pointer
               shadow-sm hover:shadow-xl hover:border-primary/30 dark:hover:border-primary/50 hover:-translate-y-1.5
               focus:outline-none focus:ring-2 focus:ring-primary focus:ring-offset-2 dark:focus:ring-offset-gray-900
               transition-all duration-300 ease-out"
        :aria-label="`LOT ${item.lot_number}: ${item.name}, 開始価格 ${formatPrice(item.starting_price)}pt`"
      >
        <!-- サムネイル画像 -->
        <div class="relative w-full aspect-[4/3] bg-gray-100 dark:bg-gray-700 overflow-hidden">
          <img
            v-if="getItemThumbnail(item) && !imageErrors[item.id]"
            :src="getItemThumbnail(item)"
            :alt="item.name"
            class="w-full h-full object-cover transition-transform duration-500 ease-out group-hover:scale-110"
            loading="lazy"
            @error="handleImageError(item.id, $event)"
          />
          <!-- プレースホルダー画像 -->
          <div
            v-else
            class="w-full h-full flex flex-col items-center justify-center text-gray-400 dark:text-gray-500 bg-gradient-to-br from-gray-100 to-gray-200 dark:from-gray-700 dark:to-gray-800"
            role="img"
            aria-label="画像なし"
          >
            <svg class="h-10 w-10 sm:h-12 sm:w-12 mb-2 opacity-50" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
            </svg>
            <span class="text-xs sm:text-sm font-medium" aria-hidden="true">画像なし</span>
          </div>

          <!-- 画像枚数バッジ（複数メディアがある場合） -->
          <div
            v-if="getMediaCount(item) > 1"
            class="absolute bottom-2 right-2 px-2.5 py-1 bg-black/70 backdrop-blur-sm text-white text-xs font-medium rounded-full flex items-center shadow-lg"
            :aria-label="`${getMediaCount(item)}枚の画像があります`"
          >
            <svg class="h-3.5 w-3.5 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
            </svg>
            <span aria-hidden="true">{{ getMediaCount(item) }}</span>
          </div>

          <!-- ホバー時のオーバーレイグラデーション -->
          <div class="absolute inset-0 bg-gradient-to-t from-black/20 via-transparent to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none"></div>
        </div>

        <!-- カード情報 -->
        <div class="p-3 sm:p-4">
          <!-- LOT番号バッジ -->
          <div class="mb-2 flex items-center flex-wrap gap-2">
            <span class="inline-flex items-center px-2.5 py-1 text-xs font-bold text-primary bg-primary/10 dark:bg-primary/20 rounded-full">
              LOT {{ item.lot_number }}
            </span>
            <!-- ステータスバッジ（オプション） -->
            <span
              v-if="item.status === 'active'"
              class="inline-flex items-center px-2.5 py-1 text-xs font-semibold text-emerald-700 dark:text-emerald-300 bg-emerald-100 dark:bg-emerald-900/40 rounded-full"
            >
              <span class="w-1.5 h-1.5 bg-emerald-500 rounded-full mr-1.5 animate-pulse"></span>
              入札中
            </span>
            <span
              v-else-if="item.status === 'ended'"
              class="inline-flex items-center px-2.5 py-1 text-xs font-semibold text-gray-600 dark:text-gray-400 bg-gray-100 dark:bg-gray-700 rounded-full"
            >
              終了
            </span>
          </div>

          <!-- 商品名 -->
          <h3 class="text-sm sm:text-base font-semibold text-gray-900 dark:text-white mb-3 line-clamp-2 min-h-[2.5rem] sm:min-h-[3rem] leading-snug group-hover:text-primary transition-colors duration-200">
            {{ item.name }}
          </h3>

          <!-- 価格セクション -->
          <div class="space-y-1.5 pt-3 border-t border-gray-100 dark:border-gray-700">
            <!-- 開始価格 -->
            <div class="flex items-baseline justify-between">
              <span class="text-xs text-gray-500 dark:text-gray-400">開始価格</span>
              <span class="text-base sm:text-lg font-bold text-gray-900 dark:text-white tabular-nums">
                {{ formatPrice(item.starting_price) }}pt
              </span>
            </div>

            <!-- 現在価格（入札がある場合） -->
            <div
              v-if="item.current_price && item.current_price > item.starting_price"
              class="flex items-baseline justify-between"
            >
              <span class="text-xs text-emerald-600 dark:text-emerald-400 font-medium">現在価格</span>
              <span class="text-base sm:text-lg font-bold text-emerald-600 dark:text-emerald-400 tabular-nums">
                {{ formatPrice(item.current_price) }}pt
              </span>
            </div>
          </div>
        </div>

        <!-- 詳細を見るヒント -->
        <div class="absolute bottom-0 left-0 right-0 p-3 bg-gradient-to-t from-black/60 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300 pointer-events-none">
          <span class="text-white text-xs font-medium flex items-center justify-center">
            <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
            </svg>
            詳細を見る
          </span>
        </div>
      </article>
    </div>

    <!-- アイテムがない場合 -->
    <div
      v-else
      class="text-center py-16 px-4"
      role="status"
      aria-label="出品アイテムがありません"
    >
      <div class="inline-flex items-center justify-center w-16 h-16 bg-gray-100 dark:bg-gray-700 rounded-full mb-4">
        <svg class="h-8 w-8 text-gray-400 dark:text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24" aria-hidden="true">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"></path>
        </svg>
      </div>
      <p class="text-gray-500 dark:text-gray-400 text-base sm:text-lg font-medium">
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

/* Card entrance animation */
.item-card {
  animation: card-entrance 0.4s ease-out;
  animation-fill-mode: both;
}

.item-card:nth-child(1) { animation-delay: 0.05s; }
.item-card:nth-child(2) { animation-delay: 0.1s; }
.item-card:nth-child(3) { animation-delay: 0.15s; }
.item-card:nth-child(4) { animation-delay: 0.2s; }
.item-card:nth-child(5) { animation-delay: 0.25s; }
.item-card:nth-child(6) { animation-delay: 0.3s; }
.item-card:nth-child(7) { animation-delay: 0.35s; }
.item-card:nth-child(8) { animation-delay: 0.4s; }

@keyframes card-entrance {
  from {
    opacity: 0;
    transform: translateY(16px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Ensure smooth transitions respect user preferences */
@media (prefers-reduced-motion: reduce) {
  .item-card {
    animation: none;
    transition: box-shadow 0.2s ease;
  }
  .item-card img {
    transition: none;
  }
}

/* Focus visible styles for better accessibility */
.item-card:focus-visible {
  outline: 2px solid hsl(var(--primary));
  outline-offset: 2px;
}

/* Tabular number display for prices */
.tabular-nums {
  font-variant-numeric: tabular-nums;
}
</style>
