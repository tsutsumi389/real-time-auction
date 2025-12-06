<template>
  <article
    class="auction-card lux-glass-strong rounded-2xl overflow-hidden border border-lux-gold/10 hover:border-lux-gold/30 transition-all duration-500 group"
    role="article"
    :aria-label="`${auction.title} auction`"
  >
    <!-- Thumbnail Image -->
    <div class="relative w-full h-52 bg-lux-noir-light overflow-hidden">
      <img
        v-if="auction.thumbnail_url && !imageError"
        :src="auction.thumbnail_url"
        :alt="auction.title"
        class="w-full h-full object-cover transition-transform duration-700 group-hover:scale-110"
        @error="handleImageError"
      />
      <div v-else class="w-full h-full flex items-center justify-center bg-lux-noir-medium">
        <div class="text-center">
          <ImageOff class="mx-auto h-12 w-12 mb-3 text-lux-silver/30" :stroke-width="1.5" />
          <p class="text-xs text-lux-silver/40 font-medium tracking-wide">NO IMAGE</p>
        </div>
      </div>

      <!-- Gradient Overlay -->
      <div class="absolute inset-0 bg-gradient-to-t from-lux-noir via-transparent to-transparent opacity-80"></div>

      <!-- Status Badge -->
      <div class="absolute top-4 right-4">
        <span :class="getStatusBadgeClasses()">
          <span class="status-indicator" :class="getStatusIndicatorClass()"></span>
          {{ getStatusLabel() }}
        </span>
      </div>

      <!-- Item Count Badge -->
      <div class="absolute bottom-4 left-4">
        <div class="flex items-center gap-2 px-3 py-1.5 rounded-lg bg-lux-noir/80 backdrop-blur-sm border border-lux-gold/20">
          <Package2 class="h-4 w-4 text-lux-gold" :stroke-width="1.5" />
          <span class="text-sm font-semibold text-lux-cream">{{ auction.item_count }}</span>
          <span class="text-xs text-lux-silver">items</span>
        </div>
      </div>
    </div>

    <!-- Card Body -->
    <div class="p-5 sm:p-6 flex flex-col">
      <!-- Title -->
      <h3 class="font-display text-xl sm:text-2xl text-lux-cream mb-2 line-clamp-2 tracking-wide group-hover:text-lux-gold transition-colors duration-300">
        {{ auction.title }}
      </h3>

      <!-- Description -->
      <p class="text-sm text-lux-silver/70 mb-5 line-clamp-2 leading-relaxed">
        {{ auction.description || 'No description available' }}
      </p>

      <!-- Meta Info -->
      <div class="space-y-3 mb-6 flex-grow">
        <!-- Start Date -->
        <div v-if="auction.started_at" class="flex items-center gap-3 text-sm">
          <div class="w-8 h-8 rounded-lg bg-lux-noir-medium flex items-center justify-center flex-shrink-0">
            <Calendar class="h-4 w-4 text-lux-gold/60" :stroke-width="1.5" />
          </div>
          <div class="min-w-0">
            <p class="text-xs text-lux-silver/50 uppercase tracking-wider mb-0.5">Started</p>
            <p class="text-lux-cream text-sm font-medium truncate">{{ formatDate(auction.started_at) }}</p>
          </div>
        </div>

        <!-- Updated Date -->
        <div class="flex items-center gap-3 text-sm">
          <div class="w-8 h-8 rounded-lg bg-lux-noir-medium flex items-center justify-center flex-shrink-0">
            <Clock class="h-4 w-4 text-lux-gold/60" :stroke-width="1.5" />
          </div>
          <div class="min-w-0">
            <p class="text-xs text-lux-silver/50 uppercase tracking-wider mb-0.5">Updated</p>
            <p class="text-lux-cream text-sm font-medium truncate">{{ formatDate(auction.updated_at) }}</p>
          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex gap-3 mt-auto">
        <button
          @click="handleViewDetails"
          class="flex-1 h-11 rounded-xl bg-lux-noir-medium border border-lux-noir-soft text-sm font-medium text-lux-cream tracking-wide hover:border-lux-gold/40 hover:bg-lux-noir-soft transition-all duration-300"
          :aria-label="`View details of ${auction.title}`"
        >
          Details
        </button>
        <button
          v-if="auction.status === 'active'"
          @click="handleJoinAuction"
          class="flex-1 h-11 rounded-xl lux-btn-gold text-sm font-semibold tracking-wider flex items-center justify-center gap-2"
          :aria-label="`Join ${auction.title} auction`"
        >
          <span class="relative flex h-2 w-2">
            <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-lux-noir opacity-75"></span>
            <span class="relative inline-flex rounded-full h-2 w-2 bg-lux-noir"></span>
          </span>
          JOIN LIVE
        </button>
      </div>
    </div>
  </article>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { Package2, Calendar, Clock, ImageOff } from 'lucide-vue-next'

const router = useRouter()

const props = defineProps({
  auction: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['join-auction', 'view-details'])

const imageError = ref(false)

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return new Intl.DateTimeFormat('en-US', {
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  }).format(date)
}

const handleImageError = () => {
  imageError.value = true
}

const getStatusBadgeClasses = () => {
  const base = 'inline-flex items-center gap-2 px-3 py-1.5 rounded-lg text-xs font-semibold tracking-wider uppercase backdrop-blur-sm'
  switch (props.auction.status) {
    case 'active':
      return `${base} bg-emerald-500/15 text-emerald-400 border border-emerald-500/30`
    case 'ended':
      return `${base} bg-lux-silver/10 text-lux-silver border border-lux-silver/20`
    case 'cancelled':
      return `${base} bg-red-500/15 text-red-400 border border-red-500/30`
    default:
      return `${base} bg-lux-silver/10 text-lux-silver border border-lux-silver/20`
  }
}

const getStatusIndicatorClass = () => {
  switch (props.auction.status) {
    case 'active':
      return 'indicator-active'
    case 'ended':
      return 'indicator-ended'
    case 'cancelled':
      return 'indicator-cancelled'
    default:
      return 'indicator-ended'
  }
}

const getStatusLabel = () => {
  switch (props.auction.status) {
    case 'active':
      return 'LIVE'
    case 'ended':
      return 'ENDED'
    case 'cancelled':
      return 'CANCELLED'
    default:
      return props.auction.status?.toUpperCase() || 'UNKNOWN'
  }
}

const handleViewDetails = () => {
  emit('view-details', props.auction)
}

const handleJoinAuction = () => {
  emit('join-auction', props.auction)
}
</script>

<style scoped>
.auction-card {
  height: 100%;
  display: flex;
  flex-direction: column;
  background: rgba(255, 255, 255, 0.03);
}

.auction-card:hover {
  box-shadow: 0 0 40px rgba(212, 175, 55, 0.08), 0 0 80px rgba(0, 0, 0, 0.3);
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

/* Status indicator dots */
.status-indicator {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  flex-shrink: 0;
}

.indicator-active {
  background-color: hsl(145 50% 50%);
  box-shadow: 0 0 8px hsl(145 50% 50% / 0.8);
  animation: pulse-glow 2s ease-in-out infinite;
}

.indicator-ended {
  background-color: hsl(220 10% 60%);
}

.indicator-cancelled {
  background-color: hsl(0 65% 50%);
}

@keyframes pulse-glow {
  0%, 100% {
    opacity: 1;
    box-shadow: 0 0 8px hsl(145 50% 50% / 0.8);
  }
  50% {
    opacity: 0.7;
    box-shadow: 0 0 12px hsl(145 50% 50% / 0.5);
  }
}

/* Luxury color utilities */
.bg-lux-noir-light {
  background-color: hsl(0 0% 8%);
}

.bg-lux-noir-medium {
  background-color: hsl(0 0% 12%);
}

.bg-lux-noir\/80 {
  background-color: hsl(0 0% 4% / 0.8);
}

.border-lux-noir-soft {
  border-color: hsl(0 0% 16%);
}

.border-lux-gold\/10 {
  border-color: hsl(43 74% 49% / 0.1);
}

.border-lux-gold\/20 {
  border-color: hsl(43 74% 49% / 0.2);
}

.border-lux-gold\/30 {
  border-color: hsl(43 74% 49% / 0.3);
}

.hover\:border-lux-gold\/30:hover {
  border-color: hsl(43 74% 49% / 0.3);
}

.hover\:border-lux-gold\/40:hover {
  border-color: hsl(43 74% 49% / 0.4);
}

.hover\:bg-lux-noir-soft:hover {
  background-color: hsl(0 0% 16%);
}

.text-lux-cream {
  color: hsl(45 30% 96%);
}

.text-lux-silver {
  color: hsl(220 10% 70%);
}

.text-lux-silver\/30 {
  color: hsl(220 10% 70% / 0.3);
}

.text-lux-silver\/40 {
  color: hsl(220 10% 70% / 0.4);
}

.text-lux-silver\/50 {
  color: hsl(220 10% 70% / 0.5);
}

.text-lux-silver\/70 {
  color: hsl(220 10% 70% / 0.7);
}

.text-lux-gold {
  color: hsl(43 74% 49%);
}

.text-lux-gold\/60 {
  color: hsl(43 74% 49% / 0.6);
}

.group-hover\:text-lux-gold:hover {
  color: hsl(43 74% 49%);
}

.bg-lux-noir {
  background-color: hsl(0 0% 4%);
}

/* Font display */
.font-display {
  font-family: 'Cormorant Garamond', Georgia, serif;
}
</style>
