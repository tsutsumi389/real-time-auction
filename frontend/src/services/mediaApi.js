/**
 * Media Management API Service
 * メディア管理関連のAPIエンドポイントを提供
 */
import axios from 'axios'
import { getToken } from './token'
import { getApiBaseUrl } from '../config/api'

// 動的にAPIベースURLを取得（ローカルネットワーク対応）
const API_BASE_URL = getApiBaseUrl()

// ファイルアップロード用のAxiosインスタンス
// タイムアウトを長めに設定（30秒）
const uploadClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Content-Type': 'multipart/form-data',
  },
})

// リクエストインターセプター
uploadClient.interceptors.request.use(
  (config) => {
    const token = getToken('admin')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// レスポンスインターセプター
uploadClient.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      const { status, data } = error.response
      return Promise.reject({
        status,
        message: data.error || 'An error occurred',
        code: data.code || null,
        details: data.details || null,
      })
    } else if (error.request) {
      return Promise.reject({
        status: 0,
        message: 'Network error. Please check your connection.',
        details: null,
      })
    } else {
      return Promise.reject({
        status: -1,
        message: error.message,
        details: null,
      })
    }
  }
)

/**
 * 商品のメディアを取得
 * @param {string} itemId - 商品ID（UUID）
 * @returns {Promise<Array>} メディア一覧
 */
export async function getItemMedia(itemId) {
  const response = await uploadClient.get(`/items/${itemId}/media`)
  return response.data
}

/**
 * 商品にメディアをアップロード
 * @param {string} itemId - 商品ID（UUID）
 * @param {File} file - アップロードするファイル
 * @param {string} mediaType - メディアタイプ（'image' または 'video'）
 * @param {function} onProgress - アップロード進捗コールバック（0-100）
 * @returns {Promise<object>} アップロードされたメディア情報
 */
export async function uploadItemMedia(itemId, file, mediaType = 'image', onProgress = null) {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('media_type', mediaType)

  const config = {}
  if (onProgress) {
    config.onUploadProgress = (progressEvent) => {
      const percentCompleted = Math.round((progressEvent.loaded * 100) / progressEvent.total)
      onProgress(percentCompleted)
    }
  }

  const response = await uploadClient.post(`/admin/items/${itemId}/media`, formData, config)
  return response.data
}

/**
 * 商品のメディアを削除
 * @param {string} itemId - 商品ID（UUID）
 * @param {number} mediaId - メディアID
 * @returns {Promise<void>}
 */
export async function deleteItemMedia(itemId, mediaId) {
  await uploadClient.delete(`/admin/items/${itemId}/media/${mediaId}`)
}

/**
 * 商品のメディア順序を変更
 * @param {string} itemId - 商品ID（UUID）
 * @param {Array<{id: number, display_order: number}>} mediaOrder - メディアの順序配列
 * @returns {Promise<object>} 更新結果
 */
export async function reorderItemMedia(itemId, mediaOrder) {
  // バックエンドの仕様に合わせて media_ids の配列として送信
  // display_order順にソートしてIDのみを抽出
  const sortedMediaIds = mediaOrder
    .sort((a, b) => a.display_order - b.display_order)
    .map((m) => m.id)

  const response = await uploadClient.put(`/admin/items/${itemId}/media/reorder`, {
    media_ids: sortedMediaIds,
  })
  return response.data
}

/**
 * ファイルサイズを人間が読みやすい形式に変換
 * @param {number} bytes - バイト数
 * @returns {string} フォーマットされたサイズ文字列
 */
export function formatFileSize(bytes) {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

/**
 * 許可されたMIMEタイプかどうかを確認
 * @param {string} mimeType - ファイルのMIMEタイプ
 * @param {string} mediaType - メディアタイプ（'image' または 'video'）
 * @returns {boolean} 許可されている場合true
 */
export function isAllowedMimeType(mimeType, mediaType = 'image') {
  const allowedImageTypes = ['image/jpeg', 'image/png', 'image/webp', 'image/gif']
  const allowedVideoTypes = ['video/mp4', 'video/quicktime', 'video/x-msvideo']

  if (mediaType === 'image') {
    return allowedImageTypes.includes(mimeType)
  } else if (mediaType === 'video') {
    return allowedVideoTypes.includes(mimeType)
  }
  return false
}

/**
 * ファイルサイズが制限内かどうかを確認
 * @param {number} fileSize - ファイルサイズ（バイト）
 * @param {string} mediaType - メディアタイプ（'image' または 'video'）
 * @returns {boolean} 制限内の場合true
 */
export function isFileSizeValid(fileSize, mediaType = 'image') {
  const maxImageSize = 5 * 1024 * 1024 // 5MB
  const maxVideoSize = 100 * 1024 * 1024 // 100MB

  if (mediaType === 'image') {
    return fileSize <= maxImageSize
  } else if (mediaType === 'video') {
    return fileSize <= maxVideoSize
  }
  return false
}

/**
 * ファイルのメディアタイプを判定
 * @param {File} file - ファイルオブジェクト
 * @returns {string|null} 'image', 'video', または null（不明な場合）
 */
export function getMediaTypeFromFile(file) {
  if (file.type.startsWith('image/')) {
    return 'image'
  } else if (file.type.startsWith('video/')) {
    return 'video'
  }
  return null
}

/**
 * バリデーションエラーメッセージを取得
 * @param {File} file - ファイルオブジェクト
 * @returns {string|null} エラーメッセージ、またはnull（バリデーション通過）
 */
export function validateFile(file) {
  const mediaType = getMediaTypeFromFile(file)

  if (!mediaType) {
    return '対応していないファイル形式です。画像（JPEG, PNG, WebP, GIF）または動画（MP4, MOV, AVI）をアップロードしてください。'
  }

  if (!isAllowedMimeType(file.type, mediaType)) {
    if (mediaType === 'image') {
      return '対応していない画像形式です。JPEG, PNG, WebP, GIF形式の画像をアップロードしてください。'
    } else {
      return '対応していない動画形式です。MP4, MOV, AVI形式の動画をアップロードしてください。'
    }
  }

  if (!isFileSizeValid(file.size, mediaType)) {
    const maxSize = mediaType === 'image' ? '5MB' : '100MB'
    return `ファイルサイズが大きすぎます。${maxSize}以下のファイルをアップロードしてください。`
  }

  return null
}

export default {
  getItemMedia,
  uploadItemMedia,
  deleteItemMedia,
  reorderItemMedia,
  formatFileSize,
  isAllowedMimeType,
  isFileSizeValid,
  getMediaTypeFromFile,
  validateFile,
}
