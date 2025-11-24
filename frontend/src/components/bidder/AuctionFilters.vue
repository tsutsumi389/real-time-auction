<template>
  <div class="auction-filters space-y-4">
    <!-- ステータスフィルタ -->
    <div class="flex flex-wrap gap-2 sm:gap-4 items-center">
      <span class="text-xs sm:text-sm font-medium text-gray-700 whitespace-nowrap">表示:</span>
      <div class="flex flex-wrap gap-2">
        <button
          v-for="status in statusOptions"
          :key="status.value"
          @click="handleStatusChange(status.value)"
          :class="getStatusButtonClasses(status.value)"
        >
          {{ status.label }}
        </button>
      </div>
    </div>

    <!-- ソートフィルタ -->
    <div class="flex flex-wrap gap-2 sm:gap-4 items-center">
      <span class="text-xs sm:text-sm font-medium text-gray-700 whitespace-nowrap">並び替え:</span>
      <select
        :value="currentSort"
        @change="handleSortChange"
        class="px-3 py-2 sm:px-4 sm:py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm sm:text-base"
      >
        <option v-for="sort in sortOptions" :key="sort.value" :value="sort.value">
          {{ sort.label }}
        </option>
      </select>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  currentStatus: {
    type: String,
    default: 'active'
  },
  currentSort: {
    type: String,
    default: 'started_at_desc'
  },
  statusOptions: {
    type: Array,
    default: () => [
      { value: 'active', label: '開催中' },
      { value: 'ended', label: '終了' },
      { value: 'cancelled', label: '中止' }
    ]
  },
  sortOptions: {
    type: Array,
    default: () => [
      { value: 'started_at_desc', label: '開始日時が新しい順' },
      { value: 'started_at_asc', label: '開始日時が古い順' },
      { value: 'updated_at_desc', label: '更新日時が新しい順' },
      { value: 'updated_at_asc', label: '更新日時が古い順' }
    ]
  }
})

const emit = defineEmits(['update:status', 'update:sort'])

const getStatusButtonClasses = (status) => {
  const base = 'px-3 py-2 sm:px-4 sm:py-2 rounded-lg text-xs sm:text-sm font-medium transition-colors whitespace-nowrap'
  const active = 'bg-blue-600 text-white'
  const inactive = 'bg-gray-100 text-gray-700 hover:bg-gray-200'

  return `${base} ${props.currentStatus === status ? active : inactive}`
}

const handleStatusChange = (status) => {
  emit('update:status', status)
}

const handleSortChange = (event) => {
  emit('update:sort', event.target.value)
}
</script>

<style scoped>
.auction-filters {
  width: 100%;
}
</style>
