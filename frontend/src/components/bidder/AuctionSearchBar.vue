<template>
  <div class="search-bar" role="search" aria-label="Search auctions">
    <div class="flex gap-3">
      <!-- Search Input -->
      <div class="relative flex-1">
        <div class="absolute left-4 top-1/2 -translate-y-1/2 pointer-events-none">
          <Search class="h-5 w-5 text-lux-silver/50" :stroke-width="1.5" />
        </div>
        <input
          v-model="localKeyword"
          type="text"
          :placeholder="placeholder"
          :disabled="loading"
          class="search-input w-full h-12 pl-12 pr-12 rounded-xl lux-input text-base placeholder:text-lux-silver/40 focus:border-lux-gold/50 focus:ring-2 focus:ring-lux-gold/10 disabled:opacity-50 disabled:cursor-not-allowed"
          aria-label="Search keywords"
          @keyup.enter="handleSearch"
        />
        <Transition
          enter-active-class="transition-all duration-200 ease-out"
          enter-from-class="opacity-0 scale-90"
          enter-to-class="opacity-100 scale-100"
          leave-active-class="transition-all duration-150 ease-in"
          leave-from-class="opacity-100 scale-100"
          leave-to-class="opacity-0 scale-90"
        >
          <button
            v-if="localKeyword"
            @click="handleClear"
            :disabled="loading"
            class="absolute right-4 top-1/2 -translate-y-1/2 p-1.5 rounded-lg text-lux-silver/60 hover:text-lux-cream hover:bg-lux-noir-soft/50 transition-all duration-200 disabled:opacity-50"
            aria-label="Clear search"
            type="button"
          >
            <X class="h-4 w-4" :stroke-width="1.5" />
          </button>
        </Transition>
      </div>

      <!-- Search Button -->
      <button
        @click="handleSearch"
        :disabled="loading"
        class="search-button h-12 px-6 rounded-xl lux-btn-gold flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
        aria-label="Execute search"
      >
        <Loader2 v-if="loading" class="h-5 w-5 animate-spin" />
        <Search v-else class="h-5 w-5" :stroke-width="1.5" />
        <span class="hidden sm:inline text-sm font-semibold tracking-wider">{{ searchButtonText }}</span>
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { Search, X, Loader2 } from 'lucide-vue-next'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: 'Search auctions by title...'
  },
  searchButtonText: {
    type: String,
    default: 'SEARCH'
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

/* Luxury Input Styling */
.search-input {
  background: hsl(0 0% 8%);
  border: 1px solid hsl(0 0% 16%);
  color: hsl(45 30% 96%);
  font-family: 'DM Sans', system-ui, sans-serif;
  transition: all 0.3s ease;
}

.search-input::placeholder {
  color: hsl(220 10% 70% / 0.4);
}

.search-input:focus {
  outline: none;
  border-color: hsl(43 74% 49% / 0.5);
  box-shadow: 0 0 0 3px hsl(43 74% 49% / 0.1);
}

/* Luxury color utilities */
.text-lux-silver\/50 {
  color: hsl(220 10% 70% / 0.5);
}

.text-lux-silver\/60 {
  color: hsl(220 10% 70% / 0.6);
}

.text-lux-cream {
  color: hsl(45 30% 96%);
}

.hover\:bg-lux-noir-soft\/50:hover {
  background-color: hsl(0 0% 16% / 0.5);
}
</style>
