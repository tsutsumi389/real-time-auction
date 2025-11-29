<template>
  <div class="bg-white rounded-lg shadow-sm p-6">
    <h2 class="text-lg font-semibold text-gray-900 mb-4">出品商品</h2>

    <!-- Tabs Navigation -->
    <div class="flex gap-2 overflow-x-auto pb-2 custom-scrollbar">
      <button
        v-for="item in items"
        :key="item.id"
        @click="handleSelect(item.id)"
        :class="[
          'flex-shrink-0 px-4 py-3 rounded-lg transition-all duration-200 border-2',
          isSelected(item.id)
            ? 'bg-blue-600 text-white border-blue-600 shadow-md'
            : 'bg-white text-gray-700 border-gray-200 hover:border-blue-300 hover:bg-blue-50'
        ]"
      >
        <div class="text-left min-w-[150px]">
          <div class="flex items-center gap-2 mb-1">
            <span class="font-semibold text-sm">
              {{ item.name }}
            </span>
          </div>
          <div class="flex items-center gap-2">
            <span
              :class="[
                'inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium',
                getStatusClass(item.status, isSelected(item.id))
              ]"
            >
              {{ getStatusLabel(item.status) }}
            </span>
            <span class="text-xs opacity-75">
              ロット {{ item.lot_number }}
            </span>
          </div>
          <div
            v-if="item.current_price > 0"
            :class="[
              'text-xs mt-1 font-semibold',
              isSelected(item.id) ? 'text-white' : 'text-gray-600'
            ]"
          >
            現在: {{ formatNumber(item.current_price) }}P
          </div>
        </div>
      </button>
    </div>

    <!-- Status Legend -->
    <div class="mt-4 pt-4 border-t border-gray-200 flex flex-wrap gap-4 text-xs text-gray-600">
      <div class="flex items-center">
        <span class="inline-block w-3 h-3 bg-yellow-100 border border-yellow-300 rounded-full mr-2"></span>
        待機中
      </div>
      <div class="flex items-center">
        <span class="inline-block w-3 h-3 bg-green-100 border border-green-300 rounded-full mr-2"></span>
        開始
      </div>
      <div class="flex items-center">
        <span class="inline-block w-3 h-3 bg-gray-100 border border-gray-300 rounded-full mr-2"></span>
        終了
      </div>
    </div>

    <!-- Item Summary Stats -->
    <div class="mt-4 grid grid-cols-3 gap-4 text-center">
      <div class="p-3 bg-gray-50 rounded-lg">
        <div class="text-2xl font-bold text-gray-900">{{ itemStats.total }}</div>
        <div class="text-xs text-gray-600 mt-1">全商品</div>
      </div>
      <div class="p-3 bg-green-50 rounded-lg">
        <div class="text-2xl font-bold text-green-600">{{ itemStats.started }}</div>
        <div class="text-xs text-gray-600 mt-1">開始済み</div>
      </div>
      <div class="p-3 bg-yellow-50 rounded-lg">
        <div class="text-2xl font-bold text-yellow-600">{{ itemStats.pending }}</div>
        <div class="text-xs text-gray-600 mt-1">待機中</div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  items: {
    type: Array,
    required: true,
    default: () => [],
  },
  selectedItemId: {
    type: [Number, String], // Support both Number and String (UUID)
    default: null,
  },
})

const emit = defineEmits(['select'])

// Computed: Item statistics
const itemStats = computed(() => {
  return {
    total: props.items.length,
    pending: props.items.filter((item) => item.status === 'pending').length,
    started: props.items.filter((item) => item.status === 'active').length,
    ended: props.items.filter((item) => item.status === 'ended').length,
  }
})

// Methods
function isSelected(itemId) {
  return props.selectedItemId === itemId
}

function handleSelect(itemId) {
  if (itemId !== props.selectedItemId) {
    emit('select', itemId)
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

function getStatusClass(status, isSelected) {
  if (isSelected) {
    // Selected tab: lighter colors for contrast with blue background
    const classes = {
      pending: 'bg-yellow-200 text-yellow-900',
      active: 'bg-green-200 text-green-900',
      ended: 'bg-gray-200 text-gray-900',
    }
    return classes[status] || 'bg-gray-200 text-gray-900'
  } else {
    // Unselected tab: normal colors
    const classes = {
      pending: 'bg-yellow-100 text-yellow-800',
      active: 'bg-green-100 text-green-800',
      ended: 'bg-gray-100 text-gray-800',
    }
    return classes[status] || 'bg-gray-100 text-gray-800'
  }
}
</script>

<style scoped>
/* Custom scrollbar for horizontal scroll */
.custom-scrollbar {
  scrollbar-width: thin;
  scrollbar-color: #cbd5e0 #f7fafc;
}

.custom-scrollbar::-webkit-scrollbar {
  height: 6px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: #f7fafc;
  border-radius: 3px;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #cbd5e0;
  border-radius: 3px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #a0aec0;
}
</style>
