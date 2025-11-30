<template>
  <div class="confetti-container pointer-events-none fixed inset-0 z-[100] overflow-hidden">
    <div
      v-for="particle in particles"
      :key="particle.id"
      class="confetti-particle absolute"
      :style="getParticleStyle(particle)"
    ></div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  active: {
    type: Boolean,
    default: false
  },
  duration: {
    type: Number,
    default: 3000
  }
})

const particles = ref([])
let animationFrameId = null
let lastTime = 0

const colors = ['#FFC700', '#FF0000', '#2E3192', '#41BBC7', '#73F018']

function createParticle() {
  return {
    id: Math.random(),
    x: Math.random() * 100, // vw
    y: -10, // start above screen
    size: Math.random() * 10 + 5,
    color: colors[Math.floor(Math.random() * colors.length)],
    speedY: Math.random() * 3 + 2,
    speedX: Math.random() * 2 - 1,
    rotation: Math.random() * 360,
    rotationSpeed: Math.random() * 10 - 5,
    opacity: 1
  }
}

function updateParticles(timestamp) {
  if (!lastTime) lastTime = timestamp
  const deltaTime = timestamp - lastTime
  lastTime = timestamp

  // Add new particles if active
  if (props.active && particles.value.length < 150) {
    if (Math.random() > 0.5) { // Control density
      particles.value.push(createParticle())
    }
  }

  // Update existing particles
  particles.value = particles.value
    .map(p => ({
      ...p,
      y: p.y + p.speedY * (deltaTime / 16), // Normalize speed
      x: p.x + p.speedX * (deltaTime / 16),
      rotation: p.rotation + p.rotationSpeed,
      speedY: p.speedY + 0.05, // Gravity
      opacity: p.y > 80 ? p.opacity - 0.02 : p.opacity // Fade out near bottom
    }))
    .filter(p => p.y < 110 && p.opacity > 0) // Remove off-screen

  if (props.active || particles.value.length > 0) {
    animationFrameId = requestAnimationFrame(updateParticles)
  }
}

function getParticleStyle(particle) {
  return {
    left: `${particle.x}vw`,
    top: `${particle.y}vh`,
    width: `${particle.size}px`,
    height: `${particle.size}px`,
    backgroundColor: particle.color,
    transform: `rotate(${particle.rotation}deg)`,
    opacity: particle.opacity
  }
}

// Watch for activation
import { watch } from 'vue'
watch(() => props.active, (newVal) => {
  if (newVal) {
    lastTime = 0
    if (!animationFrameId) {
      animationFrameId = requestAnimationFrame(updateParticles)
    }
  }
})

onUnmounted(() => {
  if (animationFrameId) {
    cancelAnimationFrame(animationFrameId)
  }
})
</script>

<style scoped>
.confetti-particle {
  border-radius: 2px;
}
</style>
