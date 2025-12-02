<template>
  <div
    :class="cn(
      'text-sm [&_p]:leading-relaxed',
      props.class
    )"
  >
    <slot />
  </div>
</template>

<script setup>
import { cn } from '@/lib/utils'

const props = defineProps({
  class: {
    type: String,
    default: ''
  }
})
</script>
