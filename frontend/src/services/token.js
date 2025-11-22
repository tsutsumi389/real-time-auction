/**
 * Token Management
 * JWTトークンのlocalStorage操作を管理
 */

const TOKEN_KEY = 'auth_token'
const TOKEN_EXPIRY_KEY = 'auth_token_expiry'

/**
 * トークンをlocalStorageに保存
 * @param {string} token - JWTトークン
 * @param {number} expiresIn - トークンの有効期限（秒）
 */
export function saveToken(token, expiresIn = 86400) {
  try {
    localStorage.setItem(TOKEN_KEY, token)

    // 有効期限を計算してタイムスタンプで保存
    const expiryTime = Date.now() + expiresIn * 1000
    localStorage.setItem(TOKEN_EXPIRY_KEY, expiryTime.toString())
  } catch (error) {
    console.error('Failed to save token:', error)
  }
}

/**
 * localStorageからトークンを取得
 * @returns {string|null} トークン（有効期限切れの場合はnull）
 */
export function getToken() {
  try {
    const token = localStorage.getItem(TOKEN_KEY)
    const expiry = localStorage.getItem(TOKEN_EXPIRY_KEY)

    if (!token || !expiry) {
      return null
    }

    // 有効期限チェック
    const now = Date.now()
    const expiryTime = parseInt(expiry, 10)

    if (now >= expiryTime) {
      // 有効期限切れの場合はトークンを削除
      removeToken()
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
 */
export function removeToken() {
  try {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(TOKEN_EXPIRY_KEY)
  } catch (error) {
    console.error('Failed to remove token:', error)
  }
}

/**
 * トークンの有効性をチェック
 * @returns {boolean} トークンが有効な場合true
 */
export function isTokenValid() {
  return getToken() !== null
}

/**
 * トークンのペイロードをデコード（Base64）
 * @returns {object|null} デコードされたペイロード
 */
export function decodeToken() {
  const token = getToken()
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
 * @returns {object|null} ユーザー情報（user_id, email, role）
 */
export function getUserFromToken() {
  const payload = decodeToken()
  if (!payload) {
    return null
  }

  return {
    adminId: payload.user_id,
    email: payload.email,
    role: payload.role,
  }
}
