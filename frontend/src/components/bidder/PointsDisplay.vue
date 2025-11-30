<template>
  <div class="inline-flex items-center gap-2 bg-white rounded-lg shadow-sm px-4 py-2">
    <svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
    <span class="text-sm text-gray-600">利用可能:</span>
    <span class="text-lg font-bold text-green-600">{{ formatNumber(points.available) }}</span>
    <span class="text-sm text-gray-500">pts</span>
  </div>
</template>

<script setup>
const props = defineProps({
  points: {
    type: Object,
    required: true,
    validator: (value) => {
      return (
        typeof value.total === 'number' &&
        typeof value.available === 'number' &&
        typeof value.reserved === 'number'
      )
    },
  },
})

// Format number with comma separator
function formatNumber(value) {
  return new Intl.NumberFormat('ja-JP').format(value)
}
</script>
