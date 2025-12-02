<template>
  <div
    :class="cn(containerClasses, props.class)"
    role="status"
    aria-busy="true"
    :aria-label="text || '読み込み中'"
  >
    <div :class="spinnerClasses"></div>
    <span v-if="text" :class="textClasses">{{ text }}</span>
    <span class="sr-only">{{ text || '読み込み中' }}</span>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { cva } from 'class-variance-authority'
import { cn } from '@/lib/utils'

const props = defineProps({
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['xs', 'sm', 'md', 'lg', 'xl', '2xl'].includes(value)
  },
  color: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'secondary', 'white', 'muted'].includes(value)
  },
  text: {
    type: String,
    default: ''
  },
  center: {
    type: Boolean,
    default: false
  },
  class: {
    type: String,
    default: ''
  }
})

const spinnerVariants = cva(
  'inline-block animate-spin rounded-full border-2 border-solid border-current border-r-transparent',
  {
    variants: {
      size: {
        xs: 'h-3 w-3',
        sm: 'h-4 w-4',
        md: 'h-6 w-6',
        lg: 'h-12 w-12',
        xl: 'h-16 w-16',
        '2xl': 'h-20 w-20'
      },
      color: {
        primary: 'text-primary',
        secondary: 'text-secondary',
        white: 'text-white',
        muted: 'text-muted-foreground'
      }
    },
    defaultVariants: {
      size: 'md',
      color: 'primary'
    }
  }
)

const textVariants = cva(
  'text-muted-foreground',
  {
    variants: {
      size: {
        xs: 'text-xs',
        sm: 'text-sm',
        md: 'text-base',
        lg: 'text-lg',
        xl: 'text-xl',
        '2xl': 'text-2xl'
      }
    },
    defaultVariants: {
      size: 'md'
    }
  }
)

const containerClasses = computed(() => {
  const base = 'inline-flex items-center gap-2'
  const center = props.center ? 'justify-center w-full' : ''
  return `${base} ${center}`.trim()
})

const spinnerClasses = computed(() => {
  return spinnerVariants({ size: props.size, color: props.color })
})

const textClasses = computed(() => {
  return textVariants({ size: props.size })
})
</script>
