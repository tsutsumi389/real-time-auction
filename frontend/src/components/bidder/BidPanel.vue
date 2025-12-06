<template>
  <div class="bid-panel-container rounded-2xl overflow-hidden relative">
    <!-- Animated Border Gradient -->
    <div class="absolute inset-0 rounded-2xl p-[1px] pointer-events-none z-0">
      <div
        :class="[
          'absolute inset-0 rounded-2xl transition-opacity duration-700',
          item.status === 'active' ? 'opacity-100' : 'opacity-40'
        ]"
        style="background: linear-gradient(135deg, rgba(212,175,55,0.4) 0%, rgba(212,175,55,0.1) 25%, rgba(212,175,55,0.3) 50%, rgba(212,175,55,0.1) 75%, rgba(212,175,55,0.4) 100%); background-size: 400% 400%; animation: gradient-flow 8s ease infinite;"
      ></div>
    </div>

    <!-- Main Content Container -->
    <div class="relative z-10 bg-gradient-to-br from-lux-noir-light via-lux-noir to-lux-noir-medium rounded-2xl overflow-hidden">
      <!-- Decorative Corner Elements -->
      <div class="absolute top-0 left-0 w-20 h-20 pointer-events-none z-20">
        <div class="absolute top-4 left-4 w-8 h-8 border-l-2 border-t-2 border-lux-gold/50 rounded-tl-lg"></div>
        <div class="absolute top-3 left-3 w-2 h-2 bg-lux-gold/60 rounded-full"></div>
      </div>
      <div class="absolute top-0 right-0 w-20 h-20 pointer-events-none z-20">
        <div class="absolute top-4 right-4 w-8 h-8 border-r-2 border-t-2 border-lux-gold/50 rounded-tr-lg"></div>
        <div class="absolute top-3 right-3 w-2 h-2 bg-lux-gold/60 rounded-full"></div>
      </div>

      <div class="grid grid-cols-1 md:grid-cols-2 gap-0">
        <!-- Left Column: Image Gallery -->
        <div class="relative">
          <ProductImageGallery
            :images="itemImages"
            :alt-text="item.name"
          />
          <!-- Status Overlay on Image -->
          <div
            v-if="item.status === 'pending'"
            class="absolute inset-0 bg-lux-noir/40 backdrop-blur-[2px] flex items-center justify-center pointer-events-none"
          >
            <div class="text-center">
              <div class="w-16 h-16 mx-auto mb-3 rounded-full border-2 border-lux-gold/40 flex items-center justify-center">
                <svg class="w-8 h-8 text-lux-gold/70" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <span class="px-4 py-2 rounded-full bg-lux-noir/80 border border-lux-gold/30 text-sm font-medium text-lux-gold tracking-wider">
                COMING SOON
              </span>
            </div>
          </div>
        </div>

        <!-- Right Column: Info & Bidding -->
        <div class="flex flex-col p-6 sm:p-8 relative">
          <!-- Subtle Background Pattern -->
          <div class="absolute inset-0 opacity-[0.02]" style="background-image: radial-gradient(circle at 1px 1px, rgba(212,175,55,0.8) 1px, transparent 0); background-size: 24px 24px;"></div>

          <div class="relative z-10 flex flex-col h-full">
            <!-- Item Header -->
            <div class="mb-6">
              <div class="flex items-center gap-3 mb-4">
                <span
                  :class="[
                    'inline-flex items-center gap-2 px-4 py-1.5 rounded-full text-xs font-semibold tracking-widest uppercase',
                    getStatusClass(item.status)
                  ]"
                >
                  <span
                    :class="[
                      'w-2 h-2 rounded-full',
                      item.status === 'active' ? 'bg-emerald-400 animate-pulse' :
                      item.status === 'pending' ? 'bg-amber-400' : 'bg-lux-silver/50'
                    ]"
                  ></span>
                  {{ getStatusLabel(item.status) }}
                </span>
                <span class="px-3 py-1 rounded-lg bg-lux-noir-medium border border-lux-gold/20 text-xs text-lux-gold font-mono tracking-wider">
                  LOT {{ item.lot_number }}
                </span>
              </div>

              <h2 class="font-display text-2xl sm:text-3xl text-lux-cream leading-tight mb-3">
                {{ item.name }}
              </h2>

              <p v-if="item.description" class="text-sm text-lux-silver/80 leading-relaxed line-clamp-3">
                {{ item.description }}
              </p>
            </div>

            <!-- Spacer -->
            <div class="flex-grow min-h-[2rem]"></div>

            <!-- Price Display Card -->
            <div class="mb-6">
              <div
                :class="[
                  'relative rounded-2xl overflow-hidden transition-all duration-500',
                  item.status === 'active' ? 'price-card-active' : 'price-card-waiting'
                ]"
              >
                <!-- Animated Background for Active State -->
                <div
                  v-if="item.status === 'active'"
                  class="absolute inset-0 bg-gradient-to-r from-lux-gold/10 via-lux-gold/5 to-lux-gold/10"
                  style="animation: shimmer-bg 3s ease-in-out infinite;"
                ></div>

                <!-- Card Content -->
                <div class="relative z-10 p-6 sm:p-8 text-center">
                  <div class="text-xs font-bold text-lux-silver/60 uppercase tracking-[0.25em] mb-3">
                    {{ item.status === 'active' ? 'Current Bid' : item.status === 'ended' ? 'Final Price' : 'Starting Price' }}
                  </div>

                  <div class="flex items-baseline justify-center gap-3">
                    <span
                      :class="[
                        'font-display font-light tracking-tight transition-all duration-500',
                        item.status === 'active' ? 'text-5xl sm:text-6xl' : 'text-4xl sm:text-5xl',
                        priceUpdated ? 'lux-shimmer scale-105' :
                        item.status === 'pending' ? 'text-lux-silver/70' : 'lux-text-gold'
                      ]"
                    >
                      {{ currentPrice > 0 ? formatNumber(currentPrice) : '—' }}
                    </span>
                    <span
                      :class="[
                        'text-lg font-medium transition-colors duration-300',
                        item.status === 'pending' ? 'text-lux-silver/40' : 'text-lux-gold/50'
                      ]"
                    >
                      pts
                    </span>
                  </div>

                  <!-- Waiting Message -->
                  <div
                    v-if="item.status === 'pending'"
                    class="mt-4 flex items-center justify-center gap-2 text-sm text-lux-silver/60"
                  >
                    <div class="flex gap-1">
                      <span class="w-1.5 h-1.5 rounded-full bg-lux-gold/50 animate-bounce" style="animation-delay: 0ms;"></span>
                      <span class="w-1.5 h-1.5 rounded-full bg-lux-gold/50 animate-bounce" style="animation-delay: 150ms;"></span>
                      <span class="w-1.5 h-1.5 rounded-full bg-lux-gold/50 animate-bounce" style="animation-delay: 300ms;"></span>
                    </div>
                    <span>オークション開始をお待ちください</span>
                  </div>
                </div>

                <!-- Bottom Accent Line -->
                <div
                  :class="[
                    'h-1 transition-all duration-500',
                    item.status === 'active' ? 'bg-gradient-to-r from-transparent via-lux-gold to-transparent' :
                    item.status === 'ended' ? 'bg-gradient-to-r from-transparent via-lux-silver/30 to-transparent' :
                    'bg-gradient-to-r from-transparent via-lux-gold/30 to-transparent'
                  ]"
                ></div>
              </div>
            </div>

            <!-- Desktop Bid Button -->
            <div class="hidden md:block mb-6">
              <button
                @click="handleBid"
                :disabled="!canBid"
                :class="[
                  'w-full py-4 px-6 rounded-xl text-base tracking-widest font-semibold transition-all duration-300',
                  canBid ? 'lux-btn-gold shadow-lg shadow-lux-gold/20' : 'bid-button-disabled'
                ]"
              >
                <span v-if="isLoading" class="flex items-center justify-center">
                  <svg class="animate-spin -ml-1 mr-3 h-5 w-5" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  入札処理中...
                </span>
                <span v-else-if="canBid" class="flex items-center justify-center gap-2">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 15l-2 5L9 9l11 4-5 2zm0 0l5 5M7.188 2.239l.777 2.897M5.136 7.965l-2.898-.777M13.95 4.05l-2.122 2.122m-5.657 5.656l-2.12 2.122" />
                  </svg>
                  入札する
                </span>
                <span v-else class="flex items-center justify-center gap-2 text-lux-silver/50">
                  <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
                  </svg>
                  {{ disabledReason || '入札できません' }}
                </span>
              </button>
            </div>

            <!-- Status Messages -->
            <div class="space-y-3">
              <!-- Winning Status -->
              <div
                v-if="isWinning"
                class="relative p-4 rounded-xl overflow-hidden winning-status-card"
              >
                <div class="absolute inset-0 bg-gradient-to-r from-lux-gold/15 via-lux-gold/5 to-lux-gold/15"></div>
                <div class="absolute inset-0 rounded-xl" style="background: linear-gradient(90deg, transparent 0%, rgba(212,175,55,0.1) 50%, transparent 100%); animation: winning-sweep 2s ease-in-out infinite;"></div>
                <div class="relative z-10 flex items-center gap-3">
                  <div class="w-10 h-10 rounded-full bg-gradient-to-br from-lux-gold to-lux-gold-dark flex items-center justify-center shadow-lg shadow-lux-gold/30">
                    <svg class="w-5 h-5 text-lux-noir" fill="currentColor" viewBox="0 0 20 20">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                    </svg>
                  </div>
                  <div>
                    <p class="text-sm font-bold text-lux-gold tracking-wide">最高額入札中</p>
                    <p class="text-xs text-lux-gold/70">現在あなたが最高額です！</p>
                  </div>
                </div>
              </div>

              <!-- Insufficient Points -->
              <div
                v-if="!hasEnoughPoints && item.status === 'active' && currentPrice > 0"
                class="p-4 rounded-xl bg-red-950/40 border border-red-500/30 flex items-center gap-3"
              >
                <div class="w-10 h-10 rounded-full bg-red-500/20 flex items-center justify-center">
                  <svg class="w-5 h-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
                  </svg>
                </div>
                <div>
                  <p class="text-sm font-semibold text-red-300">ポイント不足</p>
                  <p class="text-xs text-red-400/80">残高: {{ formatNumber(availablePoints) }} pts</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Mobile Sticky Bid Button -->
    <div class="md:hidden fixed bottom-0 left-0 right-0 z-50 safe-area-bottom">
      <div class="bg-lux-noir/95 backdrop-blur-xl border-t border-lux-gold/30 p-4 shadow-2xl">
        <div class="flex items-center justify-between gap-4 max-w-md mx-auto">
          <div class="flex-shrink-0">
            <div class="text-[10px] text-lux-silver/60 uppercase tracking-widest mb-0.5">現在価格</div>
            <div class="text-2xl font-display lux-text-gold font-light">
              {{ currentPrice > 0 ? formatNumber(currentPrice) : '—' }}
              <span class="text-xs font-body text-lux-silver/50">pts</span>
            </div>
          </div>
          <button
            @click="handleBid"
            :disabled="!canBid"
            :class="[
              'flex-1 py-4 px-6 rounded-xl text-sm tracking-widest font-semibold transition-all duration-300',
              canBid ? 'lux-btn-gold' : 'bid-button-disabled'
            ]"
          >
            <span v-if="isLoading">処理中...</span>
            <span v-else-if="canBid">入札する</span>
            <span v-else class="text-lux-silver/50">{{ disabledReason || '入札不可' }}</span>
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

/* Animated Border Gradient */
@keyframes gradient-flow {
  0%, 100% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
}

/* Price Card Styles */
.price-card-active {
  background: linear-gradient(145deg, rgba(20, 20, 20, 0.95) 0%, rgba(10, 10, 10, 0.98) 100%);
  border: 1px solid rgba(212, 175, 55, 0.4);
  box-shadow:
    0 0 40px rgba(212, 175, 55, 0.1),
    0 4px 24px rgba(0, 0, 0, 0.4),
    inset 0 1px 0 rgba(255, 255, 255, 0.05);
}

.price-card-waiting {
  background: linear-gradient(145deg, rgba(25, 25, 25, 0.9) 0%, rgba(15, 15, 15, 0.95) 100%);
  border: 1px solid rgba(212, 175, 55, 0.15);
  box-shadow:
    0 4px 20px rgba(0, 0, 0, 0.3),
    inset 0 1px 0 rgba(255, 255, 255, 0.02);
}

/* Shimmer Background Animation */
@keyframes shimmer-bg {
  0%, 100% { opacity: 0.5; transform: translateX(-100%); }
  50% { opacity: 1; transform: translateX(100%); }
}

/* Winning Status Card */
.winning-status-card {
  border: 1px solid rgba(212, 175, 55, 0.4);
  box-shadow: 0 0 20px rgba(212, 175, 55, 0.15);
}

@keyframes winning-sweep {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}

/* Disabled Button Style */
.bid-button-disabled {
  background: linear-gradient(145deg, rgba(40, 40, 40, 0.8) 0%, rgba(30, 30, 30, 0.9) 100%);
  border: 1px solid rgba(100, 100, 100, 0.2);
  color: rgba(160, 160, 160, 0.5);
  cursor: not-allowed;
}

/* Bid Panel Container */
.bid-panel-container {
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.5),
    0 0 0 1px rgba(212, 175, 55, 0.1);
}
</style>
