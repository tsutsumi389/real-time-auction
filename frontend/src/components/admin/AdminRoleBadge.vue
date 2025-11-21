<template>
  <span
    :class="[
      'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
      role === 'system_admin' ? 'bg-purple-100 text-purple-800' : 'bg-blue-100 text-blue-800',
    ]"
  >
    {{ label }}
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  role: {
    type: String,
    required: true,
    validator: (value) => ['system_admin', 'auctioneer'].includes(value),
  },
})

const label = computed(() => {
  return props.role === 'system_admin' ? 'システム管理者' : 'オークショニア'
})
</script>
