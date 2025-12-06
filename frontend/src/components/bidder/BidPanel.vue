<template>
  <div class="bg-white rounded-lg shadow-sm p-6 relative">
    <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
      <!-- Left Column: Image Gallery -->
      <div>
        <ProductImageGallery 
          :images="itemImages" 
          :alt-text="item.name" 
        />
      </div>

      <!-- Right Column: Info & Bidding -->
      <div class="flex flex-col h-full">
        <!-- Item Header -->
        <div class="mb-6">
          <div class="flex items-center gap-2 mb-2">
            <span
              :class="[
                'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
                getStatusClass(item.status)
              ]"
            >
              {{ getStatusLabel(item.status) }}
            </span>
            <span class="text-xs text-gray-500 font-mono">
              LOT: {{ item.lot_number }}
            </span>
          </div>
          <h2 class="text-2xl font-bold text-gray-900 leading-tight mb-2">{{ item.name }}</h2>
          <p v-if="item.description" class="text-gray-600 text-sm leading-relaxed">
            {{ item.description }}
          </p>
        </div>

        <div class="mt-auto space-y-6">
          <!-- Price Display (Luxury Redesigned) -->
          <div class="bg-auction-cream rounded-xl p-6 border-2 border-auction-gold-light shadow-inner-luxury text-center relative overflow-hidden">
            <!-- Background decoration -->
            <div class="absolute top-0 right-0 -mt-4 -mr-4 w-24 h-24 bg-auction-gold-light/20 rounded-full opacity-50 blur-xl"></div>
            <div class="absolute bottom-0 left-0 -mb-4 -ml-4 w-20 h-20 bg-auction-burgundy/10 rounded-full opacity-50 blur-xl"></div>
            
            <div class="relative z-10">
              <div class="text-sm font-semibold text-auction-gold-light uppercase tracking-widest mb-1">Current Price</div>
              <div class="flex items-baseline justify-center gap-1">
                <span 
                  :class="[
                    'text-5xl sm:text-6xl font-serif font-black tracking-tight transition-all duration-300 tabular-nums',
                    priceUpdated ? 'text-auction-gold-light scale-110' : 'text-auction-burgundy'
                  ]"
                >
                  {{ formatNumber(currentPrice) }}
                </span>
                <span class="text-lg font-medium text-gray-500">pts</span>
              </div>
            </div>
          </div>

          <!-- Desktop Bid Button -->
          <div class="hidden md:block">
            <button
              @click="handleBid"
              :disabled="!canBid"
              class="relative overflow-hidden w-full py-4 px-6 rounded-xl font-bold text-lg transition-all duration-200 transform"
              :class="[
                canBid
                  ? 'bg-gold-gradient text-white shadow-gold-glow hover:brightness-110 hover:shadow-luxury-xl hover:-translate-y-1 active:translate-y-0'
                  : 'bg-gray-200 text-gray-500 cursor-not-allowed shadow-none'
              ]"
            >
              <span v-if="isLoading" class="flex items-center justify-center relative z-10">
                <svg class="animate-spin h-5 w-5 mr-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Processing...
              </span>
              <span v-else-if="canBid" class="relative z-10">
                入札する
              </span>
              <span v-else class="relative z-10">
                {{ disabledReason }}
              </span>
            </button>
          </div>

          <!-- Status Messages -->
          <div class="space-y-3">
            <!-- Winning Status -->
            <div
              v-if="isWinning"
              class="p-3 bg-auction-cream border-2 border-auction-gold-light rounded-lg flex items-center animate-pulse-gold shadow-gold-glow"
            >
              <div class="flex-shrink-0 bg-gold-gradient rounded-full p-1 mr-3">
                <svg class="w-4 h-4 text-white" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
              </div>
              <div>
                <p class="text-sm font-bold text-auction-gold-light">最高入札者です!</p>
              </div>
            </div>

            <!-- Insufficient Points -->
            <div
              v-if="!hasEnoughPoints && item.status === 'active' && currentPrice > 0"
              class="p-3 bg-orange-50 border border-orange-200 rounded-lg flex items-center"
            >
              <svg class="w-5 h-5 text-orange-600 mr-3" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
              </svg>
              <p class="text-sm text-orange-800">
                ポイント不足 (残: {{ formatNumber(availablePoints) }})
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Mobile Sticky Bid Button -->
    <div class="md:hidden fixed bottom-0 left-0 right-0 p-4 bg-white border-t border-gray-200 shadow-[0_-4px_6px_-1px_rgba(0,0,0,0.1)] z-50 safe-area-bottom">
      <div class="flex items-center justify-between gap-4 max-w-md mx-auto">
        <div class="flex-shrink-0">
          <div class="text-xs text-gray-500">現在価格</div>
          <div class="text-xl font-bold text-gray-900">{{ formatNumber(currentPrice) }} <span class="text-xs font-normal">pts</span></div>
        </div>
        <button
          @click="handleBid"
          :disabled="!canBid"
          :class="[
            'flex-1 py-3 px-4 rounded-xl font-bold text-base transition-all duration-200',
            canBid
              ? 'bg-gold-gradient text-white shadow-gold-glow active:brightness-90'
              : 'bg-gray-200 text-gray-500 cursor-not-allowed shadow-none'
          ]"
        >
          <span v-if="isLoading">処理中...</span>
          <span v-else-if="canBid">入札する</span>
          <span v-else>{{ disabledReason || '入札不可' }}</span>
        </button>
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
      }, 300)
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

function handleBid(event) {
  if (props.canBid) {
    // Create ripple effect
    if (event) {
      const button = event.currentTarget
      const circle = document.createElement('span')
      const diameter = Math.max(button.clientWidth, button.clientHeight)
      const radius = diameter / 2

      const rect = button.getBoundingClientRect()
      
      circle.style.width = circle.style.height = `${diameter}px`
      circle.style.left = `${event.clientX - rect.left - radius}px`
      circle.style.top = `${event.clientY - rect.top - radius}px`
      circle.classList.add('ripple')

      const ripple = button.getElementsByClassName('ripple')[0]
      if (ripple) {
        ripple.remove()
      }

      button.appendChild(circle)
    }

    emit('bid')
  }
}

function formatNumber(value) {
  return new Intl.NumberFormat('ja-JP').format(value)
}

function getStatusLabel(status) {
  const labels = {
    pending: '待機中',
    active: '入札受付中',
    ended: '終了',
  }
  return labels[status] || status
}

function getStatusClass(status) {
  const classes = {
    pending: 'bg-yellow-100 text-yellow-800',
    active: 'bg-green-100 text-green-800',
    ended: 'bg-gray-100 text-gray-800',
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}
</script>

<style scoped>
.safe-area-bottom {
  padding-bottom: env(safe-area-inset-bottom, 1rem);
}

@keyframes pulse-slow {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.8; }
}
.animate-pulse-slow {
  animation: pulse-slow 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}

/* Ripple Effect */
.ripple {
  position: absolute;
  border-radius: 50%;
  transform: scale(0);
  animation: ripple 0.6s linear;
  background-color: rgba(255, 255, 255, 0.7);
}

@keyframes ripple {
  to {
    transform: scale(4);
    opacity: 0;
  }
}

/* Enhanced Price Animation */
@keyframes price-pop {
  0% { transform: scale(1); color: #1f2937; } /* text-gray-900 */
  50% { transform: scale(1.2); color: #2563eb; } /* text-blue-600 */
  100% { transform: scale(1); color: #1f2937; }
}

.price-pop-enter-active {
  animation: price-pop 0.4s ease-out;
}

/* Accessibility: Respect reduced motion preference */
@media (prefers-reduced-motion: reduce) {
  .animate-pulse-slow,
  .animate-pulse-gold,
  .price-pop-enter-active,
  .ripple {
    animation: none !important;
  }
  
  * {
    transition-duration: 0.01ms !important;
  }
}
</style>
