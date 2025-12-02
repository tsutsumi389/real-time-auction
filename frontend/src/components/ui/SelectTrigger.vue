<script setup>
import { computed } from 'vue'
import { SelectTrigger, useForwardProps } from 'radix-vue'
import { ChevronDown } from 'lucide-vue-next'
import { cn } from '@/lib/utils'

const props = defineProps({
  disabled: {
    type: Boolean,
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

const delegatedProps = computed(() => {
  const { class: _, ...rest } = props
  return rest
})

const forwardedProps = useForwardProps(delegatedProps)
</script>

<template>
  <SelectTrigger
    v-bind="forwardedProps"
    :class="cn(
      'flex h-10 w-full items-center justify-between rounded-md border border-input bg-background px-3 py-2 text-sm ring-offset-background placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-ring focus:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50 [&>span]:line-clamp-1',
      props.class
    )"
  >
    <slot />
    <ChevronDown class="h-4 w-4 opacity-50 shrink-0" />
  </SelectTrigger>
</template>
