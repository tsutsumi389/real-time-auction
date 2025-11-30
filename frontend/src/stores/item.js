/**
 * Item Management Store
 * 商品管理の状態管理とアクション
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import {
  getItemList,
  getItemDetail,
  createItem,
  updateItem,
  deleteItem,
} from '@/services/itemApi'

export const useItemStore = defineStore('item', () => {
  // State
  const items = ref([])
  const currentItem = ref(null)
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
    status: 'all',
    search: '',
  })

  // Actions

  /**
   * 商品一覧を取得
   * @param {object} params - クエリパラメータ
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function fetchItemList(params = {}) {
    loading.value = true
    error.value = null

    try {
      const requestParams = {
        status: params.status || filters.value.status,
        search: params.search || filters.value.search,
        page: params.page || pagination.value.currentPage,
        limit: params.limit || pagination.value.itemsPerPage,
      }

      // 空文字列のフィルタを除外（statusは除く）
      const filteredParams = Object.entries(requestParams).reduce((acc, [key, value]) => {
        if (value !== '' && value !== null && value !== undefined) {
          acc[key] = value
        }
        return acc
      }, {})

      const response = await getItemList(filteredParams)

      items.value = response.items || []
      pagination.value = {
        currentPage: response.page || 1,
        totalPages: Math.ceil((response.total || 0) / (response.limit || 20)),
        totalItems: response.total || 0,
        itemsPerPage: response.limit || 20,
      }

      loading.value = false
      return true
    } catch (err) {
      loading.value = false
      error.value = err.response?.data?.error || err.message || '商品一覧の取得に失敗しました'
      items.value = []
      return false
    }
  }

  /**
   * 商品詳細を取得
   * @param {string} itemId - 商品ID
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function fetchItemDetail(itemId) {
    loading.value = true
    error.value = null

    try {
      const response = await getItemDetail(itemId)
      currentItem.value = response
      loading.value = false
      return true
    } catch (err) {
      loading.value = false
      error.value = err.response?.data?.error || err.message || '商品詳細の取得に失敗しました'
      currentItem.value = null
      return false
    }
  }

  /**
   * 商品を新規作成
   * @param {object} data - 商品作成データ
   * @returns {Promise<object|null>} 作成された商品またはnull
   */
  async function handleCreateItem(data) {
    loading.value = true
    error.value = null

    try {
      const response = await createItem(data)
      loading.value = false
      return response
    } catch (err) {
      loading.value = false
      error.value = err.response?.data?.error || err.message || '商品の作成に失敗しました'
      return null
    }
  }

  /**
   * 商品を更新
   * @param {string} itemId - 商品ID
   * @param {object} data - 商品更新データ
   * @returns {Promise<object|null>} 更新された商品またはnull
   */
  async function handleUpdateItem(itemId, data) {
    loading.value = true
    error.value = null

    try {
      const response = await updateItem(itemId, data)
      loading.value = false
      return response
    } catch (err) {
      loading.value = false
      error.value = err.response?.data?.error || err.message || '商品の更新に失敗しました'
      return null
    }
  }

  /**
   * 商品を削除
   * @param {string} itemId - 商品ID
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function handleDeleteItem(itemId) {
    loading.value = true
    error.value = null

    try {
      await deleteItem(itemId)
      loading.value = false
      // 一覧を再取得
      await fetchItemList()
      return true
    } catch (err) {
      loading.value = false
      error.value = err.response?.data?.error || err.message || '商品の削除に失敗しました'
      return false
    }
  }

  /**
   * フィルタを設定して一覧を再取得
   * @param {object} newFilters - 新しいフィルタ条件
   */
  async function setFiltersAndFetch(newFilters) {
    filters.value = { ...filters.value, ...newFilters }
    pagination.value.currentPage = 1
    await fetchItemList()
  }

  /**
   * ページを変更して一覧を再取得
   * @param {number} page - ページ番号
   */
  async function changePage(page) {
    pagination.value.currentPage = page
    await fetchItemList()
  }

  /**
   * フィルタをリセット
   */
  async function resetFilters() {
    filters.value = {
      status: 'all',
      search: '',
    }
    pagination.value.currentPage = 1
    await fetchItemList()
  }

  /**
   * エラーをクリア
   */
  function clearError() {
    error.value = null
  }

  /**
   * 現在の商品をクリア
   */
  function clearCurrentItem() {
    currentItem.value = null
  }

  return {
    // State
    items,
    currentItem,
    pagination,
    loading,
    error,
    filters,
    // Actions
    fetchItemList,
    fetchItemDetail,
    handleCreateItem,
    handleUpdateItem,
    handleDeleteItem,
    setFiltersAndFetch,
    changePage,
    resetFilters,
    clearError,
    clearCurrentItem,
  }
})
