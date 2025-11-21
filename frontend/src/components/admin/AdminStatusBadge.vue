<template>
  <span
    :class="[
      'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
      status === 'active' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800',
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
    validator: (value) => ['active', 'inactive'].includes(value),
  },
})

const label = computed(() => {
  return props.status === 'active' ? '有効' : '停止中'
})
</script>
