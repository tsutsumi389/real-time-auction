<script setup>
import { DropdownMenuTrigger } from 'radix-vue'
import { cn } from '@/lib/utils'

const props = defineProps({
  disabled: {
    type: Boolean,
    default: false,
  },
  asChild: {
    type: Boolean,
    default: false,
  },
  class: {
    type: String,
    default: '',
  },
})
</script>

<template>
  <DropdownMenuTrigger
    :disabled="disabled"
    :as-child="asChild"
    :class="cn('', props.class)"
  >
    <slot />
  </DropdownMenuTrigger>
</template>
