<template>
  <div :class="containerClasses">
    <div :class="spinnerClasses"></div>
    <span v-if="text" :class="textClasses">{{ text }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg', 'xl'].includes(value)
  },
  text: {
    type: String,
    default: ''
  },
  center: {
    type: Boolean,
    default: false
  }
})

const containerClasses = computed(() => {
  const base = 'inline-flex items-center gap-2'
  const center = props.center ? 'justify-center w-full' : ''
  return `${base} ${center}`.trim()
})

const spinnerClasses = computed(() => {
  const base = 'animate-spin rounded-full border-b-2 border-blue-600'

  const sizes = {
    sm: 'h-4 w-4',
    md: 'h-6 w-6',
    lg: 'h-12 w-12',
    xl: 'h-16 w-16'
  }

  return `${base} ${sizes[props.size] || sizes.md}`
})

const textClasses = computed(() => {
  const sizes = {
    sm: 'text-sm',
    md: 'text-base',
    lg: 'text-lg',
    xl: 'text-xl'
  }

  return `text-gray-600 ${sizes[props.size] || sizes.md}`
})
</script>
