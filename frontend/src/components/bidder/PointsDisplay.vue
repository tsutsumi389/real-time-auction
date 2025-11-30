<template>
  <div class="inline-flex items-center gap-2 bg-white rounded-lg shadow-sm px-4 py-2">
    <span class="text-lg font-bold text-green-600">P</span>
    <span class="text-lg font-bold text-green-600">{{ formatNumber(points.available) }}</span>
    <span class="text-sm text-gray-500">pt</span>
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
  // undefined, null, NaN を 0 として扱う
  const safeValue = typeof value === 'number' && !isNaN(value) ? value : 0
  return new Intl.NumberFormat('ja-JP').format(safeValue)
}
</script>
