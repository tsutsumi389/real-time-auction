<template>
  <div class="bg-white rounded-lg shadow-sm p-6">
    <!-- Item Header -->
    <div class="mb-6">
      <h2 class="text-2xl font-bold text-gray-900 mb-2">{{ item.name }}</h2>
      <p v-if="item.description" class="text-gray-600 text-sm">
        {{ item.description }}
      </p>
      <div class="flex items-center gap-4 mt-3">
        <span
          :class="[
            'inline-flex items-center px-3 py-1 rounded-full text-sm font-medium',
            getStatusClass(item.status)
          ]"
        >
          {{ getStatusLabel(item.status) }}
        </span>
        <span class="text-sm text-gray-500">
          ロット番号: {{ item.lot_number }}
        </span>
      </div>
    </div>

    <!-- Current Price Display -->
    <div class="text-center p-8 bg-gradient-to-br from-blue-50 to-indigo-50 rounded-xl mb-6 border border-blue-100">
      <div class="text-sm text-gray-600 mb-2 font-medium">現在価格</div>
      <div
        :class="[
          'text-6xl font-bold transition-all duration-300',
          priceUpdated ? 'text-blue-600 animate-pulse' : 'text-gray-900'
        ]"
      >
        {{ formatNumber(currentPrice) }}
      </div>
      <div class="text-sm text-gray-500 mt-2">ポイント</div>
    </div>

    <!-- Bid Button -->
    <button
      @click="handleBid"
      :disabled="!canBid"
      :class="[
        'w-full py-5 px-6 rounded-xl font-bold text-xl transition-all duration-200 shadow-lg',
        canBid
          ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white hover:from-blue-700 hover:to-indigo-700 hover:shadow-xl transform hover:-translate-y-0.5 active:translate-y-0'
          : 'bg-gray-200 text-gray-500 cursor-not-allowed shadow-none'
      ]"
    >
      <span v-if="isLoading" class="flex items-center justify-center">
        <svg class="animate-spin h-6 w-6 mr-3" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        処理中...
      </span>
      <span v-else-if="canBid">
        {{ formatNumber(currentPrice) }} ポイントで入札
      </span>
      <span v-else>
        {{ disabledReason }}
      </span>
    </button>

    <!-- Additional Info -->
    <div class="mt-6 space-y-3">
      <!-- Winning Status -->
      <div
        v-if="isWinning"
        class="p-4 bg-green-50 border border-green-200 rounded-lg flex items-start"
      >
        <svg class="w-5 h-5 text-green-600 mr-3 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
        </svg>
        <div>
          <p class="text-sm font-semibold text-green-800">現在あなたが最高入札者です</p>
          <p class="text-xs text-green-600 mt-1">他の入札者が現れるまでこの価格で落札可能です</p>
        </div>
      </div>

      <!-- Insufficient Points Warning -->
      <div
        v-if="!hasEnoughPoints && item.status === 'active' && currentPrice > 0"
        class="p-4 bg-orange-50 border border-orange-200 rounded-lg flex items-start"
      >
        <svg class="w-5 h-5 text-orange-600 mr-3 mt-0.5 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
        </svg>
        <div>
          <p class="text-sm font-semibold text-orange-800">ポイントが不足しています</p>
          <p class="text-xs text-orange-600 mt-1">
            必要: {{ formatNumber(currentPrice) }} / 利用可能: {{ formatNumber(availablePoints) }}
          </p>
        </div>
      </div>

      <!-- Item Status Info -->
      <div
        v-if="item.status === 'pending'"
        class="p-4 bg-blue-50 border border-blue-200 rounded-lg"
      >
        <p class="text-sm text-blue-800">この商品はまだ開始されていません。開始をお待ちください。</p>
      </div>

      <div
        v-if="item.status === 'ended'"
        class="p-4 bg-gray-50 border border-gray-200 rounded-lg"
      >
        <p class="text-sm text-gray-800 font-semibold">この商品は終了しました</p>
        <p v-if="item.winner_id" class="text-xs text-gray-600 mt-1">
          落札者が決定しています
        </p>
      </div>

      <!-- No Price Opened Yet -->
      <div
        v-if="item.status === 'active' && currentPrice === 0"
        class="p-4 bg-yellow-50 border border-yellow-200 rounded-lg"
      >
        <p class="text-sm text-yellow-800">価格がまだ開示されていません。主催者が価格を開示するまでお待ちください。</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, computed } from 'vue'

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
    pending: '待機中',
    active: '開始',
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
@keyframes pulse {
  0%, 100% {
    transform: scale(1);
  }
  50% {
    transform: scale(1.05);
  }
}

.animate-pulse {
  animation: pulse 0.5s ease-in-out;
}
</style>
