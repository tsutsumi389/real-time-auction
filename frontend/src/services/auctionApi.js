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

/**
 * オークションを作成
 * @param {object} data - オークション作成データ
 * @param {string} data.title - オークションタイトル（必須、最大200文字）
 * @param {string} data.description - オークション説明（任意、最大2000文字）
 * @param {Array<object>} data.items - 出品物リスト（必須、最低1つ）
 * @param {string} data.items[].name - 出品物名（必須、最大200文字）
 * @param {string} data.items[].description - 出品物説明（任意、最大2000文字）
 * @param {number} data.items[].lot_number - ロット番号（必須、0以上）
 * @returns {Promise<object>} レスポンス（id, title, description, status, item_count, created_at, updated_at）
 */
export async function createAuction(data) {
  const response = await apiClient.post('/admin/auctions', data)
  return response.data
}

/**
 * オークション参加者一覧を取得（ライブ画面用）
 * @param {string} auctionId - オークションID（UUID）
 * @returns {Promise<object>} レスポンス（participants配列）
 */
export async function getParticipants(auctionId) {
  const response = await apiClient.get(`/admin/auctions/${auctionId}/participants`)
  return response.data
}

/**
 * 商品を開始
 * @param {string} itemId - 商品ID
 * @returns {Promise<object>} レスポンス（item情報）
 */
export async function startItem(itemId) {
  const response = await apiClient.post(`/admin/items/${itemId}/start`)
  return response.data
}

/**
 * 価格を開示
 * @param {string} itemId - 商品ID
 * @param {object} data - 価格開示データ
 * @param {number} data.new_price - 開示する価格
 * @returns {Promise<object>} レスポンス（price_history情報）
 */
export async function openPrice(itemId, data) {
  const response = await apiClient.post(`/admin/items/${itemId}/open-price`, data)
  return response.data
}

/**
 * 商品を終了
 * @param {string} itemId - 商品ID
 * @returns {Promise<object>} レスポンス（item情報と落札者情報）
 */
export async function endItem(itemId) {
  const response = await apiClient.post(`/admin/items/${itemId}/end`)
  return response.data
}

/**
 * 入札履歴を取得
 * @param {string} itemId - 商品ID
 * @returns {Promise<object>} レスポンス（bids配列）
 */
export async function getBidHistory(itemId) {
  const response = await apiClient.get(`/admin/items/${itemId}/bids`)
  return response.data
}

/**
 * 価格開示履歴を取得
 * @param {string} itemId - 商品ID
 * @returns {Promise<object>} レスポンス（price_history配列）
 */
export async function getPriceHistory(itemId) {
  const response = await apiClient.get(`/admin/items/${itemId}/price-history`)
  return response.data
}

/**
 * オークション詳細を取得（ライブ画面用）
 * @param {string} auctionId - オークションID（UUID）
 * @returns {Promise<object>} レスポンス（auction詳細情報とitems配列）
 */
export async function getAuctionDetail(auctionId) {
  const response = await apiClient.get(`/auctions/${auctionId}`)
  return response.data
}

/**
 * オークション詳細を取得（編集用）
 * @param {string} auctionId - オークションID（UUID）
 * @returns {Promise<object>} レスポンス（auction詳細情報と編集可否フラグ）
 */
export async function getAuctionForEdit(auctionId) {
  const response = await apiClient.get(`/admin/auctions/${auctionId}`)
  return response.data
}

/**
 * オークションを更新
 * @param {string} auctionId - オークションID（UUID）
 * @param {object} data - 更新データ
 * @param {string} data.title - オークションタイトル
 * @param {string} data.description - オークション説明
 * @returns {Promise<object>} レスポンス（更新後のオークション情報）
 */
export async function updateAuction(auctionId, data) {
  const response = await apiClient.put(`/admin/auctions/${auctionId}`, data)
  return response.data
}

/**
 * 商品を更新
 * @param {string} auctionId - オークションID（UUID）
 * @param {string} itemId - 商品ID（UUID）
 * @param {object} data - 更新データ
 * @param {string} data.name - 商品名
 * @param {string} data.description - 商品説明
 * @returns {Promise<object>} レスポンス（更新後の商品情報）
 */
export async function updateItem(auctionId, itemId, data) {
  const response = await apiClient.put(`/admin/auctions/${auctionId}/items/${itemId}`, data)
  return response.data
}

/**
 * 商品を削除
 * @param {string} auctionId - オークションID（UUID）
 * @param {string} itemId - 商品ID（UUID）
 * @returns {Promise<void>}
 */
export async function deleteItem(auctionId, itemId) {
  await apiClient.delete(`/admin/auctions/${auctionId}/items/${itemId}`)
}

/**
 * 商品を追加
 * @param {string} auctionId - オークションID（UUID）
 * @param {object} data - 商品データ
 * @param {string} data.name - 商品名
 * @param {string} data.description - 商品説明
 * @param {number} data.starting_price - 開始価格
 * @returns {Promise<object>} レスポンス（追加された商品情報）
 */
export async function addItem(auctionId, data) {
  const response = await apiClient.post(`/admin/auctions/${auctionId}/items`, data)
  return response.data
}

/**
 * 商品の順序を変更
 * @param {string} auctionId - オークションID（UUID）
 * @param {string[]} itemIds - 新しい順序の商品ID配列
 * @returns {Promise<object>} レスポンス
 */
export async function reorderItems(auctionId, itemIds) {
  const response = await apiClient.put(`/admin/auctions/${auctionId}/items/reorder`, { item_ids: itemIds })
  return response.data
}
