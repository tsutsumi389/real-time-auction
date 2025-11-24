<script setup>
import { computed } from 'vue'
import Label from '@/components/ui/Label.vue'
import Input from '@/components/ui/Input.vue'
import Button from '@/components/ui/Button.vue'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
    default: () => ({ name: '', description: '', lot_number: 0, starting_price: null })
  },
  index: {
    type: Number,
    required: true
  },
  errors: {
    type: Object,
    default: () => ({ name: '', description: '', starting_price: '' })
  },
  canMoveUp: {
    type: Boolean,
    default: false
  },
  canMoveDown: {
    type: Boolean,
    default: false
  },
  canDelete: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:modelValue', 'validate', 'move-up', 'move-down', 'delete'])

const localValue = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

function updateName(event) {
  localValue.value = { ...localValue.value, name: event.target.value }
}

function updateDescription(event) {
  localValue.value = { ...localValue.value, description: event.target.value }
}

function updateStartingPrice(event) {
  const value = event.target.value
  localValue.value = { 
    ...localValue.value, 
    starting_price: value === '' ? null : parseInt(value, 10)
  }
}

function handleNameBlur() {
  emit('validate', 'name')
}

function handleDescriptionBlur() {
  emit('validate', 'description')
}

function handleStartingPriceBlur() {
  emit('validate', 'starting_price')
}

function handleMoveUp() {
  emit('move-up')
}

function handleMoveDown() {
  emit('move-down')
}

function handleDelete() {
  emit('delete')
}
</script>

<template>
  <div class="border border-gray-200 rounded-lg p-4">
    <div class="flex justify-between items-center mb-3">
      <h3 class="font-medium text-gray-700">商品 #{{ index + 1 }}</h3>
      <div class="flex gap-2">
        <Button
          type="button"
          @click="handleMoveUp"
          :disabled="!canMoveUp"
          variant="outline"
          size="sm"
        >
          ▲
        </Button>
        <Button
          type="button"
          @click="handleMoveDown"
          :disabled="!canMoveDown"
          variant="outline"
          size="sm"
        >
          ▼
        </Button>
        <Button
          type="button"
          @click="handleDelete"
          :disabled="!canDelete"
          variant="destructive"
          size="sm"
        >
          削除
        </Button>
      </div>
    </div>

    <div class="space-y-3">
      <!-- Item Name -->
      <div>
        <Label :for="`item-name-${index}`" class="required">商品名</Label>
        <Input
          :id="`item-name-${index}`"
          :value="localValue.name"
          @input="updateName"
          type="text"
          placeholder="例: ディープインパクト産駒"
          maxlength="200"
          @blur="handleNameBlur"
          :class="{ 'border-red-500': errors.name }"
        />
        <p v-if="errors.name" class="text-red-500 text-sm mt-1">
          {{ errors.name }}
        </p>
        <p class="text-gray-500 text-sm mt-1">{{ localValue.name.length }}/200文字</p>
      </div>

      <!-- Item Description -->
      <div>
        <Label :for="`item-description-${index}`">商品説明</Label>
        <textarea
          :id="`item-description-${index}`"
          :value="localValue.description"
          @input="updateDescription"
          placeholder="商品の詳細情報を入力してください"
          rows="3"
          maxlength="2000"
          @blur="handleDescriptionBlur"
          :class="[
            'w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500',
            { 'border-red-500': errors.description }
          ]"
        ></textarea>
        <p v-if="errors.description" class="text-red-500 text-sm mt-1">
          {{ errors.description }}
        </p>
        <p class="text-gray-500 text-sm mt-1">{{ localValue.description.length }}/2000文字</p>
      </div>

      <!-- Starting Price -->
      <div>
        <Label :for="`item-starting-price-${index}`">開始価格</Label>
        <Input
          :id="`item-starting-price-${index}`"
          :value="localValue.starting_price || ''"
          @input="updateStartingPrice"
          type="number"
          placeholder="例: 5000000"
          min="1"
          @blur="handleStartingPriceBlur"
          :class="{ 'border-red-500': errors.starting_price }"
        />
        <p v-if="errors.starting_price" class="text-red-500 text-sm mt-1">
          {{ errors.starting_price }}
        </p>
        <p class="text-gray-500 text-sm mt-1">任意、単位: ポイント</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.required::after {
  content: ' *';
  color: #ef4444;
}
</style>
