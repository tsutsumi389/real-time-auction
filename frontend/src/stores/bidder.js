/**
 * Bidder Management Store
 * 入札者管理の状態管理とアクション
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getBidderList, grantPoints, getPointHistory, updateBidderStatus } from '@/services/bidderApi'

export const useBidderStore = defineStore('bidder', () => {
  // State
  const bidders = ref([])
  const pagination = ref({
    currentPage: 1,
    totalPages: 1,
    totalItems: 0,
    itemsPerPage: 20,
  })
  const loading = ref(false)
  const error = ref(null)

  // フィルタ・検索条件
  const filters = ref({
    keyword: '',
    status: ['active', 'suspended'], // デフォルト: 有効・停止中のみ
    sort: 'created_at',
    order: 'asc',
  })

  // Actions

  /**
   * 入札者一覧を取得
   * @param {object} params - クエリパラメータ
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function fetchBidderList(params = {}) {
    loading.value = true
    error.value = null

    try {
      // Build sort parameter in the format expected by backend: {field}_{direction}
      const sortField = params.sort || filters.value.sort
      const sortOrder = params.order || filters.value.order
      const sortParam = `${sortField}_${sortOrder}`

      // Build status parameter (comma-separated)
      const statusArray = params.status || filters.value.status
      const statusParam = Array.isArray(statusArray) ? statusArray.join(',') : statusArray

      const requestParams = {
        keyword: params.keyword || filters.value.keyword,
        status: statusParam,
        sort: sortParam,
        page: params.page || pagination.value.currentPage,
        limit: params.limit || pagination.value.itemsPerPage,
      }

      // 空文字列のフィルタを除外
      const filteredParams = Object.entries(requestParams).reduce((acc, [key, value]) => {
        if (value !== '' && value !== null && value !== undefined) {
          acc[key] = value
        }
        return acc
      }, {})

      const response = await getBidderList(filteredParams)

      bidders.value = response.bidders || []
      pagination.value = {
        currentPage: response.pagination.page,
        totalPages: response.pagination.total_pages,
        totalItems: response.pagination.total,
        itemsPerPage: response.pagination.limit,
      }

      loading.value = false
      return true
    } catch (err) {
      loading.value = false
      error.value = err.message || '入札者一覧の取得に失敗しました'
      return false
    }
  }

  /**
   * ポイントを付与
   * @param {string} bidderId - 入札者ID（UUID）
   * @param {number} points - 付与するポイント数
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function addPoints(bidderId, points) {
    try {
      const response = await grantPoints(bidderId, points)

      // 一覧内の該当入札者を更新
      const index = bidders.value.findIndex((bidder) => bidder.id === bidderId)
      if (index !== -1) {
        bidders.value[index] = {
          ...bidders.value[index],
          total_points: response.bidder.total_points,
          available_points: response.bidder.available_points,
          reserved_points: response.bidder.reserved_points,
          updated_at: response.bidder.updated_at,
        }
      }

      return true
    } catch (err) {
      error.value = err.message || 'ポイント付与に失敗しました'
      return false
    }
  }

  /**
   * ポイント履歴を取得
   * @param {string} bidderId - 入札者ID（UUID）
   * @param {number} page - ページ番号
   * @param {number} limit - 1ページあたりの件数
   * @returns {Promise<object|null>} 成功した場合レスポンス、失敗した場合null
   */
  async function fetchPointHistory(bidderId, page = 1, limit = 10) {
    try {
      const response = await getPointHistory(bidderId, { page, limit })
      return response
    } catch (err) {
      error.value = err.message || 'ポイント履歴の取得に失敗しました'
      return null
    }
  }

  /**
   * 入札者の状態を更新
   * @param {string} bidderId - 入札者ID（UUID）
   * @param {string} status - 新しい状態（active/suspended/deleted）
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function changeBidderStatus(bidderId, status) {
    try {
      const response = await updateBidderStatus(bidderId, status)

      // 一覧内の該当入札者を更新
      const index = bidders.value.findIndex((bidder) => bidder.id === bidderId)
      if (index !== -1) {
        bidders.value[index] = {
          ...bidders.value[index],
          status: response.status,
          updated_at: response.updated_at,
        }
      }

      return true
    } catch (err) {
      error.value = err.message || '状態の更新に失敗しました'
      return false
    }
  }

  /**
   * 検索条件を設定して一覧を取得
   * @param {object} newFilters - 新しいフィルタ条件
   */
  async function setFiltersAndFetch(newFilters) {
    filters.value = { ...filters.value, ...newFilters }
    pagination.value.currentPage = 1 // ページをリセット
    await fetchBidderList()
  }

  /**
   * ページを変更して一覧を取得
   * @param {number} page - ページ番号
   */
  async function changePage(page) {
    pagination.value.currentPage = page
    await fetchBidderList()
  }

  /**
   * ソート条件を変更して一覧を取得
   * @param {string} sort - ソート項目
   */
  async function changeSort(sort) {
    if (filters.value.sort === sort) {
      // 同じ項目の場合、順序を反転
      filters.value.order = filters.value.order === 'asc' ? 'desc' : 'asc'
    } else {
      // 異なる項目の場合、昇順に設定
      filters.value.sort = sort
      filters.value.order = 'asc'
    }
    await fetchBidderList()
  }

  /**
   * フィルタをリセット
   */
  async function resetFilters() {
    filters.value = {
      keyword: '',
      status: ['active', 'suspended'],
      sort: 'created_at',
      order: 'asc',
    }
    pagination.value.currentPage = 1
    await fetchBidderList()
  }

  /**
   * エラーをクリア
   */
  function clearError() {
    error.value = null
  }

  return {
    // State
    bidders,
    pagination,
    loading,
    error,
    filters,
    // Actions
    fetchBidderList,
    addPoints,
    fetchPointHistory,
    changeBidderStatus,
    setFiltersAndFetch,
    changePage,
    changeSort,
    resetFilters,
    clearError,
  }
})
