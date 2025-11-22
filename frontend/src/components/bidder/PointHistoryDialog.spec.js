import { describe, it, expect, vi, beforeEach } from 'vitest'
import { mount, flushPromises } from '@vue/test-utils'
import { setActivePinia, createPinia } from 'pinia'
import PointHistoryDialog from './PointHistoryDialog.vue'
import { useBidderStore } from '@/stores/bidder'
import * as bidderApi from '@/services/bidderApi'

// Mock bidderApi
vi.mock('@/services/bidderApi')

// Mock Dialog component
const mockDialogComponent = {
  template: '<div v-if="modelValue" role="dialog""><slot /></div>',
  props: ['modelValue']
}

// Mock PointTypeBadge component
const mockPointTypeBadge = {
  template: '<span class="point-type-badge-mock">{{ type }}</span>',
  props: ['type']
}

describe('PointHistoryDialog', () => {
  let bidderStore

  const mockBidder = {
    id: 'abc12345-def6-7890-abcd-ef1234567890',
    email: 'bidder1@example.com',
    display_name: '田中太郎',
  }

  const mockHistoryResponse = {
    bidder: mockBidder,
    history: [
      {
        id: 1,
        type: 'grant',
        amount: 1000,
        balance_before: 10000,
        balance_after: 11000,
        auction_id: null,
        auction_title: null,
        created_at: '2025-01-10T12:00:00Z',
      },
      {
        id: 2,
        type: 'reserve',
        amount: -5000,
        balance_before: 11000,
        balance_after: 6000,
        auction_id: 5,
        auction_title: '競走馬オークション #5',
        created_at: '2025-01-09T10:30:00Z',
      },
    ],
    pagination: {
      page: 1,
      total_pages: 1,
      total: 2,
      limit: 10,
    },
  }

  beforeEach(() => {
    setActivePinia(createPinia())
    bidderStore = useBidderStore()

    // Mock API responses
    vi.spyOn(bidderStore, 'fetchPointHistory').mockResolvedValue(mockHistoryResponse)
  })

  describe('Rendering', () => {
    it('should not render when closed', () => {
      const wrapper = mount(PointHistoryDialog, {
        props: {
          modelValue: false,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
            PointTypeBadge: mockPointTypeBadge,
          },
        },
      })

      expect(wrapper.find('[role="dialog"]').exists()).toBe(false)
    })

    it('should render when open', async () => {
      const wrapper = mount(PointHistoryDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
            PointTypeBadge: mockPointTypeBadge,
          },
        },
      })

      await flushPromises()

      expect(wrapper.find('[role="dialog"]').exists()).toBe(true)
    })

    it('should display bidder information', async () => {
      const wrapper = mount(PointHistoryDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
            PointTypeBadge: mockPointTypeBadge,
          },
        },
      })

      await flushPromises()

      expect(wrapper.text()).toContain('bidder1@example.com')
      expect(wrapper.text()).toContain('田中太郎')
    })

    it('should fetch and display point history on open', async () => {
      const wrapper = mount(PointHistoryDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
            PointTypeBadge: mockPointTypeBadge,
          },
        },
      })

      await flushPromises()

      expect(bidderStore.fetchPointHistory).toHaveBeenCalledWith(mockBidder.id, 1, 10)
      expect(wrapper.text()).toContain('1,000')
      expect(wrapper.text()).toContain('競走馬オークション #5')
    })

    it('should display loading state while fetching', async () => {
      let resolveFn
      const promise = new Promise((resolve) => {
        resolveFn = resolve
      })

      vi.spyOn(bidderStore, 'fetchPointHistory').mockReturnValue(promise)

      const wrapper = mount(PointHistoryDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
            PointTypeBadge: mockPointTypeBadge,
          },
        },
      })

      // Should show loading
      expect(wrapper.text()).toContain('読み込み中')

      resolveFn(mockHistoryResponse)
      await flushPromises()

      // Loading should be gone
      expect(wrapper.text()).not.toContain('読み込み中')
    })

    it('should display empty state when no history', async () => {
      vi.spyOn(bidderStore, 'fetchPointHistory').mockResolvedValue({
        bidder: mockBidder,
        history: [],
        pagination: {
          page: 1,
          total_pages: 0,
          total: 0,
          limit: 10,
        },
      })

      const wrapper = mount(PointHistoryDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
            PointTypeBadge: mockPointTypeBadge,
          },
        },
      })

      await flushPromises()

      expect(wrapper.text()).toContain('履歴がありません')
    })
  })

  describe('History Details', () => {
    it('should display point type badges', async () => {
      const wrapper = mount(PointHistoryDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
            PointTypeBadge: mockPointTypeBadge,
          },
        },
      })

      await flushPromises()

      expect(wrapper.text()).toContain('grant')
      expect(wrapper.text()).toContain('reserve')
    })

    it('should display points with sign and comma separator', async () => {
      const wrapper = mount(PointHistoryDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
            PointTypeBadge: mockPointTypeBadge,
          },
        },
      })

      await flushPromises()

      expect(wrapper.text()).toContain('+1,000')
      expect(wrapper.text()).toContain('-5,000')
    })

    it('should display auction title when available', async () => {
      const wrapper = mount(PointHistoryDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
            PointTypeBadge: mockPointTypeBadge,
          },
        },
      })

      await flushPromises()

      expect(wrapper.text()).toContain('競走馬オークション #5')
    })

    it('should display "-" when no auction', async () => {
      const wrapper = mount(PointHistoryDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
            PointTypeBadge: mockPointTypeBadge,
          },
        },
      })

      await flushPromises()

      const rows = wrapper.findAll('tbody tr')
      expect(rows[0].text()).toContain('-')
    })

    it('should format date correctly', async () => {
      const wrapper = mount(PointHistoryDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
            PointTypeBadge: mockPointTypeBadge,
          },
        },
      })

      await flushPromises()

      expect(wrapper.text()).toMatch(/\d{4}\/\d{2}\/\d{2}/)
    })
  })

  describe('Pagination', () => {
    it('should display pagination when multiple pages', async () => {
      vi.spyOn(bidderStore, 'fetchPointHistory').mockResolvedValue({
        ...mockHistoryResponse,
        pagination: {
          page: 1,
          total_pages: 5,
          total: 50,
          limit: 10,
        },
      })

      const wrapper = mount(PointHistoryDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
            PointTypeBadge: mockPointTypeBadge,
          },
        },
      })

      await flushPromises()

      expect(wrapper.text()).toContain('1-10')
      expect(wrapper.text()).toContain('50')
    })

    it('should fetch next page when next button is clicked', async () => {
      vi.spyOn(bidderStore, 'fetchPointHistory').mockResolvedValue({
        ...mockHistoryResponse,
        pagination: {
          page: 1,
          total_pages: 5,
          total: 50,
          limit: 10,
        },
      })

      const wrapper = mount(PointHistoryDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
            PointTypeBadge: mockPointTypeBadge,
          },
        },
      })

      await flushPromises()
      vi.clearAllMocks()

      const nextButton = wrapper.findAll('button').find((btn) => btn.text().includes('次へ'))
      await nextButton.trigger('click')

      await flushPromises()

      expect(bidderStore.fetchPointHistory).toHaveBeenCalledWith(mockBidder.id, 2, 10)
    })
  })

  describe('Error Handling', () => {
    it('should display error message when fetch fails', async () => {
      vi.spyOn(bidderStore, 'fetchPointHistory').mockRejectedValue(new Error('Network error'))

      const wrapper = mount(PointHistoryDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
            PointTypeBadge: mockPointTypeBadge,
          },
        },
      })

      await flushPromises()

      expect(wrapper.text()).toContain('履歴がありません')
    })
  })

  describe('Dialog Closing', () => {
    it('should emit update:modelValue event when close button is clicked', async () => {
      const wrapper = mount(PointHistoryDialog, {
        props: {
          modelValue: true,
          bidder: mockBidder,
        },
        global: {
          stubs: {
            
            PointTypeBadge: mockPointTypeBadge,
          },
        },
      })

      await flushPromises()

      const closeButton = wrapper.findAll('button').find((btn) => btn.text().includes('閉じる'))
      await closeButton.trigger('click')

      expect(wrapper.emitted('update:modelValue')).toBeTruthy()
      expect(wrapper.emitted('update:modelValue')[0]).toEqual([false])
    })
  })
})
