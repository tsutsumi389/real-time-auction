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
        class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-4xl sm:w-full"
      >
        <div class="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
          <div>
            <h3 class="text-lg leading-6 font-medium text-gray-900 mb-4" id="modal-title">
              ポイント履歴
            </h3>

            <!-- 入札者情報 -->
            <div class="bg-gray-50 p-3 rounded mb-4">
              <p class="text-sm text-gray-700">
                <span class="font-medium">メール:</span> {{ bidder?.email }}
              </p>
              <p class="text-sm text-gray-700 mt-1">
                <span class="font-medium">表示名:</span> {{ bidder?.display_name || '（未設定）' }}
              </p>
              <p class="text-sm text-gray-700 mt-1">
                <span class="font-medium">現在のポイント:</span> {{ formatPoints(bidder?.points || 0) }}
              </p>
            </div>

            <!-- 履歴テーブル -->
            <div v-if="loading" class="p-8 text-center">
              <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
              <p class="mt-2 text-gray-600">読み込み中...</p>
            </div>

            <div v-else-if="history.length > 0" class="overflow-x-auto">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th scope="col" class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">日時</th>
                    <th scope="col" class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">種別</th>
                    <th scope="col" class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">増減</th>
                    <th scope="col" class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">残高</th>
                    <th scope="col" class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">関連オークション</th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  <tr v-for="item in history" :key="item.id">
                    <td class="px-4 py-3 whitespace-nowrap text-sm text-gray-500">
                      {{ formatDateTime(item.created_at) }}
                    </td>
                    <td class="px-4 py-3 whitespace-nowrap text-sm">
                      <PointTypeBadge :type="item.type" />
                    </td>
                    <td class="px-4 py-3 whitespace-nowrap text-sm text-right" :class="getAmountClass(item.type)">
                      {{ formatAmount(item.amount, item.type) }}
                    </td>
                    <td class="px-4 py-3 whitespace-nowrap text-sm text-right text-gray-900">
                      {{ formatPoints(item.balance_after) }}
                    </td>
                    <td class="px-4 py-3 whitespace-nowrap text-sm text-gray-500">
                      {{ item.auction_title || '-' }}
                    </td>
                  </tr>
                </tbody>
              </table>

              <!-- ページネーション -->
              <div v-if="totalPages > 1" class="mt-4 flex items-center justify-between border-t border-gray-200 px-4 py-3">
                <div class="flex-1 flex justify-between sm:hidden">
                  <button
                    @click="changePage(currentPage - 1)"
                    :disabled="currentPage === 1"
                    class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
                  >
                    前へ
                  </button>
                  <button
                    @click="changePage(currentPage + 1)"
                    :disabled="currentPage === totalPages"
                    class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
                  >
                    次へ
                  </button>
                </div>
                <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
                  <div>
                    <p class="text-sm text-gray-700">
                      ページ <span class="font-medium">{{ currentPage }}</span> / <span class="font-medium">{{ totalPages }}</span>
                    </p>
                  </div>
                  <div>
                    <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
                      <button
                        @click="changePage(currentPage - 1)"
                        :disabled="currentPage === 1"
                        class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                      >
                        ‹
                      </button>
                      <button
                        @click="changePage(currentPage + 1)"
                        :disabled="currentPage === totalPages"
                        class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
                      >
                        ›
                      </button>
                    </nav>
                  </div>
                </div>
              </div>
            </div>

            <div v-else class="p-8 text-center text-gray-500">
              <p>ポイント履歴がありません</p>
            </div>
          </div>
        </div>
        <div class="bg-gray-50 px-4 py-3 sm:px-6 sm:flex sm:flex-row-reverse">
          <button
            @click="handleClose"
            class="w-full inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 sm:w-auto sm:text-sm"
          >
            閉じる
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useBidderStore } from '@/stores/bidder'
import PointTypeBadge from './PointTypeBadge.vue'

const props = defineProps({
  modelValue: {
    type: Boolean,
    required: true,
  },
  bidder: {
    type: Object,
    default: null,
  },
})

const emit = defineEmits(['update:modelValue'])

const bidderStore = useBidderStore()

const loading = ref(false)
const history = ref([])
const currentPage = ref(1)
const totalPages = ref(1)

// モーダルを開いたときに履歴を取得
watch(() => props.modelValue, async (newValue) => {
  if (newValue && props.bidder) {
    await fetchHistory(1)
  }
})

async function fetchHistory(page) {
  loading.value = true
  const response = await bidderStore.fetchPointHistory(props.bidder.id, page, 10)
  if (response) {
    history.value = response.history
    currentPage.value = response.pagination.page
    totalPages.value = response.pagination.total_pages
  }
  loading.value = false
}

function changePage(page) {
  fetchHistory(page)
}

function formatPoints(value) {
  return value.toLocaleString('ja-JP')
}

function formatAmount(amount, type) {
  const sign = ['grant', 'release', 'refund'].includes(type) ? '+' : ''
  return `${sign}${amount.toLocaleString('ja-JP')}`
}

function getAmountClass(type) {
  if (['grant', 'release', 'refund'].includes(type)) {
    return 'text-green-600 font-medium'
  } else {
    return 'text-red-600 font-medium'
  }
}

function formatDateTime(dateString) {
  const date = new Date(dateString)
  return date.toLocaleString('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

function handleClose() {
  emit('update:modelValue', false)
}
</script>
