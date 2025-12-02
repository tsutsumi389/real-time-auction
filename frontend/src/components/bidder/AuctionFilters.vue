<template>
  <div class="auction-filters space-y-4">
    <!-- Status filter -->
    <div class="flex flex-wrap gap-2 sm:gap-4 items-center">
      <span class="flex items-center gap-1.5 text-xs sm:text-sm font-medium text-gray-700 whitespace-nowrap">
        <Filter class="h-4 w-4 text-muted-foreground" :stroke-width="1.5" />
        表示:
      </span>
      <div class="flex flex-wrap gap-2" role="radiogroup" aria-label="オークションステータスフィルタ">
        <button
          v-for="status in statusOptions"
          :key="status.value"
          @click="handleStatusChange(status.value)"
          :class="getStatusButtonClasses(status.value)"
          role="radio"
          :aria-checked="currentStatus === status.value"
          :aria-label="`${status.label}のオークションを表示`"
        >
          {{ status.label }}
        </button>
      </div>
    </div>

    <!-- Sort filter -->
    <div class="flex flex-wrap gap-2 sm:gap-4 items-center">
      <span class="flex items-center gap-1.5 text-xs sm:text-sm font-medium text-gray-700 whitespace-nowrap">
        <ArrowUpDown class="h-4 w-4 text-muted-foreground" :stroke-width="1.5" />
        並び替え:
      </span>
      <Select :model-value="currentSort" @update:model-value="handleSortChange">
        <SelectTrigger class="w-full sm:w-[220px]" aria-label="並び替え順を選択">
          <SelectValue placeholder="並び替え" />
        </SelectTrigger>
        <SelectContent>
          <SelectItem
            v-for="sort in sortOptions"
            :key="sort.value"
            :value="sort.value"
          >
            {{ sort.label }}
          </SelectItem>
        </SelectContent>
      </Select>
    </div>
  </div>
</template>

<script setup>
import { Filter, ArrowUpDown } from 'lucide-vue-next'
import Select from '@/components/ui/Select.vue'
import SelectTrigger from '@/components/ui/SelectTrigger.vue'
import SelectContent from '@/components/ui/SelectContent.vue'
import SelectItem from '@/components/ui/SelectItem.vue'
import SelectValue from '@/components/ui/SelectValue.vue'

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
  const base = 'px-3 py-2 sm:px-4 sm:py-2 rounded-lg text-xs sm:text-sm font-medium transition-all duration-200 whitespace-nowrap focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2'
  const active = 'bg-primary text-primary-foreground shadow-sm'
  const inactive = 'bg-muted text-muted-foreground hover:bg-accent hover:text-accent-foreground'

  return `${base} ${props.currentStatus === status ? active : inactive}`
}

const handleStatusChange = (status) => {
  emit('update:status', status)
}

const handleSortChange = (value) => {
  emit('update:sort', value)
}
</script>

<style scoped>
.auction-filters {
  width: 100%;
}
</style>
