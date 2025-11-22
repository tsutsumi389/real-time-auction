/**
 * Bidder Management API Service
 * 入札者管理関連のAPIエンドポイントを提供
 */
import apiClient from './api'

/**
 * 入札者を登録
 * @param {object} data - 登録データ
 * @param {string} data.email - メールアドレス
 * @param {string} data.password - パスワード（8文字以上）
 * @param {string} data.display_name - 表示名（任意）
 * @param {number} data.initial_points - 初期ポイント（任意、0以上）
 * @returns {Promise<object>} レスポンス（id, email, display_name, status, points, created_at, updated_at）
 */
export async function registerBidder(data) {
  const response = await apiClient.post('/admin/bidders', data)
  return response.data
}

/**
 * 入札者一覧を取得
 * @param {object} params - クエリパラメータ
 * @param {string} params.keyword - 検索キーワード（メールアドレス、表示名、UUID）
 * @param {string} params.status - 状態フィルタ（active,suspended,deleted）カンマ区切り
 * @param {string} params.sort - ソートモード（id_asc/id_desc/email_asc/email_desc/points_asc/points_desc/created_at_asc/created_at_desc）
 * @param {number} params.page - ページ番号（1から開始）
 * @param {number} params.limit - 1ページあたりの件数（最大100）
 * @returns {Promise<object>} レスポンス（bidders, pagination）
 */
export async function getBidderList(params = {}) {
  const response = await apiClient.get('/admin/bidders', { params })
  return response.data
}

/**
 * ポイントを付与
 * @param {string} bidderId - 入札者ID（UUID）
 * @param {number} points - 付与するポイント数（1〜1,000,000）
 * @returns {Promise<object>} レスポンス（bidder, history）
 */
export async function grantPoints(bidderId, points) {
  const response = await apiClient.post(`/admin/bidders/${bidderId}/points`, {
    points,
  })
  return response.data
}

/**
 * ポイント履歴を取得
 * @param {string} bidderId - 入札者ID（UUID）
 * @param {object} params - クエリパラメータ
 * @param {number} params.page - ページ番号（1から開始）
 * @param {number} params.limit - 1ページあたりの件数（最大50）
 * @returns {Promise<object>} レスポンス（bidder, history, pagination）
 */
export async function getPointHistory(bidderId, params = {}) {
  const response = await apiClient.get(`/admin/bidders/${bidderId}/points/history`, { params })
  return response.data
}

/**
 * 入札者の状態を更新
 * @param {string} bidderId - 入札者ID（UUID）
 * @param {string} status - 新しい状態（active/suspended/deleted）
 * @returns {Promise<object>} レスポンス（bidder）
 */
export async function updateBidderStatus(bidderId, status) {
  const response = await apiClient.patch(`/admin/bidders/${bidderId}/status`, {
    status,
  })
  return response.data
}
