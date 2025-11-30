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
        入札者一覧に戻る
      </button>
    </div>

    <!-- ローディング状態 -->
    <div v-if="initialLoading" class="flex justify-center py-12">
      <svg class="animate-spin h-8 w-8 text-blue-600" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
    </div>

    <!-- ロードエラー -->
    <div v-else-if="loadError" class="mb-6">
      <Alert variant="destructive">
        <div class="flex items-start">
          <svg class="w-5 h-5 mr-2 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
          <span>{{ loadError }}</span>
        </div>
      </Alert>
      <div class="mt-4">
        <Button @click="handleCancel">入札者一覧に戻る</Button>
      </div>
    </div>

    <!-- コンテンツ -->
    <template v-else>
      <!-- ヘッダー -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900">入札者編集</h1>
        <p class="mt-2 text-sm text-gray-600">ID: {{ bidder?.id }}</p>
      </div>

      <!-- ステータス表示 -->
      <div class="mb-6 bg-gray-50 rounded-lg p-4">
        <div class="flex flex-wrap gap-6">
          <div>
            <span class="text-sm text-gray-500">ステータス</span>
            <div class="mt-1">
              <span
                :class="[
                  'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
                  bidder?.status === 'active' ? 'bg-green-100 text-green-800' :
                  bidder?.status === 'suspended' ? 'bg-yellow-100 text-yellow-800' :
                  'bg-red-100 text-red-800'
                ]"
              >
                {{ statusLabel(bidder?.status) }}
              </span>
            </div>
          </div>
          <div>
            <span class="text-sm text-gray-500">保有ポイント</span>
            <div class="mt-1 font-medium">{{ formatPoints(bidder?.points?.total_points || 0) }}</div>
          </div>
          <div>
            <span class="text-sm text-gray-500">利用可能ポイント</span>
            <div class="mt-1 font-medium">{{ formatPoints(bidder?.points?.available_points || 0) }}</div>
          </div>
          <div>
            <span class="text-sm text-gray-500">予約済ポイント</span>
            <div class="mt-1 font-medium">{{ formatPoints(bidder?.points?.reserved_points || 0) }}</div>
          </div>
        </div>
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

      <!-- 成功メッセージ -->
      <div v-if="successMessage" class="mb-6">
        <Alert variant="success">
          <div class="flex items-start">
            <svg class="w-5 h-5 mr-2 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            <span>{{ successMessage }}</span>
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
                placeholder="bidder@example.com"
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
                placeholder="入札者01"
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

            <!-- 新しいパスワード -->
            <div class="mb-6">
              <Label for="password" class="mb-2">新しいパスワード</Label>
              <Input
                id="password"
                v-model="formData.password"
                type="password"
                placeholder="変更する場合のみ入力"
                :class="{ 'border-red-500': errors.password }"
                @blur="validateField('password')"
                :disabled="loading"
              />
              <p v-if="errors.password" class="mt-1 text-sm text-red-600">{{ errors.password }}</p>
              <p class="mt-1 text-sm text-gray-500">変更する場合のみ入力。8文字以上で入力してください。</p>
            </div>

            <!-- パスワード（確認） -->
            <div class="mb-6">
              <Label for="confirmPassword" class="mb-2">新しいパスワード（確認）</Label>
              <Input
                id="confirmPassword"
                v-model="formData.confirmPassword"
                type="password"
                placeholder="確認のため再度入力"
                :class="{ 'border-red-500': errors.confirmPassword }"
                @blur="validateField('confirmPassword')"
                :disabled="loading"
              />
              <p v-if="errors.confirmPassword" class="mt-1 text-sm text-red-600">{{ errors.confirmPassword }}</p>
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
                保存中...
              </span>
              <span v-else>保存する</span>
            </Button>
          </div>
        </form>
      </div>
    </template>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { getBidderById, updateBidder } from '@/services/bidderApi'
import {
  validateEmail,
  validateDisplayName,
  validateOptionalPassword,
  validateOptionalConfirmPassword,
  validateBidderEditForm,
  hasErrors
} from '@/utils/validation'
import Input from '@/components/ui/Input.vue'
import Label from '@/components/ui/Label.vue'
import Button from '@/components/ui/Button.vue'
import Alert from '@/components/ui/Alert.vue'
import Separator from '@/components/ui/Separator.vue'

const router = useRouter()
const route = useRoute()

// 入札者データ
const bidder = ref(null)

// フォームデータ
const formData = reactive({
  email: '',
  display_name: '',
  password: '',
  confirmPassword: ''
})

// エラー状態
const errors = reactive({
  email: '',
  display_name: '',
  password: '',
  confirmPassword: ''
})

// フォーム全体のエラー
const formError = ref('')

// 成功メッセージ
const successMessage = ref('')

// ローディング状態
const loading = ref(false)
const initialLoading = ref(true)
const loadError = ref('')

// ステータスラベル
function statusLabel(status) {
  switch (status) {
    case 'active':
      return 'アクティブ'
    case 'suspended':
      return '停止中'
    case 'deleted':
      return '削除済み'
    default:
      return status
  }
}

// ポイントのフォーマット
function formatPoints(points) {
  return new Intl.NumberFormat('ja-JP').format(points)
}

// 初期データ取得
onMounted(async () => {
  const bidderId = route.params.id
  try {
    const data = await getBidderById(bidderId)
    bidder.value = data
    // フォームデータを初期化
    formData.email = data.email
    formData.display_name = data.display_name || ''
    formData.password = ''
    formData.confirmPassword = ''
  } catch (error) {
    if (error.status === 404) {
      loadError.value = '入札者が見つかりませんでした'
    } else {
      loadError.value = 'データの取得に失敗しました'
    }
  } finally {
    initialLoading.value = false
  }
})

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
      errors.password = validateOptionalPassword(formData.password) || ''
      // パスワード確認も再検証
      if (formData.confirmPassword || formData.password) {
        errors.confirmPassword = validateOptionalConfirmPassword(formData.password, formData.confirmPassword) || ''
      }
      break
    case 'confirmPassword':
      errors.confirmPassword = validateOptionalConfirmPassword(formData.password, formData.confirmPassword) || ''
      break
  }
}

// フォーム送信処理
async function handleSubmit() {
  // フォーム全体のバリデーション
  const validationErrors = validateBidderEditForm(formData)

  // エラーをリアクティブオブジェクトに反映
  Object.keys(errors).forEach(key => {
    errors[key] = validationErrors[key] || ''
  })

  // エラーがある場合は送信しない
  if (hasErrors(validationErrors)) {
    formError.value = '入力内容に誤りがあります。各項目を確認してください。'
    return
  }

  // メッセージをクリア
  formError.value = ''
  successMessage.value = ''
  loading.value = true

  try {
    // API呼び出し用のデータ準備
    const requestData = {
      email: formData.email,
      display_name: formData.display_name || null,
      password: formData.password || null
    }

    // API呼び出し
    const updatedBidder = await updateBidder(route.params.id, requestData)

    // 成功時: データを更新して成功メッセージを表示
    bidder.value = updatedBidder
    formData.password = ''
    formData.confirmPassword = ''
    successMessage.value = '入札者情報を更新しました'
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
    } else if (error.status === 404) {
      // 入札者が見つからない
      formError.value = '入札者が見つかりませんでした。'
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
  router.push({ name: 'bidder-list' })
}
</script>
