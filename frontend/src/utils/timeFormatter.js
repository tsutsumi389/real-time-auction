/**
 * Time Formatter Utility
 * 時刻フォーマット関連のユーティリティ関数
 */

/**
 * 相対時刻を日本語で返す（例: "3分前", "2時間前", "1日前"）
 * @param {string|Date} timestamp - タイムスタンプ（ISO8601文字列またはDateオブジェクト）
 * @returns {string} 相対時刻の日本語表示
 */
export function formatRelativeTime(timestamp) {
  if (!timestamp) {
    return ''
  }

  const date = typeof timestamp === 'string' ? new Date(timestamp) : timestamp
  const now = new Date()
  const diffMs = now - date

  // 負の値（未来の時刻）の場合
  if (diffMs < 0) {
    return '未来'
  }

  const seconds = Math.floor(diffMs / 1000)
  const minutes = Math.floor(seconds / 60)
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)
  const weeks = Math.floor(days / 7)
  const months = Math.floor(days / 30)
  const years = Math.floor(days / 365)

  if (years > 0) {
    return `${years}年前`
  } else if (months > 0) {
    return `${months}ヶ月前`
  } else if (weeks > 0) {
    return `${weeks}週間前`
  } else if (days > 0) {
    return `${days}日前`
  } else if (hours > 0) {
    return `${hours}時間前`
  } else if (minutes > 0) {
    return `${minutes}分前`
  } else if (seconds >= 10) {
    return `${seconds}秒前`
  } else {
    return 'たった今'
  }
}

/**
 * 日時を日本語の標準フォーマットで返す（例: "2025年12月3日 14:30"）
 * @param {string|Date} timestamp - タイムスタンプ（ISO8601文字列またはDateオブジェクト）
 * @returns {string} フォーマット済み日時
 */
export function formatDateTime(timestamp) {
  if (!timestamp) {
    return ''
  }

  const date = typeof timestamp === 'string' ? new Date(timestamp) : timestamp

  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')

  return `${year}年${month}月${day}日 ${hours}:${minutes}`
}

/**
 * 日付のみを日本語フォーマットで返す（例: "2025年12月3日"）
 * @param {string|Date} timestamp - タイムスタンプ（ISO8601文字列またはDateオブジェクト）
 * @returns {string} フォーマット済み日付
 */
export function formatDate(timestamp) {
  if (!timestamp) {
    return ''
  }

  const date = typeof timestamp === 'string' ? new Date(timestamp) : timestamp

  const year = date.getFullYear()
  const month = date.getMonth() + 1
  const day = date.getDate()

  return `${year}年${month}月${day}日`
}

/**
 * 時刻のみを返す（例: "14:30"）
 * @param {string|Date} timestamp - タイムスタンプ（ISO8601文字列またはDateオブジェクト）
 * @returns {string} フォーマット済み時刻
 */
export function formatTime(timestamp) {
  if (!timestamp) {
    return ''
  }

  const date = typeof timestamp === 'string' ? new Date(timestamp) : timestamp

  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')

  return `${hours}:${minutes}`
}

/**
 * 相対時刻とフル日時の両方を含むオブジェクトを返す
 * @param {string|Date} timestamp - タイムスタンプ
 * @returns {object} { relative: "3分前", full: "2025年12月3日 14:30" }
 */
export function formatTimeWithFull(timestamp) {
  return {
    relative: formatRelativeTime(timestamp),
    full: formatDateTime(timestamp),
  }
}
