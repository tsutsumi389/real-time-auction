/**
 * Validation Utilities
 * フォームバリデーション用のヘルパー関数
 */

/**
 * メールアドレスのバリデーション
 * @param {string} email - メールアドレス
 * @returns {string|null} エラーメッセージ（エラーがない場合はnull）
 */
export function validateEmail(email) {
  // 空欄チェック
  if (!email || email.trim() === '') {
    return 'メールアドレスを入力してください'
  }

  // メール形式チェック（簡易版）
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(email)) {
    return 'メールアドレスの形式が正しくありません'
  }

  // 文字数チェック
  if (email.length > 255) {
    return 'メールアドレスは255文字以内で入力してください'
  }

  return null
}

/**
 * 表示名のバリデーション
 * @param {string} displayName - 表示名
 * @returns {string|null} エラーメッセージ（エラーがない場合はnull）
 */
export function validateDisplayName(displayName) {
  // 任意項目のため空欄はOK
  if (!displayName || displayName.trim() === '') {
    return null
  }

  // 文字数チェック
  if (displayName.length > 100) {
    return '表示名は100文字以内で入力してください'
  }

  return null
}

/**
 * パスワードのバリデーション
 * @param {string} password - パスワード
 * @returns {string|null} エラーメッセージ（エラーがない場合はnull）
 */
export function validatePassword(password) {
  // 空欄チェック
  if (!password || password.trim() === '') {
    return 'パスワードを入力してください'
  }

  // 最小文字数チェック
  if (password.length < 8) {
    return 'パスワードは8文字以上で入力してください'
  }

  return null
}

/**
 * パスワード確認のバリデーション
 * @param {string} password - パスワード
 * @param {string} confirmPassword - 確認用パスワード
 * @returns {string|null} エラーメッセージ（エラーがない場合はnull）
 */
export function validateConfirmPassword(password, confirmPassword) {
  // 空欄チェック
  if (!confirmPassword || confirmPassword.trim() === '') {
    return '確認用パスワードを入力してください'
  }

  // パスワードとの一致チェック
  if (password !== confirmPassword) {
    return 'パスワードが一致しません'
  }

  return null
}

/**
 * 初期ポイントのバリデーション
 * @param {number|string} points - 初期ポイント
 * @returns {string|null} エラーメッセージ（エラーがない場合はnull）
 */
export function validateInitialPoints(points) {
  // 空欄はOK（任意項目）
  if (points === '' || points === null || points === undefined) {
    return null
  }

  // 数値変換
  const numPoints = Number(points)

  // 数値チェック
  if (isNaN(numPoints)) {
    return '数値を入力してください'
  }

  // 整数チェック
  if (!Number.isInteger(numPoints)) {
    return '整数で入力してください'
  }

  // 0以上チェック
  if (numPoints < 0) {
    return '0以上の整数を入力してください'
  }

  return null
}

/**
 * 入札者登録フォーム全体のバリデーション
 * @param {object} formData - フォームデータ
 * @param {string} formData.email - メールアドレス
 * @param {string} formData.password - パスワード
 * @param {string} formData.confirmPassword - 確認用パスワード
 * @param {string} formData.display_name - 表示名（任意）
 * @param {number|string} formData.initial_points - 初期ポイント（任意）
 * @returns {object} エラーオブジェクト（エラーがない場合は空オブジェクト）
 */
export function validateBidderRegistrationForm(formData) {
  const errors = {}

  // メールアドレス
  const emailError = validateEmail(formData.email)
  if (emailError) {
    errors.email = emailError
  }

  // 表示名
  const displayNameError = validateDisplayName(formData.display_name)
  if (displayNameError) {
    errors.display_name = displayNameError
  }

  // パスワード
  const passwordError = validatePassword(formData.password)
  if (passwordError) {
    errors.password = passwordError
  }

  // パスワード確認
  const confirmPasswordError = validateConfirmPassword(
    formData.password,
    formData.confirmPassword
  )
  if (confirmPasswordError) {
    errors.confirmPassword = confirmPasswordError
  }

  // 初期ポイント
  const pointsError = validateInitialPoints(formData.initial_points)
  if (pointsError) {
    errors.initial_points = pointsError
  }

  return errors
}

/**
 * エラーオブジェクトが空かどうかをチェック
 * @param {object} errors - エラーオブジェクト
 * @returns {boolean} エラーがある場合はtrue
 */
export function hasErrors(errors) {
  return Object.keys(errors).length > 0
}
