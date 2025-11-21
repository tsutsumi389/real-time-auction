import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import AdminTable from './AdminTable.vue'
import AdminRoleBadge from './AdminRoleBadge.vue'
import AdminStatusBadge from './AdminStatusBadge.vue'

describe('AdminTable', () => {
  const mockAdmins = [
    {
      id: 1,
      email: 'admin1@example.com',
      role: 'system_admin',
      status: 'active',
      created_at: '2025-01-01T00:00:00Z',
    },
    {
      id: 2,
      email: 'admin2@example.com',
      role: 'auctioneer',
      status: 'inactive',
      created_at: '2025-01-02T00:00:00Z',
    },
  ]

  describe('Rendering', () => {
    it('should render loading state', () => {
      const wrapper = mount(AdminTable, {
        props: {
          admins: [],
          loading: true,
        },
      })

      expect(wrapper.text()).toContain('読み込み中...')
      expect(wrapper.find('table').exists()).toBe(false)
    })

    it('should render empty state when no admins', () => {
      const wrapper = mount(AdminTable, {
        props: {
          admins: [],
          loading: false,
        },
      })

      expect(wrapper.text()).toContain('管理者が見つかりませんでした')
      expect(wrapper.find('table').exists()).toBe(false)
    })

    it('should render table with admins', () => {
      const wrapper = mount(AdminTable, {
        props: {
          admins: mockAdmins,
          loading: false,
        },
        global: {
          stubs: {
            AdminRoleBadge: true,
            AdminStatusBadge: true,
          },
        },
      })

      const table = wrapper.find('table')
      expect(table.exists()).toBe(true)

      const rows = wrapper.findAll('tbody tr')
      expect(rows).toHaveLength(2)
    })

    it('should render admin data correctly', () => {
      const wrapper = mount(AdminTable, {
        props: {
          admins: [mockAdmins[0]],
          loading: false,
        },
        global: {
          stubs: {
            AdminRoleBadge: true,
            AdminStatusBadge: true,
          },
        },
      })

      const row = wrapper.find('tbody tr')
      expect(row.text()).toContain('1')
      expect(row.text()).toContain('admin1@example.com')
    })
  })

  describe('Sorting', () => {
    it('should display sort indicator for sorted column', () => {
      const wrapper = mount(AdminTable, {
        props: {
          admins: mockAdmins,
          loading: false,
          sortField: 'email',
          sortOrder: 'asc',
        },
        global: {
          stubs: {
            AdminRoleBadge: true,
            AdminStatusBadge: true,
          },
        },
      })

      const emailHeader = wrapper.findAll('th')[1]
      expect(emailHeader.text()).toContain('メールアドレス')
      expect(emailHeader.text()).toContain('↑')
    })

    it('should display descending sort indicator', () => {
      const wrapper = mount(AdminTable, {
        props: {
          admins: mockAdmins,
          loading: false,
          sortField: 'id',
          sortOrder: 'desc',
        },
        global: {
          stubs: {
            AdminRoleBadge: true,
            AdminStatusBadge: true,
          },
        },
      })

      const idHeader = wrapper.findAll('th')[0]
      expect(idHeader.text()).toContain('↓')
    })

    it('should emit sort event when clicking column header', async () => {
      const wrapper = mount(AdminTable, {
        props: {
          admins: mockAdmins,
          loading: false,
        },
        global: {
          stubs: {
            AdminRoleBadge: true,
            AdminStatusBadge: true,
          },
        },
      })

      const emailHeader = wrapper.findAll('th')[1]
      await emailHeader.trigger('click')

      expect(wrapper.emitted('sort')).toBeTruthy()
      expect(wrapper.emitted('sort')[0]).toEqual(['email'])
    })

    it('should emit sort event for each sortable column', async () => {
      const wrapper = mount(AdminTable, {
        props: {
          admins: mockAdmins,
          loading: false,
        },
        global: {
          stubs: {
            AdminRoleBadge: true,
            AdminStatusBadge: true,
          },
        },
      })

      const headers = wrapper.findAll('th')
      const sortableColumns = ['id', 'email', 'role', 'status', 'created_at']

      for (let i = 0; i < sortableColumns.length; i++) {
        await headers[i].trigger('click')
        expect(wrapper.emitted('sort')[i]).toEqual([sortableColumns[i]])
      }
    })
  })

  describe('Actions', () => {
    it('should emit edit event when clicking edit button', async () => {
      const wrapper = mount(AdminTable, {
        props: {
          admins: mockAdmins,
          loading: false,
        },
        global: {
          stubs: {
            AdminRoleBadge: true,
            AdminStatusBadge: true,
          },
        },
      })

      const editButton = wrapper.findAll('button').find((btn) => btn.text() === '編集')
      await editButton.trigger('click')

      expect(wrapper.emitted('edit')).toBeTruthy()
      expect(wrapper.emitted('edit')[0]).toEqual([1])
    })

    it('should emit status-change event when clicking status button', async () => {
      const wrapper = mount(AdminTable, {
        props: {
          admins: mockAdmins,
          loading: false,
        },
        global: {
          stubs: {
            AdminRoleBadge: true,
            AdminStatusBadge: true,
          },
        },
      })

      const statusButton = wrapper.findAll('button').find((btn) => btn.text() === '停止')
      await statusButton.trigger('click')

      expect(wrapper.emitted('status-change')).toBeTruthy()
      expect(wrapper.emitted('status-change')[0]).toEqual([mockAdmins[0]])
    })

    it('should show "停止" button for active admin', () => {
      const wrapper = mount(AdminTable, {
        props: {
          admins: [mockAdmins[0]], // active
          loading: false,
        },
        global: {
          stubs: {
            AdminRoleBadge: true,
            AdminStatusBadge: true,
          },
        },
      })

      const statusButton = wrapper.findAll('button').find((btn) => btn.text() === '停止')
      expect(statusButton).toBeTruthy()
    })

    it('should show "復活" button for inactive admin', () => {
      const wrapper = mount(AdminTable, {
        props: {
          admins: [mockAdmins[1]], // inactive
          loading: false,
        },
        global: {
          stubs: {
            AdminRoleBadge: true,
            AdminStatusBadge: true,
          },
        },
      })

      const statusButton = wrapper.findAll('button').find((btn) => btn.text() === '復活')
      expect(statusButton).toBeTruthy()
    })

    it('should have correct aria-label for edit button', () => {
      const wrapper = mount(AdminTable, {
        props: {
          admins: [mockAdmins[0]],
          loading: false,
        },
        global: {
          stubs: {
            AdminRoleBadge: true,
            AdminStatusBadge: true,
          },
        },
      })

      const editButton = wrapper.findAll('button').find((btn) => btn.text() === '編集')
      expect(editButton.attributes('aria-label')).toBe('admin1@example.comのアカウントを編集')
    })

    it('should have correct aria-label for status button', () => {
      const wrapper = mount(AdminTable, {
        props: {
          admins: [mockAdmins[0]],
          loading: false,
        },
        global: {
          stubs: {
            AdminRoleBadge: true,
            AdminStatusBadge: true,
          },
        },
      })

      const statusButton = wrapper.findAll('button').find((btn) => btn.text() === '停止')
      expect(statusButton.attributes('aria-label')).toBe('admin1@example.comのアカウントを停止')
    })
  })

  describe('Date Formatting', () => {
    it('should format date correctly', () => {
      const wrapper = mount(AdminTable, {
        props: {
          admins: [mockAdmins[0]],
          loading: false,
        },
        global: {
          stubs: {
            AdminRoleBadge: true,
            AdminStatusBadge: true,
          },
        },
      })

      const dateCell = wrapper.findAll('td')[4]
      // 日本語ロケールでフォーマットされた日付を確認
      expect(dateCell.text()).toMatch(/\d{4}\/\d{2}\/\d{2}/)
    })
  })

  describe('Badge Components', () => {
    it('should render AdminRoleBadge with correct props', () => {
      const wrapper = mount(AdminTable, {
        props: {
          admins: [mockAdmins[0]],
          loading: false,
        },
        global: {
          components: {
            AdminRoleBadge,
            AdminStatusBadge,
          },
        },
      })

      const roleBadge = wrapper.findComponent(AdminRoleBadge)
      expect(roleBadge.exists()).toBe(true)
      expect(roleBadge.props('role')).toBe('system_admin')
    })

    it('should render AdminStatusBadge with correct props', () => {
      const wrapper = mount(AdminTable, {
        props: {
          admins: [mockAdmins[0]],
          loading: false,
        },
        global: {
          components: {
            AdminRoleBadge,
            AdminStatusBadge,
          },
        },
      })

      const statusBadge = wrapper.findComponent(AdminStatusBadge)
      expect(statusBadge.exists()).toBe(true)
      expect(statusBadge.props('status')).toBe('active')
    })
  })
})
