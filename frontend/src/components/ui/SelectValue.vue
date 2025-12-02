<script setup>
import { SelectValue } from 'radix-vue'

defineProps({
  placeholder: {
    type: String,
    default: ''
  },
  asChild: {
    type: Boolean,
    default: undefined
  }
})
</script>

<template>
  <SelectValue :placeholder="placeholder" :as-child="asChild">
    <slot />
  </SelectValue>
</template>
