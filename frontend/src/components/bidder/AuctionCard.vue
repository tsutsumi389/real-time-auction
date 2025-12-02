<template>
  <article
    class="auction-card bg-white border border-gray-200 rounded-lg overflow-hidden hover:shadow-xl hover:scale-[1.02] transition-all duration-300"
    role="article"
    :aria-label="`${auction.title}のオークション`"
  >
    <!-- Thumbnail image -->
    <div v-if="auction.thumbnail_url" class="relative w-full h-48 bg-gray-100 group">
      <img
        :src="auction.thumbnail_url"
        :alt="auction.title"
        class="w-full h-full object-cover transition-transform duration-300 group-hover:scale-105"
        @error="handleImageError"
      />
      <!-- Gradient overlay for better badge visibility -->
      <div class="absolute inset-0 bg-gradient-to-t from-black/30 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
      <div class="absolute top-3 right-3">
        <AuctionStatusBadge :status="auction.status" />
      </div>
    </div>
    <div v-else class="relative w-full h-48 bg-gray-100 flex items-center justify-center">
      <div class="text-center text-gray-400">
        <ImageOff class="mx-auto h-12 w-12 mb-2" :stroke-width="1.5" />
        <p class="text-sm">画像なし</p>
      </div>
      <div class="absolute top-3 right-3">
        <AuctionStatusBadge :status="auction.status" />
      </div>
    </div>

    <!-- Card body -->
    <div class="p-4 sm:p-6">
      <!-- Title -->
      <h3 class="text-lg sm:text-xl font-semibold text-gray-900 mb-2 line-clamp-2">
        {{ auction.title }}
      </h3>

      <!-- Description -->
      <p class="text-sm sm:text-base text-muted-foreground mb-4 line-clamp-2">
        {{ auction.description || 'オークションの説明はありません' }}
      </p>

      <!-- Meta info -->
      <div class="space-y-2.5 mb-4">
        <!-- Item count -->
        <div class="flex items-center text-sm text-gray-500">
          <Package2 class="h-4 w-4 sm:h-5 sm:w-5 mr-2 text-gray-400" :stroke-width="1.5" />
          <span>出品物: <strong>{{ auction.item_count }}</strong>点</span>
        </div>

        <!-- Start date -->
        <div v-if="auction.started_at" class="flex items-center text-sm text-gray-500">
          <Calendar class="h-4 w-4 sm:h-5 sm:w-5 mr-2 text-gray-400" :stroke-width="1.5" />
          <span>開始: {{ formatDate(auction.started_at) }}</span>
        </div>

        <!-- Updated date -->
        <div class="flex items-center text-sm text-gray-500">
          <Clock class="h-4 w-4 sm:h-5 sm:w-5 mr-2 text-gray-400" :stroke-width="1.5" />
          <span>更新: {{ formatDate(auction.updated_at) }}</span>
        </div>
      </div>

      <!-- Action buttons -->
      <div class="flex gap-2">
        <Button
          variant="default"
          class="flex-1"
          @click="handleViewDetails"
          :aria-label="`${auction.title}の詳細を見る`"
        >
          詳細を見る
        </Button>
        <Button
          v-if="auction.status === 'active'"
          variant="default"
          class="bg-green-600 hover:bg-green-700 whitespace-nowrap"
          @click="handleJoinAuction"
          :aria-label="`${auction.title}に参加する`"
        >
          参加する
        </Button>
      </div>
    </div>
  </article>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Package2, Calendar, Clock, ImageOff } from 'lucide-vue-next'
import AuctionStatusBadge from './AuctionStatusBadge.vue'
import Button from '@/components/ui/Button.vue'

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
