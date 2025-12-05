<script setup>
import { DropdownMenuItem } from 'radix-vue'
import { cn } from '@/lib/utils'

const props = defineProps({
  disabled: {
    type: Boolean,
    default: false,
  },
  textValue: {
    type: String,
    default: undefined,
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

const emits = defineEmits(['select'])
</script>

<template>
  <DropdownMenuItem
    v-bind="props"
    :class="cn(
      'relative flex cursor-pointer select-none items-center rounded-sm px-2 py-1.5 text-sm outline-none transition-colors',
      'focus:bg-gray-100 focus:text-gray-900',
      'data-[disabled]:pointer-events-none data-[disabled]:opacity-50',
      props.class
    )"
    @select="emits('select', $event)"
  >
    <slot />
  </DropdownMenuItem>
</template>
