<template>
  <div :class="cn('flex items-center space-x-2', props.class)">
    <input
      :id="props.id"
      type="radio"
      :name="radioGroup.name"
      :value="props.value"
      :checked="radioGroup.modelValue === props.value"
      @change="handleChange"
      class="h-4 w-4 rounded-full border border-primary text-primary ring-offset-background focus:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
    />
    <label
      :for="props.id"
      class="text-sm font-medium leading-none peer-disabled:cursor-not-allowed peer-disabled:opacity-70 cursor-pointer"
    >
      <slot />
    </label>
  </div>
</template>

<script setup>
import { cn } from '@/lib/utils'
import { inject } from 'vue'

const props = defineProps({
  class: {
    type: String,
    default: ''
  },
  id: {
    type: String,
    required: true
  },
  value: {
    type: String,
    required: true
  }
})

const radioGroup = inject('radioGroup')

const handleChange = () => {
  radioGroup.updateValue(props.value)
}
</script>
