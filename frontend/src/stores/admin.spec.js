import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useAdminStore } from './admin'
import * as adminApi from '@/services/adminApi'

// Mock adminApi
vi.mock('@/services/adminApi')

describe('Admin Store', () => {
  let adminStore

  beforeEach(() => {
    setActivePinia(createPinia())
    adminStore = useAdminStore()

    // Clear all mocks
    vi.clearAllMocks()
  })

  describe('Initial State', () => {
    it('should have correct initial state', () => {
      expect(adminStore.admins).toEqual([])
      expect(adminStore.loading).toBe(false)
      expect(adminStore.error).toBeNull()
      expect(adminStore.pagination).toEqual({
        currentPage: 1,
        totalPages: 1,
        totalItems: 0,
        itemsPerPage: 20,
      })
      expect(adminStore.filters).toEqual({
        search: '',
        role: '',
        status: '',
        sort: 'id',
        order: 'asc',
      })
    })
  })

  describe('fetchAdminList', () => {
    it('should fetch admin list successfully', async () => {
      const mockResponse = {
        admins: [
          { id: 1, email: 'admin1@example.com', role: 'system_admin', status: 'active' },
          { id: 2, email: 'admin2@example.com', role: 'auctioneer', status: 'active' },
        ],
        pagination: {
          current_page: 1,
          total_pages: 5,
          total_items: 100,
          items_per_page: 20,
        },
      }

      vi.mocked(adminApi.getAdminList).mockResolvedValue(mockResponse)

      const result = await adminStore.fetchAdminList()

      expect(result).toBe(true)
      expect(adminStore.admins).toEqual(mockResponse.admins)
      expect(adminStore.pagination).toEqual({
        currentPage: 1,
        totalPages: 5,
        totalItems: 100,
        itemsPerPage: 20,
      })
      expect(adminStore.error).toBeNull()
    })

    it('should handle fetch error', async () => {
      const errorMessage = 'Network error'
      vi.mocked(adminApi.getAdminList).mockRejectedValue(new Error(errorMessage))

      const result = await adminStore.fetchAdminList()

      expect(result).toBe(false)
      expect(adminStore.admins).toEqual([])
      expect(adminStore.error).toBe(errorMessage)
    })

    it('should set loading state during fetch', async () => {
      let resolveFn
      const promise = new Promise((resolve) => {
        resolveFn = resolve
      })

      vi.mocked(adminApi.getAdminList).mockReturnValue(promise)

      const fetchPromise = adminStore.fetchAdminList()

      expect(adminStore.loading).toBe(true)

      resolveFn({
        admins: [],
        pagination: {
          current_page: 1,
          total_pages: 1,
          total_items: 0,
          items_per_page: 20,
        },
      })

      await fetchPromise

      expect(adminStore.loading).toBe(false)
    })

    it('should merge params with filters', async () => {
      adminStore.filters.search = 'test'
      adminStore.filters.role = 'system_admin'

      vi.mocked(adminApi.getAdminList).mockResolvedValue({
        admins: [],
        pagination: {
          current_page: 1,
          total_pages: 1,
          total_items: 0,
          items_per_page: 20,
        },
      })

      await adminStore.fetchAdminList({ status: 'active' })

      expect(adminApi.getAdminList).toHaveBeenCalledWith(
        expect.objectContaining({
          search: 'test',
          role: 'system_admin',
          status: 'active',
        })
      )
    })

    it('should filter out empty parameters', async () => {
      vi.mocked(adminApi.getAdminList).mockResolvedValue({
        admins: [],
        pagination: {
          current_page: 1,
          total_pages: 1,
          total_items: 0,
          items_per_page: 20,
        },
      })

      await adminStore.fetchAdminList()

      const callArgs = vi.mocked(adminApi.getAdminList).mock.calls[0][0]

      // 空文字列のパラメータが除外されていることを確認
      expect(callArgs.search).toBeUndefined()
      expect(callArgs.role).toBeUndefined()
      expect(callArgs.status).toBeUndefined()
      // デフォルト値は含まれる
      expect(callArgs.sort).toBe('id')
      expect(callArgs.order).toBe('asc')
    })
  })

  describe('changeAdminStatus', () => {
    beforeEach(() => {
      adminStore.admins = [
        { id: 1, email: 'admin1@example.com', role: 'system_admin', status: 'active', updated_at: '2025-01-01T00:00:00Z' },
        { id: 2, email: 'admin2@example.com', role: 'auctioneer', status: 'active', updated_at: '2025-01-01T00:00:00Z' },
      ]
    })

    it('should update admin status successfully', async () => {
      const mockResponse = {
        admin: {
          id: 1,
          status: 'inactive',
          updated_at: '2025-01-02T00:00:00Z',
        },
      }

      vi.mocked(adminApi.updateAdminStatus).mockResolvedValue(mockResponse)

      const result = await adminStore.changeAdminStatus(1, 'inactive')

      expect(result).toBe(true)
      expect(adminStore.admins[0].status).toBe('inactive')
      expect(adminStore.admins[0].updated_at).toBe('2025-01-02T00:00:00Z')
      expect(adminStore.error).toBeNull()
    })

    it('should handle status update error', async () => {
      const errorMessage = 'Permission denied'
      vi.mocked(adminApi.updateAdminStatus).mockRejectedValue(new Error(errorMessage))

      const result = await adminStore.changeAdminStatus(1, 'inactive')

      expect(result).toBe(false)
      expect(adminStore.admins[0].status).toBe('active') // 変更されていない
      expect(adminStore.error).toBe(errorMessage)
    })

    it('should not update admins array if admin not found', async () => {
      const mockResponse = {
        admin: {
          id: 999,
          status: 'inactive',
          updated_at: '2025-01-02T00:00:00Z',
        },
      }

      vi.mocked(adminApi.updateAdminStatus).mockResolvedValue(mockResponse)

      const result = await adminStore.changeAdminStatus(999, 'inactive')

      expect(result).toBe(true)
      // 既存のadminsは変更されていない
      expect(adminStore.admins.length).toBe(2)
      expect(adminStore.admins[0].status).toBe('active')
      expect(adminStore.admins[1].status).toBe('active')
    })
  })

  describe('setFiltersAndFetch', () => {
    it('should update filters and reset page', async () => {
      adminStore.pagination.currentPage = 3

      vi.mocked(adminApi.getAdminList).mockResolvedValue({
        admins: [],
        pagination: {
          current_page: 1,
          total_pages: 1,
          total_items: 0,
          items_per_page: 20,
        },
      })

      await adminStore.setFiltersAndFetch({ search: 'test', role: 'system_admin' })

      expect(adminStore.filters.search).toBe('test')
      expect(adminStore.filters.role).toBe('system_admin')
      expect(adminStore.pagination.currentPage).toBe(1)
      expect(adminApi.getAdminList).toHaveBeenCalled()
    })
  })

  describe('changePage', () => {
    it('should change page and fetch', async () => {
      vi.mocked(adminApi.getAdminList).mockResolvedValue({
        admins: [],
        pagination: {
          current_page: 3,
          total_pages: 10,
          total_items: 200,
          items_per_page: 20,
        },
      })

      await adminStore.changePage(3)

      expect(adminStore.pagination.currentPage).toBe(3)
      expect(adminApi.getAdminList).toHaveBeenCalledWith(
        expect.objectContaining({
          page: 3,
        })
      )
    })
  })

  describe('changeSort', () => {
    it('should toggle sort order for same field', async () => {
      adminStore.filters.sort = 'email'
      adminStore.filters.order = 'asc'

      vi.mocked(adminApi.getAdminList).mockResolvedValue({
        admins: [],
        pagination: {
          current_page: 1,
          total_pages: 1,
          total_items: 0,
          items_per_page: 20,
        },
      })

      await adminStore.changeSort('email')

      expect(adminStore.filters.sort).toBe('email')
      expect(adminStore.filters.order).toBe('desc')
    })

    it('should set ascending order for different field', async () => {
      adminStore.filters.sort = 'email'
      adminStore.filters.order = 'desc'

      vi.mocked(adminApi.getAdminList).mockResolvedValue({
        admins: [],
        pagination: {
          current_page: 1,
          total_pages: 1,
          total_items: 0,
          items_per_page: 20,
        },
      })

      await adminStore.changeSort('role')

      expect(adminStore.filters.sort).toBe('role')
      expect(adminStore.filters.order).toBe('asc')
    })
  })

  describe('resetFilters', () => {
    it('should reset all filters to default', async () => {
      adminStore.filters = {
        search: 'test',
        role: 'system_admin',
        status: 'active',
        sort: 'email',
        order: 'desc',
      }
      adminStore.pagination.currentPage = 5

      vi.mocked(adminApi.getAdminList).mockResolvedValue({
        admins: [],
        pagination: {
          current_page: 1,
          total_pages: 1,
          total_items: 0,
          items_per_page: 20,
        },
      })

      await adminStore.resetFilters()

      expect(adminStore.filters).toEqual({
        search: '',
        role: '',
        status: '',
        sort: 'id',
        order: 'asc',
      })
      expect(adminStore.pagination.currentPage).toBe(1)
      expect(adminApi.getAdminList).toHaveBeenCalled()
    })
  })

  describe('clearError', () => {
    it('should clear error', async () => {
      // エラーを発生させる
      vi.mocked(adminApi.getAdminList).mockRejectedValue(new Error('Test error'))
      await adminStore.fetchAdminList()

      expect(adminStore.error).toBe('Test error')

      // エラーをクリア
      adminStore.clearError()

      expect(adminStore.error).toBeNull()
    })
  })
})
