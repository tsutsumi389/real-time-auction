<template>
  <div class="bg-white shadow rounded-lg p-6">
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
      <!-- 検索 -->
      <div class="md:col-span-2">
        <label for="keyword" class="block text-sm font-medium text-gray-700 mb-1">
          キーワード検索
        </label>
        <input
          id="keyword"
          :value="modelValue.keyword"
          @input="$emit('update:modelValue', { ...modelValue, keyword: $event.target.value })"
          type="text"
          placeholder="メールアドレス、表示名、またはIDで検索"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          @keyup.enter="$emit('search')"
        />
      </div>

      <!-- 状態フィルタ（複数選択） -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-1">状態フィルタ</label>
        <div class="space-y-2">
          <label class="flex items-center">
            <input
              type="checkbox"
              :checked="modelValue.status.includes('active')"
              @change="handleStatusToggle('active')"
              class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
            />
            <span class="ml-2 text-sm text-gray-700">有効</span>
          </label>
          <label class="flex items-center">
            <input
              type="checkbox"
              :checked="modelValue.status.includes('suspended')"
              @change="handleStatusToggle('suspended')"
              class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
            />
            <span class="ml-2 text-sm text-gray-700">停止中</span>
          </label>
          <label class="flex items-center">
            <input
              type="checkbox"
              :checked="modelValue.status.includes('deleted')"
              @change="handleStatusToggle('deleted')"
              class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
            />
            <span class="ml-2 text-sm text-gray-700">削除済み</span>
          </label>
        </div>
      </div>
    </div>

    <!-- アクションボタン -->
    <div class="mt-4 flex gap-2">
      <button
        @click="$emit('search')"
        class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-6 rounded-md transition-colors"
        :disabled="loading"
      >
        検索
      </button>
      <button
        @click="$emit('reset')"
        class="bg-gray-200 hover:bg-gray-300 text-gray-700 font-medium py-2 px-6 rounded-md transition-colors"
        :disabled="loading"
      >
        リセット
      </button>
    </div>
  </div>
</template>

<script setup>
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

const emit = defineEmits(['update:modelValue', 'search', 'reset', 'filter-change'])

function handleStatusToggle(status) {
  const currentStatus = [...props.modelValue.status]
  const index = currentStatus.indexOf(status)

  if (index > -1) {
    currentStatus.splice(index, 1)
  } else {
    currentStatus.push(status)
  }

  emit('update:modelValue', { ...props.modelValue, status: currentStatus })
  emit('filter-change')
}
</script>
