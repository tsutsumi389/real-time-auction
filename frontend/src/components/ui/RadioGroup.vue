<template>
  <div :class="cn('grid gap-2', props.class)">
    <slot />
  </div>
</template>

<script setup>
import { cn } from '@/lib/utils'
import { provide } from 'vue'

const props = defineProps({
  class: {
    type: String,
    default: ''
  },
  modelValue: {
    type: String,
    default: ''
  },
  name: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['update:modelValue'])

// Provide context for RadioGroupItem children
provide('radioGroup', {
  name: props.name,
  modelValue: props.modelValue,
  updateValue: (value) => emit('update:modelValue', value)
})
</script>
