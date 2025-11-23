<script setup>
import { computed } from 'vue'
import Card from '@/components/ui/Card.vue'
import Button from '@/components/ui/Button.vue'
import ItemForm from './ItemForm.vue'

const props = defineProps({
  modelValue: {
    type: Array,
    required: true,
    default: () => []
  },
  errors: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['update:modelValue', 'validate', 'add-item'])

const localValue = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

function updateItem(index, newValue) {
  const newItems = [...localValue.value]
  newItems[index] = newValue
  localValue.value = newItems
}

function handleValidate(index, field) {
  emit('validate', index, field)
}

function handleAddItem() {
  emit('add-item')
}

function handleMoveUp(index) {
  if (index === 0) return

  const newItems = [...localValue.value]
  const temp = newItems[index]
  newItems[index] = newItems[index - 1]
  newItems[index - 1] = temp

  // Recalculate lot numbers
  newItems.forEach((item, idx) => {
    item.lot_number = idx
  })

  localValue.value = newItems

  // Swap errors
  const newErrors = [...props.errors]
  const tempError = newErrors[index]
  newErrors[index] = newErrors[index - 1]
  newErrors[index - 1] = tempError
}

function handleMoveDown(index) {
  if (index === localValue.value.length - 1) return

  const newItems = [...localValue.value]
  const temp = newItems[index]
  newItems[index] = newItems[index + 1]
  newItems[index + 1] = temp

  // Recalculate lot numbers
  newItems.forEach((item, idx) => {
    item.lot_number = idx
  })

  localValue.value = newItems

  // Swap errors
  const newErrors = [...props.errors]
  const tempError = newErrors[index]
  newErrors[index] = newErrors[index + 1]
  newErrors[index + 1] = tempError
}

function handleDelete(index) {
  if (localValue.value.length === 1) {
    alert('最低1つの商品が必要です')
    return
  }

  if (confirm('この商品を削除してもよろしいですか?')) {
    const newItems = localValue.value.filter((_, idx) => idx !== index)

    // Recalculate lot numbers
    newItems.forEach((item, idx) => {
      item.lot_number = idx
    })

    localValue.value = newItems
  }
}
</script>

<template>
  <Card class="p-6">
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-xl font-semibold">商品情報</h2>
      <Button type="button" @click="handleAddItem" variant="outline" size="sm">
        + 商品を追加
      </Button>
    </div>

    <div class="space-y-4">
      <ItemForm
        v-for="(item, index) in localValue"
        :key="index"
        :model-value="item"
        @update:model-value="updateItem(index, $event)"
        :index="index"
        :errors="errors[index] || { name: '', description: '' }"
        @validate="(field) => handleValidate(index, field)"
        :can-move-up="index > 0"
        :can-move-down="index < localValue.length - 1"
        :can-delete="localValue.length > 1"
        @move-up="handleMoveUp(index)"
        @move-down="handleMoveDown(index)"
        @delete="handleDelete(index)"
      />
    </div>
  </Card>
</template>
