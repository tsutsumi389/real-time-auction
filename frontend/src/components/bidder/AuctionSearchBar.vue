<template>
  <div class="search-bar">
    <div class="flex gap-2 sm:gap-4">
      <input
        v-model="localKeyword"
        type="text"
        :placeholder="placeholder"
        class="flex-1 px-3 py-2 sm:px-4 sm:py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent text-sm sm:text-base"
        @keyup.enter="handleSearch"
      />
      <button
        @click="handleSearch"
        :disabled="loading"
        class="px-4 py-2 sm:px-6 sm:py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed text-sm sm:text-base font-medium whitespace-nowrap"
      >
        {{ searchButtonText }}
      </button>
      <button
        v-if="showClearButton && localKeyword"
        @click="handleClear"
        :disabled="loading"
        class="px-3 py-2 sm:px-4 sm:py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors disabled:opacity-50 disabled:cursor-not-allowed text-sm sm:text-base whitespace-nowrap"
      >
        クリア
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: 'オークションタイトルで検索...'
  },
  searchButtonText: {
    type: String,
    default: '検索'
  },
  showClearButton: {
    type: Boolean,
    default: true
  },
  loading: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'search', 'clear'])

const localKeyword = ref(props.modelValue)

// 親コンポーネントからの変更を反映
watch(() => props.modelValue, (newValue) => {
  localKeyword.value = newValue
})

const handleSearch = () => {
  emit('update:modelValue', localKeyword.value)
  emit('search', localKeyword.value)
}

const handleClear = () => {
  localKeyword.value = ''
  emit('update:modelValue', '')
  emit('clear')
}
</script>

<style scoped>
.search-bar {
  width: 100%;
}
</style>
