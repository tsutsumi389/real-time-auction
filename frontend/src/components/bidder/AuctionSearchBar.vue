<template>
  <div class="search-bar" role="search" aria-label="オークション検索">
    <div class="flex gap-2 sm:gap-4">
      <div class="relative flex-1">
        <Search class="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" :stroke-width="1.5" />
        <Input
          v-model="localKeyword"
          type="text"
          :placeholder="placeholder"
          :disabled="loading"
          class="pl-10 pr-10"
          aria-label="検索キーワード"
          @keyup.enter="handleSearch"
        />
        <button
          v-if="localKeyword"
          @click="handleClear"
          :disabled="loading"
          class="absolute right-3 top-1/2 -translate-y-1/2 p-0.5 rounded-full text-muted-foreground hover:text-foreground hover:bg-muted transition-colors disabled:opacity-50"
          aria-label="検索キーワードをクリア"
          type="button"
        >
          <X class="h-4 w-4" :stroke-width="1.5" />
        </button>
      </div>
      <Button
        @click="handleSearch"
        :disabled="loading"
        variant="default"
        aria-label="検索を実行"
      >
        <Loader2 v-if="loading" class="mr-2 h-4 w-4 animate-spin" />
        <Search v-else class="mr-2 h-4 w-4" :stroke-width="1.5" />
        {{ searchButtonText }}
      </Button>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { Search, X, Loader2 } from 'lucide-vue-next'
import Input from '@/components/ui/Input.vue'
import Button from '@/components/ui/Button.vue'

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

// Reflect changes from parent component
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
