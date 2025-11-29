<script setup>
import { computed } from 'vue'
import Card from '@/components/ui/Card.vue'

const props = defineProps({
  priceHistory: {
    type: Array,
    default: () => [],
  },
})

function formatPrice(price) {
  return new Intl.NumberFormat('ja-JP').format(price) + ' pt'
}

function formatTime(timestamp) {
  const date = new Date(timestamp)
  return date.toLocaleTimeString('ja-JP', { hour: '2-digit', minute: '2-digit', second: '2-digit' })
}

const sortedHistory = computed(() => {
  return [...props.priceHistory].sort((a, b) => new Date(b.disclosed_at) - new Date(a.disclosed_at))
})
</script>

<template>
  <Card class="p-6">
    <h3 class="text-lg font-semibold mb-4">価格開示履歴</h3>

    <div v-if="sortedHistory.length > 0" class="space-y-2 max-h-96 overflow-y-auto">
      <div
        v-for="(history, index) in sortedHistory"
        :key="history.id"
        :class="[
          'p-3 rounded-md border transition-all',
          index === 0 ? 'bg-blue-50 border-blue-200' : 'bg-white border-gray-200'
        ]"
      >
        <div class="flex justify-between items-center">
          <div class="flex items-center gap-2">
            <div
              :class="[
                'w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-semibold',
                index === 0 ? 'bg-blue-500' : 'bg-gray-400'
              ]"
            >
              {{ sortedHistory.length - index }}
            </div>
            <div>
              <div class="font-medium">{{ formatTime(history.disclosed_at) }}</div>
              <div class="text-xs text-gray-500">
                {{ history.had_bid ? '入札あり' : '入札なし' }}
              </div>
            </div>
          </div>
          <div class="text-right">
            <div class="font-semibold text-lg">{{ formatPrice(history.price) }}</div>
            <div v-if="index === 0" class="text-xs text-blue-600 font-medium">最新価格</div>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="text-center text-gray-400 py-8">
      価格開示履歴がありません
    </div>
  </Card>
</template>
