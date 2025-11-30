<script setup>
import { computed } from 'vue'
import Card from '@/components/ui/Card.vue'
import Button from '@/components/ui/Button.vue'
import EditItemForm from './EditItemForm.vue'

const props = defineProps({
  modelValue: {
    type: Array,
    required: true,
    default: () => []
  },
  errors: {
    type: Array,
    default: () => []
  },
  canEditAuction: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:modelValue', 'validate', 'add-item', 'delete-item', 'update-item'])

const localValue = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

function updateItem(index, newValue) {
  const newItems = [...localValue.value]
  newItems[index] = newValue
  localValue.value = newItems
  emit('update-item', index, newValue)
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
    item.lot_number = idx + 1
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
    item.lot_number = idx + 1
  })

  localValue.value = newItems

  // Swap errors
  const newErrors = [...props.errors]
  const tempError = newErrors[index]
  newErrors[index] = newErrors[index + 1]
  newErrors[index + 1] = tempError
}

function handleDelete(index) {
  emit('delete-item', index)
}
</script>

<template>
  <Card class="p-6">
    <div class="flex justify-between items-center mb-4">
      <h2 class="text-xl font-semibold">商品情報</h2>
      <Button
        v-if="canEditAuction"
        type="button"
        @click="handleAddItem"
        variant="outline"
        size="sm"
      >
        + 商品を追加
      </Button>
    </div>

    <div v-if="localValue.length === 0" class="text-center py-8 text-gray-500">
      商品がありません。「+ 商品を追加」をクリックして商品を追加してください。
    </div>

    <div v-else class="space-y-4">
      <EditItemForm
        v-for="(item, index) in localValue"
        :key="item.id || index"
        :model-value="item"
        @update:model-value="updateItem(index, $event)"
        :index="index"
        :errors="errors[index] || { name: '', description: '' }"
        @validate="(field) => handleValidate(index, field)"
        :can-move-up="index > 0 && canEditAuction"
        :can-move-down="index < localValue.length - 1 && canEditAuction"
        :can-delete="item.can_delete && localValue.length > 1 && canEditAuction"
        :can-edit="item.can_edit && canEditAuction"
        @move-up="handleMoveUp(index)"
        @move-down="handleMoveDown(index)"
        @delete="handleDelete(index)"
      />
    </div>
  </Card>
</template>
