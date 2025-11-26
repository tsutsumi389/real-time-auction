<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100 px-4 sm:px-6 lg:px-8">
    <Card class="w-full max-w-md p-6 sm:p-8 space-y-6">
      <!-- Header -->
      <div class="text-center space-y-2">
        <h1 class="text-2xl font-bold tracking-tight">ログイン</h1>
        <p class="text-sm text-muted-foreground">
          入札者としてログインしてオークションに参加しましょう
        </p>
      </div>

      <!-- Global Error Alert -->
      <Alert v-if="globalError" variant="destructive" class="text-sm">
        {{ globalError }}
      </Alert>

      <!-- Login Form -->
      <form @submit.prevent="handleSubmit" class="space-y-4">
        <!-- Email Field -->
        <div class="space-y-2">
          <Label for="email">メールアドレス</Label>
          <Input
            id="email"
            type="email"
            v-model="formData.email"
            placeholder="bidder@example.com"
            :class="{ 'border-destructive': errors.email }"
            @blur="validateEmail"
            @keypress.enter="handleSubmit"
            :disabled="loading"
          />
          <p v-if="errors.email" class="text-sm text-destructive">
            {{ errors.email }}
          </p>
        </div>

        <!-- Password Field -->
        <div class="space-y-2">
          <Label for="password">パスワード</Label>
          <Input
            id="password"
            type="password"
            v-model="formData.password"
            placeholder="8文字以上"
            :class="{ 'border-destructive': errors.password }"
            @blur="validatePassword"
            @keypress.enter="handleSubmit"
            :disabled="loading"
          />
          <p v-if="errors.password" class="text-sm text-destructive">
            {{ errors.password }}
          </p>
        </div>

        <!-- Submit Button -->
        <Button
          type="submit"
          class="w-full"
          :disabled="loading"
        >
          {{ loading ? 'ログイン中...' : 'ログイン' }}
        </Button>
      </form>

      <!-- Register Link -->
      <div class="text-center text-sm">
        <span class="text-muted-foreground">アカウントをお持ちでない方</span>
        <br />
        <a href="#" class="text-primary hover:underline font-medium" @click.prevent="handleRegisterClick">
          新規登録はこちら
        </a>
      </div>
    </Card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useBidderAuthStore } from '@/stores/bidderAuthStore'
import Card from '@/components/ui/Card.vue'
import Input from '@/components/ui/Input.vue'
import Label from '@/components/ui/Label.vue'
import Button from '@/components/ui/Button.vue'
import Alert from '@/components/ui/Alert.vue'

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

// Email validation regex
const EMAIL_REGEX = /^[^\s@]+@[^\s@]+\.[^\s@]+$/

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
