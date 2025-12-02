<template>
  <div
    :class="cn(skeletonVariants({ variant, size, rounded }), props.class)"
    :style="customStyle"
    role="status"
    aria-busy="true"
    aria-label="読み込み中"
  >
    <span class="sr-only">読み込み中</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { cva } from 'class-variance-authority'
import { cn } from '@/lib/utils'

const props = defineProps({
  variant: {
    type: String,
    default: 'default',
    validator: (value) => ['default', 'circular', 'rectangular'].includes(value)
  },
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg', 'full'].includes(value)
  },
  rounded: {
    type: Boolean,
    default: false
  },
  width: {
    type: String,
    default: null
  },
  height: {
    type: String,
    default: null
  },
  class: {
    type: String,
    default: ''
  }
})

const skeletonVariants = cva(
  'animate-pulse bg-muted',
  {
    variants: {
      variant: {
        default: '',
        circular: 'rounded-full',
        rectangular: 'rounded-none',
      },
      size: {
        sm: 'h-4',
        md: 'h-6',
        lg: 'h-10',
        full: 'h-full w-full',
      },
      rounded: {
        true: 'rounded-lg',
        false: '',
      },
    },
    compoundVariants: [
      {
        variant: 'default',
        rounded: false,
        class: 'rounded-md',
      },
    ],
    defaultVariants: {
      variant: 'default',
      size: 'md',
      rounded: false,
    },
  }
)

const customStyle = computed(() => {
  const style = {}
  if (props.width) {
    style.width = props.width
  }
  if (props.height) {
    style.height = props.height
  }
  return Object.keys(style).length > 0 ? style : undefined
})
</script>
