<script setup>
import { inject, computed } from 'vue'

const props = defineProps({
  value: {
    type: String,
    required: true
  }
})

const activeTab = inject('activeTab')
const setActiveTab = inject('setActiveTab')

const isActive = computed(() => activeTab.value === props.value)

function handleClick() {
  setActiveTab(props.value)
}
</script>

<template>
  <button
    type="button"
    @click="handleClick"
    :class="[
      'inline-flex items-center justify-center whitespace-nowrap rounded-sm px-3 py-1.5 text-sm font-medium ring-offset-white transition-all focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50',
      isActive
        ? 'bg-white text-gray-900 shadow-sm'
        : 'text-gray-600 hover:text-gray-900'
    ]"
  >
    <slot />
  </button>
</template>
