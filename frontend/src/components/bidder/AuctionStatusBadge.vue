<template>
  <span :class="badgeClasses">
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
  const baseClasses = 'inline-flex items-center px-3 py-1 rounded-full text-sm font-medium border'

  const statusClasses = {
    pending: 'bg-yellow-500/10 text-yellow-400 border-yellow-500/30',
    active: 'bg-emerald-500/10 text-emerald-400 border-emerald-500/30',
    ended: 'bg-lux-silver/10 text-lux-silver/80 border-lux-silver/20',
    cancelled: 'bg-red-500/10 text-red-400 border-red-500/30'
  }

  return `${baseClasses} ${statusClasses[props.status] || statusClasses.ended}`
})
</script>

<style scoped>
.bg-lux-silver\/10 {
  background-color: hsl(220 10% 70% / 0.1);
}

.text-lux-silver\/80 {
  color: hsl(220 10% 70% / 0.8);
}

.border-lux-silver\/20 {
  border-color: hsl(220 10% 70% / 0.2);
}
</style>
