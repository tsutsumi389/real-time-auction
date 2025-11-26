/**
 * Bidder Authentication API Service
 * 入札者認証関連のAPIエンドポイントを提供
 */
import bidderApiClient from './bidderApiClient'

/**
 * 入札者ログイン
 * @param {string} email - メールアドレス
 * @param {string} password - パスワード
 * @returns {Promise<object>} レスポンス（token, user）
 */
export async function bidderLogin(email, password) {
  const response = await bidderApiClient.post('/auth/bidder/login', {
    email,
    password,
  })
  return response.data
}

/**
 * 入札者ログアウト（将来実装）
 * @returns {Promise<object>} レスポンス
 */
export async function bidderLogout() {
  // 現在はバックエンドにログアウトエンドポイントがないため
  // フロントエンドでトークンを削除するだけ
  return Promise.resolve({})
}

/**
 * 現在の入札者情報を取得（将来実装）
 * @returns {Promise<object>} レスポンス（user）
 */
export async function getCurrentBidder() {
  const response = await bidderApiClient.get('/bidder/me')
  return response.data
}

/**
 * 入札者トークンの有効性を検証（将来実装）
 * @returns {Promise<boolean>} トークンが有効な場合true
 */
export async function validateBidderToken() {
  try {
    await getCurrentBidder()
    return true
  } catch (error) {
    return false
  }
}
