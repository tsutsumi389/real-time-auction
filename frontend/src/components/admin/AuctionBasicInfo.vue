<script setup>
import { computed } from 'vue'
import Card from '@/components/ui/Card.vue'
import Label from '@/components/ui/Label.vue'
import Input from '@/components/ui/Input.vue'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
    default: () => ({ title: '', description: '' })
  },
  errors: {
    type: Object,
    default: () => ({ title: '', description: '' })
  }
})

const emit = defineEmits(['update:modelValue', 'validate'])

const localValue = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

function updateTitle(event) {
  localValue.value = { ...localValue.value, title: event.target.value }
}

function updateDescription(event) {
  localValue.value = { ...localValue.value, description: event.target.value }
}

function handleTitleBlur() {
  emit('validate', 'title')
}

function handleDescriptionBlur() {
  emit('validate', 'description')
}
</script>

<template>
  <Card class="p-6">
    <h2 class="text-xl font-semibold mb-4">基本情報</h2>

    <div class="space-y-4">
      <!-- Title -->
      <div>
        <Label for="title" class="required">タイトル</Label>
        <Input
          id="title"
          :value="localValue.title"
          @input="updateTitle"
          type="text"
          placeholder="例: 2025年春季競走馬セリ"
          maxlength="200"
          @blur="handleTitleBlur"
          :class="{ 'border-red-500': errors.title }"
        />
        <p v-if="errors.title" class="text-red-500 text-sm mt-1">{{ errors.title }}</p>
        <p class="text-gray-500 text-sm mt-1">{{ localValue.title.length }}/200文字</p>
      </div>

      <!-- Description -->
      <div>
        <Label for="description">説明</Label>
        <textarea
          id="description"
          :value="localValue.description"
          @input="updateDescription"
          placeholder="オークションの概要を入力してください"
          rows="4"
          maxlength="2000"
          @blur="handleDescriptionBlur"
          :class="[
            'w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500',
            { 'border-red-500': errors.description }
          ]"
        ></textarea>
        <p v-if="errors.description" class="text-red-500 text-sm mt-1">{{ errors.description }}</p>
        <p class="text-gray-500 text-sm mt-1">{{ localValue.description.length }}/2000文字</p>
      </div>
    </div>
  </Card>
</template>

<style scoped>
.required::after {
  content: ' *';
  color: #ef4444;
}
</style>
