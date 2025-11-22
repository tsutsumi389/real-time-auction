import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import BidderTable from './BidderTable.vue'

// Mock UI components
const mockBadgeComponent = {
  template: '<span class="badge-mock"><slot /></span>',
  props: ['status', 'type']
}

describe('BidderTable', () => {
  const mockBidders = [
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
      display_name: null,
      status: 'suspended',
      total_points: 5000,
      created_at: '2025-01-02T00:00:00Z',
    },
  ]

  describe('Rendering', () => {
    it('should render loading state', () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: [],
          loading: true,
        },
      })

      expect(wrapper.text()).toContain('読み込み中...')
      expect(wrapper.find('table').exists()).toBe(false)
    })

    it('should render empty state when no bidders', () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: [],
          loading: false,
        },
      })

      expect(wrapper.text()).toContain('入札者が見つかりませんでした')
      expect(wrapper.find('table').exists()).toBe(false)
    })

    it('should render table with bidders', () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: mockBidders,
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const table = wrapper.find('table')
      expect(table.exists()).toBe(true)

      const rows = wrapper.findAll('tbody tr')
      expect(rows).toHaveLength(2)
    })

    it('should render bidder data correctly', () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: [mockBidders[0]],
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const row = wrapper.find('tbody tr')
      expect(row.text()).toContain('abc12345')
      expect(row.text()).toContain('bidder1@example.com')
      expect(row.text()).toContain('田中太郎')
      expect(row.text()).toContain('10,000')
    })

    it('should display "未設定" for null display_name', () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: [mockBidders[1]],
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const row = wrapper.find('tbody tr')
      expect(row.text()).toContain('（未設定）')
    })

    it('should display shortened UUID (first 8 characters)', () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: [mockBidders[0]],
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const row = wrapper.find('tbody tr')
      expect(row.text()).toContain('abc12345')
      expect(row.text()).not.toContain('abc12345-def6-7890-abcd-ef1234567890')
    })

    it('should format points with commas', () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: [mockBidders[0]],
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const row = wrapper.find('tbody tr')
      expect(row.text()).toContain('10,000')
    })
  })

  describe('Sorting', () => {
    it('should display sort indicator for sorted column', () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: mockBidders,
          loading: false,
          sortField: 'email',
          sortOrder: 'asc',
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const emailHeader = wrapper.findAll('th')[1]
      expect(emailHeader.text()).toContain('メールアドレス')
      expect(emailHeader.text()).toContain('↑')
    })

    it('should display descending sort indicator', () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: mockBidders,
          loading: false,
          sortField: 'id',
          sortOrder: 'desc',
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const idHeader = wrapper.findAll('th')[0]
      expect(idHeader.text()).toContain('↓')
    })

    it('should emit sort event when clicking column header', async () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: mockBidders,
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const emailHeader = wrapper.findAll('th')[1]
      await emailHeader.trigger('click')

      expect(wrapper.emitted('sort')).toBeTruthy()
      expect(wrapper.emitted('sort')[0]).toEqual(['email'])
    })
  })

  describe('Actions', () => {
    it('should emit edit event when clicking edit button', async () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: mockBidders,
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const editButton = wrapper.findAll('button').find((btn) => btn.text().includes('詳細'))
      await editButton.trigger('click')

      expect(wrapper.emitted('edit')).toBeTruthy()
      expect(wrapper.emitted('edit')[0]).toEqual(['abc12345-def6-7890-abcd-ef1234567890'])
    })

    it('should emit grant-points event when clicking grant points button', async () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: mockBidders,
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const grantButton = wrapper.findAll('button').find((btn) => btn.text().includes('pt付与'))
      await grantButton.trigger('click')

      expect(wrapper.emitted('grant-points')).toBeTruthy()
      expect(wrapper.emitted('grant-points')[0]).toEqual([mockBidders[0]])
    })

    it('should emit view-history event when clicking history button', async () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: mockBidders,
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const historyButton = wrapper.findAll('button').find((btn) => btn.text().includes('履歴'))
      await historyButton.trigger('click')

      expect(wrapper.emitted('view-history')).toBeTruthy()
      expect(wrapper.emitted('view-history')[0]).toEqual([mockBidders[0]])
    })

    it('should emit status-change event when clicking status button', async () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: mockBidders,
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const statusButton = wrapper.findAll('button').find((btn) => btn.text() === '停止')
      await statusButton.trigger('click')

      expect(wrapper.emitted('status-change')).toBeTruthy()
      expect(wrapper.emitted('status-change')[0]).toEqual([mockBidders[0]])
    })

    it('should show "停止" button for active bidder', () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: [mockBidders[0]], // active
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const statusButton = wrapper.findAll('button').find((btn) => btn.text() === '停止')
      expect(statusButton).toBeTruthy()
    })

    it('should show "復活" button for suspended bidder', () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: [mockBidders[1]], // suspended
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const statusButton = wrapper.findAll('button').find((btn) => btn.text() === '復活')
      expect(statusButton).toBeTruthy()
    })

    it('should not show grant points button for deleted bidder', () => {
      const deletedBidder = {
        ...mockBidders[0],
        status: 'deleted',
      }

      const wrapper = mount(BidderTable, {
        props: {
          bidders: [deletedBidder],
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const grantButton = wrapper.findAll('button').find((btn) => btn.text().includes('pt付与'))
      expect(grantButton).toBeFalsy()
    })

    it('should not show status change button for deleted bidder', () => {
      const deletedBidder = {
        ...mockBidders[0],
        status: 'deleted',
      }

      const wrapper = mount(BidderTable, {
        props: {
          bidders: [deletedBidder],
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const statusButton = wrapper.findAll('button').find((btn) => btn.text() === '停止' || btn.text() === '復活')
      expect(statusButton).toBeFalsy()
    })
  })

  describe('Date Formatting', () => {
    it('should format date correctly', () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: [mockBidders[0]],
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const dateCell = wrapper.findAll('td')[5]
      // 日本語ロケールでフォーマットされた日付を確認
      expect(dateCell.text()).toMatch(/\d{4}\/\d{2}\/\d{2}/)
    })
  })

  describe('Accessibility', () => {
    it('should have correct aria-label for edit button', () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: [mockBidders[0]],
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const editButton = wrapper.findAll('button').find((btn) => btn.text().includes('詳細'))
      expect(editButton.attributes('aria-label')).toContain('詳細を表示')
    })

    it('should have correct aria-label for grant points button', () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: [mockBidders[0]],
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const grantButton = wrapper.findAll('button').find((btn) => btn.text().includes('pt付与'))
      expect(grantButton.attributes('aria-label')).toContain('ポイントを付与')
    })
  })

  describe('UUID Tooltip', () => {
    it('should have title attribute with full UUID on UUID cell', () => {
      const wrapper = mount(BidderTable, {
        props: {
          bidders: [mockBidders[0]],
          loading: false,
        },
        global: {
          stubs: {
            Badge: mockBadgeComponent,
          },
        },
      })

      const idCell = wrapper.findAll('td')[0]
      expect(idCell.attributes('title')).toBe('abc12345-def6-7890-abcd-ef1234567890')
    })
  })
})
