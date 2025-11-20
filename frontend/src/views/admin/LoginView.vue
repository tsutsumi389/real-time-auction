<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-gray-50 to-gray-100 px-4 sm:px-6 lg:px-8">
    <Card class="w-full max-w-md p-6 sm:p-8 space-y-6">
      <!-- Header -->
      <div class="text-center space-y-2">
        <h1 class="text-2xl font-bold tracking-tight">管理者ログイン</h1>
        <p class="text-sm text-muted-foreground">
          システム管理者またはオークション主催者としてログインしてください
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
            placeholder="admin@example.com"
            :class="{ 'border-destructive': errors.email }"
            @blur="validateEmail"
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
    </Card>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Card from '@/components/ui/Card.vue'
import Input from '@/components/ui/Input.vue'
import Label from '@/components/ui/Label.vue'
import Button from '@/components/ui/Button.vue'
import Alert from '@/components/ui/Alert.vue'

const router = useRouter()
const authStore = useAuthStore()

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

  if (!formData.email) {
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

  if (!formData.password) {
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
  authStore.clearError()

  // Validate form
  if (!validateForm()) {
    return
  }

  loading.value = true

  try {
    const success = await authStore.login(formData.email, formData.password)

    if (success) {
      // Get redirect path from query or default to dashboard
      const redirect = router.currentRoute.value.query.redirect || '/admin/dashboard'
      router.push(redirect)
    } else {
      // Display error from store
      globalError.value = authStore.error || 'ログインに失敗しました'
    }
  } catch (error) {
    // Handle unexpected errors
    globalError.value = 'サーバーに接続できません'
    console.error('Login error:', error)
  } finally {
    loading.value = false
  }
}
</script>
