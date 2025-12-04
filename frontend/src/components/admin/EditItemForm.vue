<script setup>
import { ref, computed } from 'vue'
import Label from '@/components/ui/Label.vue'
import Input from '@/components/ui/Input.vue'
import Button from '@/components/ui/Button.vue'

const props = defineProps({
  modelValue: {
    type: Object,
    required: true,
    default: () => ({
      id: '',
      name: '',
      description: '',
      lot_number: 0,
      starting_price: null,
      can_edit: true,
      can_delete: true,
      bid_count: 0,
    }),
  },
  index: {
    type: Number,
    required: true,
  },
  errors: {
    type: Object,
    default: () => ({ name: '', description: '' }),
  },
  canMoveUp: {
    type: Boolean,
    default: false,
  },
  canMoveDown: {
    type: Boolean,
    default: false,
  },
  canDelete: {
    type: Boolean,
    default: true,
  },
  canEdit: {
    type: Boolean,
    default: true,
  },
})

const emit = defineEmits(['update:modelValue', 'validate', 'move-up', 'move-down', 'delete', 'drag-start', 'drag-over', 'drop'])

const localValue = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value),
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
    starting_price: value === '' ? null : parseInt(value, 10),
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

// Status text
const statusText = computed(() => {
  if (props.modelValue.started_at && props.modelValue.ended_at) {
    return '終了'
  } else if (props.modelValue.started_at) {
    return '進行中'
  }
  return '未開始'
})

const statusClass = computed(() => {
  if (props.modelValue.started_at && props.modelValue.ended_at) {
    return 'bg-gray-100 text-gray-800'
  } else if (props.modelValue.started_at) {
    return 'bg-green-100 text-green-800'
  }
  return 'bg-blue-100 text-blue-800'
})

// Drag and drop state
const isDragging = ref(false)

// Drag and drop handlers
function handleDragStart(event) {
  if (!props.canEdit) return
  isDragging.value = true
  event.dataTransfer.effectAllowed = 'move'
  event.dataTransfer.setData('text/plain', props.index.toString())
  emit('drag-start', props.index)
}

function handleDragEnd() {
  isDragging.value = false
}

function handleDragOver(event) {
  if (!props.canEdit) return
  event.preventDefault()
  event.dataTransfer.dropEffect = 'move'
  emit('drag-over', props.index)
}

function handleDrop(event) {
  if (!props.canEdit) return
  event.preventDefault()
  const fromIndex = parseInt(event.dataTransfer.getData('text/plain'), 10)
  emit('drop', { fromIndex, toIndex: props.index })
}
</script>

<template>
  <div
    :draggable="canEdit"
    @dragstart="handleDragStart"
    @dragend="handleDragEnd"
    @dragover="handleDragOver"
    @drop="handleDrop"
    class="border rounded-lg p-4 transition-all duration-200"
    :class="[
      canEdit ? 'border-gray-200 hover:border-blue-300 hover:shadow-sm' : 'border-gray-300 bg-gray-50',
      isDragging ? 'opacity-50 scale-95 shadow-lg border-blue-400' : ''
    ]"
  >
    <div class="flex justify-between items-center mb-3">
      <div class="flex items-center gap-3">
        <!-- Drag Handle -->
        <div v-if="canEdit" class="flex-shrink-0 text-gray-400 cursor-move">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8h16M4 16h16" />
          </svg>
        </div>
        <h3 class="font-medium text-gray-700">商品 #{{ index + 1 }}</h3>
        <span
          :class="[
            'px-2 py-0.5 text-xs font-medium rounded-full',
            statusClass,
          ]"
        >
          {{ statusText }}
        </span>
        <span v-if="modelValue.bid_count > 0" class="text-sm text-gray-500">
          (入札 {{ modelValue.bid_count }} 件)
        </span>
      </div>
      <div class="flex gap-2">
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

    <!-- Read-only warning -->
    <div
      v-if="!canEdit"
      class="mb-3 p-2 bg-yellow-50 border border-yellow-200 rounded text-sm text-yellow-700"
    >
      この商品は開始済みまたは入札があるため編集できません
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
          :disabled="!canEdit"
        />
        <p v-if="errors.name" class="text-red-500 text-sm mt-1">
          {{ errors.name }}
        </p>
        <p class="text-gray-500 text-sm mt-1">{{ localValue.name?.length || 0 }}/200文字</p>
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
            { 'border-red-500': errors.description },
            { 'bg-gray-100 cursor-not-allowed': !canEdit },
          ]"
          :disabled="!canEdit"
        ></textarea>
        <p v-if="errors.description" class="text-red-500 text-sm mt-1">
          {{ errors.description }}
        </p>
        <p class="text-gray-500 text-sm mt-1">{{ localValue.description?.length || 0 }}/2000文字</p>
      </div>

      <!-- Starting Price (read-only for existing items) -->
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
          :disabled="!canEdit"
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
