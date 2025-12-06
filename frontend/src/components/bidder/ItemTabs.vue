<template>
  <div class="lux-glass-strong rounded-2xl overflow-hidden">
    <!-- Collapsible Header -->
    <button
      @click="toggleCollapsed"
      class="w-full px-6 py-4 flex items-center justify-between hover:bg-lux-noir-light/50 transition-colors"
    >
      <div class="flex items-center gap-4">
        <h2 class="font-display text-lg text-lux-cream">Auction Items</h2>

        <!-- Progress Indicator -->
        <div class="flex items-center gap-3">
          <div class="w-32 h-1.5 bg-lux-noir-light rounded-full overflow-hidden">
            <div
              class="h-full bg-gradient-to-r from-lux-gold to-lux-gold-light rounded-full transition-all duration-700 ease-out"
              :style="{ width: `${progressPercent}%` }"
            ></div>
          </div>
          <span class="text-xs text-lux-silver tabular-nums">{{ itemStats.ended }}/{{ itemStats.total }}</span>
        </div>
      </div>

      <div class="flex items-center gap-3">
        <!-- Status Summary Pills -->
        <div class="hidden sm:flex items-center gap-2">
          <span v-if="itemStats.started > 0" class="px-2.5 py-1 rounded-full lux-badge-active text-xs font-medium">
            {{ itemStats.started }} Live
          </span>
          <span v-if="itemStats.pending > 0" class="px-2.5 py-1 rounded-full lux-badge-pending text-xs font-medium">
            {{ itemStats.pending }} Pending
          </span>
        </div>

        <!-- Collapse Arrow -->
        <div class="w-8 h-8 rounded-full bg-lux-noir-light flex items-center justify-center">
          <svg
            :class="['w-4 h-4 text-lux-silver transition-transform duration-300', isCollapsed ? '' : 'rotate-180']"
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
      enter-to-class="max-h-[400px] opacity-100"
      leave-active-class="transition-all duration-200 ease-in"
      leave-from-class="max-h-[400px] opacity-100"
      leave-to-class="max-h-0 opacity-0"
    >
      <div v-show="!isCollapsed" class="overflow-hidden">
        <div class="px-4 pb-4 border-t border-lux-gold/10">
          <!-- Item List -->
          <div class="mt-4 space-y-2 max-h-64 overflow-y-auto lux-scrollbar">
            <button
              v-for="item in sortedItems"
              :key="item.id"
              @click="handleSelect(item.id)"
              :class="[
                'w-full flex items-center justify-between px-4 py-3 rounded-xl transition-all duration-200 text-left group',
                isSelected(item.id)
                  ? 'bg-gradient-to-r from-lux-gold/15 via-lux-gold/10 to-transparent border border-lux-gold/30'
                  : 'bg-lux-noir-light/50 hover:bg-lux-noir-light border border-transparent hover:border-lux-gold/10'
              ]"
            >
              <div class="flex items-center gap-3 min-w-0">
                <!-- Lot Number -->
                <span
                  :class="[
                    'flex-shrink-0 w-10 h-10 flex items-center justify-center rounded-lg text-sm font-semibold transition-colors',
                    isSelected(item.id)
                      ? 'bg-lux-gold text-lux-noir'
                      : 'bg-lux-noir-medium text-lux-silver group-hover:text-lux-cream'
                  ]"
                >
                  {{ item.lot_number }}
                </span>

                <!-- Item Info -->
                <div class="min-w-0 flex-1">
                  <div
                    :class="[
                      'text-sm font-medium truncate transition-colors',
                      isSelected(item.id) ? 'text-lux-gold' : 'text-lux-cream'
                    ]"
                  >
                    {{ item.name }}
                  </div>
                  <div v-if="item.current_price > 0" class="text-xs text-lux-silver/60 mt-0.5">
                    Current: {{ formatNumber(item.current_price) }} pts
                  </div>
                </div>
              </div>

              <!-- Status Badge -->
              <span
                :class="[
                  'flex-shrink-0 inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium',
                  getStatusClass(item.status)
                ]"
              >
                {{ getStatusLabel(item.status) }}
              </span>
            </button>
          </div>

          <!-- Status Legend -->
          <div class="mt-4 pt-4 border-t border-lux-gold/10 flex items-center justify-center gap-6 text-[10px] text-lux-silver/60 uppercase tracking-wider">
            <div class="flex items-center gap-1.5">
              <span class="w-2 h-2 rounded-full bg-amber-400"></span>
              <span>Pending</span>
            </div>
            <div class="flex items-center gap-1.5">
              <span class="w-2 h-2 rounded-full bg-emerald-400"></span>
              <span>Live</span>
            </div>
            <div class="flex items-center gap-1.5">
              <span class="w-2 h-2 rounded-full bg-lux-silver/40"></span>
              <span>Ended</span>
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
