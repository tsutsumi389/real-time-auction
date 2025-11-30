<template>
  <div class="product-image-gallery">
    <!-- Main Image Area -->
    <div class="relative aspect-square w-full overflow-hidden rounded-xl bg-gray-100 border border-gray-200">
      <template v-if="hasImages">
        <img
          :src="currentImage"
          :alt="altText"
          class="h-full w-full object-cover object-center transition-opacity duration-300"
          :class="{ 'opacity-0': changing }"
          @load="onImageLoad"
        />
        
        <!-- Navigation Arrows (if multiple images) -->
        <div v-if="images.length > 1" class="absolute inset-0 flex items-center justify-between p-4 pointer-events-none">
          <button
            @click="prevImage"
            class="pointer-events-auto rounded-full bg-white/80 p-2 text-gray-800 shadow-sm hover:bg-white transition-colors"
            aria-label="Previous image"
          >
            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <button
            @click="nextImage"
            class="pointer-events-auto rounded-full bg-white/80 p-2 text-gray-800 shadow-sm hover:bg-white transition-colors"
            aria-label="Next image"
          >
            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>

        <!-- Image Counter badge -->
        <div v-if="images.length > 1" class="absolute bottom-4 right-4 rounded-full bg-black/50 px-3 py-1 text-xs font-medium text-white backdrop-blur-sm">
          {{ currentIndex + 1 }} / {{ images.length }}
        </div>
      </template>

      <!-- No Image Placeholder -->
      <div v-else class="flex h-full w-full items-center justify-center text-gray-400">
        <div class="text-center">
          <svg class="mx-auto h-16 w-16 mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <span class="text-sm">No Image</span>
        </div>
      </div>
    </div>

    <!-- Thumbnails (if multiple images) -->
    <div v-if="images.length > 1" class="mt-4 flex gap-2 overflow-x-auto pb-2 custom-scrollbar">
      <button
        v-for="(img, index) in images"
        :key="index"
        @click="selectImage(index)"
        class="relative h-16 w-16 flex-shrink-0 overflow-hidden rounded-lg border-2 transition-all"
        :class="currentIndex === index ? 'border-blue-600 ring-2 ring-blue-100' : 'border-transparent opacity-70 hover:opacity-100'"
      >
        <img :src="img" class="h-full w-full object-cover" alt="" />
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  images: {
    type: Array,
    default: () => []
  },
  altText: {
    type: String,
    default: 'Product Image'
  }
})

const currentIndex = ref(0)
const changing = ref(false)

const hasImages = computed(() => props.images && props.images.length > 0)
const currentImage = computed(() => hasImages.value ? props.images[currentIndex.value] : null)

// Reset index when images change
watch(() => props.images, () => {
  currentIndex.value = 0
})

function nextImage() {
  if (!hasImages.value) return
  changeImage((currentIndex.value + 1) % props.images.length)
}

function prevImage() {
  if (!hasImages.value) return
  changeImage((currentIndex.value - 1 + props.images.length) % props.images.length)
}

function selectImage(index) {
  if (currentIndex.value === index) return
  changeImage(index)
}

function changeImage(newIndex) {
  changing.value = true
  setTimeout(() => {
    currentIndex.value = newIndex
  }, 150) // Half of transition duration
}

function onImageLoad() {
  changing.value = false
}
</script>

<style scoped>
.custom-scrollbar {
  scrollbar-width: thin;
  scrollbar-color: #cbd5e0 transparent;
}
.custom-scrollbar::-webkit-scrollbar {
  height: 4px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent;
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background-color: #cbd5e0;
  border-radius: 20px;
}
</style>
