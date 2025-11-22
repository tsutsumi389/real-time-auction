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
              class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-green-100 sm:mx-0 sm:h-10 sm:w-10"
            >
              <span class="text-green-600 text-xl">+</span>
            </div>
            <div class="mt-3 text-center sm:mt-0 sm:ml-4 sm:text-left flex-1">
              <h3 class="text-lg leading-6 font-medium text-gray-900" id="modal-title">
                ポイント付与
              </h3>
              <div class="mt-4 space-y-4">
                <!-- 入札者情報 -->
                <div class="bg-gray-50 p-3 rounded">
                  <p class="text-sm text-gray-700">
                    <span class="font-medium">メール:</span> {{ bidder?.email }}
                  </p>
                  <p class="text-sm text-gray-700 mt-1">
                    <span class="font-medium">表示名:</span> {{ bidder?.display_name || '（未設定）' }}
                  </p>
                  <p class="text-sm text-gray-700 mt-1">
                    <span class="font-medium">現在のポイント:</span> {{ formatPoints(bidder?.total_points || 0) }}
                  </p>
                </div>

                <!-- 付与ポイント入力 -->
                <div>
                  <label for="points" class="block text-sm font-medium text-gray-700 mb-1">
                    付与するポイント
                  </label>
                  <input
                    id="points"
                    v-model.number="points"
                    type="number"
                    min="1"
                    max="1000000"
                    placeholder="1〜1,000,000"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-green-500"
                    :disabled="loading"
                  />
                  <p v-if="pointsError" class="mt-1 text-sm text-red-600">{{ pointsError }}</p>
                </div>

                <!-- 付与後のポイント残高プレビュー -->
                <div v-if="points > 0 && !pointsError" class="bg-blue-50 p-3 rounded">
                  <p class="text-sm text-blue-700">
                    <span class="font-medium">付与後のポイント:</span>
                    {{ formatPoints((bidder?.total_points || 0) + points) }}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
          <button
            @click="handleConfirm"
            :disabled="loading || !isValid"
            class="w-full inline-flex justify-center rounded-md border border-transparent shadow-sm px-4 py-2 bg-green-600 text-base font-medium text-white hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 sm:ml-3 sm:w-auto sm:text-sm disabled:opacity-50 disabled:cursor-not-allowed"
          >
            <span v-if="loading" class="inline-block animate-spin rounded-full h-4 w-4 border-b-2 border-white mr-2"></span>
            付与する
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
import { ref, computed, watch } from 'vue'

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

const points = ref(0)

const pointsError = computed(() => {
  if (points.value <= 0) {
    return 'ポイントは1以上の整数を入力してください'
  }
  if (points.value > 1000000) {
    return 'ポイントは1,000,000以下で入力してください'
  }
  if (!Number.isInteger(points.value)) {
    return 'ポイントは整数で入力してください'
  }
  return null
})

const isValid = computed(() => {
  return points.value > 0 && points.value <= 1000000 && Number.isInteger(points.value)
})

// モーダルを開いたときにポイントをリセット
watch(() => props.modelValue, (newValue) => {
  if (newValue) {
    points.value = 0
  }
})

function formatPoints(value) {
  return value.toLocaleString('ja-JP')
}

function handleClose() {
  if (!props.loading) {
    emit('update:modelValue', false)
  }
}

function handleConfirm() {
  if (isValid.value) {
    emit('confirm', points.value)
  }
}
</script>
