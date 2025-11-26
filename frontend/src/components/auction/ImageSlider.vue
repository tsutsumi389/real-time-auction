<script setup>
import { ref, computed } from 'vue'
import Card from '@/components/ui/Card.vue'
import Button from '@/components/ui/Button.vue'

const props = defineProps({
  media: {
    type: Array,
    default: () => [],
  },
})

const currentIndex = ref(0)

const sortedMedia = computed(() => {
  return [...props.media].sort((a, b) => a.display_order - b.display_order)
})

const currentMedia = computed(() => {
  return sortedMedia.value[currentIndex.value]
})

const hasMedia = computed(() => {
  return sortedMedia.value.length > 0
})

const hasPrevious = computed(() => {
  return currentIndex.value > 0
})

const hasNext = computed(() => {
  return currentIndex.value < sortedMedia.value.length - 1
})

function goToPrevious() {
  if (hasPrevious.value) {
    currentIndex.value--
  }
}

function goToNext() {
  if (hasNext.value) {
    currentIndex.value++
  }
}

function goToIndex(index) {
  currentIndex.value = index
}
</script>

<template>
  <Card class="p-6">
    <h3 class="text-lg font-semibold mb-4">商品画像</h3>

    <div v-if="hasMedia" class="space-y-4">
      <!-- メイン画像 -->
      <div class="relative bg-gray-100 rounded-lg overflow-hidden aspect-video">
        <img
          v-if="currentMedia.media_type === 'image'"
          :src="currentMedia.url"
          :alt="currentMedia.caption || '商品画像'"
          class="w-full h-full object-contain"
        />
        <video
          v-else-if="currentMedia.media_type === 'video'"
          :src="currentMedia.url"
          controls
          class="w-full h-full object-contain"
        >
          ブラウザが動画をサポートしていません
        </video>

        <!-- ナビゲーションボタン -->
        <div v-if="sortedMedia.length > 1" class="absolute inset-0 flex items-center justify-between p-4">
          <Button
            @click="goToPrevious"
            :disabled="!hasPrevious"
            variant="outline"
            size="sm"
            class="bg-white/80 hover:bg-white"
          >
            ‹
          </Button>
          <Button
            @click="goToNext"
            :disabled="!hasNext"
            variant="outline"
            size="sm"
            class="bg-white/80 hover:bg-white"
          >
            ›
          </Button>
        </div>

        <!-- インジケーター -->
        <div v-if="sortedMedia.length > 1" class="absolute bottom-4 left-1/2 transform -translate-x-1/2">
          <div class="flex gap-2">
            <button
              v-for="(_, index) in sortedMedia"
              :key="index"
              @click="goToIndex(index)"
              :class="[
                'w-2 h-2 rounded-full transition-all',
                index === currentIndex ? 'bg-white w-4' : 'bg-white/50'
              ]"
            />
          </div>
        </div>
      </div>

      <!-- キャプション -->
      <div v-if="currentMedia.caption" class="text-sm text-gray-600 text-center">
        {{ currentMedia.caption }}
      </div>

      <!-- サムネイル -->
      <div v-if="sortedMedia.length > 1" class="flex gap-2 overflow-x-auto">
        <button
          v-for="(item, index) in sortedMedia"
          :key="item.id"
          @click="goToIndex(index)"
          :class="[
            'flex-shrink-0 w-20 h-20 rounded-md overflow-hidden border-2 transition-all',
            index === currentIndex ? 'border-blue-500' : 'border-gray-200 hover:border-gray-300'
          ]"
        >
          <img
            v-if="item.media_type === 'image'"
            :src="item.url"
            :alt="item.caption || '商品画像'"
            class="w-full h-full object-cover"
          />
          <div
            v-else
            class="w-full h-full bg-gray-200 flex items-center justify-center text-xs text-gray-500"
          >
            動画
          </div>
        </button>
      </div>
    </div>

    <div v-else class="text-center text-gray-400 py-8">
      画像がありません
    </div>
  </Card>
</template>
