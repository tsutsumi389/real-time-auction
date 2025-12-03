/**
 * Dashboard API Service
 * ダッシュボード関連のAPIエンドポイントを提供
 */
import apiClient from './api'

/**
 * ダッシュボード統計情報を取得
 * @returns {Promise<object>} レスポンス（stats）
 */
export async function getDashboardStats() {
  const response = await apiClient.get('/admin/dashboard/stats')
  return response.data
}

/**
 * ダッシュボードの最近のアクティビティを取得
 * ロールに応じてフィルタリングされたデータが返される
 * - system_admin: 全てのアクティビティ（新規入札者を含む）
 * - auctioneer: 新規入札者を除いたアクティビティ
 * @returns {Promise<object>} レスポンス（activities）
 */
export async function getDashboardActivities() {
  const response = await apiClient.get('/admin/dashboard/activities')
  return response.data
}
