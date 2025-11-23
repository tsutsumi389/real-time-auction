<template>
  <div class="bg-white shadow rounded-lg p-6">
    <div class="grid grid-cols-1 md:grid-cols-5 gap-4">
      <!-- 状態フィルタ -->
      <div>
        <label for="status" class="block text-sm font-medium text-gray-700 mb-1"> 状態 </label>
        <select
          id="status"
          :value="modelValue.status"
          @change="handleStatusChange"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="">すべて</option>
          <option value="pending">非公開</option>
          <option value="active">公開中</option>
          <option value="ended">終了</option>
          <option value="cancelled">中止</option>
        </select>
      </div>

      <!-- 作成日フィルタ（この日以降） -->
      <div>
        <label for="created-after" class="block text-sm font-medium text-gray-700 mb-1"> 作成日（以降） </label>
        <input
          id="created-after"
          type="date"
          :value="modelValue.createdAfter"
          @change="handleCreatedAfterChange"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <!-- 更新日フィルタ（この日以前） -->
      <div>
        <label for="updated-before" class="block text-sm font-medium text-gray-700 mb-1"> 更新日（以前） </label>
        <input
          id="updated-before"
          type="date"
          :value="modelValue.updatedBefore"
          @change="handleUpdatedBeforeChange"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        />
      </div>

      <!-- ソート -->
      <div>
        <label for="sort" class="block text-sm font-medium text-gray-700 mb-1"> 並び順 </label>
        <select
          id="sort"
          :value="sortValue"
          @change="handleSortChange"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="created_at_desc">作成日（新しい順）</option>
          <option value="created_at_asc">作成日（古い順）</option>
          <option value="updated_at_desc">更新日（新しい順）</option>
          <option value="updated_at_asc">更新日（古い順）</option>
          <option value="id_desc">ID（降順）</option>
          <option value="id_asc">ID（昇順）</option>
        </select>
      </div>

      <!-- アクションボタン -->
      <div class="flex items-end">
        <button
          @click="$emit('reset')"
          class="w-full bg-gray-200 hover:bg-gray-300 text-gray-700 font-medium py-2 px-4 rounded-md transition-colors"
          :disabled="loading"
        >
          リセット
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
  },
  loading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue', 'filter-change'])

const sortValue = computed(() => {
  return `${props.modelValue.sort}_${props.modelValue.order}`
})

function handleStatusChange(event) {
  emit('update:modelValue', { ...props.modelValue, status: event.target.value })
  emit('filter-change')
}

function handleCreatedAfterChange(event) {
  emit('update:modelValue', { ...props.modelValue, createdAfter: event.target.value })
  emit('filter-change')
}

function handleUpdatedBeforeChange(event) {
  emit('update:modelValue', { ...props.modelValue, updatedBefore: event.target.value })
  emit('filter-change')
}

function handleSortChange(event) {
  const [sort, order] = event.target.value.split('_')
  const newSort = event.target.value.substring(0, event.target.value.lastIndexOf('_'))
  const newOrder = event.target.value.substring(event.target.value.lastIndexOf('_') + 1)
  emit('update:modelValue', { ...props.modelValue, sort: newSort, order: newOrder })
  emit('filter-change')
}
</script>
