/**
 * Item Management API Service
 * 商品管理関連のAPIエンドポイントを提供
 */
import apiClient from './api'

/**
 * 商品一覧を取得
 * @param {object} params - クエリパラメータ
 * @param {string} params.status - 状態フィルタ（all/assigned/unassigned）
 * @param {string} params.search - キーワード検索
 * @param {number} params.page - ページ番号（デフォルト: 1）
 * @param {number} params.limit - 1ページあたりの件数（デフォルト: 20、最大: 100）
 * @returns {Promise<object>} レスポンス（items, total, page, limit）
 */
export async function getItemList(params = {}) {
  const response = await apiClient.get('/admin/items', { params })
  return response.data
}

/**
 * 商品詳細を取得
 * @param {string} itemId - 商品ID（UUID）
 * @returns {Promise<object>} レスポンス（商品詳細情報）
 */
export async function getItemDetail(itemId) {
  const response = await apiClient.get(`/admin/items/${itemId}`)
  return response.data
}

/**
 * 商品を新規作成
 * @param {object} data - 商品作成データ
 * @param {string} data.name - 商品名（必須、最大200文字）
 * @param {string} data.description - 商品説明（任意、最大2000文字）
 * @param {number} data.starting_price - 開始価格（任意）
 * @returns {Promise<object>} レスポンス（作成された商品情報）
 */
export async function createItem(data) {
  const response = await apiClient.post('/admin/items', data)
  return response.data
}

/**
 * 商品を更新
 * @param {string} itemId - 商品ID（UUID）
 * @param {object} data - 商品更新データ
 * @param {string} data.name - 商品名（必須、最大200文字）
 * @param {string} data.description - 商品説明（任意、最大2000文字）
 * @param {number} data.starting_price - 開始価格（任意）
 * @returns {Promise<object>} レスポンス（更新された商品情報）
 */
export async function updateItem(itemId, data) {
  const response = await apiClient.put(`/admin/items/${itemId}`, data)
  return response.data
}

/**
 * 商品を削除
 * @param {string} itemId - 商品ID（UUID）
 * @returns {Promise<void>}
 */
export async function deleteItem(itemId) {
  await apiClient.delete(`/admin/items/${itemId}`)
}

/**
 * 商品をオークションに割り当て
 * @param {string} auctionId - オークションID（UUID）
 * @param {string[]} itemIds - 商品ID配列
 * @returns {Promise<object>} レスポンス
 */
export async function assignItemsToAuction(auctionId, itemIds) {
  const response = await apiClient.post(`/admin/auctions/${auctionId}/items/assign`, {
    item_ids: itemIds,
  })
  return response.data
}

/**
 * 商品をオークションから解除
 * @param {string} auctionId - オークションID（UUID）
 * @param {string} itemId - 商品ID（UUID）
 * @returns {Promise<void>}
 */
export async function unassignItemFromAuction(auctionId, itemId) {
  await apiClient.delete(`/admin/auctions/${auctionId}/items/${itemId}/unassign`)
}
