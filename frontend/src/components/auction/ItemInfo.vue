<script setup>
import { computed } from 'vue'
import Card from '@/components/ui/Card.vue'

const props = defineProps({
  item: {
    type: Object,
    default: null,
  },
})

const statusText = computed(() => {
  if (!props.item) return '-'

  switch (props.item.status) {
    case 'pending':
      return '待機中'
    case 'active':
      return '進行中'
    case 'ended':
      return '終了'
    default:
      return props.item.status
  }
})

const statusClass = computed(() => {
  if (!props.item) return ''

  switch (props.item.status) {
    case 'pending':
      return 'bg-gray-100 text-gray-700'
    case 'active':
      return 'bg-green-100 text-green-700'
    case 'ended':
      return 'bg-blue-100 text-blue-700'
    default:
      return 'bg-gray-100 text-gray-700'
  }
})

const formattedPrice = computed(() => {
  if (!props.item?.current_price) return '-'
  return new Intl.NumberFormat('ja-JP').format(props.item.current_price) + ' pt'
})

const formattedStartingPrice = computed(() => {
  if (!props.item?.starting_price) return '-'
  return new Intl.NumberFormat('ja-JP').format(props.item.starting_price) + ' pt'
})
</script>

<template>
  <Card class="p-6">
    <div v-if="item">
      <div class="flex justify-between items-start mb-4">
        <div>
          <div class="text-sm text-gray-500 mb-1">ロット番号</div>
          <div class="text-2xl font-bold">{{ item.lot_number }}</div>
        </div>
        <div :class="['px-3 py-1 rounded-full text-sm font-medium', statusClass]">
          {{ statusText }}
        </div>
      </div>

      <div class="mb-4">
        <div class="text-lg font-semibold mb-2">{{ item.name }}</div>
        <div class="text-sm text-gray-600">{{ item.description || '説明なし' }}</div>
      </div>

      <div class="grid grid-cols-2 gap-4">
        <div>
          <div class="text-sm text-gray-500 mb-1">開始価格</div>
          <div class="text-lg font-semibold">{{ formattedStartingPrice }}</div>
        </div>
        <div>
          <div class="text-sm text-gray-500 mb-1">現在価格</div>
          <div class="text-lg font-semibold text-blue-600">{{ formattedPrice }}</div>
        </div>
      </div>

      <div v-if="item.metadata" class="mt-4 pt-4 border-t">
        <div class="text-sm text-gray-500 mb-2">詳細情報</div>
        <div class="text-sm">
          <pre class="whitespace-pre-wrap text-gray-700">{{ JSON.stringify(item.metadata, null, 2) }}</pre>
        </div>
      </div>
    </div>

    <div v-else class="text-center text-gray-400 py-8">
      商品が選択されていません
    </div>
  </Card>
</template>
