/**
 * Dashboard Store
 * ダッシュボードの状態管理とアクション
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getDashboardStats, getDashboardActivities } from '@/services/dashboardApi'

export const useDashboardStore = defineStore('dashboard', () => {
  // State
  const stats = ref({
    activeAuctions: 0,
    todayBids: 0,
    totalBidders: 0,
    totalPoints: 0,
  })

  const activities = ref({
    recentBids: [],
    newBidders: [],
    endedAuctions: [],
  })

  const loading = ref(false)
  const error = ref(null)

  // Actions

  /**
   * ダッシュボード統計情報を取得
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function fetchStats() {
    loading.value = true
    error.value = null

    try {
      const response = await getDashboardStats()
      stats.value = {
        activeAuctions: response.stats.active_auctions,
        todayBids: response.stats.today_bids,
        totalBidders: response.stats.total_bidders,
        totalPoints: response.stats.total_points,
      }
      loading.value = false
      return true
    } catch (err) {
      loading.value = false
      error.value = err.message || '統計情報の取得に失敗しました'
      return false
    }
  }

  /**
   * ダッシュボードの最近のアクティビティを取得
   * ロールに応じてフィルタリングされたデータが返される
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function fetchActivities() {
    loading.value = true
    error.value = null

    try {
      const response = await getDashboardActivities()
      activities.value = {
        recentBids: response.activities.recent_bids || [],
        newBidders: response.activities.new_bidders || [],
        endedAuctions: response.activities.ended_auctions || [],
      }
      loading.value = false
      return true
    } catch (err) {
      loading.value = false
      error.value = err.message || 'アクティビティの取得に失敗しました'
      return false
    }
  }

  /**
   * 統計情報とアクティビティを一括取得
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function fetchAll() {
    loading.value = true
    error.value = null

    try {
      // 並列実行で高速化
      const [statsResult, activitiesResult] = await Promise.all([
        getDashboardStats(),
        getDashboardActivities(),
      ])

      stats.value = {
        activeAuctions: statsResult.stats.active_auctions,
        todayBids: statsResult.stats.today_bids,
        totalBidders: statsResult.stats.total_bidders,
        totalPoints: statsResult.stats.total_points,
      }

      activities.value = {
        recentBids: activitiesResult.activities.recent_bids || [],
        newBidders: activitiesResult.activities.new_bidders || [],
        endedAuctions: activitiesResult.activities.ended_auctions || [],
      }

      loading.value = false
      return true
    } catch (err) {
      loading.value = false
      error.value = err.message || 'ダッシュボードデータの取得に失敗しました'
      return false
    }
  }

  /**
   * エラーをクリア
   */
  function clearError() {
    error.value = null
  }

  /**
   * ストアをリセット
   */
  function reset() {
    stats.value = {
      activeAuctions: 0,
      todayBids: 0,
      totalBidders: 0,
      totalPoints: 0,
    }
    activities.value = {
      recentBids: [],
      newBidders: [],
      endedAuctions: [],
    }
    loading.value = false
    error.value = null
  }

  return {
    // State
    stats,
    activities,
    loading,
    error,
    // Actions
    fetchStats,
    fetchActivities,
    fetchAll,
    clearError,
    reset,
  }
})
