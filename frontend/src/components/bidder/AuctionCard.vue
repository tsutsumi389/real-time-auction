<template>
  <div class="auction-card bg-white border border-gray-200 rounded-lg overflow-hidden hover:shadow-lg transition-shadow">
    <!-- サムネイル画像 -->
    <div v-if="auction.thumbnail_url" class="relative w-full h-48 bg-gray-100">
      <img
        :src="auction.thumbnail_url"
        :alt="auction.title"
        class="w-full h-full object-cover"
        @error="handleImageError"
      />
      <div class="absolute top-3 right-3">
        <AuctionStatusBadge :status="auction.status" />
      </div>
    </div>
    <div v-else class="relative w-full h-48 bg-gray-100 flex items-center justify-center">
      <div class="text-center text-gray-400">
        <svg class="mx-auto h-12 w-12 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
        </svg>
        <p class="text-sm">画像なし</p>
      </div>
      <div class="absolute top-3 right-3">
        <AuctionStatusBadge :status="auction.status" />
      </div>
    </div>

    <!-- カード本体 -->
    <div class="p-4 sm:p-6">
      <!-- タイトル -->
      <h3 class="text-lg sm:text-xl font-semibold text-gray-900 mb-2 line-clamp-2">
        {{ auction.title }}
      </h3>

      <!-- 説明 -->
      <p class="text-sm sm:text-base text-gray-600 mb-4 line-clamp-2">
        {{ auction.description || 'オークションの説明はありません' }}
      </p>

      <!-- メタ情報 -->
      <div class="space-y-2 mb-4">
        <!-- 出品物数 -->
        <div class="flex items-center text-sm text-gray-500">
          <svg class="h-4 w-4 sm:h-5 sm:w-5 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"></path>
          </svg>
          <span>出品物: <strong>{{ auction.item_count }}</strong>点</span>
        </div>

        <!-- 開始日時 -->
        <div v-if="auction.started_at" class="flex items-center text-sm text-gray-500">
          <svg class="h-4 w-4 sm:h-5 sm:w-5 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
          </svg>
          <span>開始: {{ formatDate(auction.started_at) }}</span>
        </div>

        <!-- 更新日時 -->
        <div class="flex items-center text-sm text-gray-500">
          <svg class="h-4 w-4 sm:h-5 sm:w-5 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
          </svg>
          <span>更新: {{ formatDate(auction.updated_at) }}</span>
        </div>
      </div>

      <!-- アクションボタン -->
      <div class="flex gap-2">
        <button
          @click="handleViewDetails"
          class="flex-1 px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors text-sm sm:text-base font-medium"
        >
          詳細を見る
        </button>
        <button
          v-if="auction.status === 'active'"
          @click="handleJoinAuction"
          class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors text-sm sm:text-base font-medium whitespace-nowrap"
        >
          参加する
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import AuctionStatusBadge from './AuctionStatusBadge.vue'

const router = useRouter()

const props = defineProps({
  auction: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['join-auction'])

const imageError = ref(false)

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('ja-JP', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

const handleImageError = () => {
  imageError.value = true
}

const handleViewDetails = () => {
  router.push({ name: 'bidder-auction-detail', params: { id: props.auction.id } })
}

const handleJoinAuction = () => {
  emit('join-auction', props.auction)
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.auction-card {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.auction-card > div:last-child {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.auction-card > div:last-child > div:last-child {
  margin-top: auto;
}
</style>
