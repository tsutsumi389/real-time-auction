/**
 * Bidder Bid API Service
 * 入札者向け入札・ポイント取得APIエンドポイントを提供
 */
import bidderApiClient from './bidderApiClient'
import apiClient from './api'

/**
 * 入札者のポイント残高を取得（認証必須）
 * @returns {Promise<object>} レスポンス（points: { total, available, reserved }）
 */
export async function getPoints() {
  const response = await bidderApiClient.get('/bidder/points')
  return response.data
}

/**
 * 商品に入札（認証必須）
 * @param {number} itemId - 商品ID
 * @param {number} price - 入札価格
 * @returns {Promise<object>} レスポンス（bid, points）
 */
export async function placeBid(itemId, price) {
  const response = await bidderApiClient.post(`/bidder/items/${itemId}/bid`, {
    price,
  })
  return response.data
}

/**
 * 商品の入札履歴を取得（認証必須）
 * @param {number} itemId - 商品ID
 * @param {object} params - クエリパラメータ
 * @param {number} params.limit - 取得件数（デフォルト: 50、最大: 200）
 * @param {number} params.offset - オフセット（デフォルト: 0）
 * @returns {Promise<object>} レスポンス（bids, pagination）
 */
export async function getBidHistory(itemId, params = {}) {
  const response = await bidderApiClient.get(`/bidder/items/${itemId}/bids`, {
    params,
  })
  return response.data
}

/**
 * オークション詳細を取得（公開エンドポイント、認証不要）
 * @param {string} auctionId - オークションID（UUID）
 * @returns {Promise<object>} レスポンス（auction詳細情報とitems配列）
 */
export async function getAuctionDetail(auctionId) {
  const response = await apiClient.get(`/auctions/${auctionId}`)
  return response.data
}
