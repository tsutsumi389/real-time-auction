/**
 * Auction Management Store
 * オークション管理の状態管理とアクション
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  getAuctionList,
  startAuction,
  endAuction,
  cancelAuction,
} from '@/services/auctionApi'

export const useAuctionStore = defineStore('auction', () => {
  // State
  const auctions = ref([])
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
    status: '',
    createdAfter: '',
    updatedBefore: '',
    sort: 'created_at',
    order: 'desc',
  })

  // Actions

  /**
   * オークション一覧を取得
   * @param {object} params - クエリパラメータ
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function fetchAuctionList(params = {}) {
    loading.value = true
    error.value = null

    try {
      // Build sort parameter in the format expected by backend: {field}_{direction}
      const sortField = params.sort || filters.value.sort
      const sortOrder = params.order || filters.value.order
      const sortParam = `${sortField}_${sortOrder}`

      const requestParams = {
        keyword: params.keyword || filters.value.keyword,
        status: params.status || filters.value.status,
        created_after: params.createdAfter || filters.value.createdAfter,
        updated_before: params.updatedBefore || filters.value.updatedBefore,
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

      const response = await getAuctionList(filteredParams)

      auctions.value = response.auctions
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
      error.value = err.message || 'オークション一覧の取得に失敗しました'
      return false
    }
  }

  /**
   * オークションを開始
   * @param {string} auctionId - オークションID（UUID）
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function handleStartAuction(auctionId) {
    try {
      const response = await startAuction(auctionId)

      // 一覧内の該当オークションを更新
      const index = auctions.value.findIndex((auction) => auction.id === auctionId)
      if (index !== -1) {
        auctions.value[index] = {
          ...auctions.value[index],
          status: response.status,
          updated_at: response.updated_at,
        }
      }

      return true
    } catch (err) {
      error.value = err.message || 'オークションの開始に失敗しました'
      return false
    }
  }

  /**
   * オークションを終了
   * @param {string} auctionId - オークションID（UUID）
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function handleEndAuction(auctionId) {
    try {
      const response = await endAuction(auctionId)

      // 一覧内の該当オークションを更新
      const index = auctions.value.findIndex((auction) => auction.id === auctionId)
      if (index !== -1) {
        auctions.value[index] = {
          ...auctions.value[index],
          status: response.status,
          updated_at: response.updated_at,
        }
      }

      return true
    } catch (err) {
      error.value = err.message || 'オークションの終了に失敗しました'
      return false
    }
  }

  /**
   * オークションを中止（system_adminのみ）
   * @param {string} auctionId - オークションID（UUID）
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function handleCancelAuction(auctionId) {
    try {
      const response = await cancelAuction(auctionId)

      // 一覧内の該当オークションを更新
      const index = auctions.value.findIndex((auction) => auction.id === auctionId)
      if (index !== -1) {
        auctions.value[index] = {
          ...auctions.value[index],
          status: response.status,
          updated_at: response.updated_at,
        }
      }

      return true
    } catch (err) {
      error.value = err.message || 'オークションの中止に失敗しました'
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
    await fetchAuctionList()
  }

  /**
   * ページを変更して一覧を取得
   * @param {number} page - ページ番号
   */
  async function changePage(page) {
    pagination.value.currentPage = page
    await fetchAuctionList()
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
      // 異なる項目の場合、降順に設定（デフォルト）
      filters.value.sort = sort
      filters.value.order = 'desc'
    }
    await fetchAuctionList()
  }

  /**
   * フィルタをリセット
   */
  async function resetFilters() {
    filters.value = {
      keyword: '',
      status: '',
      createdAfter: '',
      updatedBefore: '',
      sort: 'created_at',
      order: 'desc',
    }
    pagination.value.currentPage = 1
    await fetchAuctionList()
  }

  /**
   * エラーをクリア
   */
  function clearError() {
    error.value = null
  }

  return {
    // State
    auctions,
    pagination,
    loading,
    error,
    filters,
    // Actions
    fetchAuctionList,
    handleStartAuction,
    handleEndAuction,
    handleCancelAuction,
    setFiltersAndFetch,
    changePage,
    changeSort,
    resetFilters,
    clearError,
  }
})
