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
    validator: (value) => ['active', 'suspended', 'deleted'].includes(value),
  },
})

const label = computed(() => {
  switch (props.status) {
    case 'active':
      return '有効'
    case 'suspended':
      return '停止中'
    case 'deleted':
      return '削除済み'
    default:
      return props.status
  }
})

const statusClass = computed(() => {
  switch (props.status) {
    case 'active':
      return 'bg-green-100 text-green-800'
    case 'suspended':
      return 'bg-yellow-100 text-yellow-800'
    case 'deleted':
      return 'bg-gray-100 text-gray-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
})
</script>
