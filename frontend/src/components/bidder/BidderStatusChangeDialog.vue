<template>
  <div
    v-if="modelValue"
    class="fixed z-10 inset-0 overflow-y-auto"
    aria-labelledby="modal-title"
    role="dialog"
    aria-modal="true"
  >
    <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="handleClose"></div>

      <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

      <div
        class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full"
      >
        <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
          <div class="sm:flex sm:items-start">
            <div
              :class="[
                'mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full sm:mx-0 sm:h-10 sm:w-10',
                bidder?.status === 'active' ? 'bg-red-100' : 'bg-green-100',
              ]"
            >
              <span :class="[bidder?.status === 'active' ? 'text-red-600' : 'text-green-600', 'text-xl']">
                {{ bidder?.status === 'active' ? '⚠' : '✓' }}
              </span>
            </div>
            <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left">
              <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                アカウント{{ bidder?.status === 'active' ? '停止' : '復活' }}の確認
              </h3>
              <div class="mt-2">
                <p class="text-sm text-gray-500">
                  {{ bidder?.email }} のアカウントを{{ bidder?.status === 'active' ? '停止' : '復活'
                  }}してもよろしいですか？
                </p>
                <p v-if="bidder?.status === 'active'" class="mt-2 text-sm text-red-600">
                  停止後、このアカウントではログインできなくなります。
                </p>
                <p v-else class="mt-2 text-sm text-green-600">
                  復活後、このアカウントは再びログインできるようになります。
                </p>
              </div>
            </div>
          </div>
        </div>
        <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
          <button
            @click="handleConfirm"
            :disabled="loading"
            :class="[
              bidder?.status === 'active' ? 'bg-red-600 hover:bg-red-700' : 'bg-green-600 hover:bg-green-700',
              'w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 text-base font-medium text-white focus:outline-none focus:ring-2 focus:ring-offset-2 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed',
            ]"
          >
            <span v-if="loading" class="inline-block animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></span>
            {{ bidder?.status === 'active' ? '停止する' : '復活する' }}
          </button>
          <button
            @click="handleClose"
            :disabled="loading"
            class="mt-3 w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:mt-0 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
          >
            キャンセル
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
const props = defineProps({
  modelValue: {
    type: Boolean,
    required: true,
  },
  bidder: {
    type: Object,
    default: null,
  },
  loading: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:modelValue', 'confirm'])

function handleClose() {
  if (!props.loading) {
    emit('update:modelValue', false)
  }
}

function handleConfirm() {
  emit('confirm')
}
</script>
