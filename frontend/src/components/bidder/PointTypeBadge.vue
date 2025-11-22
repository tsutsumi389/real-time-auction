<template>
  <span
    :class="[
      'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
      typeClass,
    ]"
  >
    {{ label }}
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  type: {
    type: String,
    required: true,
    validator: (value) => ['grant', 'reserve', 'release', 'consume', 'refund'].includes(value),
  },
})

const label = computed(() => {
  switch (props.type) {
    case 'grant':
      return '付与'
    case 'reserve':
      return '予約'
    case 'release':
      return '解放'
    case 'consume':
      return '消費'
    case 'refund':
      return '返金'
    default:
      return props.type
  }
})

const typeClass = computed(() => {
  switch (props.type) {
    case 'grant':
      return 'bg-green-100 text-green-800'
    case 'reserve':
      return 'bg-yellow-100 text-yellow-800'
    case 'release':
      return 'bg-blue-100 text-blue-800'
    case 'consume':
      return 'bg-red-100 text-red-800'
    case 'refund':
      return 'bg-green-100 text-green-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
})
</script>
