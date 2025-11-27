/**
 * Bidder Authentication Store
 * 入札者認証状態の管理とアクション
 * 管理者認証とは完全に分離
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { bidderLogin as apiBidderLogin, bidderLogout as apiBidderLogout } from '@/services/bidderAuth'
import { saveToken, removeToken, isTokenValid, getUserFromToken } from '@/services/token'

// 入札者用の役割識別子
const ROLE = 'bidder'

// ユーザー情報のlocalStorageキー
const USER_KEY = 'bidder_user'

/**
 * ユーザー情報をlocalStorageに保存
 */
function saveBidderUser(user) {
  try {
    localStorage.setItem(USER_KEY, JSON.stringify(user))
  } catch (error) {
    console.error('Failed to save bidder user:', error)
  }
}

/**
 * ユーザー情報をlocalStorageから取得
 */
function getBidderUser() {
  try {
    const userJson = localStorage.getItem(USER_KEY)
    if (!userJson) return null
    return JSON.parse(userJson)
  } catch (error) {
    console.error('Failed to get bidder user:', error)
    return null
  }
}

/**
 * ユーザー情報をlocalStorageから削除
 */
function removeBidderUser() {
  try {
    localStorage.removeItem(USER_KEY)
  } catch (error) {
    console.error('Failed to remove bidder user:', error)
  }
}

/**
 * トークンとユーザー情報をすべて削除
 */
function clearBidderAuth() {
  removeToken(ROLE)
  removeBidderUser()
}

export const useBidderAuthStore = defineStore('bidderAuth', () => {
  // State
  const user = ref(null)
  const loading = ref(false)
  const error = ref(null)

  // Getters
  const isAuthenticated = computed(() => {
    return user.value !== null && isTokenValid(ROLE)
  })

  const hasPoints = computed(() => {
    return user.value?.points && user.value.points.available > 0
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
      const response = await apiBidderLogin(email, password)

      // トークンを保存（入札者用）
      saveToken(response.token, ROLE)

      // ユーザー情報を設定
      const userData = {
        bidderId: response.user.id,
        email: response.user.email,
        displayName: response.user.display_name,
        userType: response.user.user_type,
        points: response.user.points || {
          total: 0,
          available: 0,
          reserved: 0,
        },
      }

      user.value = userData

      // ユーザー情報をlocalStorageに保存
      saveBidderUser(userData)

      loading.value = false
      return true
    } catch (err) {
      loading.value = false
      // エラーメッセージの処理
      let errorMessage = 'ログインに失敗しました'

      if (err.status === 401) {
        errorMessage = 'メールアドレスまたはパスワードが正しくありません'
      } else if (err.status === 403) {
        if (err.message.includes('suspended')) {
          errorMessage = 'アカウントが停止されています'
        } else if (err.message.includes('deleted')) {
          errorMessage = 'アカウントが削除されています'
        } else {
          errorMessage = err.message
        }
      } else if (err.status === 0) {
        errorMessage = 'サーバーに接続できません'
      } else if (err.message) {
        errorMessage = err.message
      }

      error.value = errorMessage
      return false
    }
  }

  /**
   * ログアウト処理
   */
  async function logout() {
    try {
      // ログアウトAPI呼び出し（ベストエフォート）
      await apiBidderLogout()
    } catch (err) {
      // エラーが発生してもローカルの状態はクリアする
      console.error('Logout API error:', err)
    } finally {
      // トークンとユーザー情報を削除
      clearBidderAuth()
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
    // トークンの有効性チェック（入札者用）
    if (!isTokenValid(ROLE)) {
      user.value = null
      return false
    }

    loading.value = true

    try {
      // まずlocalStorageからユーザー情報を取得
      let storedUser = getBidderUser()

      // localStorageにユーザー情報がない場合はトークンから復元
      if (!storedUser) {
        const tokenUser = getUserFromToken(ROLE)
        if (tokenUser) {
          storedUser = {
            bidderId: tokenUser.bidderId,
            email: tokenUser.email,
            displayName: tokenUser.displayName,
            userType: tokenUser.userType,
            points: {
              total: 0,
              available: 0,
              reserved: 0,
            },
          }
        }
      }

      if (storedUser) {
        user.value = storedUser
        loading.value = false
        return true
      } else {
        // トークンが不正な場合はクリア
        clearBidderAuth()
        user.value = null
        loading.value = false
        return false
      }

      // 将来的には、バックグラウンドでAPIによる検証を実行
      // 現在は /bidder/me エンドポイントが未実装のためスキップ
    } catch (err) {
      // トークンデコードエラーなど予期しないエラー
      console.error('Unexpected error during restore:', err)
      clearBidderAuth()
      user.value = null
      loading.value = false
      return false
    }
  }

  /**
   * ポイント情報を更新
   * @param {object} points - 新しいポイント情報
   */
  function updatePoints(points) {
    if (user.value) {
      user.value.points = points
      // localStorageも更新
      saveBidderUser(user.value)
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
    hasPoints,
    // Actions
    login,
    logout,
    restoreUser,
    updatePoints,
    clearError,
  }
})
