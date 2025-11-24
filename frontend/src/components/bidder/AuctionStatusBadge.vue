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
  const baseClasses = 'inline-flex items-center px-3 py-1 rounded-full text-sm font-medium'

  const statusClasses = {
    pending: 'bg-yellow-100 text-yellow-800',
    active: 'bg-green-100 text-green-800',
    ended: 'bg-gray-100 text-gray-800',
    cancelled: 'bg-red-100 text-red-800'
  }

  return `${baseClasses} ${statusClasses[props.status] || statusClasses.ended}`
})
</script>
