<script setup>
import { computed } from 'vue'
import { SelectScrollDownButton } from 'radix-vue'
import { ChevronDown } from 'lucide-vue-next'
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
  <SelectScrollDownButton
    :class="cn(
      'flex cursor-default items-center justify-center py-1',
      props.class
    )"
  >
    <slot>
      <ChevronDown class="h-4 w-4" />
    </slot>
  </SelectScrollDownButton>
</template>
