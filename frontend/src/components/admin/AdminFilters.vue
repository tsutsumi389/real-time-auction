<template>
  <div class="bg-white shadow rounded-lg p-6">
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <!-- 検索 -->
      <div>
        <label for="search" class="block text-sm font-medium text-gray-700 mb-1"> メールアドレス検索 </label>
        <input
          id="search"
          :value="modelValue.search"
          @input="$emit('update:modelValue', { ...modelValue, search: $event.target.value })"
          type="text"
          placeholder="例: admin@example.com"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          @keyup.enter="$emit('search')"
        />
      </div>

      <!-- ロールフィルタ -->
      <div>
        <label for="role" class="block text-sm font-medium text-gray-700 mb-1"> ロール </label>
        <select
          id="role"
          :value="modelValue.role"
          @change="handleRoleChange"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        >
          <option value="">すべて</option>
          <option value="system_admin">システム管理者</option>
          <option value="auctioneer">オークショニア</option>
        </select>
      </div>

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
          <option value="active">有効</option>
          <option value="inactive">停止中</option>
        </select>
      </div>

      <!-- アクションボタン -->
      <div class="flex items-end gap-2">
        <button
          @click="$emit('search')"
          class="flex-1 bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md transition-colors"
          :disabled="loading"
        >
          検索
        </button>
        <button
          @click="$emit('reset')"
          class="flex-1 bg-gray-200 hover:bg-gray-300 text-gray-700 font-medium py-2 px-4 rounded-md transition-colors"
          :disabled="loading"
        >
          リセット
        </button>
      </div>
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

function handleRoleChange(event) {
  emit('update:modelValue', { ...props.modelValue, role: event.target.value })
  emit('filter-change')
}

function handleStatusChange(event) {
  emit('update:modelValue', { ...props.modelValue, status: event.target.value })
  emit('filter-change')
}
</script>
