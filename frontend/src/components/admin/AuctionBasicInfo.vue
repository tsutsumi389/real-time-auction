<script setup>
import { computed } from 'vue'
import Card from '@/components/ui/Card.vue'
import Label from '@/components/ui/Label.vue'
import Input from '@/components/ui/Input.vue'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
    default: () => ({ title: '', description: '', started_at: '' })
  },
  errors: {
    type: Object,
    default: () => ({ title: '', description: '', started_at: '' })
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

function updateStartedAt(event) {
  localValue.value = { ...localValue.value, started_at: event.target.value }
}

function handleTitleBlur() {
  emit('validate', 'title')
}

function handleDescriptionBlur() {
  emit('validate', 'description')
}

function handleStartedAtBlur() {
  emit('validate', 'started_at')
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

      <!-- Started At -->
      <div>
        <Label for="started_at">開始日時</Label>
        <Input
          id="started_at"
          :value="localValue.started_at"
          @input="updateStartedAt"
          type="datetime-local"
          @blur="handleStartedAtBlur"
          :class="{ 'border-red-500': errors.started_at }"
        />
        <p v-if="errors.started_at" class="text-red-500 text-sm mt-1">{{ errors.started_at }}</p>
        <p class="text-gray-500 text-sm mt-1">任意、入札者用一覧画面での表示・ソートに使用されます</p>
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
