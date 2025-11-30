<template>
  <div
    :class="cn(
      'relative w-full rounded-lg border p-4',
      alertVariants({ variant }),
      props.class
    )"
    v-bind="$attrs"
  >
    <slot />
  </div>
</template>

<script setup>
import { cva } from 'class-variance-authority'
import { cn } from '@/lib/utils'

const props = defineProps({
  variant: {
    type: String,
    default: 'default',
    validator: (value) => ['default', 'destructive', 'success'].includes(value)
  },
  class: {
    type: String,
    default: ''
  }
})

const alertVariants = cva(
  'relative w-full rounded-lg border p-4',
  {
    variants: {
      variant: {
        default: 'bg-background text-foreground',
        destructive: 'border-destructive/50 text-destructive dark:border-destructive [&>svg]:text-destructive',
        success: 'border-green-200 bg-green-50 text-green-800 [&>svg]:text-green-500',
      },
    },
    defaultVariants: {
      variant: 'default',
    },
  }
)
</script>
