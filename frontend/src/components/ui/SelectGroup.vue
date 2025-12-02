<script setup>
import { computed } from 'vue'
import { SelectGroup } from 'radix-vue'
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
  <SelectGroup :class="cn('p-1 w-full', props.class)">
    <slot />
  </SelectGroup>
</template>
