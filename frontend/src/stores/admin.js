/**
 * Admin Management Store
 * 管理者管理の状態管理とアクション
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getAdminList, updateAdminStatus } from '@/services/adminApi'

export const useAdminStore = defineStore('admin', () => {
  // State
  const admins = ref([])
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
    search: '',
    role: '',
    status: '',
    sort: 'id',
    order: 'asc',
  })

  // Actions

  /**
   * 管理者一覧を取得
   * @param {object} params - クエリパラメータ
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function fetchAdminList(params = {}) {
    loading.value = true
    error.value = null

    try {
      const requestParams = {
        search: params.search || filters.value.search,
        role: params.role || filters.value.role,
        status: params.status || filters.value.status,
        sort: params.sort || filters.value.sort,
        order: params.order || filters.value.order,
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

      const response = await getAdminList(filteredParams)

      admins.value = response.admins
      pagination.value = {
        currentPage: response.pagination.current_page,
        totalPages: response.pagination.total_pages,
        totalItems: response.pagination.total_items,
        itemsPerPage: response.pagination.items_per_page,
      }

      loading.value = false
      return true
    } catch (err) {
      loading.value = false
      error.value = err.message || '管理者一覧の取得に失敗しました'
      return false
    }
  }

  /**
   * 管理者の状態を更新
   * @param {number} adminId - 管理者ID
   * @param {string} status - 新しい状態（active/inactive）
   * @returns {Promise<boolean>} 成功した場合true
   */
  async function changeAdminStatus(adminId, status) {
    try {
      const response = await updateAdminStatus(adminId, status)

      // 一覧内の該当管理者を更新
      const index = admins.value.findIndex((admin) => admin.id === adminId)
      if (index !== -1) {
        admins.value[index] = {
          ...admins.value[index],
          status: response.admin.status,
          updated_at: response.admin.updated_at,
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
    await fetchAdminList()
  }

  /**
   * ページを変更して一覧を取得
   * @param {number} page - ページ番号
   */
  async function changePage(page) {
    pagination.value.currentPage = page
    await fetchAdminList()
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
    await fetchAdminList()
  }

  /**
   * フィルタをリセット
   */
  async function resetFilters() {
    filters.value = {
      search: '',
      role: '',
      status: '',
      sort: 'id',
      order: 'asc',
    }
    pagination.value.currentPage = 1
    await fetchAdminList()
  }

  /**
   * エラーをクリア
   */
  function clearError() {
    error.value = null
  }

  return {
    // State
    admins,
    pagination,
    loading,
    error,
    filters,
    // Actions
    fetchAdminList,
    changeAdminStatus,
    setFiltersAndFetch,
    changePage,
    changeSort,
    resetFilters,
    clearError,
  }
})
