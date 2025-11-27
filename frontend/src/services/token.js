/**
 * Token Management
 * JWTトークンのlocalStorage操作を管理（役割別に分離）
 */

// 管理者用のトークンキー
const ADMIN_TOKEN_KEY = 'admin_auth_token'
const ADMIN_TOKEN_EXPIRY_KEY = 'admin_auth_token_expiry'

// 入札者用のトークンキー
const BIDDER_TOKEN_KEY = 'bidder_auth_token'
const BIDDER_TOKEN_EXPIRY_KEY = 'bidder_auth_token_expiry'

/**
 * トークンをlocalStorageに保存
 * @param {string} token - JWTトークン
 * @param {string} role - ユーザーの役割（'admin' または 'bidder'）
 * @param {number} expiresIn - トークンの有効期限（秒）
 */
export function saveToken(token, role = 'admin', expiresIn = 86400) {
  try {
    const tokenKey = role === 'bidder' ? BIDDER_TOKEN_KEY : ADMIN_TOKEN_KEY
    const expiryKey = role === 'bidder' ? BIDDER_TOKEN_EXPIRY_KEY : ADMIN_TOKEN_EXPIRY_KEY

    localStorage.setItem(tokenKey, token)

    // 有効期限を計算してタイムスタンプで保存
    const expiryTime = Date.now() + expiresIn * 1000
    localStorage.setItem(expiryKey, expiryTime.toString())
  } catch (error) {
    console.error('Failed to save token:', error)
  }
}

/**
 * localStorageからトークンを取得
 * @param {string} role - ユーザーの役割（'admin' または 'bidder'）
 * @returns {string|null} トークン（有効期限切れの場合はnull）
 */
export function getToken(role = 'admin') {
  try {
    const tokenKey = role === 'bidder' ? BIDDER_TOKEN_KEY : ADMIN_TOKEN_KEY
    const expiryKey = role === 'bidder' ? BIDDER_TOKEN_EXPIRY_KEY : ADMIN_TOKEN_EXPIRY_KEY

    const token = localStorage.getItem(tokenKey)
    const expiry = localStorage.getItem(expiryKey)

    if (!token || !expiry) {
      return null
    }

    // 有効期限チェック
    const now = Date.now()
    const expiryTime = parseInt(expiry, 10)

    if (now >= expiryTime) {
      // 有効期限切れの場合はトークンを削除
      removeToken(role)
      return null
    }

    return token
  } catch (error) {
    console.error('Failed to get token:', error)
    return null
  }
}

/**
 * localStorageからトークンを削除
 * @param {string} role - ユーザーの役割（'admin' または 'bidder'）
 */
export function removeToken(role = 'admin') {
  try {
    const tokenKey = role === 'bidder' ? BIDDER_TOKEN_KEY : ADMIN_TOKEN_KEY
    const expiryKey = role === 'bidder' ? BIDDER_TOKEN_EXPIRY_KEY : ADMIN_TOKEN_EXPIRY_KEY

    localStorage.removeItem(tokenKey)
    localStorage.removeItem(expiryKey)
  } catch (error) {
    console.error('Failed to remove token:', error)
  }
}

/**
 * トークンの有効性をチェック
 * @param {string} role - ユーザーの役割（'admin' または 'bidder'）
 * @returns {boolean} トークンが有効な場合true
 */
export function isTokenValid(role = 'admin') {
  return getToken(role) !== null
}

/**
 * トークンのペイロードをデコード（Base64）
 * @param {string} role - ユーザーの役割（'admin' または 'bidder'）
 * @returns {object|null} デコードされたペイロード
 */
export function decodeToken(role = 'admin') {
  const token = getToken(role)
  if (!token) {
    return null
  }

  try {
    // JWTは3つの部分に分かれている: header.payload.signature
    const parts = token.split('.')
    if (parts.length !== 3) {
      return null
    }

    // ペイロード部分をデコード
    const payload = parts[1]
    const decoded = atob(payload.replace(/-/g, '+').replace(/_/g, '/'))
    return JSON.parse(decoded)
  } catch (error) {
    console.error('Failed to decode token:', error)
    return null
  }
}

/**
 * トークンからユーザー情報を取得
 * @param {string} role - ユーザーの役割（'admin' または 'bidder'）
 * @returns {object|null} ユーザー情報（user_id, email, role）
 */
export function getUserFromToken(role = 'admin') {
  const payload = decodeToken(role)
  if (!payload) {
    return null
  }

  if (role === 'bidder') {
    return {
      bidderId: payload.user_id,
      email: payload.email,
      displayName: payload.display_name,
      userType: payload.user_type,
    }
  }

  return {
    adminId: payload.user_id,
    email: payload.email,
    role: payload.role,
  }
}
