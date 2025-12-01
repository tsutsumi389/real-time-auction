/**
 * Bidder Auction API Service
 * 入札者向けオークション一覧APIエンドポイントを提供
 */
import apiClient from './api'

/**
 * 入札者向けオークション一覧を取得（公開エンドポイント、認証不要）
 * @param {object} params - クエリパラメータ
 * @param {number} params.offset - オフセット（デフォルト: 0）
 * @param {number} params.limit - 取得件数（デフォルト: 20、最大: 100）
 * @param {string} params.keyword - タイトル検索キーワード（部分一致、ILIKE）
 * @param {string} params.status - 状態フィルタ（active/ended/cancelled、デフォルト: active）
 * @param {string} params.sort - ソート順（started_at_asc/started_at_desc/updated_at_asc/updated_at_desc、デフォルト: started_at_desc）
 * @returns {Promise<object>} レスポンス（auctions, pagination）
 */
export async function getBidderAuctionList(params = {}) {
  const response = await apiClient.get('/auctions', { params })
  return response.data
}

/**
 * オークション詳細を取得（公開エンドポイント、認証不要）
 * @param {string} id - オークションID（UUID）
 * @returns {Promise<object>} オークション詳細情報（auction, items with media）
 */
export async function getAuctionDetail(id) {
  const response = await apiClient.get(`/auctions/${id}`)
  return response.data
}
