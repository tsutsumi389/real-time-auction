/**
 * Authentication API Service
 * 認証関連のAPIエンドポイントを提供
 */
import apiClient from './api'

/**
 * 管理者ログイン
 * @param {string} email - メールアドレス
 * @param {string} password - パスワード
 * @returns {Promise<object>} レスポンス（token, admin）
 */
export async function login(email, password) {
  const response = await apiClient.post('/admin/login', {
    email,
    password,
  })
  return response.data
}

/**
 * ログアウト
 * @returns {Promise<object>} レスポンス
 */
export async function logout() {
  const response = await apiClient.post('/admin/logout')
  return response.data
}

/**
 * 現在のユーザー情報を取得
 * @returns {Promise<object>} レスポンス（admin）
 */
export async function getCurrentUser() {
  const response = await apiClient.get('/admin/me')
  return response.data
}

/**
 * トークンの有効性を検証
 * @returns {Promise<boolean>} トークンが有効な場合true
 */
export async function validateToken() {
  try {
    await getCurrentUser()
    return true
  } catch (error) {
    return false
  }
}
