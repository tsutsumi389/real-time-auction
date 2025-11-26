/**
 * Toast Notification Composable
 * トースト通知を管理するコンポーザブル
 */
import { ref } from 'vue'

const toasts = ref([])
let nextId = 0

export function useToast() {
  function addToast({ variant = 'default', title = '', description = '', duration = 5000 }) {
    const id = nextId++
    const toast = {
      id,
      variant,
      title,
      description,
    }

    toasts.value.push(toast)

    if (duration > 0) {
      setTimeout(() => {
        removeToast(id)
      }, duration)
    }

    return id
  }

  function removeToast(id) {
    const index = toasts.value.findIndex(toast => toast.id === id)
    if (index !== -1) {
      toasts.value.splice(index, 1)
    }
  }

  function success(title, description = '', duration = 5000) {
    return addToast({ variant: 'success', title, description, duration })
  }

  function error(title, description = '', duration = 5000) {
    return addToast({ variant: 'error', title, description, duration })
  }

  function warning(title, description = '', duration = 5000) {
    return addToast({ variant: 'warning', title, description, duration })
  }

  function info(title, description = '', duration = 5000) {
    return addToast({ variant: 'info', title, description, duration })
  }

  function clear() {
    toasts.value = []
  }

  return {
    toasts,
    addToast,
    removeToast,
    success,
    error,
    warning,
    info,
    clear,
  }
}
