<script setup>
import { ref, provide, computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  defaultValue: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const activeTab = computed({
  get: () => props.modelValue || props.defaultValue,
  set: (value) => emit('update:modelValue', value)
})

provide('activeTab', activeTab)

function setActiveTab(value) {
  activeTab.value = value
}

provide('setActiveTab', setActiveTab)
</script>

<template>
  <div class="w-full">
    <slot />
  </div>
</template>
