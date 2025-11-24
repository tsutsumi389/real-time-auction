/**
 * Bidder Auction Store
 * 入札者向けオークション一覧の状態管理とアクション
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getBidderAuctionList } from '@/services/bidderAuctionApi'

export const useBidderAuctionStore = defineStore('bidderAuction', () => {
  // State
  const auctions = ref([])
  const pagination = ref({
    total: 0,
    offset: 0,
    limit: 20,
    hasMore: false,
  })
  const loading = ref(false)
  const loadingMore = ref(false)
  const error = ref(null)

  // フィルタ・検索条件
  const filters = ref({
    keyword: '',
    status: 'active', // デフォルト: activeのみ
    sort: 'started_at_desc', // デフォルト: 開始日時降順（新しい順）
  })

  // Actions

  /**
   * オークション一覧を取得（初回読み込み）
   * @param {object} params - クエリパラメータ
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function fetchAuctionList(params = {}) {
    loading.value = true
    error.value = null

    try {
      const requestParams = {
        keyword: params.keyword !== undefined ? params.keyword : filters.value.keyword,
        status: params.status !== undefined ? params.status : filters.value.status,
        sort: params.sort !== undefined ? params.sort : filters.value.sort,
        offset: 0, // 初回読み込みは常にoffset=0
        limit: params.limit || pagination.value.limit,
      }

      // 空文字列のフィルタを除外（statusは除く）
      const filteredParams = Object.entries(requestParams).reduce((acc, [key, value]) => {
        if (key === 'status' || (value !== '' && value !== null && value !== undefined)) {
          acc[key] = value
        }
        return acc
      }, {})

      const response = await getBidderAuctionList(filteredParams)

      auctions.value = response.auctions || []
      pagination.value = {
        total: response.pagination?.total || 0,
        offset: response.pagination?.offset || 0,
        limit: response.pagination?.limit || 20,
        hasMore: response.pagination?.has_more || false,
      }

      loading.value = false
      return true
    } catch (err) {
      loading.value = false
      error.value = err.message || 'オークション一覧の取得に失敗しました'
      auctions.value = []
      return false
    }
  }

  /**
   * オークション一覧を追加読み込み（無限スクロール用）
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function loadMoreAuctions() {
    // すでに読み込み中、またはこれ以上データがない場合は何もしない
    if (loadingMore.value || !pagination.value.hasMore) {
      return false
    }

    loadingMore.value = true
    error.value = null

    try {
      const nextOffset = pagination.value.offset + pagination.value.limit

      const requestParams = {
        keyword: filters.value.keyword,
        status: filters.value.status,
        sort: filters.value.sort,
        offset: nextOffset,
        limit: pagination.value.limit,
      }

      // 空文字列のフィルタを除外（statusは除く）
      const filteredParams = Object.entries(requestParams).reduce((acc, [key, value]) => {
        if (key === 'status' || (value !== '' && value !== null && value !== undefined)) {
          acc[key] = value
        }
        return acc
      }, {})

      const response = await getBidderAuctionList(filteredParams)

      // 既存のリストに追加
      auctions.value = [...auctions.value, ...(response.auctions || [])]
      pagination.value = {
        total: response.pagination?.total || 0,
        offset: response.pagination?.offset || nextOffset,
        limit: response.pagination?.limit || 20,
        hasMore: response.pagination?.has_more || false,
      }

      loadingMore.value = false
      return true
    } catch (err) {
      loadingMore.value = false
      error.value = err.message || 'オークション一覧の追加読み込みに失敗しました'
      return false
    }
  }

  /**
   * フィルタを設定して一覧を再取得
   * @param {object} newFilters - 新しいフィルタ条件
   */
  async function setFiltersAndFetch(newFilters) {
    filters.value = { ...filters.value, ...newFilters }
    await fetchAuctionList()
  }

  /**
   * キーワード検索
   * @param {string} keyword - 検索キーワード
   */
  async function searchByKeyword(keyword) {
    filters.value.keyword = keyword
    await fetchAuctionList()
  }

  /**
   * ステータスフィルタ変更
   * @param {string} status - ステータス（active/ended/cancelled）
   */
  async function filterByStatus(status) {
    filters.value.status = status
    await fetchAuctionList()
  }

  /**
   * ソート順変更
   * @param {string} sort - ソート順（started_at_asc/started_at_desc/updated_at_asc/updated_at_desc）
   */
  async function changeSort(sort) {
    filters.value.sort = sort
    await fetchAuctionList()
  }

  /**
   * フィルタをリセット
   */
  async function resetFilters() {
    filters.value = {
      keyword: '',
      status: 'active',
      sort: 'started_at_desc',
    }
    await fetchAuctionList()
  }

  /**
   * エラーをクリア
   */
  function clearError() {
    error.value = null
  }

  /**
   * ストアをリセット（ページ離脱時など）
   */
  function reset() {
    auctions.value = []
    pagination.value = {
      total: 0,
      offset: 0,
      limit: 20,
      hasMore: false,
    }
    loading.value = false
    loadingMore.value = false
    error.value = null
    filters.value = {
      keyword: '',
      status: 'active',
      sort: 'started_at_desc',
    }
  }

  return {
    // State
    auctions,
    pagination,
    loading,
    loadingMore,
    error,
    filters,
    // Actions
    fetchAuctionList,
    loadMoreAuctions,
    setFiltersAndFetch,
    searchByKeyword,
    filterByStatus,
    changeSort,
    resetFilters,
    clearError,
    reset,
  }
})
