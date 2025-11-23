<template>
  <span
    :class="[
      'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
      statusClass,
    ]"
  >
    {{ label }}
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: {
    type: String,
    required: true,
    validator: (value) => ['pending', 'active', 'ended', 'cancelled'].includes(value),
  },
})

const label = computed(() => {
  switch (props.status) {
    case 'pending':
      return '非公開'
    case 'active':
      return '公開中'
    case 'ended':
      return '終了'
    case 'cancelled':
      return '中止'
    default:
      return props.status
  }
})

const statusClass = computed(() => {
  switch (props.status) {
    case 'pending':
      return 'bg-gray-100 text-gray-800'
    case 'active':
      return 'bg-green-600 text-white'
    case 'ended':
      return 'bg-blue-600 text-white'
    case 'cancelled':
      return 'bg-red-600 text-white'
    default:
      return 'bg-gray-100 text-gray-800'
  }
})
</script>
