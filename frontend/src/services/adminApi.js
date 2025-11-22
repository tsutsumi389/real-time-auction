/**
 * Admin Management API Service
 * 管理者管理関連のAPIエンドポイントを提供
 */
import apiClient from './api'

/**
 * 管理者一覧を取得
 * @param {object} params - クエリパラメータ
 * @param {string} params.search - 検索キーワード（メールアドレス）
 * @param {string} params.role - ロールフィルタ（system_admin/auctioneer）
 * @param {string} params.status - 状態フィルタ（active/inactive）
 * @param {string} params.sort - ソート項目（id/email/role/status/created_at）
 * @param {string} params.order - ソート順（asc/desc）
 * @param {number} params.page - ページ番号（1から開始）
 * @param {number} params.limit - 1ページあたりの件数
 * @returns {Promise<object>} レスポンス（admins, pagination）
 */
export async function getAdminList(params = {}) {
  const response = await apiClient.get('/admin/admins', { params })
  return response.data
}

/**
 * 管理者の状態を更新
 * @param {number} adminId - 管理者ID
 * @param {string} status - 新しい状態（active/inactive）
 * @returns {Promise<object>} レスポンス（admin）
 */
export async function updateAdminStatus(adminId, status) {
  const response = await apiClient.patch(`/admin/admins/${adminId}/status`, {
    status,
  })
  return response.data
}

/**
 * 新しい管理者を登録
 * @param {object} adminData - 管理者情報
 * @param {string} adminData.email - メールアドレス
 * @param {string} adminData.password - パスワード
 * @param {string} adminData.display_name - 表示名（任意）
 * @param {string} adminData.role - ロール（system_admin/auctioneer）
 * @returns {Promise<object>} レスポンス（admin）
 */
export async function registerAdmin(adminData) {
  const response = await apiClient.post('/admin/admins', adminData)
  return response.data
}
