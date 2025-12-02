<template>
  <div
    class="auction-card-skeleton bg-white border border-gray-200 rounded-lg overflow-hidden"
    role="status"
    aria-busy="true"
    aria-label="オークションカードを読み込み中"
  >
    <!-- Thumbnail skeleton -->
    <div class="relative w-full h-48 bg-muted">
      <Skeleton class="w-full h-full" variant="rectangular" />
      <!-- Badge skeleton -->
      <div class="absolute top-3 right-3">
        <Skeleton width="80px" height="24px" rounded />
      </div>
    </div>

    <!-- Card body -->
    <div class="p-4 sm:p-6">
      <!-- Title skeleton (2 lines) -->
      <div class="mb-2 space-y-2">
        <Skeleton width="90%" height="24px" />
        <Skeleton width="60%" height="24px" />
      </div>

      <!-- Description skeleton (2 lines) -->
      <div class="mb-4 space-y-2">
        <Skeleton width="100%" height="16px" />
        <Skeleton width="80%" height="16px" />
      </div>

      <!-- Meta info skeleton (3 lines) -->
      <div class="space-y-2.5 mb-4">
        <div class="flex items-center">
          <Skeleton width="20px" height="20px" class="mr-2" rounded />
          <Skeleton width="120px" height="16px" />
        </div>
        <div class="flex items-center">
          <Skeleton width="20px" height="20px" class="mr-2" rounded />
          <Skeleton width="150px" height="16px" />
        </div>
        <div class="flex items-center">
          <Skeleton width="20px" height="20px" class="mr-2" rounded />
          <Skeleton width="140px" height="16px" />
        </div>
      </div>

      <!-- Action buttons skeleton -->
      <div class="flex gap-2">
        <Skeleton class="flex-1" height="40px" rounded />
        <Skeleton width="100px" height="40px" rounded />
      </div>
    </div>

    <span class="sr-only">読み込み中</span>
  </div>
</template>

<script setup>
import Skeleton from '@/components/ui/Skeleton.vue'
</script>

<style scoped>
.auction-card-skeleton {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.auction-card-skeleton > div:last-child {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.auction-card-skeleton > div:last-child > div:last-child {
  margin-top: auto;
}
</style>
