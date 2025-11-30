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
  createAuction,
  getAuctionForEdit,
  updateAuction,
  updateItem,
  deleteItem,
  addItem,
  reorderItems,
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

      auctions.value = response.auctions || []
      pagination.value = {
        currentPage: response.pagination?.page || 1,
        totalPages: response.pagination?.total_pages || 1,
        totalItems: response.pagination?.total || 0,
        itemsPerPage: response.pagination?.limit || 20,
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
   * オークションを作成
   * @param {object} data - オークション作成データ
   * @param {string} data.title - オークションタイトル
   * @param {string} data.description - オークション説明
   * @param {Array<object>} data.items - 出品物リスト
   * @returns {Promise<object|null>} 成功した場合は作成されたオークション情報、失敗した場合はnull
   */
  async function handleCreateAuction(data) {
    loading.value = true
    error.value = null

    try {
      const response = await createAuction(data)

      // 一覧に新しいオークションを追加（先頭に挿入）
      auctions.value.unshift({
        id: response.id,
        title: response.title,
        description: response.description,
        status: response.status,
        item_count: response.item_count,
        created_at: response.created_at,
        updated_at: response.updated_at,
      })

      loading.value = false
      return response
    } catch (err) {
      loading.value = false
      error.value = err.message || 'オークションの作成に失敗しました'
      return null
    }
  }

  /**
   * エラーをクリア
   */
  function clearError() {
    error.value = null
  }

  // 編集用状態
  const currentAuction = ref(null)
  const isLoadingAuction = ref(false)

  /**
   * 編集用のオークション詳細を取得
   * @param {string} id - オークションID
   * @returns {Promise<object|null>} オークション情報またはnull
   */
  async function fetchAuctionForEdit(id) {
    isLoadingAuction.value = true
    error.value = null

    try {
      const response = await getAuctionForEdit(id)
      currentAuction.value = response
      isLoadingAuction.value = false
      return response
    } catch (err) {
      isLoadingAuction.value = false
      error.value = err.message || 'オークション情報の取得に失敗しました'
      currentAuction.value = null
      return null
    }
  }

  /**
   * オークションを更新
   * @param {string} id - オークションID
   * @param {object} data - 更新データ
   * @returns {Promise<object|null>} 更新後のオークション情報またはnull
   */
  async function handleUpdateAuction(id, data) {
    error.value = null

    try {
      const response = await updateAuction(id, data)

      // currentAuctionを更新
      if (currentAuction.value && currentAuction.value.id === id) {
        currentAuction.value = {
          ...currentAuction.value,
          title: response.title,
          description: response.description,
          updated_at: response.updated_at,
        }
      }

      // 一覧内の該当オークションを更新
      const index = auctions.value.findIndex((auction) => auction.id === id)
      if (index !== -1) {
        auctions.value[index] = {
          ...auctions.value[index],
          title: response.title,
          description: response.description,
          updated_at: response.updated_at,
        }
      }

      return response
    } catch (err) {
      error.value = err.message || 'オークションの更新に失敗しました'
      return null
    }
  }

  /**
   * 商品を更新
   * @param {string} auctionId - オークションID
   * @param {string} itemId - 商品ID
   * @param {object} data - 更新データ
   * @returns {Promise<object|null>} 更新後の商品情報またはnull
   */
  async function handleUpdateItem(auctionId, itemId, data) {
    error.value = null

    try {
      const response = await updateItem(auctionId, itemId, data)

      // currentAuctionのitemsを更新
      if (currentAuction.value && currentAuction.value.id === auctionId) {
        const itemIndex = currentAuction.value.items.findIndex((item) => item.id === itemId)
        if (itemIndex !== -1) {
          currentAuction.value.items[itemIndex] = {
            ...currentAuction.value.items[itemIndex],
            name: response.name,
            description: response.description,
            updated_at: response.updated_at,
          }
        }
      }

      return response
    } catch (err) {
      error.value = err.message || '商品の更新に失敗しました'
      return null
    }
  }

  /**
   * 商品を削除
   * @param {string} auctionId - オークションID
   * @param {string} itemId - 商品ID
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function handleDeleteItem(auctionId, itemId) {
    error.value = null

    try {
      await deleteItem(auctionId, itemId)

      // currentAuctionのitemsから削除
      if (currentAuction.value && currentAuction.value.id === auctionId) {
        currentAuction.value.items = currentAuction.value.items.filter((item) => item.id !== itemId)
        // lot_numberを再計算
        currentAuction.value.items.forEach((item, index) => {
          item.lot_number = index + 1
        })
      }

      return true
    } catch (err) {
      error.value = err.message || '商品の削除に失敗しました'
      return false
    }
  }

  /**
   * 商品を追加
   * @param {string} auctionId - オークションID
   * @param {object} data - 商品データ
   * @returns {Promise<object|null>} 追加された商品情報またはnull
   */
  async function handleAddItem(auctionId, data) {
    error.value = null

    try {
      const response = await addItem(auctionId, data)

      // currentAuctionのitemsに追加
      if (currentAuction.value && currentAuction.value.id === auctionId) {
        currentAuction.value.items.push({
          ...response,
          can_edit: true,
          can_delete: true,
          bid_count: 0,
        })
      }

      return response
    } catch (err) {
      error.value = err.message || '商品の追加に失敗しました'
      return null
    }
  }

  /**
   * 商品の順序を変更
   * @param {string} auctionId - オークションID
   * @param {string[]} itemIds - 新しい順序の商品ID配列
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function handleReorderItems(auctionId, itemIds) {
    error.value = null

    try {
      await reorderItems(auctionId, itemIds)

      // currentAuctionのitemsを再ソート
      if (currentAuction.value && currentAuction.value.id === auctionId) {
        const itemMap = new Map(currentAuction.value.items.map((item) => [item.id, item]))
        currentAuction.value.items = itemIds.map((id, index) => {
          const item = itemMap.get(id)
          return { ...item, lot_number: index + 1 }
        })
      }

      return true
    } catch (err) {
      error.value = err.message || '商品順序の変更に失敗しました'
      return false
    }
  }

  /**
   * 現在のオークションをクリア
   */
  function clearCurrentAuction() {
    currentAuction.value = null
  }

  return {
    // State
    auctions,
    pagination,
    loading,
    error,
    filters,
    currentAuction,
    isLoadingAuction,
    // Actions
    fetchAuctionList,
    handleStartAuction,
    handleEndAuction,
    handleCancelAuction,
    handleCreateAuction,
    setFiltersAndFetch,
    changePage,
    changeSort,
    resetFilters,
    clearError,
    // Edit actions
    fetchAuctionForEdit,
    handleUpdateAuction,
    handleUpdateItem,
    handleDeleteItem,
    handleAddItem,
    handleReorderItems,
    clearCurrentAuction,
  }
})
