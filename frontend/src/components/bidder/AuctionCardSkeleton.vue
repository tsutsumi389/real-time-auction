<template>
  <div
    class="auction-card-skeleton lux-glass-strong rounded-2xl overflow-hidden border border-lux-gold/10"
    role="status"
    aria-busy="true"
    aria-label="Loading auction card"
  >
    <!-- Thumbnail Skeleton -->
    <div class="relative w-full h-52 bg-lux-noir-medium">
      <div class="skeleton-shimmer w-full h-full"></div>
      <!-- Status Badge Skeleton -->
      <div class="absolute top-4 right-4">
        <div class="skeleton-shimmer w-16 h-7 rounded-lg"></div>
      </div>
      <!-- Item Count Badge Skeleton -->
      <div class="absolute bottom-4 left-4">
        <div class="skeleton-shimmer w-24 h-8 rounded-lg"></div>
      </div>
    </div>

    <!-- Card Body -->
    <div class="p-5 sm:p-6 flex flex-col">
      <!-- Title Skeleton -->
      <div class="mb-3 space-y-2">
        <div class="skeleton-shimmer h-7 w-[90%] rounded"></div>
        <div class="skeleton-shimmer h-7 w-[60%] rounded"></div>
      </div>

      <!-- Description Skeleton -->
      <div class="mb-5 space-y-2">
        <div class="skeleton-shimmer h-4 w-full rounded"></div>
        <div class="skeleton-shimmer h-4 w-[75%] rounded"></div>
      </div>

      <!-- Meta Info Skeleton -->
      <div class="space-y-3 mb-6 flex-grow">
        <!-- Start Date -->
        <div class="flex items-center gap-3">
          <div class="skeleton-shimmer w-8 h-8 rounded-lg flex-shrink-0"></div>
          <div class="flex-1 space-y-1.5">
            <div class="skeleton-shimmer h-3 w-16 rounded"></div>
            <div class="skeleton-shimmer h-4 w-32 rounded"></div>
          </div>
        </div>
        <!-- Updated Date -->
        <div class="flex items-center gap-3">
          <div class="skeleton-shimmer w-8 h-8 rounded-lg flex-shrink-0"></div>
          <div class="flex-1 space-y-1.5">
            <div class="skeleton-shimmer h-3 w-16 rounded"></div>
            <div class="skeleton-shimmer h-4 w-32 rounded"></div>
          </div>
        </div>
      </div>

      <!-- Button Skeleton -->
      <div class="flex gap-3 mt-auto">
        <div class="skeleton-shimmer flex-1 h-11 rounded-xl"></div>
        <div class="skeleton-shimmer flex-1 h-11 rounded-xl"></div>
      </div>
    </div>

    <span class="sr-only">Loading</span>
  </div>
</template>

<script setup>
</script>

<style scoped>
.auction-card-skeleton {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: rgba(255, 255, 255, 0.03);
}

/* Skeleton shimmer animation */
.skeleton-shimmer {
  background: linear-gradient(
    90deg,
    hsl(0 0% 12%) 0%,
    hsl(0 0% 16%) 50%,
    hsl(0 0% 12%) 100%
  );
  background-size: 200% 100%;
  animation: shimmer 1.5s ease-in-out infinite;
}

@keyframes shimmer {
  0% {
    background-position: 200% 0;
  }
  100% {
    background-position: -200% 0;
  }
}

/* Luxury color utilities */
.bg-lux-noir-medium {
  background-color: hsl(0 0% 12%);
}

.border-lux-gold\/10 {
  border-color: hsl(43 74% 49% / 0.1);
}
</style>
