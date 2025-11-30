<script setup>
import { ref, watch } from 'vue'
import Button from '@/components/ui/Button.vue'

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: '削除の確認',
  },
  message: {
    type: String,
    default: 'この操作は取り消せません。本当に削除しますか？',
  },
  confirmText: {
    type: String,
    default: '削除',
  },
  cancelText: {
    type: String,
    default: 'キャンセル',
  },
  isDeleting: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['confirm', 'cancel'])

const isVisible = ref(props.isOpen)

watch(
  () => props.isOpen,
  (newVal) => {
    isVisible.value = newVal
  }
)

function handleConfirm() {
  emit('confirm')
}

function handleCancel() {
  emit('cancel')
}

function handleBackdropClick(event) {
  if (event.target === event.currentTarget) {
    handleCancel()
  }
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="isVisible"
        class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
        @click="handleBackdropClick"
      >
        <div class="bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-6">
          <!-- Header -->
          <div class="flex items-center gap-3 mb-4">
            <div class="flex-shrink-0 w-10 h-10 bg-red-100 rounded-full flex items-center justify-center">
              <svg
                class="w-6 h-6 text-red-600"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                />
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-gray-900">{{ title }}</h3>
          </div>

          <!-- Content -->
          <p class="text-gray-600 mb-6">{{ message }}</p>

          <!-- Actions -->
          <div class="flex justify-end gap-3">
            <Button
              type="button"
              variant="outline"
              @click="handleCancel"
              :disabled="isDeleting"
            >
              {{ cancelText }}
            </Button>
            <Button
              type="button"
              variant="destructive"
              @click="handleConfirm"
              :disabled="isDeleting"
            >
              {{ isDeleting ? '削除中...' : confirmText }}
            </Button>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.2s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-enter-active .bg-white,
.modal-leave-active .bg-white {
  transition: transform 0.2s ease;
}

.modal-enter-from .bg-white,
.modal-leave-to .bg-white {
  transform: scale(0.95);
}
</style>
