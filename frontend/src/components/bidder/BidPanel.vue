<template>
  <div class="lux-glass-strong rounded-2xl overflow-hidden relative">
    <!-- Subtle Corner Accents -->
    <div class="absolute top-0 left-0 w-12 h-12 border-l border-t border-lux-gold/20 rounded-tl-2xl pointer-events-none z-10"></div>
    <div class="absolute top-0 right-0 w-12 h-12 border-r border-t border-lux-gold/20 rounded-tr-2xl pointer-events-none z-10"></div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-0">
      <!-- Left Column: Image Gallery -->
      <div class="relative">
        <ProductImageGallery
          :images="itemImages"
          :alt-text="item.name"
        />
      </div>

      <!-- Right Column: Info & Bidding -->
      <div class="flex flex-col p-6 sm:p-8">
        <!-- Item Header -->
        <div class="mb-6">
          <div class="flex items-center gap-3 mb-4">
            <span
              :class="[
                'inline-flex items-center px-3 py-1 rounded-full text-xs font-medium tracking-wide',
                getStatusClass(item.status)
              ]"
            >
              {{ getStatusLabel(item.status) }}
            </span>
            <span class="text-xs text-lux-silver font-mono tracking-wider">
              LOT {{ item.lot_number }}
            </span>
          </div>

          <h2 class="font-display text-2xl sm:text-3xl text-lux-cream leading-tight mb-3">
            {{ item.name }}
          </h2>

          <p v-if="item.description" class="text-sm text-lux-silver leading-relaxed line-clamp-3">
            {{ item.description }}
          </p>
        </div>

        <!-- Spacer -->
        <div class="flex-grow"></div>

        <!-- Price Display -->
        <div class="mb-6">
          <div class="relative rounded-xl bg-lux-noir-light border border-lux-gold/20 p-6 text-center overflow-hidden">
            <!-- Background Glow -->
            <div class="absolute inset-0 bg-gradient-to-br from-lux-gold/5 via-transparent to-lux-gold/5 pointer-events-none"></div>

            <div class="relative z-10">
              <div class="text-xs font-semibold text-lux-silver uppercase tracking-[0.2em] mb-2">
                Current Price
              </div>
              <div class="flex items-baseline justify-center gap-2">
                <span
                  :class="[
                    'font-display text-5xl sm:text-6xl font-light tracking-tight transition-all duration-300',
                    priceUpdated ? 'lux-shimmer scale-105' : 'lux-text-gold'
                  ]"
                >
                  {{ formatNumber(currentPrice) }}
                </span>
                <span class="text-lg font-medium text-lux-silver/70">pts</span>
              </div>
            </div>
          </div>
        </div>

        <!-- Desktop Bid Button -->
        <div class="hidden md:block mb-6">
          <button
            @click="handleBid"
            :disabled="!canBid"
            class="w-full py-4 px-6 rounded-xl text-base lux-btn-gold tracking-widest"
          >
            <span v-if="isLoading" class="flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Processing...
            </span>
            <span v-else-if="canBid">Place Bid</span>
            <span v-else class="text-lux-silver/50">{{ disabledReason || 'Bidding Unavailable' }}</span>
          </button>
        </div>

        <!-- Status Messages -->
        <div class="space-y-3">
          <!-- Winning Status -->
          <div
            v-if="isWinning"
            class="relative p-4 rounded-xl overflow-hidden lux-winning-glow"
          >
            <div class="absolute inset-0 bg-gradient-to-r from-lux-gold/10 via-lux-gold/5 to-lux-gold/10"></div>
            <div class="relative z-10 flex items-center gap-3">
              <div class="w-8 h-8 rounded-full bg-lux-gold/20 flex items-center justify-center">
                <svg class="w-4 h-4 text-lux-gold" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
              </div>
              <div>
                <p class="text-sm font-semibold text-lux-gold">Highest Bidder</p>
                <p class="text-xs text-lux-gold/70">You are currently winning!</p>
              </div>
            </div>
          </div>

          <!-- Insufficient Points -->
          <div
            v-if="!hasEnoughPoints && item.status === 'active' && currentPrice > 0"
            class="p-4 rounded-xl bg-red-950/30 border border-red-500/20 flex items-center gap-3"
          >
            <div class="w-8 h-8 rounded-full bg-red-500/10 flex items-center justify-center">
              <svg class="w-4 h-4 text-red-400" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
              </svg>
            </div>
            <div>
              <p class="text-sm font-medium text-red-300">Insufficient Points</p>
              <p class="text-xs text-red-400/70">Balance: {{ formatNumber(availablePoints) }} pts</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Mobile Sticky Bid Button -->
    <div class="md:hidden fixed bottom-0 left-0 right-0 z-50 safe-area-bottom">
      <div class="lux-glass border-t border-lux-gold/20 p-4">
        <div class="flex items-center justify-between gap-4 max-w-md mx-auto">
          <div class="flex-shrink-0">
            <div class="text-xs text-lux-silver uppercase tracking-wider">Current</div>
            <div class="text-xl font-display lux-text-gold">
              {{ formatNumber(currentPrice) }}
              <span class="text-xs font-body text-lux-silver/60">pts</span>
            </div>
          </div>
          <button
            @click="handleBid"
            :disabled="!canBid"
            class="flex-1 py-3.5 px-6 rounded-xl text-sm lux-btn-gold tracking-widest"
          >
            <span v-if="isLoading">Processing...</span>
            <span v-else-if="canBid">Place Bid</span>
            <span v-else class="text-lux-silver/50">{{ disabledReason || 'Unavailable' }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'
import ProductImageGallery from './ProductImageGallery.vue'

const props = defineProps({
  item: {
    type: Object,
    required: true,
  },
  currentPrice: {
    type: Number,
    required: true,
  },
  availablePoints: {
    type: Number,
    required: true,
  },
  isLoading: {
    type: Boolean,
    default: false,
  },
  canBid: {
    type: Boolean,
    required: true,
  },
  disabledReason: {
    type: String,
    default: '',
  },
  isWinning: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['bid'])

// Local state for price update animation
const priceUpdated = ref(false)

// Watch for price changes to trigger animation
watch(
  () => props.currentPrice,
  (newPrice, oldPrice) => {
    if (newPrice !== oldPrice && oldPrice !== undefined) {
      priceUpdated.value = true
      setTimeout(() => {
        priceUpdated.value = false
      }, 1000)
    }
  }
)

// Computed
const hasEnoughPoints = computed(() => {
  return props.availablePoints >= props.currentPrice
})

// Extract images from item (supports media array from API)
const itemImages = computed(() => {
  // Check for media array from API (preferred format)
  if (props.item.media && Array.isArray(props.item.media)) {
    // Filter to only images and map to URLs (prefer thumbnail for gallery)
    return props.item.media
      .filter(m => m.media_type === 'image')
      .sort((a, b) => a.display_order - b.display_order)
      .map(m => ({
        url: m.url,
        thumbnail: m.thumbnail_url || m.url
      }))
  }
  // Legacy support for images array
  if (props.item.images && Array.isArray(props.item.images)) {
    return props.item.images.map(url => ({ url, thumbnail: url }))
  }
  // Legacy support for single image_url
  if (props.item.image_url) {
    return [{ url: props.item.image_url, thumbnail: props.item.image_url }]
  }
  return []
})

function handleBid() {
  if (props.canBid) {
    emit('bid')
  }
}

function formatNumber(value) {
  return new Intl.NumberFormat('ja-JP').format(value)
}

function getStatusLabel(status) {
  const labels = {
    pending: 'Waiting',
    active: 'Live',
    ended: 'Ended',
  }
  return labels[status] || status
}

function getStatusClass(status) {
  const classes = {
    pending: 'lux-badge-pending',
    active: 'lux-badge-active',
    ended: 'lux-badge-ended',
  }
  return classes[status] || 'lux-badge-ended'
}
</script>

<style scoped>
.safe-area-bottom {
  padding-bottom: env(safe-area-inset-bottom, 1rem);
}

.line-clamp-3 {
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
