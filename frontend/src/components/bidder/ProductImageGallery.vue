<template>
  <div class="product-image-gallery">
    <!-- Main Image Area -->
    <div class="relative aspect-square w-full overflow-hidden" style="background: hsl(0 0% 8%)">
      <template v-if="hasImages">
        <!-- Main Image -->
        <img
          :src="currentImage"
          :alt="altText"
          class="h-full w-full object-cover object-center transition-all duration-500"
          :class="{ 'opacity-0 scale-105': changing }"
          @load="onImageLoad"
        />

        <!-- Gradient Overlay -->
        <div class="absolute inset-0 bg-gradient-to-t from-lux-noir/60 via-transparent to-transparent pointer-events-none"></div>

        <!-- Navigation Arrows (if multiple images) -->
        <div v-if="images.length > 1" class="absolute inset-0 flex items-center justify-between p-4 pointer-events-none">
          <button
            @click="prevImage"
            class="pointer-events-auto w-10 h-10 rounded-full bg-lux-noir/60 backdrop-blur-sm border border-lux-gold/20 text-lux-cream hover:bg-lux-noir/80 hover:border-lux-gold/40 transition-all duration-200 flex items-center justify-center group"
            aria-label="Previous image"
          >
            <svg class="w-5 h-5 group-hover:-translate-x-0.5 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 19l-7-7 7-7" />
            </svg>
          </button>
          <button
            @click="nextImage"
            class="pointer-events-auto w-10 h-10 rounded-full bg-lux-noir/60 backdrop-blur-sm border border-lux-gold/20 text-lux-cream hover:bg-lux-noir/80 hover:border-lux-gold/40 transition-all duration-200 flex items-center justify-center group"
            aria-label="Next image"
          >
            <svg class="w-5 h-5 group-hover:translate-x-0.5 transition-transform" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 5l7 7-7 7" />
            </svg>
          </button>
        </div>

        <!-- Image Counter Badge -->
        <div v-if="images.length > 1" class="absolute bottom-4 right-4 flex items-center gap-2">
          <!-- Dot Indicators -->
          <div class="flex items-center gap-1.5 px-3 py-1.5 rounded-full bg-lux-noir/60 backdrop-blur-sm border border-lux-gold/20">
            <span
              v-for="(_, index) in images"
              :key="index"
              @click="selectImage(index)"
              :class="[
                'w-1.5 h-1.5 rounded-full cursor-pointer transition-all duration-200',
                currentIndex === index ? 'bg-lux-gold w-4' : 'bg-lux-silver/40 hover:bg-lux-silver/60'
              ]"
            ></span>
          </div>
        </div>

        <!-- Image Number -->
        <div v-if="images.length > 1" class="absolute bottom-4 left-4">
          <span class="px-3 py-1.5 rounded-full bg-lux-noir/60 backdrop-blur-sm border border-lux-gold/20 text-xs font-medium text-lux-cream tabular-nums">
            {{ currentIndex + 1 }} / {{ images.length }}
          </span>
        </div>
      </template>

      <!-- No Image Placeholder -->
      <div v-else class="flex h-full w-full items-center justify-center">
        <div class="text-center">
          <div class="w-20 h-20 mx-auto mb-4 rounded-full bg-lux-noir-medium border border-lux-gold/10 flex items-center justify-center">
            <svg class="w-10 h-10 text-lux-silver/30" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
          </div>
          <p class="text-sm text-lux-silver/50 font-medium">No Image Available</p>
        </div>
      </div>
    </div>

    <!-- Thumbnails (if multiple images) -->
    <div v-if="images.length > 1" class="mt-3 px-2">
      <div class="flex gap-2 overflow-x-auto lux-scrollbar pb-2">
        <button
          v-for="(img, index) in images"
          :key="index"
          @click="selectImage(index)"
          :class="[
            'relative flex-shrink-0 w-16 h-16 rounded-lg overflow-hidden transition-all duration-200',
            currentIndex === index
              ? 'ring-2 ring-lux-gold ring-offset-2 ring-offset-lux-noir'
              : 'opacity-50 hover:opacity-80'
          ]"
        >
          <img :src="getThumbnailUrl(img)" class="w-full h-full object-cover" alt="" />
          <div v-if="currentIndex === index" class="absolute inset-0 bg-lux-gold/10"></div>
        </button>
      </div>
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

// Get the full URL for main display
const currentImage = computed(() => {
  if (!hasImages.value) return null
  const img = props.images[currentIndex.value]
  // Support both string URLs and object with url/thumbnail
  return typeof img === 'string' ? img : img.url
})

// Get thumbnail URL for thumbnails list
function getThumbnailUrl(img) {
  if (typeof img === 'string') return img
  return img.thumbnail || img.url
}

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
  }, 200) // Half of transition duration
}

function onImageLoad() {
  changing.value = false
}
</script>

<style scoped>
.product-image-gallery {
  position: relative;
  background: hsl(0 0% 4%);
}
</style>
