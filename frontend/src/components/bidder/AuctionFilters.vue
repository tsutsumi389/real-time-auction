<template>
  <div class="auction-filters lux-glass-strong rounded-xl p-4 sm:p-5 border border-lux-gold/10">
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <!-- Status Filter -->
      <div class="flex flex-wrap items-center gap-3">
        <span class="flex items-center gap-2 text-xs sm:text-sm font-medium text-lux-silver whitespace-nowrap">
          <Filter class="h-4 w-4 text-lux-gold/60" :stroke-width="1.5" />
          <span class="hidden sm:inline">表示:</span>
        </span>
        <div class="flex flex-wrap gap-2" role="radiogroup" aria-label="オークションステータスフィルタ">
          <button
            v-for="status in statusOptions"
            :key="status.value"
            @click="handleStatusChange(status.value)"
            :class="getStatusButtonClasses(status.value)"
            role="radio"
            :aria-checked="currentStatus === status.value"
            :aria-label="`${status.label}のオークションを表示`"
          >
            <span class="status-dot" :class="getStatusDotClass(status.value)"></span>
            {{ status.label }}
          </button>
        </div>
      </div>

      <!-- Sort Filter -->
      <div class="flex items-center gap-3">
        <span class="flex items-center gap-2 text-xs sm:text-sm font-medium text-lux-silver whitespace-nowrap">
          <ArrowUpDown class="h-4 w-4 text-lux-gold/60" :stroke-width="1.5" />
          <span class="hidden sm:inline">並び替え:</span>
        </span>
        <div class="relative">
          <select
            :value="currentSort"
            @change="handleSortChange($event.target.value)"
            class="sort-select appearance-none w-full sm:w-[200px] h-10 pl-4 pr-10 rounded-lg bg-lux-noir-light border border-lux-noir-soft text-sm text-lux-cream cursor-pointer focus:outline-none focus:border-lux-gold/50 focus:ring-2 focus:ring-lux-gold/10 transition-all duration-300"
            aria-label="並び替え順を選択"
          >
            <option
              v-for="sort in sortOptions"
              :key="sort.value"
              :value="sort.value"
              class="bg-lux-noir text-lux-cream"
            >
              {{ sort.label }}
            </option>
          </select>
          <div class="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none">
            <ChevronDown class="h-4 w-4 text-lux-silver/60" :stroke-width="1.5" />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { Filter, ArrowUpDown, ChevronDown } from 'lucide-vue-next'

const props = defineProps({
  currentStatus: {
    type: String,
    default: 'active'
  },
  currentSort: {
    type: String,
    default: 'started_at_desc'
  },
  statusOptions: {
    type: Array,
    default: () => [
      { value: 'active', label: '開催中' },
      { value: 'ended', label: '終了' },
      { value: 'cancelled', label: '中止' }
    ]
  },
  sortOptions: {
    type: Array,
    default: () => [
      { value: 'started_at_desc', label: '開始日時が新しい順' },
      { value: 'started_at_asc', label: '開始日時が古い順' },
      { value: 'updated_at_desc', label: '更新日時が新しい順' },
      { value: 'updated_at_asc', label: '更新日時が古い順' }
    ]
  }
})

const emit = defineEmits(['update:status', 'update:sort'])

const getStatusButtonClasses = (status) => {
  const isActive = props.currentStatus === status
  return [
    'status-btn flex items-center gap-2 px-4 py-2 rounded-lg text-xs sm:text-sm font-medium transition-all duration-300 whitespace-nowrap',
    isActive
      ? 'bg-lux-gold/15 text-lux-gold border border-lux-gold/40 shadow-[0_0_15px_rgba(212,175,55,0.15)]'
      : 'bg-lux-noir-light/50 text-lux-silver border border-lux-noir-soft hover:border-lux-gold/30 hover:text-lux-cream'
  ].join(' ')
}

const getStatusDotClass = (status) => {
  if (status === 'active') return 'dot-active'
  if (status === 'ended') return 'dot-ended'
  return 'dot-cancelled'
}

const handleStatusChange = (status) => {
  emit('update:status', status)
}

const handleSortChange = (value) => {
  emit('update:sort', value)
}
</script>

<style scoped>
.auction-filters {
  width: 100%;
}

/* Status dot indicator */
.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.dot-active {
  background-color: hsl(145 50% 50%);
  box-shadow: 0 0 8px hsl(145 50% 50% / 0.6);
}

.dot-ended {
  background-color: hsl(220 10% 60%);
}

.dot-cancelled {
  background-color: hsl(0 65% 50%);
}

/* Sort select styling */
.sort-select {
  background-color: hsl(0 0% 8%);
  border-color: hsl(0 0% 16%);
}

.sort-select option {
  background-color: hsl(0 0% 8%);
  color: hsl(45 30% 96%);
  padding: 8px;
}

/* Luxury color utilities */
.bg-lux-noir-light {
  background-color: hsl(0 0% 8%);
}

.bg-lux-noir-light\/50 {
  background-color: hsl(0 0% 8% / 0.5);
}

.border-lux-noir-soft {
  border-color: hsl(0 0% 16%);
}

.text-lux-silver {
  color: hsl(220 10% 70%);
}

.text-lux-cream {
  color: hsl(45 30% 96%);
}

.text-lux-gold {
  color: hsl(43 74% 49%);
}

.bg-lux-gold\/15 {
  background-color: hsl(43 74% 49% / 0.15);
}

.border-lux-gold\/40 {
  border-color: hsl(43 74% 49% / 0.4);
}

.hover\:border-lux-gold\/30:hover {
  border-color: hsl(43 74% 49% / 0.3);
}

.hover\:text-lux-cream:hover {
  color: hsl(45 30% 96%);
}

.text-lux-gold\/60 {
  color: hsl(43 74% 49% / 0.6);
}

.text-lux-silver\/60 {
  color: hsl(220 10% 70% / 0.6);
}

.border-lux-gold\/10 {
  border-color: hsl(43 74% 49% / 0.1);
}

.focus\:border-lux-gold\/50:focus {
  border-color: hsl(43 74% 49% / 0.5);
}

.focus\:ring-lux-gold\/10:focus {
  --tw-ring-color: hsl(43 74% 49% / 0.1);
}
</style>
