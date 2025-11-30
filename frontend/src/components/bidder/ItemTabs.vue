<template>
  <div class="bg-white rounded-lg shadow-sm">
    <!-- Collapsible Header -->
    <button
      @click="toggleCollapsed"
      class="w-full px-4 py-3 flex items-center justify-between hover:bg-gray-50 transition-colors rounded-lg"
    >
      <div class="flex items-center gap-3">
        <h2 class="text-sm font-semibold text-gray-700">出品商品</h2>
        <!-- Progress Indicator -->
        <div class="flex items-center gap-2">
          <div class="w-24 h-1.5 bg-gray-200 rounded-full overflow-hidden">
            <div 
              class="h-full bg-auction-gold rounded-full transition-all duration-500"
              :style="{ width: `${progressPercent}%` }"
            ></div>
          </div>
          <span class="text-xs text-gray-500">{{ itemStats.ended }}/{{ itemStats.total }}</span>
        </div>
      </div>
      <div class="flex items-center gap-2">
        <!-- Status Summary Pills -->
        <div class="hidden sm:flex items-center gap-1.5 text-xs">
          <span v-if="itemStats.started > 0" class="px-2 py-0.5 rounded-full bg-green-100 text-green-700 font-medium">
            {{ itemStats.started }}件 開催中
          </span>
          <span v-if="itemStats.pending > 0" class="px-2 py-0.5 rounded-full bg-yellow-100 text-yellow-700 font-medium">
            {{ itemStats.pending }}件 待機
          </span>
        </div>
        <!-- Collapse Arrow -->
        <svg 
          :class="['w-5 h-5 text-gray-400 transition-transform duration-200', isCollapsed ? '' : 'rotate-180']"
          fill="none" 
          stroke="currentColor" 
          viewBox="0 0 24 24"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </div>
    </button>

    <!-- Collapsible Content -->
    <div 
      :class="[
        'overflow-hidden transition-all duration-300 ease-in-out',
        isCollapsed ? 'max-h-0' : 'max-h-[500px]'
      ]"
    >
      <div class="px-4 pb-4">
        <!-- Item List (Vertical) -->
        <div class="space-y-1.5 max-h-64 overflow-y-auto custom-scrollbar">
          <button
            v-for="item in sortedItems"
            :key="item.id"
            @click="handleSelect(item.id)"
            :class="[
              'w-full flex items-center justify-between px-3 py-2 rounded-lg transition-all duration-150 text-left',
              isSelected(item.id)
                ? 'bg-auction-gold/10 border border-auction-gold/30 ring-1 ring-auction-gold/20'
                : 'bg-gray-50 hover:bg-gray-100 border border-transparent'
            ]"
          >
            <div class="flex items-center gap-3 min-w-0">
              <!-- Lot Number -->
              <span 
                :class="[
                  'flex-shrink-0 w-8 h-8 flex items-center justify-center rounded-full text-xs font-bold',
                  isSelected(item.id) 
                    ? 'bg-auction-gold text-white' 
                    : 'bg-gray-200 text-gray-600'
                ]"
              >
                {{ item.lot_number }}
              </span>
              <!-- Item Info -->
              <div class="min-w-0 flex-1">
                <div 
                  :class="[
                    'text-sm font-medium truncate',
                    isSelected(item.id) ? 'text-auction-gold-dark' : 'text-gray-800'
                  ]"
                >
                  {{ item.name }}
                </div>
                <div v-if="item.current_price > 0" class="text-xs text-gray-500">
                  現在: {{ formatNumber(item.current_price) }}P
                </div>
              </div>
            </div>
            <!-- Status Badge -->
            <span
              :class="[
                'flex-shrink-0 inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium',
                getStatusClass(item.status)
              ]"
            >
              {{ getStatusLabel(item.status) }}
            </span>
          </button>
        </div>

        <!-- Status Legend (Compact) -->
        <div class="mt-3 pt-3 border-t border-gray-100 flex items-center justify-center gap-4 text-xs text-gray-500">
          <div class="flex items-center gap-1">
            <span class="w-2 h-2 bg-yellow-400 rounded-full"></span>
            <span>待機</span>
          </div>
          <div class="flex items-center gap-1">
            <span class="w-2 h-2 bg-green-500 rounded-full"></span>
            <span>開催中</span>
          </div>
          <div class="flex items-center gap-1">
            <span class="w-2 h-2 bg-gray-400 rounded-full"></span>
            <span>終了</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref } from 'vue'

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

// Local state
const isCollapsed = ref(true) // Start collapsed

// Computed: Sorted items (active first, then pending, then ended)
const sortedItems = computed(() => {
  const statusOrder = { active: 0, pending: 1, ended: 2 }
  return [...props.items].sort((a, b) => {
    const orderA = statusOrder[a.status] ?? 3
    const orderB = statusOrder[b.status] ?? 3
    if (orderA !== orderB) return orderA - orderB
    return (a.lot_number || 0) - (b.lot_number || 0)
  })
})

// Computed: Item statistics
const itemStats = computed(() => {
  return {
    total: props.items.length,
    pending: props.items.filter((item) => item.status === 'pending').length,
    started: props.items.filter((item) => item.status === 'active').length,
    ended: props.items.filter((item) => item.status === 'ended').length,
  }
})

// Computed: Progress percentage
const progressPercent = computed(() => {
  if (itemStats.value.total === 0) return 0
  return Math.round((itemStats.value.ended / itemStats.value.total) * 100)
})

// Methods
function toggleCollapsed() {
  isCollapsed.value = !isCollapsed.value
}

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
    pending: '待機',
    active: '開催中',
    ended: '終了',
  }
  return labels[status] || status
}

function getStatusClass(status) {
  const classes = {
    pending: 'bg-yellow-100 text-yellow-700',
    active: 'bg-green-100 text-green-700',
    ended: 'bg-gray-100 text-gray-600',
  }
  return classes[status] || 'bg-gray-100 text-gray-600'
}
</script>

<style scoped>
/* Custom scrollbar for vertical scroll */
.custom-scrollbar {
  scrollbar-width: thin;
  scrollbar-color: #cbd5e0 #f7fafc;
}

.custom-scrollbar::-webkit-scrollbar {
  width: 4px;
}

.custom-scrollbar::-webkit-scrollbar-track {
  background: #f7fafc;
  border-radius: 2px;
}

.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #cbd5e0;
  border-radius: 2px;
}

.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #a0aec0;
}
</style>