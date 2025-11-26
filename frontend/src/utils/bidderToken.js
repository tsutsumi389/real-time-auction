/**
 * Bidder Token Management
 * 入札者用JWTトークンのlocalStorage操作を管理
 * 管理者トークンとは完全に分離
 */

const TOKEN_KEY = 'bidder_token'
const TOKEN_EXPIRY_KEY = 'bidder_token_expiry'
const USER_KEY = 'bidder_user'

/**
 * トークンをlocalStorageに保存
 * @param {string} token - JWTトークン
 * @param {number} expiresIn - トークンの有効期限（秒）
 */
export function saveBidderToken(token, expiresIn = 86400) {
  try {
    localStorage.setItem(TOKEN_KEY, token)

    // 有効期限を計算してタイムスタンプで保存
    const expiryTime = Date.now() + expiresIn * 1000
    localStorage.setItem(TOKEN_EXPIRY_KEY, expiryTime.toString())
  } catch (error) {
    console.error('Failed to save bidder token:', error)
  }
}

/**
 * localStorageからトークンを取得
 * @returns {string|null} トークン（有効期限切れの場合はnull）
 */
export function getBidderToken() {
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
      removeBidderToken()
      return null
    }

    return token
  } catch (error) {
    console.error('Failed to get bidder token:', error)
    return null
  }
}

/**
 * localStorageからトークンを削除
 */
export function removeBidderToken() {
  try {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(TOKEN_EXPIRY_KEY)
  } catch (error) {
    console.error('Failed to remove bidder token:', error)
  }
}

/**
 * トークンの有効性をチェック
 * @returns {boolean} トークンが有効な場合true
 */
export function isBidderTokenValid() {
  return getBidderToken() !== null
}

/**
 * トークンのペイロードをデコード（Base64）
 * @returns {object|null} デコードされたペイロード
 */
export function decodeBidderToken() {
  const token = getBidderToken()
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
    console.error('Failed to decode bidder token:', error)
    return null
  }
}

/**
 * トークンからユーザー情報を取得
 * @returns {object|null} ユーザー情報（user_id, email, display_name, user_type）
 */
export function getBidderFromToken() {
  const payload = decodeBidderToken()
  if (!payload) {
    return null
  }

  return {
    bidderId: payload.user_id,
    email: payload.email,
    displayName: payload.display_name,
    userType: payload.user_type,
  }
}

/**
 * ユーザー情報をlocalStorageに保存
 * @param {object} user - ユーザー情報
 */
export function saveBidderUser(user) {
  try {
    localStorage.setItem(USER_KEY, JSON.stringify(user))
  } catch (error) {
    console.error('Failed to save bidder user:', error)
  }
}

/**
 * ユーザー情報をlocalStorageから取得
 * @returns {object|null} ユーザー情報
 */
export function getBidderUser() {
  try {
    const userJson = localStorage.getItem(USER_KEY)
    if (!userJson) {
      return null
    }
    return JSON.parse(userJson)
  } catch (error) {
    console.error('Failed to get bidder user:', error)
    return null
  }
}

/**
 * ユーザー情報をlocalStorageから削除
 */
export function removeBidderUser() {
  try {
    localStorage.removeItem(USER_KEY)
  } catch (error) {
    console.error('Failed to remove bidder user:', error)
  }
}

/**
 * トークンとユーザー情報をすべて削除
 */
export function clearBidderAuth() {
  removeBidderToken()
  removeBidderUser()
}
