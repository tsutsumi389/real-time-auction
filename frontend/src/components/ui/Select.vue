<script setup>
import { SelectRoot, useForwardPropsEmits } from 'radix-vue'
import { computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: undefined
  },
  defaultValue: {
    type: String,
    default: undefined
  },
  defaultOpen: {
    type: Boolean,
    default: false
  },
  open: {
    type: Boolean,
    default: undefined
  },
  dir: {
    type: String,
    default: undefined
  },
  name: {
    type: String,
    default: undefined
  },
  autocomplete: {
    type: String,
    default: undefined
  },
  disabled: {
    type: Boolean,
    default: false
  },
  required: {
    type: Boolean,
    default: false
  }
})

const emits = defineEmits(['update:modelValue', 'update:open'])

const forwarded = useForwardPropsEmits(props, emits)
</script>

<template>
  <SelectRoot v-bind="forwarded">
    <slot />
  </SelectRoot>
</template>
