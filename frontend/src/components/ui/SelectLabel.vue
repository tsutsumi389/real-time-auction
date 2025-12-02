<script setup>
import { computed } from 'vue'
import { SelectLabel } from 'radix-vue'
import { cn } from '@/lib/utils'

const props = defineProps({
  for: {
    type: String,
    default: undefined
  },
  asChild: {
    type: Boolean,
    default: undefined
  },
  class: {
    type: String,
    default: ''
  }
})
</script>

<template>
  <SelectLabel :class="cn('py-1.5 pl-8 pr-2 text-sm font-semibold', props.class)">
    <slot />
  </SelectLabel>
</template>
