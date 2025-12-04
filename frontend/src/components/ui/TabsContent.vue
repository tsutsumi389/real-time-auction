<script setup>
import { inject, computed } from 'vue'

const props = defineProps({
  value: {
    type: String,
    required: true
  }
})

const activeTab = inject('activeTab')

const isActive = computed(() => activeTab.value === props.value)
</script>

<template>
  <div
    v-if="isActive"
    class="mt-2 ring-offset-white focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-blue-500 focus-visible:ring-offset-2"
  >
    <slot />
  </div>
</template>
