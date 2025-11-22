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

    <!-- ヘッダー -->
    <div class="mb-8">
      <h1 class="text-3xl font-bold text-gray-900">新規管理者登録</h1>
    </div>

    <!-- エラー表示エリア -->
    <div v-if="formError" class="mb-6">
      <Alert variant="destructive">
        <div class="flex items-start">
          <svg class="w-5 h-5 mr-2 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
          <span>{{ formError }}</span>
        </div>
      </Alert>
    </div>

    <!-- 登録フォーム -->
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

        <!-- 認証情報セクション -->
        <div>
          <h2 class="text-lg font-semibold text-gray-900 mb-1">認証情報</h2>
          <Separator class="mb-6" />

          <!-- パスワード -->
          <div class="mb-6">
            <Label for="password" class="mb-2">
              パスワード <span class="text-red-500">*</span>
            </Label>
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
            <Label for="password_confirm" class="mb-2">
              パスワード（確認） <span class="text-red-500">*</span>
            </Label>
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
            <p class="mt-1 text-sm text-gray-500">確認のため、もう一度入力してください。</p>
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
            <RadioGroup v-model="formData.role" name="role" class="space-y-4">
              <RadioGroupItem id="role-system-admin" value="system_admin">
                <div>
                  <div class="font-medium">システム管理者（system_admin）</div>
                  <div class="text-sm text-gray-500">全権限。ユーザー管理、ポイント管理が可能。</div>
                </div>
              </RadioGroupItem>
              <RadioGroupItem id="role-auctioneer" value="auctioneer">
                <div>
                  <div class="font-medium">主催者（auctioneer）</div>
                  <div class="text-sm text-gray-500">オークション管理のみ。ユーザー管理は不可。</div>
                </div>
              </RadioGroupItem>
            </RadioGroup>
            <p v-if="errors.role" class="mt-2 text-sm text-red-600">{{ errors.role }}</p>
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
              登録中...
            </span>
            <span v-else>登録する</span>
          </Button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { registerAdmin } from '@/services/adminApi'
import {
  validateEmail,
  validateDisplayName,
  validatePassword,
  validatePasswordConfirm,
  validateRole,
  validateAdminRegistrationForm,
  hasNoErrors
} from '@/lib/validation'
import Input from '@/components/ui/Input.vue'
import Label from '@/components/ui/Label.vue'
import Button from '@/components/ui/Button.vue'
import Alert from '@/components/ui/Alert.vue'
import Separator from '@/components/ui/Separator.vue'
import RadioGroup from '@/components/ui/RadioGroup.vue'
import RadioGroupItem from '@/components/ui/RadioGroupItem.vue'

const router = useRouter()

// フォームデータ
const formData = reactive({
  email: '',
  display_name: '',
  password: '',
  password_confirm: '',
  role: ''
})

// エラー状態
const errors = reactive({
  email: '',
  display_name: '',
  password: '',
  password_confirm: '',
  role: ''
})

// フォーム全体のエラー
const formError = ref('')

// ローディング状態
const loading = ref(false)

// フィールド単位のバリデーション
function validateField(fieldName) {
  // エラーをクリア
  errors[fieldName] = ''

  switch (fieldName) {
    case 'email':
      errors.email = validateEmail(formData.email) || ''
      break
    case 'display_name':
      errors.display_name = validateDisplayName(formData.display_name) || ''
      break
    case 'password':
      errors.password = validatePassword(formData.password) || ''
      // パスワード確認も再検証
      if (formData.password_confirm) {
        errors.password_confirm = validatePasswordConfirm(formData.password, formData.password_confirm) || ''
      }
      break
    case 'password_confirm':
      errors.password_confirm = validatePasswordConfirm(formData.password, formData.password_confirm) || ''
      break
    case 'role':
      errors.role = validateRole(formData.role) || ''
      break
  }
}

// フォーム送信処理
async function handleSubmit() {
  // フォーム全体のバリデーション
  const validationErrors = validateAdminRegistrationForm(formData)

  // エラーをリアクティブオブジェクトに反映
  Object.keys(errors).forEach(key => {
    errors[key] = validationErrors[key] || ''
  })

  // エラーがある場合は送信しない
  if (!hasNoErrors(validationErrors)) {
    formError.value = '入力内容に誤りがあります。各項目を確認してください。'
    return
  }

  // フォームエラーをクリア
  formError.value = ''
  loading.value = true

  try {
    // API呼び出し
    await registerAdmin({
      email: formData.email,
      password: formData.password,
      display_name: formData.display_name || undefined,
      role: formData.role
    })

    // 成功時: 管理者一覧にリダイレクト
    router.push({ name: 'admin-list' })
  } catch (error) {
    // エラーハンドリング
    if (error.status === 409) {
      // メールアドレス重複エラー
      errors.email = 'このメールアドレスは既に登録されています'
      formError.value = 'メールアドレスが既に使用されています。'
    } else if (error.status === 400) {
      // バリデーションエラー
      formError.value = error.message || '入力内容に誤りがあります。'
    } else if (error.status === 403) {
      // 権限エラー
      formError.value = 'この操作を行う権限がありません。'
    } else {
      // その他のエラー
      formError.value = 'サーバーエラーが発生しました。しばらくしてから再度お試しください。'
    }
  } finally {
    loading.value = false
  }
}

// キャンセル処理
function handleCancel() {
  router.push({ name: 'admin-list' })
}
</script>
