/**
 * Form Validation Utilities
 * フォームバリデーション関数を提供
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

  // 255文字制限チェック
  if (email.length > 255) {
    return 'メールアドレスは255文字以内で入力してください'
  }

  // メール形式チェック
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(email)) {
    return 'メールアドレスの形式が正しくありません'
  }

  return null
}

/**
 * 表示名のバリデーション
 * @param {string} displayName - 表示名
 * @returns {string|null} エラーメッセージ（エラーがない場合はnull）
 */
export function validateDisplayName(displayName) {
  // 任意項目なので空欄の場合はエラーなし
  if (!displayName || displayName.trim() === '') {
    return null
  }

  // 100文字制限チェック
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

  // 8文字以上チェック
  if (password.length < 8) {
    return 'パスワードは8文字以上で入力してください'
  }

  return null
}

/**
 * 確認用パスワードのバリデーション
 * @param {string} password - パスワード
 * @param {string} confirmPassword - 確認用パスワード
 * @returns {string|null} エラーメッセージ（エラーがない場合はnull）
 */
export function validatePasswordConfirm(password, confirmPassword) {
  // 空欄チェック
  if (!confirmPassword || confirmPassword.trim() === '') {
    return '確認用パスワードを入力してください'
  }

  // パスワード一致チェック
  if (password !== confirmPassword) {
    return 'パスワードが一致しません'
  }

  return null
}

/**
 * ロールのバリデーション
 * @param {string} role - ロール
 * @returns {string|null} エラーメッセージ（エラーがない場合はnull）
 */
export function validateRole(role) {
  // 空欄チェック
  if (!role || role.trim() === '') {
    return 'ロールを選択してください'
  }

  // 有効なロール値チェック
  const validRoles = ['system_admin', 'auctioneer']
  if (!validRoles.includes(role)) {
    return 'ロールの値が正しくありません'
  }

  return null
}

/**
 * 管理者登録フォーム全体のバリデーション
 * @param {object} formData - フォームデータ
 * @param {string} formData.email - メールアドレス
 * @param {string} formData.display_name - 表示名
 * @param {string} formData.password - パスワード
 * @param {string} formData.password_confirm - 確認用パスワード
 * @param {string} formData.role - ロール
 * @returns {object} エラーオブジェクト（エラーがない場合は空オブジェクト）
 */
export function validateAdminRegistrationForm(formData) {
  const errors = {}

  // 各フィールドのバリデーション
  const emailError = validateEmail(formData.email)
  if (emailError) {
    errors.email = emailError
  }

  const displayNameError = validateDisplayName(formData.display_name)
  if (displayNameError) {
    errors.display_name = displayNameError
  }

  const passwordError = validatePassword(formData.password)
  if (passwordError) {
    errors.password = passwordError
  }

  const passwordConfirmError = validatePasswordConfirm(formData.password, formData.password_confirm)
  if (passwordConfirmError) {
    errors.password_confirm = passwordConfirmError
  }

  const roleError = validateRole(formData.role)
  if (roleError) {
    errors.role = roleError
  }

  return errors
}

/**
 * 管理者編集用パスワードのバリデーション（任意）
 * @param {string} password - パスワード
 * @returns {string|null} エラーメッセージ（エラーがない場合はnull）
 */
export function validateOptionalPassword(password) {
  // 任意項目なので空欄の場合はエラーなし
  if (!password || password.trim() === '') {
    return null
  }

  // 8文字以上チェック
  if (password.length < 8) {
    return 'パスワードは8文字以上で入力してください'
  }

  return null
}

/**
 * 管理者編集用確認パスワードのバリデーション（パスワードが入力されている場合のみ必須）
 * @param {string} password - パスワード
 * @param {string} confirmPassword - 確認用パスワード
 * @returns {string|null} エラーメッセージ（エラーがない場合はnull）
 */
export function validateOptionalPasswordConfirm(password, confirmPassword) {
  // パスワードが空欄の場合
  if (!password || password.trim() === '') {
    // 確認パスワードが入力されている場合は警告
    if (confirmPassword && confirmPassword.trim() !== '') {
      return '新しいパスワードを入力してください'
    }
    return null
  }

  // パスワードが入力されている場合、確認も必須
  if (!confirmPassword || confirmPassword.trim() === '') {
    return '確認用パスワードを入力してください'
  }

  // パスワード一致チェック
  if (password !== confirmPassword) {
    return 'パスワードが一致しません'
  }

  return null
}

/**
 * 状態のバリデーション
 * @param {string} status - 状態
 * @returns {string|null} エラーメッセージ（エラーがない場合はnull）
 */
export function validateStatus(status) {
  // 空欄チェック
  if (!status || status.trim() === '') {
    return '状態を選択してください'
  }

  // 有効な状態値チェック
  const validStatuses = ['active', 'suspended']
  if (!validStatuses.includes(status)) {
    return '状態の値が正しくありません'
  }

  return null
}

/**
 * 管理者編集フォーム全体のバリデーション
 * @param {object} formData - フォームデータ
 * @param {string} formData.email - メールアドレス
 * @param {string} formData.display_name - 表示名
 * @param {string} formData.password - パスワード（任意）
 * @param {string} formData.password_confirm - 確認用パスワード
 * @param {string} formData.role - ロール
 * @param {string} formData.status - 状態
 * @returns {object} エラーオブジェクト（エラーがない場合は空オブジェクト）
 */
export function validateAdminEditForm(formData) {
  const errors = {}

  // 各フィールドのバリデーション
  const emailError = validateEmail(formData.email)
  if (emailError) {
    errors.email = emailError
  }

  const displayNameError = validateDisplayName(formData.display_name)
  if (displayNameError) {
    errors.display_name = displayNameError
  }

  const passwordError = validateOptionalPassword(formData.password)
  if (passwordError) {
    errors.password = passwordError
  }

  const passwordConfirmError = validateOptionalPasswordConfirm(formData.password, formData.password_confirm)
  if (passwordConfirmError) {
    errors.password_confirm = passwordConfirmError
  }

  const roleError = validateRole(formData.role)
  if (roleError) {
    errors.role = roleError
  }

  const statusError = validateStatus(formData.status)
  if (statusError) {
    errors.status = statusError
  }

  return errors
}

/**
 * エラーオブジェクトが空かどうかをチェック
 * @param {object} errors - エラーオブジェクト
 * @returns {boolean} エラーがない場合true
 */
export function hasNoErrors(errors) {
  return Object.keys(errors).length === 0
}
