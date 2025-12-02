<script setup>
import { computed } from 'vue'
import { SelectScrollUpButton } from 'radix-vue'
import { ChevronUp } from 'lucide-vue-next'
import { cn } from '@/lib/utils'

const props = defineProps({
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
  <SelectScrollUpButton
    :class="cn(
      'flex cursor-default items-center justify-center py-1',
      props.class
    )"
  >
    <slot>
      <ChevronUp class="h-4 w-4" />
    </slot>
  </SelectScrollUpButton>
</template>
