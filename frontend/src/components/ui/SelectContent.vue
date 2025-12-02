<script setup>
import { computed } from 'vue'
import {
  SelectContent,
  SelectPortal,
  SelectViewport,
  useForwardPropsEmits
} from 'radix-vue'
import { cn } from '@/lib/utils'
import SelectScrollUpButton from './SelectScrollUpButton.vue'
import SelectScrollDownButton from './SelectScrollDownButton.vue'

const props = defineProps({
  forceMount: {
    type: Boolean,
    default: undefined
  },
  position: {
    type: String,
    default: 'popper',
    validator: (value) => ['item-aligned', 'popper'].includes(value)
  },
  bodyLock: {
    type: Boolean,
    default: undefined
  },
  side: {
    type: String,
    default: 'bottom',
    validator: (value) => ['top', 'right', 'bottom', 'left'].includes(value)
  },
  sideOffset: {
    type: Number,
    default: 0
  },
  align: {
    type: String,
    default: 'start',
    validator: (value) => ['start', 'center', 'end'].includes(value)
  },
  alignOffset: {
    type: Number,
    default: 0
  },
  avoidCollisions: {
    type: Boolean,
    default: true
  },
  collisionBoundary: {
    type: [Object, Array],
    default: () => []
  },
  collisionPadding: {
    type: [Number, Object],
    default: 0
  },
  arrowPadding: {
    type: Number,
    default: 0
  },
  sticky: {
    type: String,
    default: 'partial',
    validator: (value) => ['partial', 'always'].includes(value)
  },
  hideWhenDetached: {
    type: Boolean,
    default: false
  },
  updatePositionStrategy: {
    type: String,
    default: 'optimized',
    validator: (value) => ['optimized', 'always'].includes(value)
  },
  prioritizePosition: {
    type: Boolean,
    default: false
  },
  class: {
    type: String,
    default: ''
  }
})

const emits = defineEmits(['closeAutoFocus', 'escapeKeyDown', 'pointerDownOutside'])

const delegatedProps = computed(() => {
  const { class: _, ...rest } = props
  return rest
})

const forwarded = useForwardPropsEmits(delegatedProps, emits)
</script>

<template>
  <SelectPortal>
    <SelectContent
      v-bind="{ ...forwarded, ...$attrs }"
      :class="cn(
        'relative z-50 max-h-96 min-w-32 overflow-hidden rounded-md border bg-popover text-popover-foreground shadow-md data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2',
        position === 'popper' && 'data-[side=bottom]:translate-y-1 data-[side=left]:-translate-x-1 data-[side=right]:translate-x-1 data-[side=top]:-translate-y-1',
        props.class
      )"
    >
      <SelectScrollUpButton />
      <SelectViewport
        :class="cn(
          'p-1',
          position === 'popper' && 'h-[var(--radix-select-trigger-height)] w-full min-w-[var(--radix-select-trigger-width)]'
        )"
      >
        <slot />
      </SelectViewport>
      <SelectScrollDownButton />
    </SelectContent>
  </SelectPortal>
</template>
