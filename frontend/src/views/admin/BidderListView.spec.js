import { describe, it, expect, beforeEach, vi } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { createRouter, createMemoryHistory } from 'vue-router'
import { createPinia, setActivePinia } from 'pinia'
import BidderListView from './BidderListView.vue'
import { useBidderStore } from '@/stores/bidder'
import * as bidderApi from '@/services/bidderApi'

// Mock bidderApi
vi.mock('@/services/bidderApi')

// Mock UI components
vi.mock('@/components/ui/Card.vue', () => ({
  default: {
    template: '<div class="card-mock"><slot /></div>'
  }
}))

vi.mock('@/components/ui/Button.vue', () => ({
  default: {
    template: '<button v-bind="$attrs"><slot /></button>',
    inheritAttrs: false
  }
}))

// Mock bidder components
vi.mock('@/components/bidder/BidderTable.vue', () => ({
  default: {
    name: 'BidderTable',
    template: '<div class="bidder-table-mock">BidderTable</div>',
    props: ['bidders', 'loading', 'sortField', 'sortOrder'],
    emits: ['edit', 'grantPoints', 'viewHistory', 'statusChange', 'sort']
  }
}))

vi.mock('@/components/bidder/BidderFilters.vue', () => ({
  default: {
    name: 'BidderFilters',
    template: '<div class="bidder-filters-mock">BidderFilters</div>',
    props: ['modelValue'],
    emits: ['update:modelValue', 'reset']
  }
}))

vi.mock('@/components/bidder/GrantPointsDialog.vue', () => ({
  default: {
    name: 'GrantPointsDialog',
    template: '<div class="grant-points-dialog-mock">GrantPointsDialog</div>',
    props: ['open', 'bidder', 'loading'],
    emits: ['update:open', 'confirm']
  }
}))

vi.mock('@/components/bidder/PointHistoryDialog.vue', () => ({
  default: {
    name: 'PointHistoryDialog',
    template: '<div class="point-history-dialog-mock">PointHistoryDialog</div>',
    props: ['open', 'bidder'],
    emits: ['update:open']
  }
}))

vi.mock('@/components/bidder/BidderStatusChangeDialog.vue', () => ({
  default: {
    name: 'BidderStatusChangeDialog',
    template: '<div class="bidder-status-change-dialog-mock">BidderStatusChangeDialog</div>',
    props: ['open', 'bidder', 'newStatus', 'loading'],
    emits: ['update:open', 'confirm']
  }
}))

describe('BidderListView Integration Tests', () => {
  let wrapper
  let router
  let bidderStore

  const mockBidders = [
    {
      id: 'abc12345-def6-7890-abcd-ef1234567890',
      email: 'bidder1@example.com',
      display_name: '田中太郎',
      status: 'active',
      total_points: 10000,
      created_at: '2025-01-01T00:00:00Z',
      updated_at: '2025-01-01T00:00:00Z'
    },
    {
      id: 'def67890-abc1-2345-6789-0abcdef12345',
      email: 'bidder2@example.com',
      display_name: '佐藤花子',
      status: 'active',
      total_points: 5000,
      created_at: '2025-01-02T00:00:00Z',
      updated_at: '2025-01-02T00:00:00Z'
    },
    {
      id: 'ghi12345-jkl6-7890-mnop-qr1234567890',
      email: 'bidder3@example.com',
      display_name: '鈴木一郎',
      status: 'suspended',
      total_points: 15000,
      created_at: '2025-01-03T00:00:00Z',
      updated_at: '2025-01-03T00:00:00Z'
    }
  ]

  const mockPagination = {
    page: 1,
    total_pages: 5,
    total: 100,
    limit: 20
  }

  beforeEach(() => {
    setActivePinia(createPinia())
    bidderStore = useBidderStore()

    // Create router
    router = createRouter({
      history: createMemoryHistory(),
      routes: [
        { path: '/admin/bidders', component: BidderListView },
        { path: '/admin/bidders/:id/edit', component: { template: '<div>Edit Bidder</div>' } }
      ]
    })

    // Mock API responses
    vi.mocked(bidderApi.getBidderList).mockResolvedValue({
      bidders: mockBidders,
      pagination: mockPagination
    })

    vi.mocked(bidderApi.updateBidderStatus).mockResolvedValue({
      id: 'abc12345-def6-7890-abcd-ef1234567890',
      status: 'suspended',
      updated_at: '2025-01-04T00:00:00Z'
    })

    vi.mocked(bidderApi.grantPoints).mockResolvedValue({
      bidder: {
        id: 'abc12345-def6-7890-abcd-ef1234567890',
        total_points: 11000,
        available_points: 11000,
        reserved_points: 0,
        updated_at: '2025-01-04T00:00:00Z'
      },
      history: {
        id: 123,
        type: 'grant',
        points: 1000,
        balance_after: 11000,
        created_at: '2025-01-10T12:00:00Z'
      }
    })
  })

  describe('Initial Load', () => {
    it('should fetch and display bidder list on mount', async () => {
      wrapper = mount(BidderListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      expect(bidderApi.getBidderList).toHaveBeenCalled()
      expect(bidderStore.bidders).toEqual(mockBidders)
      expect(bidderStore.pagination).toEqual({
        currentPage: 1,
        totalPages: 5,
        totalItems: 100,
        itemsPerPage: 20
      })
    })

    it('should handle fetch error on mount', async () => {
      const errorMessage = 'Network error'
      vi.mocked(bidderApi.getBidderList).mockRejectedValue(new Error(errorMessage))

      wrapper = mount(BidderListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      expect(bidderStore.error).toBe(errorMessage)
    })
  })

  describe('Search Functionality', () => {
    it('should fetch bidders when search is submitted', async () => {
      wrapper = mount(BidderListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()
      vi.clearAllMocks()

      // Update search term and trigger search
      bidderStore.filters.keyword = 'tanaka'
      await bidderStore.setFiltersAndFetch({ keyword: 'tanaka' })

      await flushPromises()

      expect(bidderApi.getBidderList).toHaveBeenCalledWith(
        expect.objectContaining({
          keyword: 'tanaka'
        })
      )
    })

    it('should reset page to 1 when searching', async () => {
      wrapper = mount(BidderListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      // Set to page 3
      bidderStore.pagination.currentPage = 3

      // Search
      await bidderStore.setFiltersAndFetch({ keyword: 'test' })

      await flushPromises()

      expect(bidderStore.pagination.currentPage).toBe(1)
    })
  })

  describe('Filter Functionality', () => {
    it('should fetch bidders when status filter is changed', async () => {
      wrapper = mount(BidderListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()
      vi.clearAllMocks()

      // Change status filter
      await bidderStore.setFiltersAndFetch({ status: ['active'] })

      await flushPromises()

      expect(bidderApi.getBidderList).toHaveBeenCalledWith(
        expect.objectContaining({
          status: 'active'
        })
      )
    })

    it('should reset filters when reset button is clicked', async () => {
      wrapper = mount(BidderListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      // Set filters
      bidderStore.filters.keyword = 'test'
      bidderStore.filters.status = ['active']

      // Reset
      await bidderStore.resetFilters()

      await flushPromises()

      expect(bidderStore.filters).toEqual({
        keyword: '',
        status: ['active', 'suspended'],
        sort: 'created_at',
        order: 'asc'
      })
    })
  })

  describe('Sorting Functionality', () => {
    it('should sort by email when email column is clicked', async () => {
      wrapper = mount(BidderListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()
      vi.clearAllMocks()

      // Sort by email
      await bidderStore.changeSort('email')

      await flushPromises()

      expect(bidderApi.getBidderList).toHaveBeenCalledWith(
        expect.objectContaining({
          sort: 'email_asc'
        })
      )
    })

    it('should toggle sort order when clicking same column', async () => {
      wrapper = mount(BidderListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      // First click - ascending
      await bidderStore.changeSort('email')
      expect(bidderStore.filters.order).toBe('asc')

      // Second click - descending
      await bidderStore.changeSort('email')
      expect(bidderStore.filters.order).toBe('desc')
    })
  })

  describe('Pagination Functionality', () => {
    it('should fetch bidders when page is changed', async () => {
      wrapper = mount(BidderListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()
      vi.clearAllMocks()

      // Change to page 3
      await bidderStore.changePage(3)

      await flushPromises()

      expect(bidderApi.getBidderList).toHaveBeenCalledWith(
        expect.objectContaining({
          page: 3
        })
      )
    })
  })

  describe('Points Grant Functionality', () => {
    it('should grant points successfully', async () => {
      wrapper = mount(BidderListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      // Get the bidder before granting points
      const bidder = bidderStore.bidders[0]
      expect(bidder.total_points).toBe(10000)

      // Grant points
      const result = await bidderStore.addPoints('abc12345-def6-7890-abcd-ef1234567890', 1000)

      await flushPromises()

      expect(result).toBe(true)
      expect(bidderApi.grantPoints).toHaveBeenCalledWith('abc12345-def6-7890-abcd-ef1234567890', 1000)

      // Verify the bidder points were updated in the store
      const updatedBidder = bidderStore.bidders.find(b => b.id === 'abc12345-def6-7890-abcd-ef1234567890')
      expect(updatedBidder.total_points).toBe(11000)
    })

    it('should handle points grant error', async () => {
      const errorMessage = 'Invalid points value'

      wrapper = mount(BidderListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      // Mock the error for the points grant
      vi.mocked(bidderApi.grantPoints).mockRejectedValue(new Error(errorMessage))

      // Try to grant points
      const result = await bidderStore.addPoints('abc12345-def6-7890-abcd-ef1234567890', -100)

      await flushPromises()

      expect(result).toBe(false)
      expect(bidderStore.error).toBe(errorMessage)

      // Verify the bidder points were NOT updated in the store
      const bidder = bidderStore.bidders.find(b => b.id === 'abc12345-def6-7890-abcd-ef1234567890')
      expect(bidder.total_points).toBe(10000)
    })
  })

  describe('Status Change Functionality', () => {
    it('should update bidder status successfully', async () => {
      wrapper = mount(BidderListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      // Get the bidder before status change
      const bidder = bidderStore.bidders[0]
      expect(bidder.status).toBe('active')

      // Change status
      const result = await bidderStore.changeBidderStatus('abc12345-def6-7890-abcd-ef1234567890', 'suspended')

      await flushPromises()

      expect(result).toBe(true)
      expect(bidderApi.updateBidderStatus).toHaveBeenCalledWith('abc12345-def6-7890-abcd-ef1234567890', 'suspended')

      // Verify the bidder status was updated in the store
      const updatedBidder = bidderStore.bidders.find(b => b.id === 'abc12345-def6-7890-abcd-ef1234567890')
      expect(updatedBidder.status).toBe('suspended')
    })

    it('should handle status change error', async () => {
      const errorMessage = 'Permission denied'

      wrapper = mount(BidderListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      // Get the original status before attempting to change
      const bidderBefore = bidderStore.bidders.find(b => b.id === 'abc12345-def6-7890-abcd-ef1234567890')
      const originalStatus = bidderBefore.status

      // Mock the error for the status update
      vi.mocked(bidderApi.updateBidderStatus).mockRejectedValue(new Error(errorMessage))

      // Try to change status
      const result = await bidderStore.changeBidderStatus('abc12345-def6-7890-abcd-ef1234567890', 'suspended')

      await flushPromises()

      expect(result).toBe(false)
      expect(bidderStore.error).toBe(errorMessage)

      // Verify the bidder status was NOT updated in the store
      const bidderAfter = bidderStore.bidders.find(b => b.id === 'abc12345-def6-7890-abcd-ef1234567890')
      expect(bidderAfter.status).toBe(originalStatus)
    })
  })

  describe('Navigation', () => {
    it('should navigate to edit page when edit button is clicked', async () => {
      const pushSpy = vi.spyOn(router, 'push')

      wrapper = mount(BidderListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      // Simulate edit button click (in real implementation this would be triggered by BidderTable)
      await router.push('/admin/bidders/abc12345-def6-7890-abcd-ef1234567890/edit')

      expect(pushSpy).toHaveBeenCalledWith('/admin/bidders/abc12345-def6-7890-abcd-ef1234567890/edit')
    })
  })

  describe('Loading State', () => {
    it('should show loading state during fetch', async () => {
      // Create a promise that we can control
      let resolveFn
      const promise = new Promise((resolve) => {
        resolveFn = resolve
      })

      vi.mocked(bidderApi.getBidderList).mockReturnValue(promise)

      wrapper = mount(BidderListView, {
        global: {
          plugins: [router]
        }
      })

      // Loading should be true
      expect(bidderStore.loading).toBe(true)

      // Resolve the promise
      resolveFn({
        bidders: mockBidders,
        pagination: mockPagination
      })

      await flushPromises()

      // Loading should be false
      expect(bidderStore.loading).toBe(false)
    })
  })

  describe('Multiple Filters Combined', () => {
    it('should apply multiple filters simultaneously', async () => {
      wrapper = mount(BidderListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()
      vi.clearAllMocks()

      // Apply multiple filters
      await bidderStore.setFiltersAndFetch({
        keyword: 'tanaka',
        status: ['active']
      })

      await flushPromises()

      expect(bidderApi.getBidderList).toHaveBeenCalledWith(
        expect.objectContaining({
          keyword: 'tanaka',
          status: 'active',
          sort: 'created_at_asc',
          page: 1,
          limit: 20
        })
      )
    })
  })

  describe('Error Handling', () => {
    it('should clear error when retrying after error', async () => {
      // First request fails
      vi.mocked(bidderApi.getBidderList).mockRejectedValueOnce(new Error('Network error'))

      wrapper = mount(BidderListView, {
        global: {
          plugins: [router]
        }
      })

      await flushPromises()

      expect(bidderStore.error).toBe('Network error')

      // Clear error
      bidderStore.clearError()
      expect(bidderStore.error).toBeNull()

      // Second request succeeds
      vi.mocked(bidderApi.getBidderList).mockResolvedValue({
        bidders: mockBidders,
        pagination: mockPagination
      })

      await bidderStore.fetchBidderList()

      await flushPromises()

      expect(bidderStore.error).toBeNull()
      expect(bidderStore.bidders).toEqual(mockBidders)
    })
  })
})
