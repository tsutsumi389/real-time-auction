<template>
  <span :class="badgeClasses">
    <span v-if="props.status === 'active'" class="inline-block w-2 h-2 bg-auction-green-racing rounded-full mr-2 animate-pulse"></span>
    {{ label }}
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: {
    type: String,
    required: true,
    validator: (value) => ['pending', 'active', 'ended', 'cancelled'].includes(value)
  }
})

const label = computed(() => {
  const labels = {
    pending: '準備中',
    active: '開催中',
    ended: '終了',
    cancelled: '中止'
  }
  return labels[props.status] || props.status
})

const badgeClasses = computed(() => {
  const baseClasses = 'inline-flex items-center px-3 py-1 rounded-full text-sm font-medium shadow-sm'

  const statusClasses = {
    pending: 'bg-amber-100 text-amber-800',
    active: 'bg-auction-green-racing/10 text-auction-green-racing border border-auction-green-racing/30',
    ended: 'bg-muted text-muted-foreground',
    cancelled: 'bg-red-100 text-red-800'
  }

  return `${baseClasses} ${statusClasses[props.status] || statusClasses.ended}`
})
</script>

<style scoped>
/* Accessibility: Respect reduced motion preference */
@media (prefers-reduced-motion: reduce) {
  .animate-pulse {
    animation: none !important;
  }
}
</style>
