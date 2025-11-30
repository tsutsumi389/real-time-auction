<template>
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- 戻るリンク -->
    <div class="mb-6">
      <button
        @click="handleCancel"
        class="inline-flex items-center text-sm text-gray-600 hover:text-gray-900"
      >
        <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        管理者一覧に戻る
      </button>
    </div>

    <!-- ローディング表示 -->
    <div v-if="initialLoading" class="flex justify-center items-center py-12">
      <LoadingSpinner />
      <span class="ml-2 text-gray-600">読み込み中...</span>
    </div>

    <!-- 読み込みエラー表示 -->
    <div v-else-if="loadError" class="max-w-2xl">
      <Alert variant="destructive">
        <div class="flex items-start">
          <svg class="w-5 h-5 mr-2 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
          <span>{{ loadError }}</span>
        </div>
      </Alert>
      <div class="mt-4">
        <Button @click="handleCancel" variant="outline">管理者一覧に戻る</Button>
      </div>
    </div>

    <!-- 編集フォーム -->
    <template v-else>
      <!-- ヘッダー -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">管理者編集</h1>
        <p class="mt-2 text-sm text-gray-600">ID: {{ adminId }}</p>
      </div>

      <!-- エラー表示エリア -->
      <div v-if="formError" class="mb-6 max-w-2xl">
        <Alert variant="destructive">
          <div class="flex items-start">
            <svg class="w-5 h-5 mr-2 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
            </svg>
            <span>{{ formError }}</span>
          </div>
        </Alert>
      </div>

      <!-- 自己編集警告 -->
      <div v-if="isSelfEdit" class="mb-6 max-w-2xl">
        <Alert>
          <div class="flex items-start">
            <svg class="w-5 h-5 mr-2 mt-0.5 text-amber-500" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
            </svg>
            <span class="text-amber-700">自分自身のアカウントを編集しています。ロールと状態は変更できません。</span>
          </div>
        </Alert>
      </div>

      <!-- 編集フォーム -->
      <div class="max-w-2xl">
        <form @submit.prevent="handleSubmit" class="space-y-8">
          <!-- 基本情報セクション -->
          <div>
            <h2 class="text-lg font-semibold text-gray-900 mb-1">基本情報</h2>
            <Separator class="mb-6" />

            <!-- メールアドレス -->
            <div class="mb-6">
              <Label for="email" class="mb-2">
                メールアドレス <span class="text-red-500">*</span>
              </Label>
              <Input
                id="email"
                v-model="formData.email"
                type="email"
                placeholder="admin@example.com"
                :class="{ 'border-red-500': errors.email }"
                @blur="validateField('email')"
                :disabled="loading"
              />
              <p v-if="errors.email" class="mt-1 text-sm text-red-600">{{ errors.email }}</p>
            </div>

            <!-- 表示名 -->
            <div class="mb-6">
              <Label for="display_name" class="mb-2">表示名</Label>
              <Input
                id="display_name"
                v-model="formData.display_name"
                type="text"
                placeholder="システム管理者"
                :class="{ 'border-red-500': errors.display_name }"
                @blur="validateField('display_name')"
                :disabled="loading"
              />
              <p v-if="errors.display_name" class="mt-1 text-sm text-red-600">{{ errors.display_name }}</p>
              <p class="mt-1 text-sm text-gray-500">任意。未入力の場合はメールアドレスが使用されます。</p>
            </div>
          </div>

          <!-- パスワード変更セクション -->
          <div>
            <h2 class="text-lg font-semibold text-gray-900 mb-1">パスワード変更</h2>
            <Separator class="mb-6" />
            <p class="text-sm text-gray-500 mb-4">パスワードを変更する場合のみ入力してください。空欄の場合は現在のパスワードが維持されます。</p>

            <!-- パスワード -->
            <div class="mb-6">
              <Label for="password" class="mb-2">新しいパスワード</Label>
              <Input
                id="password"
                v-model="formData.password"
                type="password"
                placeholder="8文字以上"
                :class="{ 'border-red-500': errors.password }"
                @blur="validateField('password')"
                :disabled="loading"
              />
              <p v-if="errors.password" class="mt-1 text-sm text-red-600">{{ errors.password }}</p>
              <p class="mt-1 text-sm text-gray-500">8文字以上で入力してください。</p>
            </div>

            <!-- パスワード（確認） -->
            <div class="mb-6">
              <Label for="password_confirm" class="mb-2">新しいパスワード（確認）</Label>
              <Input
                id="password_confirm"
                v-model="formData.password_confirm"
                type="password"
                placeholder="確認のため再度入力"
                :class="{ 'border-red-500': errors.password_confirm }"
                @blur="validateField('password_confirm')"
                :disabled="loading"
              />
              <p v-if="errors.password_confirm" class="mt-1 text-sm text-red-600">{{ errors.password_confirm }}</p>
            </div>
          </div>

          <!-- 権限設定セクション -->
          <div>
            <h2 class="text-lg font-semibold text-gray-900 mb-1">権限設定</h2>
            <Separator class="mb-6" />

            <!-- ロール -->
            <div class="mb-6">
              <Label class="mb-3">
                ロール <span class="text-red-500">*</span>
              </Label>
              <RadioGroup v-model="formData.role" name="role" class="space-y-4" :disabled="isSelfEdit">
                <RadioGroupItem id="role-system-admin" value="system_admin" :disabled="isSelfEdit">
                  <div>
                    <div class="font-medium" :class="{ 'text-gray-400': isSelfEdit }">システム管理者（system_admin）</div>
                    <div class="text-sm" :class="isSelfEdit ? 'text-gray-400' : 'text-gray-500'">全権限。ユーザー管理、ポイント管理が可能。</div>
                  </div>
                </RadioGroupItem>
                <RadioGroupItem id="role-auctioneer" value="auctioneer" :disabled="isSelfEdit">
                  <div>
                    <div class="font-medium" :class="{ 'text-gray-400': isSelfEdit }">主催者（auctioneer）</div>
                    <div class="text-sm" :class="isSelfEdit ? 'text-gray-400' : 'text-gray-500'">オークション管理のみ。ユーザー管理は不可。</div>
                  </div>
                </RadioGroupItem>
              </RadioGroup>
              <p v-if="errors.role" class="mt-2 text-sm text-red-600">{{ errors.role }}</p>
              <p v-if="isSelfEdit" class="mt-2 text-sm text-amber-600">自分自身のロールは変更できません。</p>
            </div>
          </div>

          <!-- アカウント状態セクション -->
          <div>
            <h2 class="text-lg font-semibold text-gray-900 mb-1">アカウント状態</h2>
            <Separator class="mb-6" />

            <!-- 状態 -->
            <div class="mb-6">
              <Label class="mb-3">
                状態 <span class="text-red-500">*</span>
              </Label>
              <RadioGroup v-model="formData.status" name="status" class="space-y-4" :disabled="isSelfEdit">
                <RadioGroupItem id="status-active" value="active" :disabled="isSelfEdit">
                  <div>
                    <div class="font-medium" :class="{ 'text-gray-400': isSelfEdit }">有効（active）</div>
                    <div class="text-sm" :class="isSelfEdit ? 'text-gray-400' : 'text-gray-500'">通常のログインとすべての操作が可能です。</div>
                  </div>
                </RadioGroupItem>
                <RadioGroupItem id="status-suspended" value="suspended" :disabled="isSelfEdit">
                  <div>
                    <div class="font-medium" :class="{ 'text-gray-400': isSelfEdit }">停止（suspended）</div>
                    <div class="text-sm" :class="isSelfEdit ? 'text-gray-400' : 'text-gray-500'">ログインができなくなります。再度有効にすることで復帰できます。</div>
                  </div>
                </RadioGroupItem>
              </RadioGroup>
              <p v-if="errors.status" class="mt-2 text-sm text-red-600">{{ errors.status }}</p>
              <p v-if="isSelfEdit" class="mt-2 text-sm text-amber-600">自分自身の状態は変更できません。</p>
            </div>
          </div>

          <!-- ボタンエリア -->
          <div class="flex items-center justify-end space-x-4 pt-6 border-t">
            <Button
              type="button"
              variant="outline"
              @click="handleCancel"
              :disabled="loading"
            >
              キャンセル
            </Button>
            <Button
              type="submit"
              :disabled="loading"
              class="min-w-[120px]"
            >
              <span v-if="loading" class="flex items-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                更新中...
              </span>
              <span v-else>更新する</span>
            </Button>
          </div>
        </form>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAdminStore } from '@/stores/admin'
import { useAuthStore } from '@/stores/auth'
import {
  validateEmail,
  validateDisplayName,
  validateOptionalPassword,
  validateOptionalPasswordConfirm,
  validateRole,
  validateStatus,
  validateAdminEditForm,
  hasNoErrors
} from '@/lib/validation'
import Input from '@/components/ui/Input.vue'
import Label from '@/components/ui/Label.vue'
import Button from '@/components/ui/Button.vue'
import Alert from '@/components/ui/Alert.vue'
import Separator from '@/components/ui/Separator.vue'
import RadioGroup from '@/components/ui/RadioGroup.vue'
import RadioGroupItem from '@/components/ui/RadioGroupItem.vue'
import LoadingSpinner from '@/components/ui/LoadingSpinner.vue'

const route = useRoute()
const router = useRouter()
const adminStore = useAdminStore()
const authStore = useAuthStore()

// Admin ID from route params
const adminId = computed(() => Number(route.params.id))

// Check if editing own account
const isSelfEdit = computed(() => {
  return authStore.user?.adminId === adminId.value
})

// Form data
const formData = reactive({
  email: '',
  display_name: '',
  password: '',
  password_confirm: '',
  role: '',
  status: ''
})

// Error states
const errors = reactive({
  email: '',
  display_name: '',
  password: '',
  password_confirm: '',
  role: '',
  status: ''
})

// Form-level error
const formError = ref('')

// Loading states
const initialLoading = ref(true)
const loading = ref(false)
const loadError = ref('')

// Field validation
function validateField(fieldName) {
  // Clear error
  errors[fieldName] = ''

  switch (fieldName) {
    case 'email':
      errors.email = validateEmail(formData.email) || ''
      break
    case 'display_name':
      errors.display_name = validateDisplayName(formData.display_name) || ''
      break
    case 'password':
      errors.password = validateOptionalPassword(formData.password) || ''
      // Re-validate password confirmation
      if (formData.password_confirm || formData.password) {
        errors.password_confirm = validateOptionalPasswordConfirm(formData.password, formData.password_confirm) || ''
      }
      break
    case 'password_confirm':
      errors.password_confirm = validateOptionalPasswordConfirm(formData.password, formData.password_confirm) || ''
      break
    case 'role':
      if (!isSelfEdit.value) {
        errors.role = validateRole(formData.role) || ''
      }
      break
    case 'status':
      if (!isSelfEdit.value) {
        errors.status = validateStatus(formData.status) || ''
      }
      break
  }
}

// Load admin data
async function loadAdmin() {
  initialLoading.value = true
  loadError.value = ''

  const admin = await adminStore.fetchAdmin(adminId.value)

  if (!admin) {
    loadError.value = adminStore.error || '管理者が見つかりません'
    initialLoading.value = false
    return
  }

  // Populate form with admin data
  formData.email = admin.email || ''
  formData.display_name = admin.display_name || ''
  formData.role = admin.role || ''
  formData.status = admin.status || ''
  formData.password = ''
  formData.password_confirm = ''

  initialLoading.value = false
}

// Form submission
async function handleSubmit() {
  // Form validation
  const validationErrors = validateAdminEditForm(formData)

  // If self-edit, remove role and status errors (they are not editable)
  if (isSelfEdit.value) {
    delete validationErrors.role
    delete validationErrors.status
  }

  // Update error reactive object
  Object.keys(errors).forEach(key => {
    errors[key] = validationErrors[key] || ''
  })

  // If there are errors, stop submission
  if (!hasNoErrors(validationErrors)) {
    formError.value = '入力内容に誤りがあります。各項目を確認してください。'
    return
  }

  // Clear form error
  formError.value = ''
  loading.value = true

  try {
    // Build update request
    const updateData = {
      email: formData.email,
      display_name: formData.display_name || '',
      role: formData.role,
      status: formData.status
    }

    // Include password only if provided
    if (formData.password) {
      updateData.password = formData.password
    }

    // Call update API
    const result = await adminStore.updateAdmin(adminId.value, updateData)

    if (result.success) {
      // Success: redirect to admin list
      router.push({ name: 'admin-list' })
    } else {
      // Handle errors
      const err = result.error
      const status = err.response?.status
      const errorMessage = err.response?.data?.error || err.message

      if (status === 409) {
        // Email already exists
        errors.email = 'このメールアドレスは既に登録されています'
        formError.value = 'メールアドレスが既に使用されています。'
      } else if (status === 403) {
        // Permission errors
        if (errorMessage.includes('own role')) {
          formError.value = '自分自身のロールは変更できません。'
        } else if (errorMessage.includes('suspend own')) {
          formError.value = '自分自身のアカウントを停止することはできません。'
        } else if (errorMessage.includes('last system admin')) {
          formError.value = '最後のシステム管理者を降格または停止することはできません。'
        } else {
          formError.value = 'この操作を行う権限がありません。'
        }
      } else if (status === 400) {
        // Validation error
        formError.value = errorMessage || '入力内容に誤りがあります。'
      } else if (status === 404) {
        formError.value = '管理者が見つかりません。'
      } else {
        // Other errors
        formError.value = 'サーバーエラーが発生しました。しばらくしてから再度お試しください。'
      }
    }
  } finally {
    loading.value = false
  }
}

// Cancel handler
function handleCancel() {
  router.push({ name: 'admin-list' })
}

// Lifecycle
onMounted(() => {
  loadAdmin()
})

onUnmounted(() => {
  adminStore.clearCurrentAdmin()
})
</script>
