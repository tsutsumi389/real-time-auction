<template>
  <div class="item-tabs-container rounded-2xl overflow-hidden">
    <!-- Collapsible Header -->
    <button
      @click="toggleCollapsed"
      class="w-full px-6 py-5 flex items-center justify-between hover:bg-lux-noir-light/30 transition-all duration-300 group"
    >
      <div class="flex items-center gap-4">
        <!-- Icon -->
        <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-lux-gold/20 to-lux-gold/5 border border-lux-gold/30 flex items-center justify-center group-hover:border-lux-gold/50 transition-colors">
          <svg class="w-5 h-5 text-lux-gold" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
          </svg>
        </div>

        <div>
          <h2 class="font-display text-lg text-lux-cream">出品アイテム</h2>
          <p class="text-[10px] text-lux-silver/50 uppercase tracking-widest">Auction Items</p>
        </div>

        <!-- Progress Indicator -->
        <div class="hidden sm:flex items-center gap-3 ml-4 pl-4 border-l border-lux-gold/20">
          <div class="w-28 h-2 bg-lux-noir-medium rounded-full overflow-hidden border border-lux-noir-soft">
            <div
              class="h-full bg-gradient-to-r from-lux-gold to-lux-gold-light rounded-full transition-all duration-700 ease-out relative"
              :style="{ width: `${progressPercent}%` }"
            >
              <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/30 to-transparent" style="animation: shimmer-progress 2s ease-in-out infinite;"></div>
            </div>
          </div>
          <span class="text-xs text-lux-gold font-semibold tabular-nums">{{ itemStats.ended }}/{{ itemStats.total }}</span>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <!-- Status Summary Pills -->
        <div class="hidden sm:flex items-center gap-2">
          <span v-if="itemStats.started > 0" class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-emerald-500/15 border border-emerald-500/30 text-xs font-semibold text-emerald-400">
            <span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"></span>
            {{ itemStats.started }} Live
          </span>
          <span v-if="itemStats.pending > 0" class="inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg bg-amber-500/15 border border-amber-500/30 text-xs font-semibold text-amber-400">
            {{ itemStats.pending }} Pending
          </span>
        </div>

        <!-- Collapse Arrow -->
        <div class="w-9 h-9 rounded-xl bg-lux-noir-medium border border-lux-gold/20 flex items-center justify-center group-hover:border-lux-gold/40 transition-all">
          <svg
            :class="['w-4 h-4 text-lux-gold transition-transform duration-300', isCollapsed ? '' : 'rotate-180']"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
          </svg>
        </div>
      </div>
    </button>

    <!-- Collapsible Content -->
    <Transition
      enter-active-class="transition-all duration-300 ease-out"
      enter-from-class="max-h-0 opacity-0"
      enter-to-class="max-h-[500px] opacity-100"
      leave-active-class="transition-all duration-200 ease-in"
      leave-from-class="max-h-[500px] opacity-100"
      leave-to-class="max-h-0 opacity-0"
    >
      <div v-show="!isCollapsed" class="overflow-hidden">
        <div class="px-4 pb-4 border-t border-lux-gold/15 bg-gradient-to-b from-lux-noir-light/30 to-transparent">
          <!-- Item List -->
          <div class="mt-4 space-y-2 max-h-72 overflow-y-auto lux-scrollbar pr-1">
            <button
              v-for="item in sortedItems"
              :key="item.id"
              @click="handleSelect(item.id)"
              :class="[
                'w-full flex items-center justify-between px-4 py-3.5 rounded-xl transition-all duration-200 text-left group',
                isSelected(item.id)
                  ? 'item-row-selected'
                  : 'item-row-default'
              ]"
            >
              <div class="flex items-center gap-3 min-w-0">
                <!-- Lot Number -->
                <span
                  :class="[
                    'flex-shrink-0 w-11 h-11 flex items-center justify-center rounded-xl text-sm font-bold transition-all duration-300',
                    isSelected(item.id)
                      ? 'bg-gradient-to-br from-lux-gold to-lux-gold-dark text-lux-noir shadow-lg shadow-lux-gold/30'
                      : item.status === 'active'
                      ? 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'
                      : 'bg-lux-noir-medium text-lux-silver/70 border border-lux-noir-soft group-hover:text-lux-cream group-hover:border-lux-gold/20'
                  ]"
                >
                  {{ item.lot_number }}
                </span>

                <!-- Item Info -->
                <div class="min-w-0 flex-1">
                  <div
                    :class="[
                      'text-sm font-semibold truncate transition-colors',
                      isSelected(item.id) ? 'text-lux-gold' : 'text-lux-cream'
                    ]"
                  >
                    {{ item.name }}
                  </div>
                  <div class="flex items-center gap-2 mt-1">
                    <span
                      v-if="item.current_price > 0"
                      :class="[
                        'text-xs tabular-nums',
                        isSelected(item.id) ? 'text-lux-gold/70' : 'text-lux-silver/50'
                      ]"
                    >
                      {{ formatNumber(item.current_price) }} pts
                    </span>
                    <span v-else class="text-xs text-lux-silver/40">
                      未開始
                    </span>
                  </div>
                </div>
              </div>

              <!-- Status Badge -->
              <span
                :class="[
                  'flex-shrink-0 inline-flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold',
                  getStatusClass(item.status)
                ]"
              >
                <span
                  v-if="item.status === 'active'"
                  class="w-1.5 h-1.5 rounded-full bg-current animate-pulse"
                ></span>
                {{ getStatusLabel(item.status) }}
              </span>
            </button>
          </div>

          <!-- Status Legend -->
          <div class="mt-4 pt-4 border-t border-lux-gold/10 flex items-center justify-center gap-8 text-[10px] text-lux-silver/50 uppercase tracking-wider">
            <div class="flex items-center gap-2">
              <span class="w-2.5 h-2.5 rounded-full bg-amber-400/80"></span>
              <span>待機中</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="w-2.5 h-2.5 rounded-full bg-emerald-400 animate-pulse"></span>
              <span>開催中</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="w-2.5 h-2.5 rounded-full bg-lux-silver/40"></span>
              <span>終了</span>
            </div>
          </div>
        </div>
      </div>
    </Transition>
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
    type: [Number, String],
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
    pending: 'Pending',
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
/* Container Styling */
.item-tabs-container {
  background: hsl(0 0% 4%);
  border: 1px solid rgba(212, 175, 55, 0.2);
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.5),
    0 0 0 1px rgba(212, 175, 55, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.03);
}

/* Item Row Styles */
.item-row-selected {
  background: linear-gradient(135deg, rgba(212, 175, 55, 0.15) 0%, rgba(212, 175, 55, 0.05) 100%);
  border: 1px solid rgba(212, 175, 55, 0.4);
  box-shadow:
    0 0 20px rgba(212, 175, 55, 0.1),
    inset 0 1px 0 rgba(255, 255, 255, 0.05);
}

.item-row-default {
  background: linear-gradient(135deg, rgba(25, 25, 25, 0.6) 0%, rgba(20, 20, 20, 0.8) 100%);
  border: 1px solid rgba(60, 60, 60, 0.2);
}

.item-row-default:hover {
  background: linear-gradient(135deg, rgba(35, 35, 35, 0.7) 0%, rgba(25, 25, 25, 0.9) 100%);
  border-color: rgba(212, 175, 55, 0.2);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.02);
}

/* Progress Shimmer Animation */
@keyframes shimmer-progress {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}
</style>
