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
      saveToken(response.token, response.expires_in)

      // ユーザー情報を設定
      user.value = {
        adminId: response.admin.id,
        email: response.admin.email,
        role: response.admin.role,
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
      // トークンからユーザー情報を取得
      const tokenUser = getUserFromToken()
      if (tokenUser) {
        user.value = tokenUser
      }

      // APIでユーザー情報を検証
      const response = await getCurrentUser()
      user.value = {
        adminId: response.admin.id,
        email: response.admin.email,
        role: response.admin.role,
      }

      loading.value = false
      return true
    } catch (err) {
      // トークンが無効な場合はクリア
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
