<script setup>
import { computed } from 'vue'
import Card from '@/components/ui/Card.vue'

const props = defineProps({
  bids: {
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

const sortedBids = computed(() => {
  return [...props.bids].sort((a, b) => new Date(b.bid_at) - new Date(a.bid_at))
})
</script>

<template>
  <Card class="p-6">
    <h3 class="text-lg font-semibold mb-4">入札履歴</h3>

    <div v-if="sortedBids.length > 0" class="space-y-2 max-h-96 overflow-y-auto">
      <TransitionGroup
        enter-active-class="transition-all duration-500 ease-out"
        enter-from-class="opacity-0 -translate-y-4 scale-95"
        enter-to-class="opacity-100 translate-y-0 scale-100"
        leave-active-class="transition-all duration-300 ease-in"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
        move-class="transition-transform duration-300"
      >
        <div
          v-for="(bid, index) in sortedBids"
          :key="bid.id"
          :class="[
            'p-3 rounded-md border transition-all',
            index === 0 ? 'bg-green-50 border-green-200 animate-pulse-subtle' : 'bg-white border-gray-200'
          ]"
        >
        <div class="flex justify-between items-center">
          <div class="flex items-center gap-2">
            <div
              :class="[
                'w-8 h-8 rounded-full flex items-center justify-center text-white text-sm font-semibold',
                bid.is_winning ? 'bg-green-500' : 'bg-gray-400'
              ]"
            >
              {{ index + 1 }}
            </div>
            <div>
              <div class="font-medium">{{ bid.bidder_display_name }}</div>
              <div class="text-xs text-gray-500">{{ formatTime(bid.bid_at) }}</div>
            </div>
          </div>
          <div class="text-right">
            <div class="font-semibold text-lg">{{ formatPrice(bid.price) }}</div>
            <div v-if="bid.is_winning" class="text-xs text-green-600 font-medium">最高入札</div>
          </div>
        </div>
        </div>
      </TransitionGroup>
    </div>

    <div v-else class="text-center text-gray-400 py-8">
      入札履歴がありません
    </div>
  </Card>
</template>
