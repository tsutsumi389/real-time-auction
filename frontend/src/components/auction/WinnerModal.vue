<script setup>
import Dialog from '@/components/ui/Dialog.vue'
import Button from '@/components/ui/Button.vue'

const props = defineProps({
  open: {
    type: Boolean,
    default: false,
  },
  winner: {
    type: Object,
    default: null,
  },
  item: {
    type: Object,
    default: null,
  },
})

const emit = defineEmits(['close'])

function formatPrice(price) {
  if (!price) return '-'
  return new Intl.NumberFormat('ja-JP').format(price) + ' pt'
}
</script>

<template>
  <Dialog :open="open" @update:open="val => !val && emit('close')" title="落札者発表">
    <div v-if="winner && item" class="space-y-6">
      <!-- 落札情報 -->
      <div class="text-center">
        <div class="mb-4">
          <div class="inline-flex items-center justify-center w-16 h-16 bg-green-100 rounded-full mb-4">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              class="w-8 h-8 text-green-600"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </div>
          <h3 class="text-2xl font-bold text-gray-900 mb-2">落札成立</h3>
          <p class="text-sm text-gray-500">おめでとうございます!</p>
        </div>
      </div>

      <!-- 商品情報 -->
      <div class="bg-gray-50 rounded-lg p-4">
        <div class="text-sm text-gray-500 mb-1">落札商品</div>
        <div class="font-semibold text-lg mb-2">{{ item.name }}</div>
        <div class="text-sm text-gray-600">ロット番号: {{ item.lot_number }}</div>
      </div>

      <!-- 落札者情報 -->
      <div class="bg-green-50 rounded-lg p-4">
        <div class="text-sm text-gray-500 mb-2">落札者</div>
        <div class="flex items-center gap-3 mb-3">
          <div class="w-12 h-12 bg-green-200 rounded-full flex items-center justify-center text-green-700 font-bold text-lg">
            {{ winner.display_name?.charAt(0).toUpperCase() }}
          </div>
          <div>
            <div class="font-semibold text-lg">{{ winner.display_name }}</div>
            <div class="text-xs text-gray-500">ID: {{ winner.id?.slice(0, 8) }}</div>
          </div>
        </div>
        <div class="border-t border-green-200 pt-3">
          <div class="flex justify-between items-center">
            <span class="text-sm text-gray-600">落札価格</span>
            <span class="text-2xl font-bold text-green-600">{{ formatPrice(winner.winning_price) }}</span>
          </div>
        </div>
      </div>

      <!-- 入札統計 -->
      <div class="grid grid-cols-2 gap-4">
        <div class="text-center p-3 bg-gray-50 rounded-lg">
          <div class="text-sm text-gray-500 mb-1">総入札数</div>
          <div class="text-xl font-bold">{{ winner.bid_count || 0 }}</div>
        </div>
        <div class="text-center p-3 bg-gray-50 rounded-lg">
          <div class="text-sm text-gray-500 mb-1">競合者数</div>
          <div class="text-xl font-bold">{{ winner.competitor_count || 0 }}</div>
        </div>
      </div>
    </div>

    <template #footer="{ close }">
      <Button @click="close" variant="default" class="w-full">
        閉じる
      </Button>
    </template>
  </Dialog>
</template>
