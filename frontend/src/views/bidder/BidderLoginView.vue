<template>
  <div class="theme-luxury min-h-screen relative overflow-hidden flex items-center justify-center px-4 sm:px-6 lg:px-8">
    <!-- Ambient Background -->
    <div class="absolute inset-0 bg-lux-noir">
      <!-- Gradient Orbs -->
      <div class="absolute top-0 left-1/4 w-[600px] h-[600px] bg-lux-gold/5 rounded-full blur-[120px] animate-pulse-slow"></div>
      <div class="absolute bottom-0 right-1/4 w-[500px] h-[500px] bg-lux-gold/3 rounded-full blur-[100px] animate-pulse-slow" style="animation-delay: 1s;"></div>

      <!-- Subtle Grid Pattern -->
      <div class="absolute inset-0 opacity-[0.02]" style="background-image: linear-gradient(rgba(212,175,55,0.3) 1px, transparent 1px), linear-gradient(90deg, rgba(212,175,55,0.3) 1px, transparent 1px); background-size: 60px 60px;"></div>

      <!-- Noise Texture -->
      <div class="absolute inset-0 lux-noise"></div>
    </div>

    <!-- Login Card -->
    <div class="relative z-10 w-full max-w-md lux-fade-in">
      <!-- Decorative Top Line -->
      <div class="h-px w-full bg-gradient-to-r from-transparent via-lux-gold/50 to-transparent mb-8"></div>

      <div class="lux-glass-strong rounded-2xl p-8 sm:p-10 space-y-8 shadow-2xl">
        <!-- Subtle Corner Accents -->
        <div class="absolute top-0 left-0 w-16 h-16 border-l border-t border-lux-gold/20 rounded-tl-2xl pointer-events-none"></div>
        <div class="absolute bottom-0 right-0 w-16 h-16 border-r border-b border-lux-gold/20 rounded-br-2xl pointer-events-none"></div>

        <!-- Header -->
        <div class="text-center space-y-4">
          <!-- Logo/Brand Mark -->
          <div class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-lux-noir-light border border-lux-gold/30 mb-2">
            <svg class="w-8 h-8 text-lux-gold" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>

          <div>
            <h1 class="font-display text-3xl sm:text-4xl font-light tracking-wide text-lux-cream">
              Welcome Back
            </h1>
            <p class="mt-3 text-sm text-lux-silver tracking-wide">
              プレミアムオークションへのログイン
            </p>
          </div>
        </div>

        <!-- Global Error Alert -->
        <div
          v-if="globalError"
          class="p-4 rounded-lg bg-red-950/50 border border-red-500/30 text-red-300 text-sm flex items-center gap-3"
        >
          <svg class="w-5 h-5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          {{ globalError }}
        </div>

        <!-- Login Form -->
        <form @submit.prevent="handleSubmit" class="space-y-6">
          <!-- Email Field -->
          <div class="space-y-2">
            <label for="email" class="block text-sm font-medium text-lux-pearl tracking-wide">
              メールアドレス
            </label>
            <div class="relative">
              <input
                id="email"
                type="email"
                v-model="formData.email"
                placeholder="your@email.com"
                class="w-full px-4 py-3.5 rounded-lg lux-input text-lux-cream placeholder:text-lux-silver/40 focus:ring-2 focus:ring-lux-gold/20"
                :class="{ 'border-red-500/50 focus:border-red-500/50': errors.email }"
                @blur="validateEmail"
                @keypress.enter="handleSubmit"
                :disabled="loading"
              />
              <div class="absolute inset-y-0 right-0 flex items-center pr-4 pointer-events-none">
                <svg class="w-5 h-5 text-lux-silver/50" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                </svg>
              </div>
            </div>
            <p v-if="errors.email" class="text-sm text-red-400 flex items-center gap-1.5">
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
              </svg>
              {{ errors.email }}
            </p>
          </div>

          <!-- Password Field -->
          <div class="space-y-2">
            <label for="password" class="block text-sm font-medium text-lux-pearl tracking-wide">
              パスワード
            </label>
            <div class="relative">
              <input
                id="password"
                :type="showPassword ? 'text' : 'password'"
                v-model="formData.password"
                placeholder="••••••••"
                class="w-full px-4 py-3.5 rounded-lg lux-input text-lux-cream placeholder:text-lux-silver/40 focus:ring-2 focus:ring-lux-gold/20"
                :class="{ 'border-red-500/50 focus:border-red-500/50': errors.password }"
                @blur="validatePassword"
                @keypress.enter="handleSubmit"
                :disabled="loading"
              />
              <button
                type="button"
                @click="showPassword = !showPassword"
                class="absolute inset-y-0 right-0 flex items-center pr-4 text-lux-silver/50 hover:text-lux-gold transition-colors"
              >
                <svg v-if="!showPassword" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                </svg>
                <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.88 9.88l-3.29-3.29m7.532 7.532l3.29 3.29M3 3l3.59 3.59m0 0A9.953 9.953 0 0112 5c4.478 0 8.268 2.943 9.543 7a10.025 10.025 0 01-4.132 5.411m0 0L21 21" />
                </svg>
              </button>
            </div>
            <p v-if="errors.password" class="text-sm text-red-400 flex items-center gap-1.5">
              <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
              </svg>
              {{ errors.password }}
            </p>
          </div>

          <!-- Submit Button -->
          <button
            type="submit"
            class="w-full py-4 rounded-lg lux-btn-gold text-base tracking-widest"
            :disabled="loading"
          >
            <span v-if="loading" class="inline-flex items-center justify-center">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Authenticating...
            </span>
            <span v-else>Sign In</span>
          </button>
        </form>

        <!-- Divider -->
        <div class="relative">
          <div class="absolute inset-0 flex items-center">
            <div class="w-full border-t border-lux-noir-soft"></div>
          </div>
          <div class="relative flex justify-center text-sm">
            <span class="px-4 bg-transparent text-lux-silver/60">または</span>
          </div>
        </div>

        <!-- Register Link -->
        <div class="text-center">
          <p class="text-sm text-lux-silver">
            アカウントをお持ちでない方
          </p>
          <a
            href="#"
            class="inline-flex items-center mt-2 text-sm font-medium text-lux-gold hover:text-lux-gold-light transition-colors group"
            @click.prevent="handleRegisterClick"
          >
            新規登録はこちら
            <svg class="w-4 h-4 ml-1 transform group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 8l4 4m0 0l-4 4m4-4H3" />
            </svg>
          </a>
        </div>
      </div>

      <!-- Decorative Bottom Line -->
      <div class="h-px w-full bg-gradient-to-r from-transparent via-lux-gold/30 to-transparent mt-8"></div>

      <!-- Footer Text -->
      <p class="text-center mt-6 text-xs text-lux-silver/40 tracking-wider">
        PREMIUM AUCTION EXPERIENCE
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useBidderAuthStore } from '@/stores/bidderAuthStore'

const router = useRouter()
const bidderAuthStore = useBidderAuthStore()

// Form data
const formData = reactive({
  email: '',
  password: ''
})

// Error states
const errors = reactive({
  email: '',
  password: ''
})

const globalError = ref('')
const loading = ref(false)
const showPassword = ref(false)

// Email validation regex
const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

// Apply theme class to body
onMounted(() => {
  document.body.classList.add('theme-luxury')
})

onUnmounted(() => {
  document.body.classList.remove('theme-luxury')
})

/**
 * Validate email field
 */
function validateEmail() {
  errors.email = ''

  if (!formData.email || formData.email.trim() === '') {
    errors.email = 'メールアドレスを入力してください'
    return false
  }

  if (!EMAIL_REGEX.test(formData.email)) {
    errors.email = 'メールアドレスの形式が正しくありません'
    return false
  }

  return true
}

/**
 * Validate password field
 */
function validatePassword() {
  errors.password = ''

  if (!formData.password || formData.password.trim() === '') {
    errors.password = 'パスワードを入力してください'
    return false
  }

  if (formData.password.length < 8) {
    errors.password = 'パスワードは8文字以上で入力してください'
    return false
  }

  return true
}

/**
 * Validate entire form
 */
function validateForm() {
  const isEmailValid = validateEmail()
  const isPasswordValid = validatePassword()

  return isEmailValid && isPasswordValid
}

/**
 * Handle form submission
 */
async function handleSubmit() {
  // Clear previous errors
  globalError.value = ''
  bidderAuthStore.clearError()

  // Validate form
  if (!validateForm()) {
    return
  }

  loading.value = true

  try {
    const success = await bidderAuthStore.login(formData.email, formData.password)

    if (success) {
      // Get redirect path from query or default to auction list
      const redirect = router.currentRoute.value.query.redirect || '/auctions'
      router.push(redirect)
    } else {
      // Display error from store
      globalError.value = bidderAuthStore.error || 'ログインに失敗しました'
      // Clear password field on error (security best practice)
      formData.password = ''
    }
  } catch (error) {
    // Handle unexpected errors
    globalError.value = 'サーバーに接続できません'
    console.error('Bidder login error:', error)
    // Clear password field on error
    formData.password = ''
  } finally {
    loading.value = false
  }
}

/**
 * Handle register link click
 */
function handleRegisterClick() {
  // 将来実装: 新規登録画面へ遷移
  alert('新規登録機能は現在実装中です')
}
</script>

<style scoped>
@keyframes pulse-slow {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.5; }
}

.animate-pulse-slow {
  animation: pulse-slow 8s ease-in-out infinite;
}
</style>
