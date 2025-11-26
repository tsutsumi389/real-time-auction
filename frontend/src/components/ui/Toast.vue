<script setup>
import { computed } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'default',
    validator: (value) => ['default', 'success', 'error', 'warning', 'info'].includes(value),
  },
  title: {
    type: String,
    default: '',
  },
  description: {
    type: String,
    default: '',
  },
})

const emit = defineEmits(['close'])

const variantClasses = computed(() => {
  const variants = {
    default: 'bg-white border-gray-200',
    success: 'bg-green-50 border-green-200',
    error: 'bg-red-50 border-red-200',
    warning: 'bg-yellow-50 border-yellow-200',
    info: 'bg-blue-50 border-blue-200',
  }
  return variants[props.variant]
})

const iconClasses = computed(() => {
  const icons = {
    default: 'text-gray-500',
    success: 'text-green-500',
    error: 'text-red-500',
    warning: 'text-yellow-500',
    info: 'text-blue-500',
  }
  return icons[props.variant]
})

const icon = computed(() => {
  const icons = {
    default: '●',
    success: '✓',
    error: '✕',
    warning: '⚠',
    info: 'ℹ',
  }
  return icons[props.variant]
})
</script>

<template>
  <div
    :class="[
      'relative flex w-full max-w-md items-start gap-3 rounded-lg border p-4 shadow-lg transition-all',
      variantClasses
    ]"
  >
    <div :class="['text-xl', iconClasses]">{{ icon }}</div>
    <div class="flex-1">
      <div v-if="title" class="font-semibold text-sm mb-1">{{ title }}</div>
      <div v-if="description" class="text-sm text-gray-600">{{ description }}</div>
      <slot />
    </div>
    <button
      @click="emit('close')"
      class="ml-auto text-gray-400 hover:text-gray-600 transition-colors"
    >
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="16"
        height="16"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <line x1="18" y1="6" x2="6" y2="18"></line>
        <line x1="6" y1="6" x2="18" y2="18"></line>
      </svg>
    </button>
  </div>
</template>
