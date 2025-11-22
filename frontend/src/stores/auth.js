/**
 * Authentication Store
 * 認証状態の管理とアクション
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as apiLogin, logout as apiLogout, getCurrentUser } from '@/services/auth'
import { saveToken, removeToken, isTokenValid, getUserFromToken } from '@/services/token'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const loading = ref(false)
  const error = ref(null)

  // Getters
  const isAuthenticated = computed(() => {
    return user.value !== null && isTokenValid()
  })

  const isSystemAdmin = computed(() => {
    return user.value?.role === 'system_admin'
  })

  const isAuctioneer = computed(() => {
    return user.value?.role === 'auctioneer'
  })

  // Actions

  /**
   * ログイン処理
   * @param {string} email - メールアドレス
   * @param {string} password - パスワード
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function login(email, password) {
    loading.value = true
    error.value = null

    try {
      // ログインAPI呼び出し
      const response = await apiLogin(email, password)

      // トークンを保存
      saveToken(response.token)

      // ユーザー情報を設定
      user.value = {
        adminId: response.user.id,
        email: response.user.email,
        role: response.user.role,
      }

      loading.value = false
      return true
    } catch (err) {
      loading.value = false
      error.value = err.message || 'ログインに失敗しました'
      return false
    }
  }

  /**
   * ログアウト処理
   */
  async function logout() {
    try {
      // ログアウトAPI呼び出し（ベストエフォート）
      await apiLogout()
    } catch (err) {
      // エラーが発生してもローカルの状態はクリアする
      console.error('Logout API error:', err)
    } finally {
      // トークンとユーザー情報を削除
      removeToken()
      user.value = null
      error.value = null
    }
  }

  /**
   * トークンからユーザー情報を復元
   * ページリロード時などに使用
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function restoreUser() {
    // トークンの有効性チェック
    if (!isTokenValid()) {
      user.value = null
      return false
    }

    loading.value = true

    try {
      // まずトークンからユーザー情報を取得して即座に設定
      const tokenUser = getUserFromToken()

      if (tokenUser) {
        user.value = tokenUser
      } else {
        // トークンが不正な場合はクリア
        removeToken()
        user.value = null
        loading.value = false
        return false
      }

      // バックグラウンドでAPIによる検証を実行
      // ネットワークエラーの場合でもトークンが有効であれば継続利用
      try {
        const response = await getCurrentUser()
        user.value = {
          adminId: response.user.id,
          email: response.user.email,
          role: response.user.role,
        }
      } catch (apiError) {
        // APIエラー（401/403）の場合のみトークンを削除
        if (apiError.response && (apiError.response.status === 401 || apiError.response.status === 403)) {
          removeToken()
          user.value = null
          loading.value = false
          return false
        }
        // その他のエラー（ネットワークエラーなど）はトークンの情報を継続使用
        console.warn('Failed to verify user with API, using token data:', apiError.message)
      }

      loading.value = false
      return true
    } catch (err) {
      // トークンデコードエラーなど予期しないエラー
      console.error('Unexpected error during restore:', err)
      removeToken()
      user.value = null
      loading.value = false
      return false
    }
  }

  /**
   * エラーをクリア
   */
  function clearError() {
    error.value = null
  }

  return {
    // State
    user,
    loading,
    error,
    // Getters
    isAuthenticated,
    isSystemAdmin,
    isAuctioneer,
    // Actions
    login,
    logout,
    restoreUser,
    clearError,
  }
})
