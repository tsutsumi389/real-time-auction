/**
 * Auction Management API Service
 * オークション管理関連のAPIエンドポイントを提供
 */
import apiClient from './api'

/**
 * オークション一覧を取得
 * @param {object} params - クエリパラメータ
 * @param {number} params.page - ページ番号（デフォルト: 1）
 * @param {number} params.limit - 1ページあたりの件数（デフォルト: 20、最大: 100）
 * @param {string} params.keyword - タイトル・説明検索キーワード
 * @param {string} params.status - 状態フィルタ（pending/active/ended/cancelled）
 * @param {string} params.created_after - 作成日フィルタ（YYYY-MM-DD、この日以降）
 * @param {string} params.updated_before - 更新日フィルタ（YYYY-MM-DD、この日以前）
 * @param {string} params.sort - ソート順（created_at_asc/created_at_desc/updated_at_asc/updated_at_desc/id_asc/id_desc）
 * @returns {Promise<object>} レスポンス（auctions, pagination）
 */
export async function getAuctionList(params = {}) {
  const response = await apiClient.get('/admin/auctions', { params })
  return response.data
}

/**
 * オークションを開始
 * @param {string} auctionId - オークションID（UUID）
 * @returns {Promise<object>} レスポンス（auction）
 */
export async function startAuction(auctionId) {
  const response = await apiClient.post(`/admin/auctions/${auctionId}/start`)
  return response.data
}

/**
 * オークションを終了
 * @param {string} auctionId - オークションID（UUID）
 * @returns {Promise<object>} レスポンス（auction）
 */
export async function endAuction(auctionId) {
  const response = await apiClient.post(`/admin/auctions/${auctionId}/end`)
  return response.data
}

/**
 * オークションを中止（system_adminのみ）
 * @param {string} auctionId - オークションID（UUID）
 * @returns {Promise<object>} レスポンス（auction）
 */
export async function cancelAuction(auctionId) {
  const response = await apiClient.post(`/admin/auctions/${auctionId}/cancel`)
  return response.data
}
