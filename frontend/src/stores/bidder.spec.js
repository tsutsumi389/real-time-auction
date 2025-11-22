import { describe, it, expect, beforeEach, vi } from 'vitest'
import { setActivePinia, createPinia } from 'pinia'
import { useBidderStore } from './bidder'
import * as bidderApi from '@/services/bidderApi'

// Mock bidderApi
vi.mock('@/services/bidderApi')

describe('Bidder Store', () => {
  let bidderStore

  beforeEach(() => {
    setActivePinia(createPinia())
    bidderStore = useBidderStore()

    // Clear all mocks
    vi.clearAllMocks()
  })

  describe('Initial State', () => {
    it('should have correct initial state', () => {
      expect(bidderStore.bidders).toEqual([])
      expect(bidderStore.loading).toBe(false)
      expect(bidderStore.error).toBeNull()
      expect(bidderStore.pagination).toEqual({
        currentPage: 1,
        totalPages: 1,
        totalItems: 0,
        itemsPerPage: 20,
      })
      expect(bidderStore.filters).toEqual({
        keyword: '',
        status: ['active', 'suspended'],
        sort: 'created_at',
        order: 'asc',
      })
    })
  })

  describe('fetchBidderList', () => {
    it('should fetch bidder list successfully', async () => {
      const mockResponse = {
        bidders: [
          {
            id: 'abc12345-def6-7890-abcd-ef1234567890',
            email: 'bidder1@example.com',
            display_name: '田中太郎',
            status: 'active',
            total_points: 10000,
            created_at: '2025-01-01T00:00:00Z',
          },
          {
            id: 'def67890-abc1-2345-6789-0abcdef12345',
            email: 'bidder2@example.com',
            display_name: '佐藤花子',
            status: 'active',
            total_points: 5000,
            created_at: '2025-01-02T00:00:00Z',
          },
        ],
        pagination: {
          page: 1,
          total_pages: 5,
          total: 100,
          limit: 20,
        },
      }

      vi.mocked(bidderApi.getBidderList).mockResolvedValue(mockResponse)

      const result = await bidderStore.fetchBidderList()

      expect(result).toBe(true)
      expect(bidderStore.bidders).toEqual(mockResponse.bidders)
      expect(bidderStore.pagination).toEqual({
        currentPage: 1,
        totalPages: 5,
        totalItems: 100,
        itemsPerPage: 20,
      })
      expect(bidderStore.error).toBeNull()
    })

    it('should handle fetch error', async () => {
      const errorMessage = 'Network error'
      vi.mocked(bidderApi.getBidderList).mockRejectedValue(new Error(errorMessage))

      const result = await bidderStore.fetchBidderList()

      expect(result).toBe(false)
      expect(bidderStore.bidders).toEqual([])
      expect(bidderStore.error).toBe(errorMessage)
    })

    it('should set loading state during fetch', async () => {
      let resolveFn
      const promise = new Promise((resolve) => {
        resolveFn = resolve
      })

      vi.mocked(bidderApi.getBidderList).mockReturnValue(promise)

      const fetchPromise = bidderStore.fetchBidderList()

      expect(bidderStore.loading).toBe(true)

      resolveFn({
        bidders: [],
        pagination: {
          page: 1,
          total_pages: 1,
          total: 0,
          limit: 20,
        },
      })

      await fetchPromise

      expect(bidderStore.loading).toBe(false)
    })

    it('should build correct sort parameter', async () => {
      vi.mocked(bidderApi.getBidderList).mockResolvedValue({
        bidders: [],
        pagination: {
          page: 1,
          total_pages: 1,
          total: 0,
          limit: 20,
        },
      })

      await bidderStore.fetchBidderList({ sort: 'email', order: 'desc' })

      expect(bidderApi.getBidderList).toHaveBeenCalledWith(
        expect.objectContaining({
          sort: 'email_desc',
        })
      )
    })

    it('should build correct status parameter from array', async () => {
      vi.mocked(bidderApi.getBidderList).mockResolvedValue({
        bidders: [],
        pagination: {
          page: 1,
          total_pages: 1,
          total: 0,
          limit: 20,
        },
      })

      await bidderStore.fetchBidderList({ status: ['active', 'suspended'] })

      expect(bidderApi.getBidderList).toHaveBeenCalledWith(
        expect.objectContaining({
          status: 'active,suspended',
        })
      )
    })

    it('should filter out empty parameters', async () => {
      bidderStore.filters.keyword = ''

      vi.mocked(bidderApi.getBidderList).mockResolvedValue({
        bidders: [],
        pagination: {
          page: 1,
          total_pages: 1,
          total: 0,
          limit: 20,
        },
      })

      await bidderStore.fetchBidderList()

      const callArgs = vi.mocked(bidderApi.getBidderList).mock.calls[0][0]

      // 空文字列のパラメータが除外されていることを確認
      expect(callArgs.keyword).toBeUndefined()
      // デフォルト値は含まれる
      expect(callArgs.status).toBe('active,suspended')
      expect(callArgs.sort).toBe('created_at_asc')
    })
  })

  describe('addPoints', () => {
    beforeEach(() => {
      bidderStore.bidders = [
        {
          id: 'abc12345-def6-7890-abcd-ef1234567890',
          email: 'bidder1@example.com',
          display_name: '田中太郎',
          status: 'active',
          total_points: 10000,
          available_points: 8000,
          reserved_points: 2000,
          created_at: '2025-01-01T00:00:00Z',
          updated_at: '2025-01-01T00:00:00Z',
        },
      ]
    })

    it('should add points successfully', async () => {
      const mockResponse = {
        bidder: {
          id: 'abc12345-def6-7890-abcd-ef1234567890',
          total_points: 11000,
          available_points: 9000,
          reserved_points: 2000,
          updated_at: '2025-01-02T00:00:00Z',
        },
      }

      vi.mocked(bidderApi.grantPoints).mockResolvedValue(mockResponse)

      const result = await bidderStore.addPoints('abc12345-def6-7890-abcd-ef1234567890', 1000)

      expect(result).toBe(true)
      expect(bidderStore.bidders[0].total_points).toBe(11000)
      expect(bidderStore.bidders[0].available_points).toBe(9000)
      expect(bidderStore.bidders[0].updated_at).toBe('2025-01-02T00:00:00Z')
      expect(bidderStore.error).toBeNull()
    })

    it('should handle points grant error', async () => {
      const errorMessage = 'Invalid points value'
      vi.mocked(bidderApi.grantPoints).mockRejectedValue(new Error(errorMessage))

      const result = await bidderStore.addPoints('abc12345-def6-7890-abcd-ef1234567890', -100)

      expect(result).toBe(false)
      expect(bidderStore.bidders[0].total_points).toBe(10000) // 変更されていない
      expect(bidderStore.error).toBe(errorMessage)
    })

    it('should not update bidders array if bidder not found', async () => {
      const mockResponse = {
        bidder: {
          id: 'nonexistent-id',
          total_points: 11000,
          available_points: 11000,
          reserved_points: 0,
          updated_at: '2025-01-02T00:00:00Z',
        },
      }

      vi.mocked(bidderApi.grantPoints).mockResolvedValue(mockResponse)

      const result = await bidderStore.addPoints('nonexistent-id', 1000)

      expect(result).toBe(true)
      // 既存のbiddersは変更されていない
      expect(bidderStore.bidders.length).toBe(1)
      expect(bidderStore.bidders[0].total_points).toBe(10000)
    })
  })

  describe('fetchPointHistory', () => {
    it('should fetch point history successfully', async () => {
      const mockResponse = {
        bidder: {
          id: 'abc12345-def6-7890-abcd-ef1234567890',
          email: 'bidder1@example.com',
          display_name: '田中太郎',
        },
        history: [
          {
            id: 1,
            type: 'grant',
            points: 1000,
            balance_before: 10000,
            balance_after: 11000,
            created_at: '2025-01-10T12:00:00Z',
          },
        ],
        pagination: {
          page: 1,
          total_pages: 1,
          total: 1,
          limit: 10,
        },
      }

      vi.mocked(bidderApi.getPointHistory).mockResolvedValue(mockResponse)

      const result = await bidderStore.fetchPointHistory('abc12345-def6-7890-abcd-ef1234567890', 1, 10)

      expect(result).toEqual(mockResponse)
      expect(bidderApi.getPointHistory).toHaveBeenCalledWith('abc12345-def6-7890-abcd-ef1234567890', {
        page: 1,
        limit: 10,
      })
    })

    it('should handle point history fetch error', async () => {
      const errorMessage = 'Network error'
      vi.mocked(bidderApi.getPointHistory).mockRejectedValue(new Error(errorMessage))

      const result = await bidderStore.fetchPointHistory('abc12345-def6-7890-abcd-ef1234567890')

      expect(result).toBeNull()
      expect(bidderStore.error).toBe(errorMessage)
    })
  })

  describe('changeBidderStatus', () => {
    beforeEach(() => {
      bidderStore.bidders = [
        {
          id: 'abc12345-def6-7890-abcd-ef1234567890',
          email: 'bidder1@example.com',
          display_name: '田中太郎',
          status: 'active',
          total_points: 10000,
          created_at: '2025-01-01T00:00:00Z',
          updated_at: '2025-01-01T00:00:00Z',
        },
      ]
    })

    it('should update bidder status successfully', async () => {
      const mockResponse = {
        id: 'abc12345-def6-7890-abcd-ef1234567890',
        status: 'suspended',
        updated_at: '2025-01-02T00:00:00Z',
      }

      vi.mocked(bidderApi.updateBidderStatus).mockResolvedValue(mockResponse)

      const result = await bidderStore.changeBidderStatus('abc12345-def6-7890-abcd-ef1234567890', 'suspended')

      expect(result).toBe(true)
      expect(bidderStore.bidders[0].status).toBe('suspended')
      expect(bidderStore.bidders[0].updated_at).toBe('2025-01-02T00:00:00Z')
      expect(bidderStore.error).toBeNull()
    })

    it('should handle status update error', async () => {
      const errorMessage = 'Permission denied'
      vi.mocked(bidderApi.updateBidderStatus).mockRejectedValue(new Error(errorMessage))

      const result = await bidderStore.changeBidderStatus('abc12345-def6-7890-abcd-ef1234567890', 'suspended')

      expect(result).toBe(false)
      expect(bidderStore.bidders[0].status).toBe('active') // 変更されていない
      expect(bidderStore.error).toBe(errorMessage)
    })
  })

  describe('setFiltersAndFetch', () => {
    it('should update filters and reset page', async () => {
      bidderStore.pagination.currentPage = 3

      vi.mocked(bidderApi.getBidderList).mockResolvedValue({
        bidders: [],
        pagination: {
          page: 1,
          total_pages: 1,
          total: 0,
          limit: 20,
        },
      })

      await bidderStore.setFiltersAndFetch({ keyword: 'test', status: ['active'] })

      expect(bidderStore.filters.keyword).toBe('test')
      expect(bidderStore.filters.status).toEqual(['active'])
      expect(bidderStore.pagination.currentPage).toBe(1)
      expect(bidderApi.getBidderList).toHaveBeenCalled()
    })
  })

  describe('changePage', () => {
    it('should change page and fetch', async () => {
      vi.mocked(bidderApi.getBidderList).mockResolvedValue({
        bidders: [],
        pagination: {
          page: 3,
          total_pages: 10,
          total: 200,
          limit: 20,
        },
      })

      await bidderStore.changePage(3)

      expect(bidderStore.pagination.currentPage).toBe(3)
      expect(bidderApi.getBidderList).toHaveBeenCalledWith(
        expect.objectContaining({
          page: 3,
        })
      )
    })
  })

  describe('changeSort', () => {
    it('should toggle sort order for same field', async () => {
      bidderStore.filters.sort = 'email'
      bidderStore.filters.order = 'asc'

      vi.mocked(bidderApi.getBidderList).mockResolvedValue({
        bidders: [],
        pagination: {
          page: 1,
          total_pages: 1,
          total: 0,
          limit: 20,
        },
      })

      await bidderStore.changeSort('email')

      expect(bidderStore.filters.sort).toBe('email')
      expect(bidderStore.filters.order).toBe('desc')
    })

    it('should set ascending order for different field', async () => {
      bidderStore.filters.sort = 'email'
      bidderStore.filters.order = 'desc'

      vi.mocked(bidderApi.getBidderList).mockResolvedValue({
        bidders: [],
        pagination: {
          page: 1,
          total_pages: 1,
          total: 0,
          limit: 20,
        },
      })

      await bidderStore.changeSort('created_at')

      expect(bidderStore.filters.sort).toBe('created_at')
      expect(bidderStore.filters.order).toBe('asc')
    })
  })

  describe('resetFilters', () => {
    it('should reset all filters to default', async () => {
      bidderStore.filters = {
        keyword: 'test',
        status: ['active'],
        sort: 'email',
        order: 'desc',
      }
      bidderStore.pagination.currentPage = 5

      vi.mocked(bidderApi.getBidderList).mockResolvedValue({
        bidders: [],
        pagination: {
          page: 1,
          total_pages: 1,
          total: 0,
          limit: 20,
        },
      })

      await bidderStore.resetFilters()

      expect(bidderStore.filters).toEqual({
        keyword: '',
        status: ['active', 'suspended'],
        sort: 'created_at',
        order: 'asc',
      })
      expect(bidderStore.pagination.currentPage).toBe(1)
      expect(bidderApi.getBidderList).toHaveBeenCalled()
    })
  })

  describe('clearError', () => {
    it('should clear error', async () => {
      // エラーを発生させる
      vi.mocked(bidderApi.getBidderList).mockRejectedValue(new Error('Test error'))
      await bidderStore.fetchBidderList()

      expect(bidderStore.error).toBe('Test error')

      // エラーをクリア
      bidderStore.clearError()

      expect(bidderStore.error).toBeNull()
    })
  })
})
